package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/errs"
)

type PurchaseItemRepository interface {
	Create(*context.Context, *entities.PurchaseItem) *errs.XError
	Update(*context.Context, *entities.PurchaseItem) *errs.XError
	Get(*context.Context, uint) (*entities.PurchaseItem, *errs.XError)
	GetByPurchaseId(*context.Context, uint) ([]entities.PurchaseItem, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type purchaseItemRepository struct {
	GormDAL
}

func ProvidePurchaseItemRepository(customDB GormDAL) PurchaseItemRepository {
	return &purchaseItemRepository{GormDAL: customDB}
}

func (pir *purchaseItemRepository) Create(ctx *context.Context, item *entities.PurchaseItem) *errs.XError {
	res := pir.WithDB(ctx).Create(item)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save purchase item", res.Error)
	}
	return nil
}

func (pir *purchaseItemRepository) Update(ctx *context.Context, item *entities.PurchaseItem) *errs.XError {
	return pir.GormDAL.Update(ctx, *item)
}

func (pir *purchaseItemRepository) Get(ctx *context.Context, id uint) (*entities.PurchaseItem, *errs.XError) {
	item := entities.PurchaseItem{}
	res := pir.WithDB(ctx).Model(entities.PurchaseItem{}).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Preload("Purchase").
		Preload("Product").
		Preload("Product.Category").
		Find(&item, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find purchase item", res.Error)
	}
	return &item, nil
}

func (pir *purchaseItemRepository) GetByPurchaseId(ctx *context.Context, purchaseId uint) ([]entities.PurchaseItem, *errs.XError) {
	var items []entities.PurchaseItem
	res := pir.WithDB(ctx).Model(entities.PurchaseItem{}).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("purchase_id = ?", purchaseId).
		Preload("Product").
		Preload("Product.Category").
		Find(&items)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find purchase items", res.Error)
	}
	return items, nil
}

func (pir *purchaseItemRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	item := &entities.PurchaseItem{Model: &entities.Model{ID: id, IsActive: false}}
	return pir.GormDAL.Delete(ctx, item)
}
