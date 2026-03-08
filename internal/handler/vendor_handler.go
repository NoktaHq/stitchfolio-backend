package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type VendorHandler struct {
	vendorSvc service.VendorService
	resp      response.Response
	dataResp  response.DataResponse
}

func ProvideVendorHandler(svc service.VendorService) *VendorHandler {
	return &VendorHandler{vendorSvc: svc}
}

// SaveVendor saves a new vendor.
func (h VendorHandler) SaveVendor(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var vendor requestModel.Vendor
	if err := ctx.Bind(&vendor); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	if errr := h.vendorSvc.SaveVendor(&context, vendor); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// UpdateVendor updates an existing vendor.
func (h VendorHandler) UpdateVendor(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var vendor requestModel.Vendor
	if err := ctx.Bind(&vendor); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(ctx.Param("id"))
	if errr := h.vendorSvc.UpdateVendor(&context, vendor, uint(id)); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get returns a vendor by id.
func (h VendorHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	vendor, errr := h.vendorSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(vendor).FormatAndSend(&context, ctx, http.StatusOK)
}

// GetAllVendors returns all vendors with optional search.
func (h VendorHandler) GetAllVendors(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)
	vendors, errr := h.vendorSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(vendors).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete soft-deletes a vendor.
func (h VendorHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	if errr := h.vendorSvc.Delete(&context, uint(id)); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.resp.SuccessResponse("Delete success").FormatAndSend(&context, ctx, http.StatusOK)
}

// AutocompleteVendor returns vendors for autocomplete.
func (h VendorHandler) AutocompleteVendor(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)
	vendors, errr := h.vendorSvc.AutocompleteVendor(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(vendors).FormatAndSend(&context, ctx, http.StatusOK)
}
