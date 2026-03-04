package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	requesModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/imkarthi24/sf-backend/internal/utils/validator"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type UserHandler struct {
	userSvc  service.UserService
	resp     response.Response
	dataResp response.DataResponse
}

func ProvideUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{userSvc: svc}
}

var _ = (*responseModel.Response)(nil) // used by swagger comments

// Save User
//
//	@Summary		Save User
//	@Description	Saves an instance of User
//	@Tags			User
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			user	body		requestModel.User	true	"user"
//	@Router			/user [post]
func (h UserHandler) SaveUser(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var user requesModel.User
	err := ctx.Bind(&user)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	if ok, err := validator.ValidateUser(user); !ok {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.userSvc.SaveUser(&context, user)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)

}

// Update User
//
//	@Summary		Update User
//	@Description	Updates an instance of User
//	@Tags			User
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			user	body		requestModel.User	true	"user"
//	@Param			id		path		int					true	"User id"
//	@Router			/user/{id} [put]
func (h UserHandler) UpdateUser(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var user requesModel.User
	err := ctx.Bind(&user)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.userSvc.UpdateUser(&context, user, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)

}

// Login User
//
//	@Summary		Login User
//	@Description	Saves an instance of User
//	@Tags			User
//	@Accept			json
//	@Failure		400			{object}	responseModel.Response
//	@Param			userLogin	body		requestModel.Login	true	"login"
//	@Router			/user/login [post]
func (h UserHandler) Login(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	var userLogin requesModel.Login
	err := ctx.Bind(&userLogin)
	if err != nil || util.IsNilOrEmptyString(&userLogin.Email) || util.IsNilOrEmptyString(&userLogin.Password) {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	login, errr := h.userSvc.Login(&context, userLogin)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(login).FormatAndSend(&context, ctx, http.StatusOK)

}

// Get User
//
//	@Summary		Get a specific  User
//	@Description	Get an instance of User
//	@Tags			User
//	@Accept			json
//	@Success		200	{object}	responseModel.User
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"User id"
//	@Router			/user/{id} [get]
func (h UserHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	user, errr := h.userSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(user).FormatAndSend(&context, ctx, http.StatusOK)

}

// Get all active users
//
//	@Summary		Get all active users
//	@Description	Get all active users
//	@Tags			User
//	@Accept			json
//	@Success		200		{object}	responseModel.User
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/user [get]
func (h UserHandler) GetAllUsers(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	login, errr := h.userSvc.GetAllUsers(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(login).FormatAndSend(&context, ctx, http.StatusOK)

}

// Delete a User
//
//	@Summary		Delete User
//	@Description	Deletes an instance of User
//	@Tags			User
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"user id"
//
//	@Router			/user/{id} [delete]
func (h UserHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.userSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)

}

// Reset User Password
//
//	@Summary		Reset User Password
//	@Description	Reset User Password
//	@Tags			User
//	@Success		200		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Param			email	query		string	true	"email"
//
//	@Router			/user/forgot-password [POST]
func (h UserHandler) ForgotPassword(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	mail := ctx.Query("email")
	err := h.userSvc.ForgotPassword(&context, mail)
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.DefaultSuccessResponse().FormatAndSend(&context, ctx, http.StatusOK)

}

// Update User Password
//
//	@Summary		Update User Password
//	@Description	Update User Password
//	@Tags			User
//	@Success		200			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Param			resetString	query		string	true	"resetString"
//	@Param			password	query		string	true	"password"
//
//	@Router			/user/reset-password [POST]
func (h UserHandler) ResetPassword(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	password := ctx.Query("password")
	resetString := ctx.Query("resetString")

	err := h.userSvc.ResetPassword(&context, resetString, password)
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Password Update Success").FormatAndSend(&context, ctx, http.StatusOK)

}

// Refresh JWT Token
//
//	@Summary		Refresh JWT Token
//	@Description	Refresh JWT Token
//	@Tags			User
//	@Accept			json
//	@Router			/user/refresh-token [get]
func (h UserHandler) RefreshToken(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	jwt, err := h.userSvc.RefreshToken(&context)
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(jwt).FormatAndSend(&context, ctx, http.StatusOK)

}

// Get  Users for AutoComplete
//
//	@Summary		Get  Users for AutoComplete
//	@Description	Get  Users for AutoComplete
//	@Produce		json
//	@Success		200	{object}	entities.User
//	@Tags			User
//	@Router			/user/autocomplete [get]
//	@Param			name	query	string	true	"name"
//	@Param			role	query	string	true	"role"
func (h UserHandler) GetUsersForAutoComplete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	name := ctx.Query("name")
	role := ctx.Query("role")
	roles := util.SplitNonEmpty(role, ",")
	if !util.IsNilOrEmptyString(&name) && len(name) < 3 {
		h.dataResp.FailureResponse(nil, errs.INCORRECT_PARAMETER).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	name = util.EncloseWithSingleQuote(name)
	users, errr := h.userSvc.GetUsersForAutoComplete(&context, name, roles)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.dataResp.DefaultSuccessResponse(users).FormatAndSend(&context, ctx, http.StatusOK)

}

// Save User Config
//
//	@Summary		Save User Config
//	@Description	Saves an instance of User Config
//	@Tags			UserConfig
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			userConfig	body		requestModel.UserConfig	true	"userconfig"
//
//	@Router			/user/config [post]
func (h UserHandler) SaveUserConfig(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	var userConfig requesModel.UserConfig
	err := ctx.Bind(&userConfig)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.userSvc.SaveUserConfig(&context, userConfig)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)

}

