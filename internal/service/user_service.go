package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/constants"
	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	"github.com/imkarthi24/sf-backend/internal/model/models"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/service/email"
	"github.com/loop-kar/pixie/util"
	"github.com/thoas/go-funk"
)

type UserService interface {
	SaveUser(*context.Context, requestModel.User) *errs.XError
	UpdateUser(*context.Context, requestModel.User, uint) *errs.XError
	Login(*context.Context, requestModel.Login) (string, *errs.XError)
	Get(*context.Context, uint) (*responseModel.User, *errs.XError)
	GetAllUsers(ctx *context.Context, search string) ([]responseModel.User, *errs.XError) // Used in Browse for Users
	Delete(*context.Context, uint) *errs.XError
	ForgotPassword(*context.Context, string) *errs.XError
	ResetPassword(*context.Context, string, string) *errs.XError
	RefreshToken(*context.Context) (string, *errs.XError)
	GetUsersForAutoComplete(ctx *context.Context, name string, role []string) ([]responseModel.UserAutoComplete, *errs.XError)
	UpdateChannel(*context.Context, uint, uint) *errs.XError
	SwitchUserChannel(ctx *context.Context, id uint) (string, *errs.XError)

	//User Config
	SaveUserConfig(ctx *context.Context, config requestModel.UserConfig) *errs.XError
	UpdateUserConfig(ctx *context.Context, config requestModel.UserConfig, id uint) *errs.XError
	GetUserConfig(ctx *context.Context, userId uint) (*entities.UserConfig, *errs.XError)

	//UserChannelDetail
	AddUserChannelDetail(ctx *context.Context, userChannleDet []requestModel.UserChannelDetail) *errs.XError
	UpdateUserChannelDetail(ctx *context.Context, userChannleDet requestModel.UserChannelDetail, id uint) *errs.XError
	DeleteUserChannelDetail(ctx *context.Context, id uint) *errs.XError
	GetUserChannelDetails(ctx *context.Context, userId uint) ([]responseModel.UserChannelDetail, *errs.XError)
}

type userService struct {
	userRepo    repository.UserRepository
	channelRepo repository.ChannelRepository
	mapper      mapper.Mapper
	config      config.AppConfig
	respMapper  mapper.ResponseMapper
	emailSvc    email.EmailService
	notifSvc    NotificationService
}

func ProvideUserService(repo repository.UserRepository, channelRepo repository.ChannelRepository, mapper mapper.Mapper, config config.AppConfig, respMapper mapper.ResponseMapper, emailSvc email.EmailService, notifSvc NotificationService) UserService {
	return userService{
		userRepo:    repo,
		channelRepo: channelRepo,
		mapper:      mapper,
		config:      config,
		respMapper:  respMapper,
		emailSvc:    emailSvc,
		notifSvc:    notifSvc,
	}
}

func (svc userService) SaveUser(ctx *context.Context, user requestModel.User) *errs.XError {

	dbUser, err := svc.mapper.User(user)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save User", err)
	}

	generatedPassword := util.GeneratePassword()
	dbUser.Password = util.HashPassword(generatedPassword, svc.config.Server.SecretKey)

	errr := svc.userRepo.Create(ctx, dbUser)
	if errr != nil {
		return errr
	}

	// Adding current channel to UserChannelDetails
	errr = svc.AddUserChannelDetail(ctx, []requestModel.UserChannelDetail{
		{
			UserID:    dbUser.ID,
			IsActive:  true,
			ChannelId: dbUser.ChannelId,
		},
	})
	if errr != nil {
		return errr
	}

	// Set reset link for "set up your account" email (async send via notification task)
	resetString := util.GenerateRandom(12)
	if errr = svc.userRepo.SetPasswordReset(ctx, dbUser.ID, resetString); errr != nil {
		return errr
	}

	// Enqueue email asynchronously; mail failure will not affect DB (notification task sends later)
	return svc.sendUserCreatedEmail(ctx, *dbUser, resetString)
}

