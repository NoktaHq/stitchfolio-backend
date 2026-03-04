package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type ProductService interface {
	SaveProduct(*context.Context, requestModel.Product) *errs.XError
	UpdateProduct(*context.Context, requestModel.Product, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Product, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Product, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	AutocompleteProduct(*context.Context, string) ([]responseModel.ProductAutoComplete, *errs.XError)
	GetBySKU(*context.Context, string) (*responseModel.Product, *errs.XError)
	GetLowStockProducts(*context.Context) ([]responseModel.Product, *errs.XError)
}

type productService struct {
	productRepo   repository.ProductRepository
	inventoryRepo repository.InventoryRepository
	mapper        mapper.Mapper
	respMapper    mapper.ResponseMapper
}

func ProvideProductService(
	repo repository.ProductRepository,
	inventoryRepo repository.InventoryRepository,
	mapper mapper.Mapper,
	respMapper mapper.ResponseMapper,
) ProductService {
	return productService{
		productRepo:   repo,
		inventoryRepo: inventoryRepo,
		mapper:        mapper,
		respMapper:    respMapper,
	}
}

func (svc productService) SaveProduct(ctx *context.Context, product requestModel.Product) *errs.XError {
	dbProduct, err := svc.mapper.Product(product)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save product", err)
	}

	errr := svc.productRepo.Create(ctx, dbProduct)
	if errr != nil {
		return errr
	}

	// Auto-create inventory entry for the product
	inventory := &entities.Inventory{
		Model:             &entities.Model{IsActive: true},
		ProductId:         dbProduct.ID,
		Quantity:          0,
		LowStockThreshold: product.LowStockThreshold,
	}

	errr = svc.inventoryRepo.Create(ctx, inventory)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc productService) UpdateProduct(ctx *context.Context, product requestModel.Product, id uint) *errs.XError {
	dbProduct, err := svc.mapper.Product(product)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update product and threshold", err)
	}

	dbProduct.ID = id
	errr := svc.productRepo.Update(ctx, dbProduct)
	if errr != nil {
		return errr
	}

	errr = svc.inventoryRepo.UpdateThreshold(ctx, id, product.LowStockThreshold)
	if errr != nil {
		return errr
	}

	return nil
}

func (svc productService) Get(ctx *context.Context, id uint) (*responseModel.Product, *errs.XError) {
	product, err := svc.productRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedProduct, mapErr := svc.respMapper.Product(product)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Product data", mapErr)
	}

	return mappedProduct, nil
}

func (svc productService) GetAll(ctx *context.Context, search string) ([]responseModel.Product, *errs.XError) {
	products, err := svc.productRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedProducts, mapErr := svc.respMapper.Products(products)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Product data", mapErr)
	}

	return mappedProducts, nil
}

func (svc productService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.productRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (svc productService) AutocompleteProduct(ctx *context.Context, search string) ([]responseModel.ProductAutoComplete, *errs.XError) {
	products, err := svc.productRepo.AutocompleteProduct(ctx, search)
	if err != nil {
		return nil, err
	}

	res := make([]responseModel.ProductAutoComplete, 0)
	for _, product := range products {
		currentStock := 0
		isLowStock := false
		if product.Inventory != nil {
			currentStock = product.Inventory.Quantity
			isLowStock = product.Inventory.IsLowStock()
		}

		res = append(res, responseModel.ProductAutoComplete{
			ID:           product.ID,
			Name:         product.Name,
			SKU:          product.SKU,
			CurrentStock: currentStock,
			IsLowStock:   isLowStock,
		})
	}

	return res, nil
}

func (svc productService) GetBySKU(ctx *context.Context, sku string) (*responseModel.Product, *errs.XError) {
	product, err := svc.productRepo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}

	mappedProduct, mapErr := svc.respMapper.Product(product)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Product data", mapErr)
	}

	return mappedProduct, nil
}

func (svc productService) GetLowStockProducts(ctx *context.Context) ([]responseModel.Product, *errs.XError) {
	products, err := svc.productRepo.GetLowStockProducts(ctx)
	if err != nil {
		return nil, err
	}

	mappedProducts, mapErr := svc.respMapper.Products(products)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Product data", mapErr)
	}

	return mappedProducts, nil
}
