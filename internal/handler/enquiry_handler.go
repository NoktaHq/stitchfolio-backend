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

type EnquiryHandler struct {
	enquirySvc service.EnquiryService
	resp       response.Response
	dataResp   response.DataResponse
}

func ProvideEnquiryHandler(svc service.EnquiryService) *EnquiryHandler {
	return &EnquiryHandler{enquirySvc: svc}
}

// Save Enquiry
//
//	@Summary		Save Enquiry
//	@Description	Saves an instance of Enquiry
//	@Tags			Enquiry
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			enquiry	body		requestModel.Enquiry	true	"enquiry"
//	@Router			/enquiry [post]
func (h EnquiryHandler) SaveEnquiry(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var enquiry requesModel.Enquiry
	err := ctx.Bind(&enquiry)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.enquirySvc.SaveEnquiry(&context, enquiry)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update Enquiry
//
//	@Summary		Update Enquiry
//	@Description	Updates an instance of Enquiry
//	@Tags			Enquiry
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			enquiry	body		requestModel.Enquiry	true	"enquiry"
//	@Param			id		path		int						true	"Enquiry id"
//	@Router			/enquiry/{id} [put]
func (h EnquiryHandler) UpdateEnquiry(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var enquiry requesModel.Enquiry
	err := ctx.Bind(&enquiry)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.enquirySvc.UpdateEnquiry(&context, enquiry, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get Enquiry
//
//	@Summary		Get a specific Enquiry
//	@Description	Get an instance of Enquiry
//	@Tags			Enquiry
//	@Accept			json
//	@Success		200	{object}	responseModel.Enquiry
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Enquiry id"
//	@Router			/enquiry/{id} [get]
func (h EnquiryHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	enquiry, errr := h.enquirySvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(enquiry).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active enquiries
//
//	@Summary		Get all active enquiries
//	@Description	Get all active enquiries
//	@Tags			Enquiry
//	@Accept			json
//	@Success		200		{object}	responseModel.Enquiry
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/enquiry [get]
func (h EnquiryHandler) GetAllEnquiries(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	enquiries, errr := h.enquirySvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(enquiries).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete an Enquiry
//
//	@Summary		Delete Enquiry
//	@Description	Deletes an instance of Enquiry
//	@Tags			Enquiry
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"enquiry id"
//
//	@Router			/enquiry/{id} [delete]
func (h EnquiryHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.enquirySvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}

// Update Enquiry and Customer
//
//	@Summary		Update Enquiry and Customer
//	@Description	Updates both enquiry and customer data together
//	@Tags			Enquiry
//	@Accept			json
//	@Success		200		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			enquiry	body		requestModel.UpdateEnquiryAndCustomer	true	"enquiry and customer"
//	@Param			id		path		int										true	"Enquiry id"
//	@Router			/enquiry/{id}/customer [put]
func (h EnquiryHandler) UpdateEnquiryAndCustomer(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var request requesModel.UpdateEnquiryAndCustomer
	err := ctx.Bind(&request)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.enquirySvc.UpdateEnquiryAndCustomer(&context, request, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusOK)
}