func (svc userService) UpdateUser(ctx *context.Context, user requestModel.User, id uint) *errs.XError {
	dbUser, err := svc.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	updateUser := dbUser

	//Updating only the necessary fields
	updateUser.Model = &entities.Model{ID: id, IsActive: user.IsActive}
	updateUser.FirstName = user.FirstName
	updateUser.LastName = user.LastName
	updateUser.Extension = user.Extension
	updateUser.PhoneNumber = user.PhoneNumber
	updateUser.Email = user.Email
	updateUser.Role = entities.RoleType(user.Role)
	updateUser.IsLoginDisabled = user.IsLoginDisabled
	updateUser.IsLoggedIn = user.IsLoggedIn
	updateUser.ResetPasswordString = user.ResetPasswordString

	errr := svc.userRepo.Update(ctx, updateUser)
	if errr != nil {
		return errr
	}

	//Updating the UserChannelDetails seperately
	for _, chDet := range user.UserChannelDetails {
		chDet.UserID = id
		errr := svc.UpdateUserChannelDetail(ctx, chDet, chDet.ID)
		if errr != nil {
			return errr
		}
	}
	return nil

}

func (svc userService) Login(ctx *context.Context, login requestModel.Login) (string, *errs.XError) {
	user, err := svc.userRepo.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return "", err
	}

	//login disabled
	if user.IsLoginDisabled {
		return "", errs.NewXError(errs.INVALID, errs.LOGIN_DISABLED, nil)

	}

	//User not present
	if user.Email != login.Email || !user.IsActive {
		return "", errs.NewXError(errs.INVALID, errs.INVALID_USER, nil)
	}

	if !util.IsPasswordMatching(login.Password, user.Password, svc.config.Server.SecretKey) {
		user.LoginFailureCounter = user.LoginFailureCounter + 1
		if user.LoginFailureCounter > 5 {
			user.IsLoginDisabled = true
		}

		err := svc.userRepo.Update(ctx, user)
		if err != nil {
			return "", err
		}

		return "", errs.NewXError(errs.INVALID, errs.INVALID_CREDS, nil)
	}

	//Login success
	loginTime := util.GetLocalTime()
	user.LastLoginTime = &loginTime
	user.IsLoggedIn = true
	user.LoginFailureCounter = 0

	err = svc.userRepo.Update(ctx, user)
	if err != nil {
		return "", err
	}

	accessibleChannelIds := make([]uint, 0)
	funk.ForEach(user.UserChannelDetails, func(x entities.UserChannelDetail) {
		accessibleChannelIds = append(accessibleChannelIds, x.UserChannelID)
	})

	channel, err := svc.channelRepo.Get(ctx, accessibleChannelIds[0])
	if err != nil {
		return "", err
	}

	jwtResponse := models.Session{
		Email:                 user.Email,
		Role:                  user.Role,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		UserId:                &user.ID,
		ChannelId:             channel.ID,
		ChannelName:           channel.Name,
		AccessibleLocationIds: accessibleChannelIds,
	}

	jwtToken, errr := util.GenerateJWT(svc.config.Server.JwtSecretKey, svc.config.Server.JwtExpiryMinutes, util.StructToMap(jwtResponse))
	if errr != nil {
		return "", errs.NewXError(errs.INTERNAL, errs.JWT_ERROR, errr)

	}

	return jwtToken, nil
}

func (svc userService) GetAllUsers(ctx *context.Context, search string) ([]responseModel.User, *errs.XError) {

	users, err := svc.userRepo.GetAllUsers(ctx, search)
	if err != nil {
		return nil, err
	}

	return svc.respMapper.UserBrowse(users), nil

}

func (svc userService) Get(ctx *context.Context, id uint) (*responseModel.User, *errs.XError) {

	user, err := svc.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedUser, mapErr := svc.respMapper.User(user)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map User data", mapErr)
	}
	mappedUser.UserChannelDetails, err = svc.GetUserChannelDetails(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return mappedUser, nil

}

