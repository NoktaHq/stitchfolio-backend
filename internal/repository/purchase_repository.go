package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type PurchaseRepository interface {
	Create(*context.Context, *entities.Purchase) *errs.XError
	Update(*context.Context, *entities.Purchase) *errs.XError
	Get(*context.Context, uint) (*entities.Purchase, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Purchase, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type purchaseRepository struct {
	GormDAL
}

func ProvidePurchaseRepository(customDB GormDAL) PurchaseRepository {
	return &purchaseRepository{GormDAL: customDB}
}

func (pr *purchaseRepository) Create(ctx *context.Context, purchase *entities.Purchase) *errs.XError {
	res := pr.WithDB(ctx).Create(purchase)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save purchase", res.Error)
	}
	return nil
}

func (pr *purchaseRepository) Update(ctx *context.Context, purchase *entities.Purchase) *errs.XError {
	return pr.GormDAL.Update(ctx, *purchase)
}

func (pr *purchaseRepository) Get(ctx *context.Context, id uint) (*entities.Purchase, *errs.XError) {
	purchase := entities.Purchase{}
	res := pr.WithDB(ctx).Model(entities.Purchase{}).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Preload("Vendor").
		Preload("PurchaseItems").
		Preload("PurchaseItems.Product").
		Preload("PurchaseItems.Product.Category").
		Find(&purchase, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find purchase", res.Error)
	}
	return &purchase, nil
}

func (pr *purchaseRepository) GetAll(ctx *context.Context, search string) ([]entities.Purchase, *errs.XError) {
	var purchases []entities.Purchase
	res := pr.WithDB(ctx).Model(entities.Purchase{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.ILike(search, "purchase_number")).
		Scopes(db.Paginate(ctx)).
		Preload("Vendor").
		Preload("PurchaseItems").
		Preload("PurchaseItems.Product").
		Order("purchase_date DESC").
		Find(&purchases)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find purchases", res.Error)
	}
	return purchases, nil
}

func (pr *purchaseRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	purchase := &entities.Purchase{Model: &entities.Model{ID: id, IsActive: false}}
	return pr.GormDAL.Delete(ctx, purchase)
}
