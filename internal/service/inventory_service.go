package service

import (
	"context"
	"fmt"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
)

type InventoryService interface {
	Get(*context.Context, uint) (*responseModel.Inventory, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Inventory, *errs.XError)
	GetByProductId(*context.Context, uint) (*responseModel.Inventory, *errs.XError)
	UpdateThreshold(*context.Context, requestModel.Inventory, uint) *errs.XError
	GetLowStockItems(*context.Context) ([]responseModel.LowStockItem, *errs.XError)

	// Stock movement operations
	RecordStockMovement(*context.Context, requestModel.StockMovementRequest) (*responseModel.StockMovementResponse, *errs.XError)
}

type inventoryService struct {
	inventoryRepo    repository.InventoryRepository
	inventoryLogRepo repository.InventoryLogRepository
	productRepo      repository.ProductRepository
	mapper           mapper.Mapper
	respMapper       mapper.ResponseMapper
}

func ProvideInventoryService(
	repo repository.InventoryRepository,
	logRepo repository.InventoryLogRepository,
	productRepo repository.ProductRepository,
	mapper mapper.Mapper,
	respMapper mapper.ResponseMapper,
) InventoryService {
	return inventoryService{
		inventoryRepo:    repo,
		inventoryLogRepo: logRepo,
		productRepo:      productRepo,
		mapper:           mapper,
		respMapper:       respMapper,
	}
}

func (svc inventoryService) Get(ctx *context.Context, id uint) (*responseModel.Inventory, *errs.XError) {
	inventory, err := svc.inventoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedInventory, mapErr := svc.respMapper.Inventory(inventory)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Inventory data", mapErr)
	}

	return mappedInventory, nil
}

func (svc inventoryService) GetAll(ctx *context.Context, search string) ([]responseModel.Inventory, *errs.XError) {
	inventories, err := svc.inventoryRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedInventories, mapErr := svc.respMapper.Inventories(inventories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Inventory data", mapErr)
	}

	return mappedInventories, nil
}

func (svc inventoryService) GetByProductId(ctx *context.Context, productId uint) (*responseModel.Inventory, *errs.XError) {
	inventory, err := svc.inventoryRepo.GetByProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	mappedInventory, mapErr := svc.respMapper.Inventory(inventory)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Inventory data", mapErr)
	}

	return mappedInventory, nil
}

func (svc inventoryService) UpdateThreshold(ctx *context.Context, inventory requestModel.Inventory, id uint) *errs.XError {
	// Get current inventory
	currentInventory, err := svc.inventoryRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	// Update only the threshold
	errr := svc.inventoryRepo.UpdateThreshold(ctx, currentInventory.ProductId, inventory.LowStockThreshold)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc inventoryService) GetLowStockItems(ctx *context.Context) ([]responseModel.LowStockItem, *errs.XError) {
	inventories, err := svc.inventoryRepo.GetLowStockItems(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]responseModel.LowStockItem, 0)
	for _, inv := range inventories {
		categoryName := ""
		if inv.Product != nil && inv.Product.Category != nil {
			categoryName = inv.Product.Category.Name
		}

		productName := ""
		productSKU := ""
		if inv.Product != nil {
			productName = inv.Product.Name
			productSKU = inv.Product.SKU
		}

		res = append(res, responseModel.LowStockItem{
			ProductId:         inv.ProductId,
			ProductName:       productName,
			ProductSKU:        productSKU,
			CurrentStock:      inv.Quantity,
			LowStockThreshold: inv.LowStockThreshold,
			CategoryName:      categoryName,
		})
	}

	return res, nil
}

// RecordStockMovement handles all stock movements (IN, OUT, ADJUST) with business rules
func (svc inventoryService) RecordStockMovement(ctx *context.Context, request requestModel.StockMovementRequest) (*responseModel.StockMovementResponse, *errs.XError) {
	// Validation
	if request.Quantity <= 0 {
		return nil, errs.NewXError(errs.INVALID_REQUEST, "Quantity must be greater than 0", nil)
	}

	changeType := entities.InventoryLogChangeType(request.ChangeType)
	if changeType != entities.InventoryLogChangeTypeIN &&
		changeType != entities.InventoryLogChangeTypeOUT &&
		changeType != entities.InventoryLogChangeTypeADJUST {
		return nil, errs.NewXError(errs.INVALID_REQUEST, "Invalid change type. Must be IN, OUT, or ADJUST", nil)
	}

	// Get current inventory
	inventory, err := svc.inventoryRepo.GetByProductId(ctx, request.ProductId)
	if err != nil {
		return nil, errs.NewXError(errs.INVALID_REQUEST, "Product inventory not found", err)
	}

	previousStock := inventory.Quantity

	// Calculate new stock based on change type
	var newStock int
	var netChange int

	switch changeType {
	case entities.InventoryLogChangeTypeIN:
		newStock = previousStock + request.Quantity
		netChange = request.Quantity

	case entities.InventoryLogChangeTypeOUT:
		netChange = -request.Quantity
		newStock = previousStock + netChange

		// Prevent negative stock unless admin override
		if newStock < 0 && !request.AdminOverride {
			return nil, errs.NewXError(
				errs.INVALID_REQUEST,
				fmt.Sprintf("Insufficient stock. Available: %d, Requested: %d", previousStock, request.Quantity),
				nil,
			)
		}

	case entities.InventoryLogChangeTypeADJUST:
		// For ADJUST, the quantity can be positive (add) or negative (remove)
		// We treat the request.Quantity as the adjustment amount
		if request.Quantity > 0 {
			netChange = request.Quantity
		} else {
			netChange = -request.Quantity // Make it negative for removal
		}
		newStock = previousStock + netChange
	}

	// Create inventory log entry
	logEntry := &entities.InventoryLog{
		Model:      &entities.Model{IsActive: true},
		ProductId:  request.ProductId,
		ChangeType: changeType,
		Quantity:   request.Quantity,
		Reason:     request.Reason,
		Notes:      request.Notes,
		LoggedAt:   util.GetLocalTime(),
	}

	errr := svc.inventoryLogRepo.Create(ctx, logEntry)
	if errr != nil {
		return nil, errs.NewXError(errs.DATABASE, "Failed to create inventory log", errr)
	}

	// Update inventory quantity
	errr = svc.inventoryRepo.UpdateQuantity(ctx, request.ProductId, newStock)
	if errr != nil {
		return nil, errs.NewXError(errs.DATABASE, "Failed to update inventory quantity", errr)
	}

	// Return response
	response := &responseModel.StockMovementResponse{
		Success:       true,
		Message:       fmt.Sprintf("Stock %s recorded successfully", request.ChangeType),
		ProductId:     request.ProductId,
		PreviousStock: previousStock,
		NewStock:      newStock,
		ChangeAmount:  netChange,
	}

	return response, nil
}
