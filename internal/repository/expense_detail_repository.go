package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/errs"
)

type ExpenseDetailRepository interface {
	Create(*context.Context, *entities.ExpenseDetail) *errs.XError
	Update(*context.Context, *entities.ExpenseDetail) *errs.XError
	Get(*context.Context, uint) (*entities.ExpenseDetail, *errs.XError)
	GetByExpenseId(*context.Context, uint) ([]entities.ExpenseDetail, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type expenseDetailRepository struct {
	GormDAL
}

func ProvideExpenseDetailRepository(dal GormDAL) ExpenseDetailRepository {
	return &expenseDetailRepository{GormDAL: dal}
}

func (edr *expenseDetailRepository) Create(ctx *context.Context, detail *entities.ExpenseDetail) *errs.XError {
	res := edr.WithDB(ctx).Create(detail)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save expense detail", res.Error)
	}
	return nil
}

func (edr *expenseDetailRepository) Update(ctx *context.Context, detail *entities.ExpenseDetail) *errs.XError {
	return edr.GormDAL.Update(ctx, *detail)
}

func (edr *expenseDetailRepository) Get(ctx *context.Context, id uint) (*entities.ExpenseDetail, *errs.XError) {
	detail := entities.ExpenseDetail{}
	res := edr.WithDB(ctx).Model(entities.ExpenseDetail{}).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Find(&detail, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find expense detail", res.Error)
	}
	return &detail, nil
}

func (edr *expenseDetailRepository) GetByExpenseId(ctx *context.Context, expenseId uint) ([]entities.ExpenseDetail, *errs.XError) {
	var details []entities.ExpenseDetail
	res := edr.WithDB(ctx).Model(entities.ExpenseDetail{}).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("expense_id = ?", expenseId).
		Find(&details)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find expense details", res.Error)
	}
	return details, nil
}

func (edr *expenseDetailRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	detail := &entities.ExpenseDetail{Model: &entities.Model{ID: id, IsActive: false}}
	err := edr.GormDAL.Delete(ctx, detail)
	if err != nil {
		return err
	}
	return nil
}
