package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type DressTypeRepository interface {
	Create(*context.Context, *entities.DressType) *errs.XError
	Update(*context.Context, *entities.DressType) *errs.XError
	Get(*context.Context, uint) (*entities.DressType, *errs.XError)
	GetAll(*context.Context, string) ([]entities.DressType, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type dressTypeRepository struct {
	GormDAL
}

func ProvideDressTypeRepository(dal GormDAL) DressTypeRepository {
	return &dressTypeRepository{GormDAL: dal}
}

func (dtr *dressTypeRepository) Create(ctx *context.Context, dressType *entities.DressType) *errs.XError {
	res := dtr.WithDB(ctx).Create(&dressType)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save dress type", res.Error)
	}
	return nil
}

func (dtr *dressTypeRepository) Update(ctx *context.Context, dressType *entities.DressType) *errs.XError {
	return dtr.GormDAL.Update(ctx, *dressType)
}

func (dtr *dressTypeRepository) Get(ctx *context.Context, id uint) (*entities.DressType, *errs.XError) {
	dressType := entities.DressType{}
	res := dtr.WithDB(ctx).
		Model(dressType).
		Scopes(scopes.WithAuditInfo()).
		Find(&dressType, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find dress type", res.Error)
	}
	return &dressType, nil
}

func (dtr *dressTypeRepository) GetAll(ctx *context.Context, search string) ([]entities.DressType, *errs.XError) {
	var dressTypes []entities.DressType
	res := dtr.WithDB(ctx).Model(entities.DressType{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.ILike(search, "name")).
		Scopes(db.Paginate(ctx)).
		Find(&dressTypes)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find dress types", res.Error)
	}
	return dressTypes, nil
}

func (dtr *dressTypeRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	dressType := &entities.DressType{Model: &entities.Model{ID: id, IsActive: false}}
	err := dtr.GormDAL.Delete(ctx, dressType)
	if err != nil {
		return err
	}
	return nil
}
