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

type EnquiryHistoryHandler struct {
	enquiryHistorySvc service.EnquiryHistoryService
	resp              response.Response
	dataResp          response.DataResponse
}

func ProvideEnquiryHistoryHandler(svc service.EnquiryHistoryService) *EnquiryHistoryHandler {
	return &EnquiryHistoryHandler{enquiryHistorySvc: svc}
}

// Save EnquiryHistory
//
//	@Summary		Save EnquiryHistory
//	@Description	Saves an instance of EnquiryHistory
//	@Tags			EnquiryHistory
//	@Accept			json
//	@Success		201				{object}	responseModel.Response
//	@Failure		400				{object}	responseModel.Response
//	@Failure		501				{object}	responseModel.Response
//	@Param			enquiryHistory	body		requestModel.EnquiryHistory	true	"enquiryHistory"
//	@Router			/enquiry-history [post]
func (h EnquiryHistoryHandler) SaveEnquiryHistory(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var enquiryHistory requesModel.EnquiryHistory
	err := ctx.Bind(&enquiryHistory)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.enquiryHistorySvc.SaveEnquiryHistory(&context, enquiryHistory)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Get EnquiryHistory
//
//	@Summary		Get a specific EnquiryHistory
//	@Description	Get an instance of EnquiryHistory
//	@Tags			EnquiryHistory
//	@Accept			json
//	@Success		200	{object}	responseModel.EnquiryHistory
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"EnquiryHistory id"
//	@Router			/enquiry-history/{id} [get]
func (h EnquiryHistoryHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	enquiryHistory, errr := h.enquiryHistorySvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(enquiryHistory).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active enquiry histories
//
//	@Summary		Get all active enquiry histories
//	@Description	Get all active enquiry histories
//	@Tags			EnquiryHistory
//	@Accept			json
//	@Success		200		{object}	responseModel.EnquiryHistory
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/enquiry-history [get]
func (h EnquiryHistoryHandler) GetAllEnquiryHistories(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	enquiryHistories, errr := h.enquiryHistorySvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(enquiryHistories).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get enquiry histories by enquiry id
//
//	@Summary		Get enquiry histories by enquiry id
//	@Description	Get enquiry histories by enquiry id
//	@Tags			EnquiryHistory
//	@Accept			json
//	@Success		200			{object}	responseModel.EnquiryHistory
//	@Failure		400			{object}	responseModel.DataResponse
//	@Param			enquiryId	path		int	true	"enquiry id"
//	@Router			/enquiry-history/enquiry/{enquiryId} [get]
func (h EnquiryHistoryHandler) GetByEnquiryId(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	enquiryId, _ := strconv.Atoi(ctx.Param("enquiryId"))

	enquiryHistories, errr := h.enquiryHistorySvc.GetByEnquiryId(&context, uint(enquiryId))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(enquiryHistories).FormatAndSend(&context, ctx, http.StatusOK)
}
