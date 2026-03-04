package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type MeasurementHistoryRepository interface {
	Create(*context.Context, *entities.MeasurementHistory) *errs.XError
	Get(*context.Context, uint) (*entities.MeasurementHistory, *errs.XError)
	GetAll(*context.Context, string) ([]entities.MeasurementHistory, *errs.XError)
	GetByMeasurementId(*context.Context, uint) ([]entities.MeasurementHistory, *errs.XError)
}

type measurementHistoryRepository struct {
	GormDAL
}

func ProvideMeasurementHistoryRepository(customDB GormDAL) MeasurementHistoryRepository {
	return &measurementHistoryRepository{GormDAL: customDB}
}

func (mhr *measurementHistoryRepository) Create(ctx *context.Context, measurementHistory *entities.MeasurementHistory) *errs.XError {
	res := mhr.WithDB(ctx).Create(&measurementHistory)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save measurement history", res.Error)
	}
	return nil
}

func (mhr *measurementHistoryRepository) Get(ctx *context.Context, id uint) (*entities.MeasurementHistory, *errs.XError) {
	measurementHistory := entities.MeasurementHistory{}
	res := mhr.WithDB(ctx).
		Preload("Measurement").
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Find(&measurementHistory, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement history", res.Error)
	}
	return &measurementHistory, nil
}

func (mhr *measurementHistoryRepository) GetAll(ctx *context.Context, search string) ([]entities.MeasurementHistory, *errs.XError) {
	var measurementHistories []entities.MeasurementHistory
	res := mhr.WithDB(ctx).Model(entities.MeasurementHistory{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(db.Paginate(ctx)).
		Preload("Measurement").
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Find(&measurementHistories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement histories", res.Error)
	}
	return measurementHistories, nil
}

func (mhr *measurementHistoryRepository) GetByMeasurementId(ctx *context.Context, measurementId uint) ([]entities.MeasurementHistory, *errs.XError) {
	var measurementHistories []entities.MeasurementHistory
	res := mhr.WithDB(ctx).
		Where("measurement_id = ?", measurementId).
		Scopes(scopes.IsActive()).
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Order("performed_at DESC").
		Find(&measurementHistories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement histories by measurement id", res.Error)
	}
	return measurementHistories, nil
}
