package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	requesModel "github.com/imkarthi24/sf-backend/internal/model/request"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type CustomerHandler struct {
	customerSvc service.CustomerService
	resp        response.Response
	dataResp    response.DataResponse
}

func ProvideCustomerHandler(svc service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerSvc: svc}
}

// Save Customer
//
//	@Summary		Save Customer
//	@Description	Saves an instance of Customer
//	@Tags			Customer
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			customer	body		requestModel.Customer	true	"customer"
//	@Router			/customer [post]
func (h CustomerHandler) SaveCustomer(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var customer requesModel.Customer
	err := ctx.Bind(&customer)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.customerSvc.SaveCustomer(&context, customer)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update Customer
//
//	@Summary		Update Customer
//	@Description	Updates an instance of Customer
//	@Tags			Customer
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			customer	body		requestModel.Customer	true	"customer"
//	@Param			id			path		int						true	"Customer id"
//	@Router			/customer/{id} [put]
func (h CustomerHandler) UpdateCustomer(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var customer requesModel.Customer
	err := ctx.Bind(&customer)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.customerSvc.UpdateCustomer(&context, customer, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get Customer
//
//	@Summary		Get a specific Customer
//	@Description	Get an instance of Customer
//	@Tags			Customer
//	@Accept			json
//	@Success		200	{object}	responseModel.Customer
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Customer id"
//	@Router			/customer/{id} [get]
func (h CustomerHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	customer, errr := h.customerSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(customer).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active customers
//
//	@Summary		Get all active customers
//	@Description	Get all active customers
//	@Tags			Customer
//	@Accept			json
//	@Success		200		{object}	responseModel.Customer
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/customer [get]
func (h CustomerHandler) GetAllCustomers(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	customers, errr := h.customerSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(customers).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete a Customer
//
//	@Summary		Delete Customer
//	@Description	Deletes an instance of Customer
//	@Tags			Customer
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"customer id"
//
//	@Router			/customer/{id} [delete]
func (h CustomerHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.customerSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}

// Autocomplete for customers
//
//	@Summary		Autocomplete for customers
//	@Description	Autocomplete for customers
//	@Tags			Customer
//	@Accept			json
//	@Success		200		{object}	responseModel.CustomerAutoComplete
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/customer/autocomplete [get]
func (h CustomerHandler) AutocompleteCustomer(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	customers, errr := h.customerSvc.AutocompleteCustomer(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(customers).FormatAndSend(&context, ctx, http.StatusOK)
}
