package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	"github.com/imkarthi24/sf-backend/internal/model/models"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type ExpenseTrackerService interface {
	SaveExpenseTracker(*context.Context, requestModel.ExpenseTracker) *errs.XError
	UpdateExpenseTracker(*context.Context, requestModel.ExpenseTracker, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.ExpenseTracker, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.ExpenseTracker, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type expenseTrackerService struct {
	expenseTrackerRepo repository.ExpenseTrackerRepository
	fileStoreSvc       FileStoreService
	entityDocumentSvc  EntityDocumentService
	mapper             mapper.Mapper
	respMapper         mapper.ResponseMapper
}

func ProvideExpenseTrackerService(repo repository.ExpenseTrackerRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper, fileStoreSvc FileStoreService, entityDocumentSvc EntityDocumentService) ExpenseTrackerService {
	return expenseTrackerService{
		expenseTrackerRepo: repo,
		fileStoreSvc:       fileStoreSvc,
		entityDocumentSvc:  entityDocumentSvc,
		mapper:             mapper,
		respMapper:        respMapper,
	}
}

func (svc expenseTrackerService) SaveExpenseTracker(ctx *context.Context, expenseTracker requestModel.ExpenseTracker) *errs.XError {
	dbExpenseTracker, err := svc.mapper.ExpenseTracker(expenseTracker)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save expense tracker", err)
	}

	errr := svc.expenseTrackerRepo.Create(ctx, dbExpenseTracker)
	if errr != nil {
		return errr
	}

	for _, f := range expenseTracker.Files {
		if xerr := svc.confirmTempFileForExpense(ctx, f, dbExpenseTracker.ID); xerr != nil {
			return xerr
		}
	}

	return svc.expenseTrackerRepo.RecalculateAndUpdateBalance(ctx, dbExpenseTracker.ID)
}

func (svc expenseTrackerService) UpdateExpenseTracker(ctx *context.Context, expenseTracker requestModel.ExpenseTracker, id uint) *errs.XError {
	dbExpenseTracker, err := svc.mapper.ExpenseTracker(expenseTracker)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update expense tracker", err)
	}

	dbExpenseTracker.ID = id
	errr := svc.expenseTrackerRepo.Update(ctx, dbExpenseTracker)
	if errr != nil {
		return errr
	}

	for _, f := range expenseTracker.Files {
		if xerr := svc.confirmTempFileForExpense(ctx, f, id); xerr != nil {
			return xerr
		}
	}

	return svc.expenseTrackerRepo.RecalculateAndUpdateBalance(ctx, id)
}

func (svc expenseTrackerService) Get(ctx *context.Context, id uint) (*responseModel.ExpenseTracker, *errs.XError) {
	expenseTracker, err := svc.expenseTrackerRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedExpenseTracker, mapErr := svc.respMapper.ExpenseTracker(expenseTracker)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map ExpenseTracker data", mapErr)
	}

	if xerr := svc.fillExpenseFiles(ctx, mappedExpenseTracker); xerr != nil {
		return nil, xerr
	}

	return mappedExpenseTracker, nil
}

func (svc expenseTrackerService) GetAll(ctx *context.Context, search string) ([]responseModel.ExpenseTracker, *errs.XError) {
	expenseTrackers, err := svc.expenseTrackerRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedExpenseTrackers, mapErr := svc.respMapper.ExpenseTrackers(expenseTrackers)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map ExpenseTracker data", mapErr)
	}

	for i := range mappedExpenseTrackers {
		if xerr := svc.fillExpenseFiles(ctx, &mappedExpenseTrackers[i]); xerr != nil {
			return nil, xerr
		}
	}

	return mappedExpenseTrackers, nil
}

func (svc expenseTrackerService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.expenseTrackerRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

const entityNameExpense = "Expense"

// confirmTempFileForExpense saves entity_document and confirms temp upload (moves file, updates file_store_metadata).
func (svc expenseTrackerService) confirmTempFileForExpense(ctx *context.Context, confirmFile models.ConfirmFile, expenseId uint) *errs.XError {
	entityDoc := requestModel.EntityDocuments{
		IsActive:     true,
		Type:         confirmFile.Kind,
		DocumentType: confirmFile.Kind,
		EntityName:   entityNameExpense,
		EntityId:     expenseId,
		Description:  confirmFile.Description,
	}
	if xerr := svc.entityDocumentSvc.SaveEntityDocument(ctx, entityDoc); xerr != nil {
		return xerr
	}
	newKey := svc.fileStoreSvc.GetFileKey(ctx, entityNameExpense, expenseId, confirmFile.Kind)
	newUrl, xerr := svc.fileStoreSvc.ConfirmTempUpload(ctx, confirmFile.FileKey, newKey)
	if xerr != nil {
		return xerr
	}
	return svc.fileStoreSvc.UpdateEntityIdAndKey(ctx, confirmFile.Id, expenseId, entityNameExpense, newKey, newUrl)
}

// fillExpenseFiles populates Files for the expense using entity_document + file_store_metadata (Expense/id/type).
func (svc expenseTrackerService) fillExpenseFiles(ctx *context.Context, expense *responseModel.ExpenseTracker) *errs.XError {
	files, xerr := svc.entityDocumentSvc.GetEntityDocumentsByEntity(ctx, expense.ID, entities.Entity_Expense, entities.EntityDocumentsType(""))
	if xerr != nil {
		return xerr
	}
	expense.Files = files
	return nil
}