func (svc userService) Delete(ctx *context.Context, id uint) *errs.XError {

	err := svc.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (svc userService) sendUserCreatedEmail(ctx *context.Context, user entities.User, resetString string) *errs.XError {

	siteUrl := utils.GetSiteURL(svc.config.Site)
	passwordResetPath := constants.PASSWORD_RESET_UI_PATH
	resetLink := fmt.Sprintf("%s%s?resetString=%s", siteUrl, passwordResetPath, resetString)

	organization := "Stitchfolio"
	if user.ChannelId != 0 {
		channel, err := svc.channelRepo.Get(ctx, user.ChannelId)
		if err == nil {
			organization = channel.Name
		}
	}

	var templateName string
	subject := "Welcome to Stitchfolio"

	if user.Role == entities.SUPERADMIN {
		templateName = constants.SUPER_ADMIN_ONBOARD_HTML_TEMPLATE
		subject = "Welcome to Stitchfolio - Super Admin Onboarding"
	} else {
		templateName = constants.EMPLOYEE_CREATION_HTML_TEMPLATE
		subject = fmt.Sprintf("Welcome to %s", organization)
	}

	recipients := []string{user.Email}
	templateMap := map[string]string{
		"**name**":         user.FirstName,
		"**email**":        user.Email,
		"**organization**": organization,
		"**reset_link**":   resetLink,
	}
	if user.Role != entities.SUPERADMIN {
		templateMap["**department**"] = user.Department
	}

	notif := requestModel.EmaiNotification{
		Notification: &requestModel.Notification{
			SourceEntity: "user",
			EntityId:     user.ID,
		},
		EmailContent: &email.EmailContent{
			To:                   recipients,
			Subject:              subject,
			HtmlTemplateFileName: &templateName,
			TemplateValueMap:     templateMap,
		},
	}
	return svc.notifSvc.CreateEmailNotification(ctx, notif)
}

func (svc userService) sendPasswordResetMail(ctx *context.Context, user entities.User, resetString string) *errs.XError {

	siteUrl := utils.GetSiteURL(svc.config.Site)
	passwordResetPath := constants.PASSWORD_RESET_UI_PATH
	resetLink := fmt.Sprintf("%s%s?resetString=%s", siteUrl, passwordResetPath, resetString)

	templateName := constants.FORGOT_PASSWORD_HTML_TEMPLATE
	subject := "Reset Your Password - Stitchfolio"
	recipients := []string{user.Email}
	templateMap := map[string]string{
		"**name**":       user.FirstName,
		"**reset_link**": resetLink,
	}

	notif := requestModel.EmaiNotification{
		Notification: &requestModel.Notification{
			SourceEntity: "user",
			EntityId:     user.ID,
		},
		EmailContent: &email.EmailContent{
			To:                   recipients,
			Subject:              subject,
			HtmlTemplateFileName: &templateName,
			TemplateValueMap:     templateMap,
		},
	}
	return svc.notifSvc.CreateEmailNotification(ctx, notif)
}

func (svc userService) sendPasswordResetSuccessMail(ctx *context.Context, user entities.User, resetString string) *errs.XError {

	templateName := constants.PASSWORD_RESET_CONFIRMATION_HTML
	subject := "Password Reset Successful - Stitchfolio"
	recipients := []string{user.Email}
	templateMap := map[string]string{
		"**name**": user.FirstName,
	}

	notif := requestModel.EmaiNotification{
		Notification: &requestModel.Notification{
			SourceEntity: "user",
			EntityId:     user.ID,
		},
		EmailContent: &email.EmailContent{
			To:                   recipients,
			Subject:              subject,
			HtmlTemplateFileName: &templateName,
			TemplateValueMap:     templateMap,
		},
	}
	return svc.notifSvc.CreateEmailNotification(ctx, notif)
}

func (svc userService) ForgotPassword(ctx *context.Context, mail string) *errs.XError {

	user, err := svc.userRepo.GetUserByEmail(ctx, mail)
	if err != nil {
		return err
	}

	if user.Model == nil {
		return errs.NewXError(errs.INVALID_USER, "Invalid User email", nil)
	}

	resetString := util.GenerateRandom(12)

	err = svc.userRepo.SetPasswordReset(ctx, user.ID, resetString)
	if err != nil {
		return err
	}

	// Enqueue email asynchronously; mail failure will not affect DB
	return svc.sendPasswordResetMail(ctx, *user, resetString)
}

func (svc userService) ResetPassword(ctx *context.Context, resetString, password string) *errs.XError {

	user, err := svc.userRepo.ResetPassword(ctx, resetString, util.HashPassword(password, svc.config.Server.SecretKey))
	if err != nil {
		return err
	}

	// Enqueue confirmation email asynchronously; mail failure will not affect DB
	return svc.sendPasswordResetSuccessMail(ctx, *user, resetString)
}

func (svc userService) RefreshToken(ctx *context.Context) (string, *errs.XError) {

	session := utils.GetSession(ctx)
	if session == nil {
		return "", errs.NewXError(errs.SMTP_ERROR, "Unable to get user session", nil)

	}

	user, err := svc.Get(ctx, *session.UserId)
	if err != nil {
		return "", err
	}

	channel, err := svc.channelRepo.Get(ctx, session.ChannelId)
	if err != nil {
		return "", err
	}

	accessibleLocationIds := make([]uint, 0)
	funk.ForEach(user.UserChannelDetails, func(x responseModel.UserChannelDetail) {
		accessibleLocationIds = append(accessibleLocationIds, x.ChannelID)
	})

	jwtResponse := models.Session{
		Email:                 user.Email,
		Role:                  session.Role,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		UserId:                &user.ID,
		ChannelId:             channel.ChannelId,
		ChannelName:           channel.Name,
		AccessibleLocationIds: accessibleLocationIds,
	}

	jwtToken, errr := util.GenerateJWT(svc.config.Server.JwtSecretKey, svc.config.Server.JwtExpiryMinutes, util.StructToMap(jwtResponse))
	if errr != nil {
		return "", errs.NewXError(errs.INTERNAL, errs.JWT_ERROR, errr)

	}

	return jwtToken, nil
}

func (svc userService) GetUsersForAutoComplete(ctx *context.Context, name string, role []string) ([]responseModel.UserAutoComplete, *errs.XError) {

	users, err := svc.userRepo.GetUsersForAutoComplete(ctx, name, role)
	if err != nil {
		return nil, err
	}

	result := make([]responseModel.UserAutoComplete, 0)

	funk.ForEach(users, func(usr entities.User) {
		result = append(result, responseModel.UserAutoComplete{
			UserID:    usr.ID,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
		})
	})

	return result, nil
}

func (svc userService) UpdateChannel(ctx *context.Context, userId uint, channelId uint) *errs.XError {
	return svc.userRepo.UpdateChannel(ctx, userId, channelId)
}

func (svc userService) GetUserConfig(ctx *context.Context, userId uint) (*entities.UserConfig, *errs.XError) {
	userConfig, err := svc.userRepo.GetUserConfig(ctx, userId)
	if err != nil {
		return nil, err
	}

	return userConfig, nil
}

func (svc userService) SaveUserConfig(ctx *context.Context, config requestModel.UserConfig) *errs.XError {

	userConfig := entities.UserConfig{
		Model:  &entities.Model{IsActive: config.IsActive},
		Config: config.Config,
		UserID: config.UserID,
	}
	errr := svc.userRepo.CreateUserConfig(ctx, &userConfig)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc userService) UpdateUserConfig(ctx *context.Context, config requestModel.UserConfig, id uint) *errs.XError {
	userConfig := entities.UserConfig{
		Model:  &entities.Model{IsActive: true, ID: id},
		Config: config.Config,
		UserID: config.UserID,
	}

	errr := svc.userRepo.UpdateUserConfig(ctx, &userConfig)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc userService) AddUserChannelDetail(ctx *context.Context, userChannleDetails []requestModel.UserChannelDetail) *errs.XError {
	for _, userChannleDet := range userChannleDetails {
		userChannelDetail := entities.UserChannelDetail{
			Model:         &entities.Model{IsActive: true},
			UserID:        userChannleDet.UserID,
			UserChannelID: userChannleDet.ChannelId,
		}
		errr := svc.userRepo.CreateUserChannelDetail(ctx, &userChannelDetail)
		if errr != nil {
			return errr
		}
	}

	return nil
}

// UpdateUserChannelDetail implements UserService.
func (svc userService) UpdateUserChannelDetail(ctx *context.Context, userChannleDet requestModel.UserChannelDetail, id uint) *errs.XError {
	userChannelDetail := entities.UserChannelDetail{
		Model:         &entities.Model{IsActive: userChannleDet.IsActive, ID: id},
		UserID:        userChannleDet.UserID,
		UserChannelID: userChannleDet.ChannelId,
	}
	errr := svc.userRepo.UpdateUserChannelDetail(ctx, &userChannelDetail)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc userService) DeleteUserChannelDetail(ctx *context.Context, id uint) *errs.XError {
	err := svc.userRepo.DeleteUserChannelDetail(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (svc userService) SwitchUserChannel(ctx *context.Context, channelId uint) (string, *errs.XError) {

	session := utils.GetSession(ctx)
	if session == nil {
		return "", errs.NewXError(errs.SMTP_ERROR, "Unable to get user session", nil)

	}

	user, err := svc.Get(ctx, *session.UserId)
	if err != nil {
		return "", err
	}

	accessibleLocationIds := make([]uint, 0)
	funk.ForEach(user.UserChannelDetails, func(x responseModel.UserChannelDetail) {
		accessibleLocationIds = append(accessibleLocationIds, x.ChannelID)
	})

	if !funk.ContainsUInt(accessibleLocationIds, channelId) {
		return "", errs.NewXError(errs.INSUFFICIENT_ACCESS, "Channel Access Denied", nil).SetCode(http.StatusUnauthorized)
	}

	channel, err := svc.channelRepo.Get(ctx, channelId)
	if err != nil {
		return "", err
	}

	jwtResponse := models.Session{
		Email:                 user.Email,
		Role:                  entities.RoleType(user.Role),
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		UserId:                &user.ID,
		ChannelId:             channel.ChannelId,
		ChannelName:           channel.Name,
		AccessibleLocationIds: accessibleLocationIds,
	}

	jwtToken, errr := util.GenerateJWT(svc.config.Server.JwtSecretKey, svc.config.Server.JwtExpiryMinutes, util.StructToMap(jwtResponse))
	if errr != nil {
		return "", errs.NewXError(errs.INTERNAL, errs.JWT_ERROR, errr)

	}

	return jwtToken, nil
}

func (svc userService) GetUserChannelDetails(ctx *context.Context, userId uint) ([]responseModel.UserChannelDetail, *errs.XError) {

	channels, err := svc.userRepo.GetUserAccessibleChannels(ctx, userId)
	if err != nil {
		return nil, err
	}
	res := make([]responseModel.UserChannelDetail, 0)

	for _, ch := range channels {
		if ch.UserChannel == nil {
			//This means the channel is inactive , even though an active
			//record for this corresponding channel is present
			//in UserChannelDetails. So skip this to prevent nil deref
			continue
		}
		res = append(res, responseModel.UserChannelDetail{
			ID:        ch.ID,
			ChannelID: ch.UserChannelID,
			Name:      ch.UserChannel.Name,
		})
	}

	return res, nil
}
