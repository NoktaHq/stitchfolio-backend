package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
)

type PurchaseService interface {
	SavePurchase(*context.Context, requestModel.Purchase) *errs.XError
	UpdatePurchase(*context.Context, requestModel.Purchase, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Purchase, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Purchase, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	ReceivePurchase(*context.Context, uint) (*responseModel.ReceivePurchaseResponse, *errs.XError)
}

type purchaseService struct {
	purchaseRepo     repository.PurchaseRepository
	purchaseItemRepo repository.PurchaseItemRepository
	inventoryRepo    repository.InventoryRepository
	inventoryLogRepo repository.InventoryLogRepository
	mapper           mapper.Mapper
	respMapper       mapper.ResponseMapper
}

func ProvidePurchaseService(
	purchaseRepo repository.PurchaseRepository,
	purchaseItemRepo repository.PurchaseItemRepository,
	inventoryRepo repository.InventoryRepository,
	inventoryLogRepo repository.InventoryLogRepository,
	m mapper.Mapper,
	respMapper mapper.ResponseMapper,
) PurchaseService {
	return purchaseService{
		purchaseRepo:     purchaseRepo,
		purchaseItemRepo: purchaseItemRepo,
		inventoryRepo:    inventoryRepo,
		inventoryLogRepo: inventoryLogRepo,
		mapper:           m,
		respMapper:       respMapper,
	}
}

func (svc purchaseService) SavePurchase(ctx *context.Context, purchase requestModel.Purchase) *errs.XError {
	ent, err := svc.mapper.Purchase(purchase)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save purchase", err)
	}
	if ent.Status == "" {
		ent.Status = entities.PurchaseStatusDRAFT
	}
	return svc.purchaseRepo.Create(ctx, ent)
}

func (svc purchaseService) UpdatePurchase(ctx *context.Context, purchase requestModel.Purchase, id uint) *errs.XError {
	ent, err := svc.mapper.Purchase(purchase)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update purchase", err)
	}
	ent.ID = id
	return svc.purchaseRepo.Update(ctx, ent)
}

func (svc purchaseService) Get(ctx *context.Context, id uint) (*responseModel.Purchase, *errs.XError) {
	purchase, err := svc.purchaseRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.Purchase(purchase)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Purchase data", mapErr)
	}
	return mapped, nil
}

func (svc purchaseService) GetAll(ctx *context.Context, search string) ([]responseModel.Purchase, *errs.XError) {
	purchases, err := svc.purchaseRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.Purchases(purchases)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Purchase data", mapErr)
	}
	return mapped, nil
}

func (svc purchaseService) Delete(ctx *context.Context, id uint) *errs.XError {
	return svc.purchaseRepo.Delete(ctx, id)
}

// ReceivePurchase receives all pending quantities for the purchase: for each line adds (QuantityOrdered - QuantityReceived)
// to inventory, creates InventoryLog with SourceType PURCHASE and SourceId = PurchaseItem.ID, then sets QuantityReceived = QuantityOrdered.
func (svc purchaseService) ReceivePurchase(ctx *context.Context, purchaseId uint) (*responseModel.ReceivePurchaseResponse, *errs.XError) {
	purchase, err := svc.purchaseRepo.Get(ctx, purchaseId)
	if err != nil {
		return nil, err
	}
	if purchase.Status == entities.PurchaseStatusCANCELLED {
		return nil, errs.NewXError(errs.INVALID_REQUEST, "Cannot receive a cancelled purchase", nil)
	}

	var receivedLines int
	for i := range purchase.PurchaseItems {
		item := &purchase.PurchaseItems[i]
		toReceive := item.QuantityOrdered - item.QuantityReceived
		if toReceive <= 0 {
			continue
		}

		inv, invErr := svc.inventoryRepo.GetByProductId(ctx, item.ProductId)
		if invErr != nil {
			return nil, errs.NewXError(errs.DATABASE, "Failed to get inventory for product", invErr)
		}
		newQty := inv.Quantity + toReceive

		logEntry := &entities.InventoryLog{
			Model:      &entities.Model{IsActive: true},
			ProductId:  item.ProductId,
			ChangeType: entities.InventoryLogChangeTypeIN,
			Quantity:   toReceive,
			Reason:     "Purchase receipt",
			Notes:      "Purchase #" + purchase.PurchaseNumber,
			LoggedAt:   util.GetLocalTime(),
			SourceType: entities.InventoryLogSourceTypePURCHASE,
			SourceId:   &item.ID,
		}
		if logErr := svc.inventoryLogRepo.Create(ctx, logEntry); logErr != nil {
			return nil, errs.NewXError(errs.DATABASE, "Failed to create inventory log", logErr)
		}
		if updErr := svc.inventoryRepo.UpdateQuantity(ctx, item.ProductId, newQty); updErr != nil {
			return nil, errs.NewXError(errs.DATABASE, "Failed to update inventory quantity", updErr)
		}

		item.QuantityReceived = item.QuantityOrdered
		if updErr := svc.purchaseItemRepo.Update(ctx, item); updErr != nil {
			return nil, errs.NewXError(errs.DATABASE, "Failed to update purchase item", updErr)
		}
		receivedLines++
	}

	allReceived := true
	for _, item := range purchase.PurchaseItems {
		if item.QuantityReceived < item.QuantityOrdered {
			allReceived = false
			break
		}
	}
	newStatus := entities.PurchaseStatusPARTIALLY_RECEIVED
	if allReceived {
		newStatus = entities.PurchaseStatusRECEIVED
	}
	purchase.Status = newStatus
	if updateErr := svc.purchaseRepo.Update(ctx, purchase); updateErr != nil {
		return nil, updateErr
	}

	return &responseModel.ReceivePurchaseResponse{
		Success:        true,
		Message:        "Purchase received successfully",
		PurchaseId:     purchaseId,
		LinesReceived:  receivedLines,
		NewStatus:      string(newStatus),
	}, nil
}
