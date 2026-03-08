package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type VendorRepository interface {
	Create(*context.Context, *entities.Vendor) *errs.XError
	Update(*context.Context, *entities.Vendor) *errs.XError
	Get(*context.Context, uint) (*entities.Vendor, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Vendor, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	AutocompleteVendor(*context.Context, string) ([]entities.Vendor, *errs.XError)
}

type vendorRepository struct {
	GormDAL
}

func ProvideVendorRepository(customDB GormDAL) VendorRepository {
	return &vendorRepository{GormDAL: customDB}
}

func (vr *vendorRepository) Create(ctx *context.Context, vendor *entities.Vendor) *errs.XError {
	res := vr.WithDB(ctx).Create(vendor)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save vendor", res.Error)
	}
	return nil
}

func (vr *vendorRepository) Update(ctx *context.Context, vendor *entities.Vendor) *errs.XError {
	return vr.GormDAL.Update(ctx, *vendor)
}

func (vr *vendorRepository) Get(ctx *context.Context, id uint) (*entities.Vendor, *errs.XError) {
	vendor := entities.Vendor{}
	res := vr.WithDB(ctx).Model(entities.Vendor{}).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Find(&vendor, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find vendor", res.Error)
	}
	return &vendor, nil
}

func (vr *vendorRepository) GetAll(ctx *context.Context, search string) ([]entities.Vendor, *errs.XError) {
	var vendors []entities.Vendor
	res := vr.WithDB(ctx).Model(entities.Vendor{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.ILike(search, "name", "contact_person", "email")).
		Scopes(db.Paginate(ctx)).
		Find(&vendors)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find vendors", res.Error)
	}
	return vendors, nil
}

func (vr *vendorRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	vendor := &entities.Vendor{Model: &entities.Model{ID: id, IsActive: false}}
	return vr.GormDAL.Delete(ctx, vendor)
}

func (vr *vendorRepository) AutocompleteVendor(ctx *context.Context, search string) ([]entities.Vendor, *errs.XError) {
	var vendors []entities.Vendor
	res := vr.WithDB(ctx).Model(entities.Vendor{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "name", "contact_person", "email")).
		Select("id", "name").
		Find(&vendors)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find vendors for autocomplete", res.Error)
	}
	return vendors, nil
}
