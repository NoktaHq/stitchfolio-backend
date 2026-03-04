package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type CategoryRepository interface {
	Create(*context.Context, *entities.Category) *errs.XError
	Update(*context.Context, *entities.Category) *errs.XError
	Get(*context.Context, uint) (*entities.Category, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Category, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	AutocompleteCategory(*context.Context, string) ([]entities.Category, *errs.XError)
}

type categoryRepository struct {
	GormDAL
}

func ProvideCategoryRepository(customDB GormDAL) CategoryRepository {
	return &categoryRepository{GormDAL: customDB}
}

func (cr *categoryRepository) Create(ctx *context.Context, category *entities.Category) *errs.XError {
	res := cr.WithDB(ctx).Create(&category)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save category", res.Error)
	}
	return nil
}

func (cr *categoryRepository) Update(ctx *context.Context, category *entities.Category) *errs.XError {
	return cr.GormDAL.Update(ctx, *category)
}

func (cr *categoryRepository) Get(ctx *context.Context, id uint) (*entities.Category, *errs.XError) {
	category := entities.Category{}
	res := cr.WithDB(ctx).Model(category).
		Scopes(scopes.WithAuditInfo()).
		Preload("Products").
		Find(&category, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find category", res.Error)
	}
	return &category, nil
}

func (cr *categoryRepository) GetAll(ctx *context.Context, search string) ([]entities.Category, *errs.XError) {
	var categories []entities.Category
	res := cr.WithDB(ctx).Model(entities.Category{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.ILike(search, "name")).
		Scopes(db.Paginate(ctx)).
		Find(&categories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find categories", res.Error)
	}
	return categories, nil
}

func (cr *categoryRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	category := &entities.Category{Model: &entities.Model{ID: id, IsActive: false}}
	err := cr.GormDAL.Delete(ctx, category)
	if err != nil {
		return err
	}
	return nil
}

func (cr *categoryRepository) AutocompleteCategory(ctx *context.Context, search string) ([]entities.Category, *errs.XError) {
	var categories []entities.Category
	res := cr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "name")).
		Select("id", "name").
		Find(&categories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find categories for autocomplete", res.Error)
	}
	return categories, nil
}
