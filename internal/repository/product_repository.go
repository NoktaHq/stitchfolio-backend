package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type ProductRepository interface {
	Create(*context.Context, *entities.Product) *errs.XError
	Update(*context.Context, *entities.Product) *errs.XError
	Get(*context.Context, uint) (*entities.Product, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Product, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	AutocompleteProduct(*context.Context, string) ([]entities.Product, *errs.XError)
	GetBySKU(*context.Context, string) (*entities.Product, *errs.XError)
	GetLowStockProducts(*context.Context) ([]entities.Product, *errs.XError)
}

type productRepository struct {
	GormDAL
}

func ProvideProductRepository(customDB GormDAL) ProductRepository {
	return &productRepository{GormDAL: customDB}
}

func (pr *productRepository) Create(ctx *context.Context, product *entities.Product) *errs.XError {
	res := pr.WithDB(ctx).Create(&product)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save product", res.Error)
	}
	return nil
}

func (pr *productRepository) Update(ctx *context.Context, product *entities.Product) *errs.XError {
	return pr.GormDAL.Update(ctx, *product)
}

func (pr *productRepository) Get(ctx *context.Context, id uint) (*entities.Product, *errs.XError) {
	product := entities.Product{}
	res := pr.WithDB(ctx).Model(product).
		Preload("Category").
		Preload("Inventory").
		Find(&product, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find product", res.Error)
	}
	return &product, nil
}

func (pr *productRepository) GetAll(ctx *context.Context, search string) ([]entities.Product, *errs.XError) {
	var products []entities.Product
	res := pr.WithDB(ctx).Model(entities.Product{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "name", "sku", "description")).
		Scopes(db.Paginate(ctx)).
		Scopes(scopes.WithAuditInfo()).
		Preload("Category").
		Preload("Inventory").
		Find(&products)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find products", res.Error)
	}
	return products, nil
}

func (pr *productRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	product := &entities.Product{Model: &entities.Model{ID: id, IsActive: false}}
	err := pr.GormDAL.Delete(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (pr *productRepository) AutocompleteProduct(ctx *context.Context, search string) ([]entities.Product, *errs.XError) {
	var products []entities.Product
	res := pr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "name", "sku")).
		Select("id", "name", "sku").
		Preload("Inventory").
		Find(&products)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find products for autocomplete", res.Error)
	}
	return products, nil
}

func (pr *productRepository) GetBySKU(ctx *context.Context, sku string) (*entities.Product, *errs.XError) {
	product := entities.Product{}
	res := pr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("sku = ?", sku).
		Preload("Category").
		Preload("Inventory").
		First(&product)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find product by SKU", res.Error)
	}
	return &product, nil
}

func (pr *productRepository) GetLowStockProducts(ctx *context.Context) ([]entities.Product, *errs.XError) {
	var products []entities.Product
	res := pr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Joins("INNER JOIN \"stich\".\"Inventories\" ON \"stich\".\"Inventories\".product_id = \"stich\".\"Products\".id").
		Where("\"stich\".\"Inventories\".quantity <= \"stich\".\"Inventories\".low_stock_threshold").
		Preload("Category").
		Preload("Inventory").
		Find(&products)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find low stock products", res.Error)
	}
	return products, nil
}
