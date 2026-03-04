package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type OrderHistoryRepository interface {
	Create(*context.Context, *entities.OrderHistory) *errs.XError
	Get(*context.Context, uint) (*entities.OrderHistory, *errs.XError)
	GetAll(*context.Context, string) ([]entities.OrderHistory, *errs.XError)
	GetByOrderId(*context.Context, uint) ([]entities.OrderHistory, *errs.XError)
}

type orderHistoryRepository struct {
	GormDAL
}

func ProvideOrderHistoryRepository(customDB GormDAL) OrderHistoryRepository {
	return &orderHistoryRepository{GormDAL: customDB}
}

func (ohr *orderHistoryRepository) Create(ctx *context.Context, orderHistory *entities.OrderHistory) *errs.XError {
	res := ohr.WithDB(ctx).Create(&orderHistory)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save order history", res.Error)
	}
	return nil
}

func (ohr *orderHistoryRepository) Get(ctx *context.Context, id uint) (*entities.OrderHistory, *errs.XError) {
	orderHistory := entities.OrderHistory{}
	res := ohr.WithDB(ctx).
		Preload("Order").
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Find(&orderHistory, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find order history", res.Error)
	}
	return &orderHistory, nil
}

func (ohr *orderHistoryRepository) GetAll(ctx *context.Context, search string) ([]entities.OrderHistory, *errs.XError) {
	var orderHistories []entities.OrderHistory
	res := ohr.WithDB(ctx).Model(entities.OrderHistory{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(db.Paginate(ctx)).
		Preload("Order").
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Find(&orderHistories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find order histories", res.Error)
	}
	return orderHistories, nil
}

func (ohr *orderHistoryRepository) GetByOrderId(ctx *context.Context, orderId uint) ([]entities.OrderHistory, *errs.XError) {
	var orderHistories []entities.OrderHistory
	res := ohr.WithDB(ctx).
		Where("order_id = ?", orderId).
		Scopes(scopes.IsActive()).
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name")).
		Order("performed_at DESC").
		Find(&orderHistories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find order histories by order id", res.Error)
	}
	return orderHistories, nil
}
