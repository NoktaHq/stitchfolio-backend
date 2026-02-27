package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type ExpenseDetailService interface {
	Save(*context.Context, requestModel.ExpenseDetail, uint) *errs.XError
	Update(*context.Context, requestModel.ExpenseDetail, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.ExpenseDetail, *errs.XError)
	GetByExpenseId(*context.Context, uint) ([]responseModel.ExpenseDetail, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type expenseDetailService struct {
	expenseDetailRepo repository.ExpenseDetailRepository
	mapper            mapper.Mapper
	respMapper        mapper.ResponseMapper
}

func ProvideExpenseDetailService(repo repository.ExpenseDetailRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) ExpenseDetailService {
	return &expenseDetailService{
		expenseDetailRepo: repo,
		mapper:            mapper,
		respMapper:        respMapper,
	}
}

func (svc *expenseDetailService) Save(ctx *context.Context, req requestModel.ExpenseDetail, expenseId uint) *errs.XError {
	req.ExpenseId = expenseId
	ent, err := svc.mapper.ExpenseDetail(req)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save expense detail", err)
	}
	return svc.expenseDetailRepo.Create(ctx, ent)
}

func (svc *expenseDetailService) Update(ctx *context.Context, req requestModel.ExpenseDetail, id uint) *errs.XError {
	ent, err := svc.mapper.ExpenseDetail(req)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update expense detail", err)
	}
	ent.ID = id
	return svc.expenseDetailRepo.Update(ctx, ent)
}

func (svc *expenseDetailService) Get(ctx *context.Context, id uint) (*responseModel.ExpenseDetail, *errs.XError) {
	detail, err := svc.expenseDetailRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.ExpenseDetail(detail)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map expense detail", mapErr)
	}
	return mapped, nil
}

func (svc *expenseDetailService) GetByExpenseId(ctx *context.Context, expenseId uint) ([]responseModel.ExpenseDetail, *errs.XError) {
	details, err := svc.expenseDetailRepo.GetByExpenseId(ctx, expenseId)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.ExpenseDetails(details)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map expense details", mapErr)
	}
	if mapped == nil {
		return []responseModel.ExpenseDetail{}, nil
	}
	return mapped, nil
}

func (svc *expenseDetailService) Delete(ctx *context.Context, id uint) *errs.XError {
	return svc.expenseDetailRepo.Delete(ctx, id)
}
