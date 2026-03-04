package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type EnquiryHistoryRepository interface {
	Create(*context.Context, *entities.EnquiryHistory) *errs.XError
	Get(*context.Context, uint) (*entities.EnquiryHistory, *errs.XError)
	GetAll(*context.Context, string) ([]entities.EnquiryHistory, *errs.XError)
	GetByEnquiryId(*context.Context, uint) ([]entities.EnquiryHistory, *errs.XError)
}

type enquiryHistoryRepository struct {
	GormDAL
}

func ProvideEnquiryHistoryRepository(dal GormDAL) EnquiryHistoryRepository {
	return &enquiryHistoryRepository{GormDAL: dal}
}

func (ehr *enquiryHistoryRepository) Create(ctx *context.Context, enquiryHistory *entities.EnquiryHistory) *errs.XError {
	res := ehr.WithDB(ctx).Create(&enquiryHistory)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save enquiry history", res.Error)
	}
	return nil
}

func (ehr *enquiryHistoryRepository) Get(ctx *context.Context, id uint) (*entities.EnquiryHistory, *errs.XError) {
	enquiryHistory := entities.EnquiryHistory{}
	res := ehr.WithDB(ctx).
		Preload("Employee", scopes.SelectFields("first_name", "last_name")).
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Find(&enquiryHistory, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find enquiry history", res.Error)
	}
	return &enquiryHistory, nil
}

func (ehr *enquiryHistoryRepository) GetAll(ctx *context.Context, search string) ([]entities.EnquiryHistory, *errs.XError) {
	var enquiryHistories []entities.EnquiryHistory
	res := ehr.WithDB(ctx).Model(entities.EnquiryHistory{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(db.Paginate(ctx)).
		Preload("Employee", scopes.SelectFields("first_name", "last_name")).
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Find(&enquiryHistories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find enquiry histories", res.Error)
	}
	return enquiryHistories, nil
}

func (ehr *enquiryHistoryRepository) GetByEnquiryId(ctx *context.Context, enquiryId uint) ([]entities.EnquiryHistory, *errs.XError) {
	var enquiryHistories []entities.EnquiryHistory
	res := ehr.WithDB(ctx).
		Where("enquiry_id = ?", enquiryId).
		Scopes(scopes.IsActive()).
		Preload("Employee", scopes.SelectFields("first_name", "last_name")).
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Order("performed_at DESC").
		Find(&enquiryHistories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find enquiry histories by enquiry id", res.Error)
	}
	return enquiryHistories, nil
}
