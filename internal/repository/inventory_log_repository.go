package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type InventoryLogRepository interface {
	Create(*context.Context, *entities.InventoryLog) *errs.XError
	Get(*context.Context, uint) (*entities.InventoryLog, *errs.XError)
	GetAll(*context.Context, string) ([]entities.InventoryLog, *errs.XError)
	GetByProductId(*context.Context, uint) ([]entities.InventoryLog, *errs.XError)
	GetByChangeType(*context.Context, entities.InventoryLogChangeType) ([]entities.InventoryLog, *errs.XError)
	GetByDateRange(*context.Context, string, string) ([]entities.InventoryLog, *errs.XError)
}

type inventoryLogRepository struct {
	GormDAL
}

func ProvideInventoryLogRepository(customDB GormDAL) InventoryLogRepository {
	return &inventoryLogRepository{GormDAL: customDB}
}

func (ilr *inventoryLogRepository) Create(ctx *context.Context, log *entities.InventoryLog) *errs.XError {
	res := ilr.WithDB(ctx).Create(&log)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to create inventory log", res.Error)
	}
	return nil
}

func (ilr *inventoryLogRepository) Get(ctx *context.Context, id uint) (*entities.InventoryLog, *errs.XError) {
	log := entities.InventoryLog{}
	res := ilr.WithDB(ctx).
		Preload("Product").
		Preload("Product.Category").
		Find(&log, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventory log", res.Error)
	}
	return &log, nil
}

func (ilr *inventoryLogRepository) GetAll(ctx *context.Context, search string) ([]entities.InventoryLog, *errs.XError) {
	var logs []entities.InventoryLog
	res := ilr.WithDB(ctx).Table(entities.InventoryLog{}.TableNameForQuery()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(db.Paginate(ctx)).
		Preload("Product").
		Preload("Product.Category").
		Order("logged_at DESC").
		Find(&logs)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventory logs", res.Error)
	}
	return logs, nil
}

func (ilr *inventoryLogRepository) GetByProductId(ctx *context.Context, productId uint) ([]entities.InventoryLog, *errs.XError) {
	var logs []entities.InventoryLog
	res := ilr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.WithAuditInfo()).
		Where("product_id = ?", productId).
		Preload("Product").
		Order("logged_at DESC").
		Find(&logs)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventory logs for product", res.Error)
	}
	return logs, nil
}

func (ilr *inventoryLogRepository) GetByChangeType(ctx *context.Context, changeType entities.InventoryLogChangeType) ([]entities.InventoryLog, *errs.XError) {
	var logs []entities.InventoryLog
	res := ilr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("change_type = ?", changeType).
		Preload("Product").
		Order("logged_at DESC").
		Find(&logs)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventory logs by change type", res.Error)
	}
	return logs, nil
}

func (ilr *inventoryLogRepository) GetByDateRange(ctx *context.Context, startDate string, endDate string) ([]entities.InventoryLog, *errs.XError) {
	var logs []entities.InventoryLog
	query := ilr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Preload("Product").
		Order("logged_at DESC")

	if startDate != "" {
		query = query.Where("logged_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("logged_at <= ?", endDate)
	}

	res := query.Find(&logs)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find inventory logs by date range", res.Error)
	}
	return logs, nil
}
