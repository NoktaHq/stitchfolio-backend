package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type PurchaseItemService interface {
	Save(*context.Context, requestModel.PurchaseItem, uint) *errs.XError
	Update(*context.Context, requestModel.PurchaseItem, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.PurchaseItem, *errs.XError)
	GetByPurchaseId(*context.Context, uint) ([]responseModel.PurchaseItem, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type purchaseItemService struct {
	purchaseItemRepo repository.PurchaseItemRepository
	purchaseRepo     repository.PurchaseRepository
	mapper           mapper.Mapper
	respMapper       mapper.ResponseMapper
}

func ProvidePurchaseItemService(
	purchaseItemRepo repository.PurchaseItemRepository,
	purchaseRepo repository.PurchaseRepository,
	m mapper.Mapper,
	respMapper mapper.ResponseMapper,
) PurchaseItemService {
	return purchaseItemService{
		purchaseItemRepo: purchaseItemRepo,
		purchaseRepo:     purchaseRepo,
		mapper:           m,
		respMapper:       respMapper,
	}
}

func (svc purchaseItemService) Save(ctx *context.Context, req requestModel.PurchaseItem, purchaseId uint) *errs.XError {
	req.PurchaseId = purchaseId
	ent, err := svc.mapper.PurchaseItem(req)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save purchase item", err)
	}
	return svc.purchaseItemRepo.Create(ctx, ent)
}

func (svc purchaseItemService) Update(ctx *context.Context, req requestModel.PurchaseItem, id uint) *errs.XError {
	ent, err := svc.mapper.PurchaseItem(req)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update purchase item", err)
	}
	ent.ID = id
	return svc.purchaseItemRepo.Update(ctx, ent)
}

func (svc purchaseItemService) Get(ctx *context.Context, id uint) (*responseModel.PurchaseItem, *errs.XError) {
	item, err := svc.purchaseItemRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.PurchaseItem(item)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map PurchaseItem data", mapErr)
	}
	return mapped, nil
}

func (svc purchaseItemService) GetByPurchaseId(ctx *context.Context, purchaseId uint) ([]responseModel.PurchaseItem, *errs.XError) {
	items, err := svc.purchaseItemRepo.GetByPurchaseId(ctx, purchaseId)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.PurchaseItems(items)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map PurchaseItem data", mapErr)
	}
	if mapped == nil {
		return []responseModel.PurchaseItem{}, nil
	}
	return mapped, nil
}

func (svc purchaseItemService) Delete(ctx *context.Context, id uint) *errs.XError {
	return svc.purchaseItemRepo.Delete(ctx, id)
}
