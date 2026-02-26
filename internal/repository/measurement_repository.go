package repository

import (
	"context"
	"errors"

	"github.com/imkarthi24/sf-backend/internal/entities"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"

	"gorm.io/gorm"
)

type MeasurementRepository interface {
	Create(*context.Context, *entities.Measurement) *errs.XError
	BatchCreate(*context.Context, []*entities.Measurement) *errs.XError
	Update(*context.Context, *entities.Measurement) *errs.XError
	BatchUpdate(*context.Context, []*entities.Measurement) *errs.XError
	Get(*context.Context, uint) (*entities.Measurement, *errs.XError)
	GetByPersonIdAndDressTypeId(*context.Context, uint, uint) (*entities.Measurement, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.MeasurementBrowse, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type measurementRepository struct {
	GormDAL
}

func ProvideMeasurementRepository(customDB GormDAL) MeasurementRepository {
	return &measurementRepository{GormDAL: customDB}
}

func (mr *measurementRepository) Create(ctx *context.Context, measurement *entities.Measurement) *errs.XError {
	res := mr.WithDB(ctx).Create(&measurement)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save measurement", res.Error)
	}
	return nil
}

func (mr *measurementRepository) BatchCreate(ctx *context.Context, measurements []*entities.Measurement) *errs.XError {
	if len(measurements) == 0 {
		return nil
	}

	res := mr.WithDB(ctx).CreateInBatches(measurements, 100)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to batch save measurements", res.Error)
	}
	return nil
}

func (mr *measurementRepository) Update(ctx *context.Context, measurement *entities.Measurement) *errs.XError {
	return mr.GormDAL.Update(ctx, *measurement)
}

func (mr *measurementRepository) BatchUpdate(ctx *context.Context, measurements []*entities.Measurement) *errs.XError {
	if len(measurements) == 0 {
		return nil
	}

	for _, measurement := range measurements {
		if measurement.ID == 0 {
			continue
		}
		err := mr.GormDAL.Update(ctx, *measurement)
		if err != nil {
			return errs.NewXError(errs.DATABASE, "Unable to batch update measurements", err)
		}
	}

	return nil
}

func (mr *measurementRepository) Get(ctx *context.Context, id uint) (*entities.Measurement, *errs.XError) {
	measurement := entities.Measurement{}
	res := mr.WithDB(ctx).Model(measurement).
		Preload("Person").
		Preload("DressType").
		Preload("TakenBy", scopes.SelectFields("first_name", "last_name")).
		Scopes(scopes.WithAuditInfo()).
		Find(&measurement, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement", res.Error)
	}
	return &measurement, nil
}

func (mr *measurementRepository) GetByPersonIdAndDressTypeId(ctx *context.Context, personId uint, dressTypeId uint) (*entities.Measurement, *errs.XError) {
	measurement := entities.Measurement{}
	res := mr.WithDB(ctx).
		Where("person_id = ? AND dress_type_id = ?", personId, dressTypeId).
		First(&measurement)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewXError(errs.DATABASE, "Unable to find measurement", res.Error)
	}
	return &measurement, nil
}
func (mr *measurementRepository) GetAll(
	ctx *context.Context,
	search string,
) ([]responseModel.MeasurementBrowse, *errs.XError) {

	var result []responseModel.MeasurementBrowse

	// Subquery: latest measurement PER PERSON
	latestSubQuery := mr.WithDB(ctx).
		Table(entities.Measurement{}.TableNameForQuery()).
		Select(`
			DISTINCT ON (person_id)
			id,
			person_id,
			updated_at,
			taken_by_id
		`).
		Scopes(scopes.IsActive(), scopes.Channel("E")).
		Order(`person_id, updated_at DESC`)

	query := mr.WithDB(ctx).
		Table(entities.Measurement{}.TableNameForQuery()).
		Select(`
			latest.id,
			E.is_active,
			E.person_id,
			P.customer_id,
			CONCAT(P.first_name, ' ', P.last_name) AS person_name,
			latest.updated_at,
			CONCAT(U.first_name, ' ', U.last_name) AS updated_by,
			STRING_AGG(DISTINCT DT.name, ', ' ORDER BY DT.name) AS dress_types
		`).
		Joins(`INNER JOIN (?) latest ON latest.person_id = E.person_id`, latestSubQuery).
		Joins(`INNER JOIN "stich"."DressTypes" DT ON DT.id = E.dress_type_id`).
		Joins(`INNER JOIN "stich"."Persons" P ON P.id = E.person_id`).
		Joins(`INNER JOIN "stich"."Users" U ON U.id = latest.taken_by_id`).
		Scopes(scopes.IsActive("E"), scopes.Channel("E")).
		Group(`
			latest.id,
			latest.updated_at,
			E.is_active,
			E.person_id,
			P.customer_id,
			P.first_name,
			P.last_name,
			U.first_name,
			U.last_name
		`).
		Order(`latest.updated_at DESC`).
		Scopes(db.Paginate(ctx))

	// 🔍 Optional Search
	if !util.IsNilOrEmptyString(&search) {
		formatted := util.EncloseWithSymbol(search, "%")
		query = query.
			Joins(`INNER JOIN "stich"."Customers" C ON C.id = P.customer_id`).
			Where(`
				C.first_name ILIKE ? OR
				C.last_name ILIKE ? OR
				CONCAT(C.first_name, ' ', C.last_name) ILIKE ? OR
				CONCAT(P.first_name, ' ', P.last_name) ILIKE ? OR
				C.phone_number ILIKE ? OR
				C.whatsapp_number ILIKE ? 
			`, formatted, formatted, formatted, formatted, formatted, formatted)
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, errs.NewXError(errs.DATABASE, "UNABLE_TO_FIND_MEASUREMENTS", err)
	}

	return result, nil
}

func (mr *measurementRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	measurement := &entities.Measurement{Model: &entities.Model{ID: id, IsActive: false}}
	err := mr.GormDAL.Delete(ctx, measurement)
	if err != nil {
		return err
	}
	return nil
}
