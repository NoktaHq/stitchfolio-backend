package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type EnquiryRepository interface {
	Create(*context.Context, *entities.Enquiry) *errs.XError
	Update(*context.Context, *entities.Enquiry) *errs.XError
	UpdateEnquiryAndCustomer(*context.Context, *entities.Enquiry, *entities.Customer) *errs.XError
	Get(*context.Context, uint) (*entities.Enquiry, *errs.XError)
	GetAll(*context.Context, string) ([]entities.Enquiry, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type enquiryRepository struct {
	GormDAL
}

func ProvideEnquiryRepository(customDB GormDAL) EnquiryRepository {
	return &enquiryRepository{GormDAL: customDB}
}

func (er *enquiryRepository) Create(ctx *context.Context, enquiry *entities.Enquiry) *errs.XError {
	res := er.WithDB(ctx).Create(&enquiry)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save enquiry", res.Error)
	}
	return nil
}

func (er *enquiryRepository) Update(ctx *context.Context, enquiry *entities.Enquiry) *errs.XError {
	return er.GormDAL.Update(ctx, *enquiry)
}

func (er *enquiryRepository) UpdateEnquiryAndCustomer(ctx *context.Context, enquiry *entities.Enquiry, customer *entities.Customer) *errs.XError {
	// Update customer first
	if customer != nil && customer.ID != 0 {
		customerErr := er.GormDAL.Update(ctx, *customer)
		if customerErr != nil {
			return customerErr
		}
	}

	// Then update enquiry
	return er.GormDAL.Update(ctx, *enquiry)
}

func (er *enquiryRepository) Get(ctx *context.Context, id uint) (*entities.Enquiry, *errs.XError) {
	enquiry := entities.Enquiry{}
	res := er.WithDB(ctx).Preload("Customer").Find(&enquiry, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find enquiry", res.Error)
	}
	return &enquiry, nil
}

func (er *enquiryRepository) GetAll(ctx *context.Context, search string) ([]entities.Enquiry, *errs.XError) {
	var enquiries []entities.Enquiry
	res := er.WithDB(ctx).Model(entities.Enquiry{}).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "subject", "notes", "status")).
		Scopes(scopes.WithAuditInfo()).
		Scopes(db.Paginate(ctx)).
		Preload("Customer", scopes.SelectFields("first_name", "last_name", "phone_number")).
		Find(&enquiries)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find enquiries", res.Error)
	}
	return enquiries, nil
}

func (er *enquiryRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	enquiry := &entities.Enquiry{Model: &entities.Model{ID: id, IsActive: false}}
	err := er.GormDAL.Delete(ctx, enquiry)
	if err != nil {
		return err
	}
	return nil
}
