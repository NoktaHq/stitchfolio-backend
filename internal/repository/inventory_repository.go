package repository

import (
	"context"
	"time"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type InventoryRepository interface {
	Create(*context.Context, *entities.Inventory) *errs.XError
	Update(*context.Context, *entities.Inventory) *errs.XError
	Get(*context.Context, uint) (*entities.Inventory, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Inventory, *errs.XError)
	GetByProductId(*context.Context, uint) (*entities.Inventory, *errs.XError)
	UpdateQuantity(*context.Context, uint, int) *errs.XError
	GetLowStockItems(*context.Context) ([]entities.Inventory, *errs.XError)
	UpdateThreshold(*context.Context, uint, int) *errs.XError
}

type inventoryRepository struct {
	GormDAL
}

func ProvideInventoryRepository(customDB GormDAL) InventoryRepository {
	return &inventoryRepository{GormDAL: customDB}
}

func (ir *inventoryRepository) Create(ctx *context.Context, inventory *entities.Inventory) *errs.XError {
	res := ir.WithDB(ctx).Create(&inventory)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to create inventory", res.Error)
	}
	return nil
}

func (ir *inventoryRepository) Update(ctx *context.Context, inventory *entities.Inventory) *errs.XError {
	return ir.GormDAL.Update(ctx, *inventory)
}

func (ir *inventoryRepository) Get(ctx *context.Context, id uint) (*entities.Inventory, *errs.XError) {
	inventory := entities.Inventory{}
	res := ir.WithDB(ctx).
		Preload("Product").
		Preload("Product.Category").
		Find(&inventory, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventory", res.Error)
	}
	return &inventory, nil
}

func (ir *inventoryRepository) GetAll(ctx *context.Context, search string) ([]entities.Inventory, *errs.XError) {
	var inventories []entities.Inventory
	res := ir.WithDB(ctx).Model(entities.Inventory{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(db.Paginate(ctx)).
		Preload("Product").
		Preload("Product.Category").
		Find(&inventories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventories", res.Error)
	}
	return inventories, nil
}

func (ir *inventoryRepository) GetByProductId(ctx *context.Context, productId uint) (*entities.Inventory, *errs.XError) {
	inventory := entities.Inventory{}
	res := ir.WithDB(ctx).
		Where("product_id = ?", productId).
		Preload("Product").
		First(&inventory)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventory for product", res.Error)
	}
	return &inventory, nil
}

func (ir *inventoryRepository) UpdateQuantity(ctx *context.Context, productId uint, newQuantity int) *errs.XError {
	res := ir.WithDB(ctx).
		Model(&entities.Inventory{}).
		Where("product_id = ?", productId).
		Updates(map[string]interface{}{
			"quantity":   newQuantity,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to update inventory quantity", res.Error)
	}
	return nil
}

func (ir *inventoryRepository) GetLowStockItems(ctx *context.Context) ([]entities.Inventory, *errs.XError) {
	var inventories []entities.Inventory
	res := ir.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("quantity <= low_stock_threshold").
		Preload("Product").
		Preload("Product.Category").
		Find(&inventories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find low stock items", res.Error)
	}
	return inventories, nil
}

func (ir *inventoryRepository) UpdateThreshold(ctx *context.Context, productId uint, threshold int) *errs.XError {
	res := ir.WithDB(ctx).
		Model(&entities.Inventory{}).
		Where("product_id = ?", productId).
		Update("low_stock_threshold", threshold)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to update low stock threshold", res.Error)
	}
	return nil
}
