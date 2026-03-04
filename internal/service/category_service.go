package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type CategoryService interface {
	SaveCategory(*context.Context, requestModel.Category) *errs.XError
	UpdateCategory(*context.Context, requestModel.Category, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Category, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Category, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	AutocompleteCategory(*context.Context, string) ([]responseModel.CategoryAutoComplete, *errs.XError)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	mapper       mapper.Mapper
	respMapper   mapper.ResponseMapper
}

func ProvideCategoryService(
	repo repository.CategoryRepository,
	mapper mapper.Mapper,
	respMapper mapper.ResponseMapper,
) CategoryService {
	return categoryService{
		categoryRepo: repo,
		mapper:       mapper,
		respMapper:   respMapper,
	}
}

func (svc categoryService) SaveCategory(ctx *context.Context, category requestModel.Category) *errs.XError {
	dbCategory, err := svc.mapper.Category(category)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save category", err)
	}

	errr := svc.categoryRepo.Create(ctx, dbCategory)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc categoryService) UpdateCategory(ctx *context.Context, category requestModel.Category, id uint) *errs.XError {
	dbCategory, err := svc.mapper.Category(category)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update category", err)
	}

	dbCategory.ID = id
	errr := svc.categoryRepo.Update(ctx, dbCategory)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc categoryService) Get(ctx *context.Context, id uint) (*responseModel.Category, *errs.XError) {
	category, err := svc.categoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedCategory, mapErr := svc.respMapper.Category(category)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Category data", mapErr)
	}

	return mappedCategory, nil
}

func (svc categoryService) GetAll(ctx *context.Context, search string) ([]responseModel.Category, *errs.XError) {
	categories, err := svc.categoryRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedCategories, mapErr := svc.respMapper.Categories(categories)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Category data", mapErr)
	}

	return mappedCategories, nil
}

func (svc categoryService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.categoryRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (svc categoryService) AutocompleteCategory(ctx *context.Context, search string) ([]responseModel.CategoryAutoComplete, *errs.XError) {
	categories, err := svc.categoryRepo.AutocompleteCategory(ctx, search)
	if err != nil {
		return nil, err
	}

	res := make([]responseModel.CategoryAutoComplete, 0)
	for _, category := range categories {
		res = append(res, responseModel.CategoryAutoComplete{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return res, nil
}
