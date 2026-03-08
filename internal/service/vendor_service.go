package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type VendorService interface {
	SaveVendor(*context.Context, requestModel.Vendor) *errs.XError
	UpdateVendor(*context.Context, requestModel.Vendor, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Vendor, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Vendor, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	AutocompleteVendor(*context.Context, string) ([]responseModel.VendorAutoComplete, *errs.XError)
}

type vendorService struct {
	vendorRepo repository.VendorRepository
	mapper     mapper.Mapper
	respMapper mapper.ResponseMapper
}

func ProvideVendorService(
	repo repository.VendorRepository,
	m mapper.Mapper,
	respMapper mapper.ResponseMapper,
) VendorService {
	return vendorService{
		vendorRepo: repo,
		mapper:     m,
		respMapper: respMapper,
	}
}

func (svc vendorService) SaveVendor(ctx *context.Context, vendor requestModel.Vendor) *errs.XError {
	ent, err := svc.mapper.Vendor(vendor)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save vendor", err)
	}
	return svc.vendorRepo.Create(ctx, ent)
}

func (svc vendorService) UpdateVendor(ctx *context.Context, vendor requestModel.Vendor, id uint) *errs.XError {
	ent, err := svc.mapper.Vendor(vendor)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update vendor", err)
	}
	ent.ID = id
	return svc.vendorRepo.Update(ctx, ent)
}

func (svc vendorService) Get(ctx *context.Context, id uint) (*responseModel.Vendor, *errs.XError) {
	vendor, err := svc.vendorRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.Vendor(vendor)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Vendor data", mapErr)
	}
	return mapped, nil
}

func (svc vendorService) GetAll(ctx *context.Context, search string) ([]responseModel.Vendor, *errs.XError) {
	vendors, err := svc.vendorRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}
	mapped, mapErr := svc.respMapper.Vendors(vendors)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Vendor data", mapErr)
	}
	return mapped, nil
}

func (svc vendorService) Delete(ctx *context.Context, id uint) *errs.XError {
	return svc.vendorRepo.Delete(ctx, id)
}

func (svc vendorService) AutocompleteVendor(ctx *context.Context, search string) ([]responseModel.VendorAutoComplete, *errs.XError) {
	vendors, err := svc.vendorRepo.AutocompleteVendor(ctx, search)
	if err != nil {
		return nil, err
	}
	res := make([]responseModel.VendorAutoComplete, 0, len(vendors))
	for _, v := range vendors {
		res = append(res, responseModel.VendorAutoComplete{ID: v.ID, Name: v.Name})
	}
	return res, nil
}
