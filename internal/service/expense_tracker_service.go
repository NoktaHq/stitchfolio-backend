package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
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
	mapper             mapper.Mapper
	respMapper         mapper.ResponseMapper
}

func ProvideExpenseTrackerService(repo repository.ExpenseTrackerRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) ExpenseTrackerService {
	return expenseTrackerService{
		expenseTrackerRepo: repo,
		mapper:             mapper,
		respMapper:         respMapper,
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

	return mappedExpenseTrackers, nil
}

func (svc expenseTrackerService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.expenseTrackerRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
