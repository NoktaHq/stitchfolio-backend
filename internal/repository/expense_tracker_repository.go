package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/loop-kar/pixie/constants"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
)

type ExpenseTrackerRepository interface {
	Create(*context.Context, *entities.Expense) *errs.XError
	Update(*context.Context, *entities.Expense) *errs.XError
	Get(*context.Context, uint) (*entities.Expense, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Expense, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	RecalculateAndUpdateBalance(*context.Context, uint) *errs.XError
}

type expenseTrackerRepository struct {
	GormDAL
}

func ProvideExpenseTrackerRepository(dal GormDAL) ExpenseTrackerRepository {
	return &expenseTrackerRepository{GormDAL: dal}
}

func (etr *expenseTrackerRepository) Create(ctx *context.Context, expenseTracker *entities.Expense) *errs.XError {
	res := etr.WithDB(ctx).Create(&expenseTracker)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save expense tracker", res.Error)
	}
	return nil
}

func (etr *expenseTrackerRepository) Update(ctx *context.Context, expenseTracker *entities.Expense) *errs.XError {
	return etr.GormDAL.Update(ctx, *expenseTracker)
}

func (etr *expenseTrackerRepository) Get(ctx *context.Context, id uint) (*entities.Expense, *errs.XError) {
	expenseTracker := entities.Expense{}
	res := etr.WithDB(ctx).Model(expenseTracker).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Preload("ExpenseDetails").
		Find(&expenseTracker, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find expense tracker", res.Error)
	}
	return &expenseTracker, nil
}

func (etr *expenseTrackerRepository) GetAll(ctx *context.Context, search string) ([]entities.Expense, *errs.XError) {
	var expenseTrackers []entities.Expense

	filterValue := util.ReadValueFromContext(ctx, constants.FILTER_KEY)
	var filter string
	if filterValue != nil {
		filter = filterValue.(string)
	}

	res := etr.WithDB(ctx).Model(entities.Expense{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.GetExpenseTrackers_Search(search)).
		Scopes(scopes.GetExpenseTrackers_Filter(filter)).
		Scopes(db.Paginate(ctx)).
		Scopes(scopes.WithAuditInfo()).
		Preload("ExpenseDetails").
		Find(&expenseTrackers)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find expense trackers", res.Error)
	}
	return expenseTrackers, nil
}

func (etr *expenseTrackerRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	expenseTracker := &entities.Expense{Model: &entities.Model{ID: id, IsActive: false}}
	err := etr.GormDAL.Delete(ctx, expenseTracker)
	if err != nil {
		return err
	}
	return nil
}

// RecalculateAndUpdateBalance sets Expense.Balance = Expense.Price - Sum(ExpenseDetail.Price) for the given expense.
// Call after add/update/delete of expense details (or expense price) to keep balance correct.
func (etr *expenseTrackerRepository) RecalculateAndUpdateBalance(ctx *context.Context, expenseId uint) *errs.XError {
	var expense entities.Expense
	res := etr.WithDB(ctx).Model(entities.Expense{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Select("id", "price").
		First(&expense, expenseId)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to find expense for balance recalc", res.Error)
	}
	var sumResult struct{ Total float64 }
	res = etr.WithDB(ctx).Model(entities.ExpenseDetail{}).
		Where("expense_id = ?", expenseId).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Select("COALESCE(SUM(price), 0) as total").
		Scan(&sumResult)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to sum expense details for balance", res.Error)
	}
	balance := expense.Price - sumResult.Total
	upd := map[string]interface{}{"balance": balance}
	if session := utils.GetSession(ctx); session != nil {
		upd["updated_by_id"] = session.UserId
	}
	res = etr.WithDB(ctx).Model(entities.Expense{}).Where("id = ?", expenseId).Updates(upd)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to update expense balance", res.Error)
	}
	return nil
}
