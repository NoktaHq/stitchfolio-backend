package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/constants"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
)

type TaskRepository interface {
	Create(*context.Context, *entities.Task) *errs.XError
	Update(*context.Context, *entities.Task) *errs.XError
	Get(*context.Context, uint) (*entities.Task, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Task, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type taskRepository struct {
	GormDAL
}

func ProvideTaskRepository(dal GormDAL) TaskRepository {
	return &taskRepository{GormDAL: dal}
}

func (tr *taskRepository) Create(ctx *context.Context, task *entities.Task) *errs.XError {
	res := tr.WithDB(ctx).Create(&task)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save task", res.Error)
	}
	return nil
}

func (tr *taskRepository) Update(ctx *context.Context, task *entities.Task) *errs.XError {
	return tr.GormDAL.Update(ctx, *task)
}

func (tr *taskRepository) Get(ctx *context.Context, id uint) (*entities.Task, *errs.XError) {
	task := entities.Task{}
	res := tr.WithDB(ctx).Model(task).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Find(&task, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find task", res.Error)
	}
	return &task, nil
}

func (tr *taskRepository) GetAll(ctx *context.Context, search string) ([]entities.Task, *errs.XError) {
	var tasks []entities.Task

	filterValue := util.ReadValueFromContext(ctx, constants.FILTER_KEY)
	var filter string
	if filterValue != nil {
		filter = filterValue.(string)
	}

	res := tr.WithDB(ctx).Model(entities.Task{}).
		Scopes(scopes.Channel(), scopes.IsActive(), scopes.TasksForCurrentUser()).
		Scopes(scopes.GetTasks_Search(search)).
		Scopes(scopes.GetTasks_Filter(filter)).
		Scopes(scopes.WithAuditInfo()).
		Scopes(db.Paginate(ctx)).
		Find(&tasks)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find tasks", res.Error)
	}
	return tasks, nil
}

func (tr *taskRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	task := &entities.Task{Model: &entities.Model{ID: id, IsActive: false}}
	err := tr.GormDAL.Delete(ctx, task)
	if err != nil {
		return err
	}
	return nil
}
