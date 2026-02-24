package service

import (
	"context"
	"strings"
	"time"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
)

type OrderService interface {
	SaveOrder(*context.Context, requestModel.Order) *errs.XError
	UpdateOrder(*context.Context, requestModel.Order, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Order, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Order, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type orderService struct {
	orderRepo        repository.OrderRepository
	orderHistoryRepo repository.OrderHistoryRepository
	mapper           mapper.Mapper
	respMapper       mapper.ResponseMapper
}

func ProvideOrderService(repo repository.OrderRepository, orderHistoryRepo repository.OrderHistoryRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) OrderService {
	return orderService{
		orderRepo:        repo,
		orderHistoryRepo: orderHistoryRepo,
		mapper:           mapper,
		respMapper:       respMapper,
	}
}

func (svc orderService) SaveOrder(ctx *context.Context, order requestModel.Order) *errs.XError {
	dbOrder, err := svc.mapper.Order(order)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save order", err)
	}

	// Set TakenById to the current user if it's not provided in the request
	if order.OrderTakenById == nil {
		userID := utils.GetUserId(ctx)
		dbOrder.OrderTakenById = &userID
	}

	errr := svc.orderRepo.Create(ctx, dbOrder)
	if errr != nil {
		return errr
	}

	// Record order history for CREATED action
	errr = svc.recordOrderHistory(ctx, dbOrder.ID, entities.OrderHistoryActionCreated, nil, nil, nil, nil)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc orderService) UpdateOrder(ctx *context.Context, order requestModel.Order, id uint) *errs.XError {
	// Get the old order before updating
	oldOrder, err := svc.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	dbOrder, mapErr := svc.mapper.Order(order)
	if mapErr != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update order", mapErr)
	}

	// Set TakenById to the current user if it's not provided in the request
	if order.OrderTakenById == nil {
		userID := utils.GetUserId(ctx)
		dbOrder.OrderTakenById = &userID
	}

	dbOrder.ID = id
	errr := svc.orderRepo.Update(ctx, dbOrder)
	if errr != nil {
		return errr
	}

	// Determine changed fields
	var changedFields []string
	if oldOrder.Status != dbOrder.Status {
		changedFields = append(changedFields, entities.OrderChangeFieldStatus)
	}
	if !timeEqual(oldOrder.ExpectedDeliveryDate, dbOrder.ExpectedDeliveryDate) {
		changedFields = append(changedFields, entities.OrderChangeFieldExpectedDeliveryDate)
	}
	if !timeEqual(oldOrder.DeliveredDate, dbOrder.DeliveredDate) {
		changedFields = append(changedFields, entities.OrderChangeFieldDeliveredDate)
	}

	changedFieldsStr := strings.Join(changedFields, ",")

	// Record order history for UPDATED action with old values
	errr = svc.recordOrderHistory(ctx, id, entities.OrderHistoryActionUpdated, &oldOrder.Status, oldOrder.ExpectedDeliveryDate, oldOrder.DeliveredDate, &changedFieldsStr)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc orderService) Get(ctx *context.Context, id uint) (*responseModel.Order, *errs.XError) {
	order, err := svc.orderRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedOrder, mapErr := svc.respMapper.Order(order)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Order data", mapErr)
	}

	return mappedOrder, nil
}

func (svc orderService) GetAll(ctx *context.Context, search string) ([]responseModel.Order, *errs.XError) {
	orders, err := svc.orderRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedOrders, mapErr := svc.respMapper.Orders(orders)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Order data", mapErr)
	}

	return mappedOrders, nil
}

func (svc orderService) Delete(ctx *context.Context, id uint) *errs.XError {
	// Get the old order before deleting
	oldOrder, err := svc.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	err = svc.orderRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Record order history for DELETED action with old values
	err = svc.recordOrderHistory(ctx, id, entities.OrderHistoryActionDeleted, &oldOrder.Status, oldOrder.ExpectedDeliveryDate, oldOrder.DeliveredDate, nil)
	if err != nil {
		return err
	}

	return nil
}

// recordOrderHistory creates an order history record
func (svc orderService) recordOrderHistory(ctx *context.Context, orderId uint, action entities.OrderHistoryAction, oldStatus *entities.OrderStatus, oldExpectedDeliveryDate *time.Time, oldDeliveredDate *time.Time, changedFields *string) *errs.XError {
	userID := utils.GetUserId(ctx)
	performedAt := util.GetLocalTime()

	history := &entities.OrderHistory{
		Model:         &entities.Model{IsActive: true},
		Action:        action,
		OrderId:       orderId,
		PerformedAt:   performedAt,
		PerformedById: userID,
	}

	// Set changed fields if provided
	if changedFields != nil {
		history.ChangedFields = *changedFields
	}

	// Set old values for tracking
	if oldStatus != nil {
		history.Status = oldStatus
	}
	if oldExpectedDeliveryDate != nil {
		history.ExpectedDeliveryDate = oldExpectedDeliveryDate
	}
	if oldDeliveredDate != nil {
		history.DeliveredDate = oldDeliveredDate
	}

	return svc.orderHistoryRepo.Create(ctx, history)
}

// timeEqual compares two *time.Time values, handling nil cases
func timeEqual(t1, t2 *time.Time) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.Equal(*t2)
}
