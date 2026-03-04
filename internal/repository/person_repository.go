package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type PersonRepository interface {
	Create(*context.Context, *entities.Person) *errs.XError
	Update(*context.Context, *entities.Person) *errs.XError
	Get(*context.Context, uint) (*entities.Person, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Person, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	GetByCustomerId(*context.Context, uint) ([]entities.Person, *errs.XError)
}

type personRepository struct {
	GormDAL
}

func ProvidePersonRepository(customDB GormDAL) PersonRepository {
	return &personRepository{GormDAL: customDB}
}

func (pr *personRepository) Create(ctx *context.Context, person *entities.Person) *errs.XError {
	res := pr.WithDB(ctx).Create(&person)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save person", res.Error)
	}
	return nil
}

func (pr *personRepository) Update(ctx *context.Context, person *entities.Person) *errs.XError {
	return pr.GormDAL.Update(ctx, *person)
}

func (pr *personRepository) Get(ctx *context.Context, id uint) (*entities.Person, *errs.XError) {
	person := entities.Person{}
	res := pr.WithDB(ctx).
		Preload("Customer").
		Preload("Measurements", scopes.IsActive()).
		Preload("Measurements.DressType", scopes.SelectFields("name")).
		Find(&person, id)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find person", res.Error)
	}
	return &person, nil
}

func (pr *personRepository) GetAll(ctx *context.Context, search string) ([]entities.Person, *errs.XError) {
	var persons []entities.Person
	res := pr.WithDB(ctx).Model(entities.Person{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "first_name", "last_name")).
		Scopes(db.Paginate(ctx)).
		Preload("Customer").
		Find(&persons)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find persons", res.Error)
	}
	return persons, nil
}

func (pr *personRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	person := &entities.Person{Model: &entities.Model{ID: id, IsActive: false}}
	err := pr.GormDAL.Delete(ctx, person)
	if err != nil {
		return err
	}
	return nil
}

func (pr *personRepository) GetByCustomerId(ctx *context.Context, customerId uint) ([]entities.Person, *errs.XError) {
	var persons []entities.Person
	res := pr.WithDB(ctx).
		Where("customer_id = ?", customerId).
		Scopes(scopes.IsActive()).
		Preload("Measurements").
		Find(&persons)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find persons by customer id", res.Error)
	}
	return persons, nil
}