// Update User Config
//
//	@Summary		Update User Config
//	@Description	Update an instance of User Config
//	@Tags			UserConfig
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			userConfig	body		requestModel.UserConfig	true	"userconfig"
//	@Param			id			path		int						true	"userconfigid"
//	@Router			/user/config/{id} [put]
func (h UserHandler) UpdateUserConfig(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	var userConfig requesModel.UserConfig
	err := ctx.Bind(&userConfig)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.userSvc.UpdateUserConfig(&context, userConfig, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusCreated)

}

// Get User Config
//
//	@Summary		Get a specific  UserConfig
//	@Description	Get an instance of UserConfig
//	@Tags			UserConfig
//	@Accept			json
//	@Success		200	{object}	entities.UserConfig
//	@Failure		400	{object}	responseModel.DataResponse
//
//	@Param			id	query		int	true	"userId"
//	@Router			/user/config/ [get]
func (h UserHandler) GetUserConfig(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Query("userId"))

	user, errr := h.userSvc.GetUserConfig(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(user).FormatAndSend(&context, ctx, http.StatusOK)

}

// Save User Channel Detail
//
//	@Summary		Save User Channel Detail
//	@Description	Saves an instance of User Channel Detail
//	@Tags			UserChannelDetail
//	@Accept			json
//	@Success		201					{object}	responseModel.Response
//	@Failure		400					{object}	responseModel.Response
//	@Failure		501					{object}	responseModel.Response
//	@Param			userChannelDetail	body		[]requestModel.UserChannelDetail	true	"userChannelDetail"
//
//	@Router			/user/channel [post]
func (h UserHandler) SaveUserChannelDetail(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	var userChanDet []requesModel.UserChannelDetail
	err := ctx.Bind(&userChanDet)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.userSvc.AddUserChannelDetail(&context, userChanDet)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)

}

// Update User Channel Detail
//
//	@Summary		Update User Channel Detail
//	@Description	Update an instance of User Channel Detail
//	@Tags			UserChannelDetail
//	@Accept			json
//	@Success		201					{object}	responseModel.Response
//	@Failure		400					{object}	responseModel.Response
//	@Failure		501					{object}	responseModel.Response
//	@Param			userChannelDetail	body		requestModel.UserChannelDetail	true	"userconfig"
//	@Param			id					path		int								true	"userChannelDetailId"
//	@Router			/user/channel/{id} [put]
func (h UserHandler) UpdateUserChannelDetail(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	var userChanDetail requesModel.UserChannelDetail
	err := ctx.Bind(&userChanDetail)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.userSvc.UpdateUserChannelDetail(&context, userChanDetail, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusCreated)

}

// Delete a User Channel Detail
//
//	@Summary		Delete User Channel Detail
//	@Description	Deletes an instance of User Channel Detail
//	@Tags			UserChannelDetail
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"user channel detail id"
//
//	@Router			/user/channel/{id} [delete]
func (h UserHandler) DeleteUserChannelDetail(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.userSvc.DeleteUserChannelDetail(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)

}

// Switch Channel
//
//	@Summary		Refresh JWT Token and Switches the Channel
//	@Description	Refresh JWT Token and switches the user to the requested channel
//	@Tags			User
//	@Accept			json
//	@Param			id	path	int	true	"Channel id"
//	@Router			/user/switch-channel [get]
func (h UserHandler) SwitchChannel(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	channelId, _ := strconv.Atoi(ctx.Param("id"))
	jwt, err := h.userSvc.SwitchUserChannel(&context, uint(channelId))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(jwt).FormatAndSend(&context, ctx, http.StatusOK)

}

// Get Accessible Channels
//
//	@Summary		Gets the channel the user can access
//	@Description	Gets the channel the user can access, to be used for channel switcher
//	@Tags			User
//	@Accept			json
//	@Param			userId	query	int	true	"User id"
//	@Router			/user/channel/accessible [get]
func (h UserHandler) GetUserAccessibleChannels(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	userId, _ := strconv.Atoi(ctx.Query("userId"))
	jwt, err := h.userSvc.GetUserChannelDetails(&context, uint(userId))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(jwt).FormatAndSend(&context, ctx, http.StatusOK)

}
