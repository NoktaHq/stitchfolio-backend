package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type InventoryLogService interface {
	Get(*context.Context, uint) (*responseModel.InventoryLog, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.InventoryLog, *errs.XError)
	GetByProductId(*context.Context, uint) ([]responseModel.InventoryLog, *errs.XError)
	GetByChangeType(*context.Context, string) ([]responseModel.InventoryLog, *errs.XError)
	GetByDateRange(*context.Context, string, string) ([]responseModel.InventoryLog, *errs.XError)
}

type inventoryLogService struct {
	inventoryLogRepo repository.InventoryLogRepository
	inventoryRepo    repository.InventoryRepository
	respMapper       mapper.ResponseMapper
}

func ProvideInventoryLogService(
	repo repository.InventoryLogRepository,
	inventoryRepo repository.InventoryRepository,
	respMapper mapper.ResponseMapper,
) InventoryLogService {
	return inventoryLogService{
		inventoryLogRepo: repo,
		inventoryRepo:    inventoryRepo,
		respMapper:       respMapper,
	}
}

func (svc inventoryLogService) Get(ctx *context.Context, id uint) (*responseModel.InventoryLog, *errs.XError) {
	log, err := svc.inventoryLogRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedLog, mapErr := svc.respMapper.InventoryLog(log)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map InventoryLog data", mapErr)
	}

	return mappedLog, nil
}

func (svc inventoryLogService) GetAll(ctx *context.Context, search string) ([]responseModel.InventoryLog, *errs.XError) {
	logs, err := svc.inventoryLogRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedLogs, mapErr := svc.respMapper.InventoryLogs(logs)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map InventoryLog data", mapErr)
	}

	return mappedLogs, nil
}

func (svc inventoryLogService) GetByProductId(ctx *context.Context, productId uint) ([]responseModel.InventoryLog, *errs.XError) {
	logs, err := svc.inventoryLogRepo.GetByProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	mappedLogs, mapErr := svc.respMapper.InventoryLogs(logs)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map InventoryLog data", mapErr)
	}

	return mappedLogs, nil
}

func (svc inventoryLogService) GetByChangeType(ctx *context.Context, changeType string) ([]responseModel.InventoryLog, *errs.XError) {
	logs, err := svc.inventoryLogRepo.GetByChangeType(ctx, entities.InventoryLogChangeType(changeType))
	if err != nil {
		return nil, err
	}

	mappedLogs, mapErr := svc.respMapper.InventoryLogs(logs)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map InventoryLog data", mapErr)
	}

	return mappedLogs, nil
}

func (svc inventoryLogService) GetByDateRange(ctx *context.Context, startDate string, endDate string) ([]responseModel.InventoryLog, *errs.XError) {
	logs, err := svc.inventoryLogRepo.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	mappedLogs, mapErr := svc.respMapper.InventoryLogs(logs)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map InventoryLog data", mapErr)
	}

	return mappedLogs, nil
}
