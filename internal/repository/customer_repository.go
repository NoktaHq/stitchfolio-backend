package repository

import (
	"context"
	"errors"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(*context.Context, *entities.Customer) *errs.XError
	Update(*context.Context, *entities.Customer) *errs.XError
	Get(*context.Context, uint) (*entities.Customer, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Customer, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	GetByPhoneNumber(*context.Context, string) (*entities.Customer, *errs.XError)
	AutocompleteCustomer(*context.Context, string) ([]entities.Customer, *errs.XError)
}

type customerRepository struct {
	GormDAL
}

func ProvideCustomerRepository(customDB GormDAL) CustomerRepository {
	return &customerRepository{GormDAL: customDB}
}

func (cr *customerRepository) Create(ctx *context.Context, customer *entities.Customer) *errs.XError {
	res := cr.WithDB(ctx).Create(&customer)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save customer", res.Error)
	}
	return nil
}

func (cr *customerRepository) Update(ctx *context.Context, customer *entities.Customer) *errs.XError {
	return cr.GormDAL.Update(ctx, *customer)
}

func (cr *customerRepository) Get(ctx *context.Context, id uint) (*entities.Customer, *errs.XError) {
	customer := entities.Customer{}
	res := cr.WithDB(ctx).Model(customer).
		Scopes(scopes.WithAuditInfo()).
		Preload("Persons").
		Preload("Persons.Measurements").
		Preload("Persons.Measurements.DressType").
		Preload("Enquiries").
		Preload("Orders").
		Find(&customer, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find customer", res.Error)
	}
	return &customer, nil
}

func (cr *customerRepository) GetAll(ctx *context.Context, search string) ([]entities.Customer, *errs.XError) {
	var customers []entities.Customer
	res := cr.WithDB(ctx).Model(entities.Customer{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "first_name", "last_name", "email", "phone_number")).
		Scopes(db.Paginate(ctx)).
		Find(&customers)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find customers", res.Error)
	}
	return customers, nil
}

func (cr *customerRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	customer := &entities.Customer{Model: &entities.Model{ID: id, IsActive: false}}
	err := cr.GormDAL.Delete(ctx, customer)
	if err != nil {
		return err
	}
	return nil
}

func (cr *customerRepository) GetByPhoneNumber(ctx *context.Context, phoneNumber string) (*entities.Customer, *errs.XError) {
	if phoneNumber == "" {
		return nil, nil
	}
	customer := entities.Customer{}
	res := cr.WithDB(ctx).Where("phone_number = ?", phoneNumber).First(&customer)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewXError(errs.DATABASE, "Unable to find customer by phone number", res.Error)
	}
	return &customer, nil
}

func (cr *customerRepository) AutocompleteCustomer(ctx *context.Context, search string) ([]entities.Customer, *errs.XError) {
	var customers []entities.Customer
	res := cr.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "first_name", "last_name", "email", "phone_number")).
		Select("id", "first_name", "last_name", "email", "phone_number").
		Find(&customers)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find customers for autocomplete", res.Error)
	}
	return customers, nil
}
