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

type OrderHistoryHandler struct {
	orderHistorySvc service.OrderHistoryService
	resp            response.Response
	dataResp        response.DataResponse
}

func ProvideOrderHistoryHandler(svc service.OrderHistoryService) *OrderHistoryHandler {
	return &OrderHistoryHandler{orderHistorySvc: svc}
}

// Save OrderHistory
//
//	@Summary		Save OrderHistory
//	@Description	Saves an instance of OrderHistory
//	@Tags			OrderHistory
//	@Accept			json
//	@Success		201				{object}	responseModel.Response
//	@Failure		400				{object}	responseModel.Response
//	@Failure		501				{object}	responseModel.Response
//	@Param			orderHistory	body		requestModel.OrderHistory	true	"orderHistory"
//	@Router			/order-history [post]
func (h OrderHistoryHandler) SaveOrderHistory(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var orderHistory requesModel.OrderHistory
	err := ctx.Bind(&orderHistory)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.orderHistorySvc.SaveOrderHistory(&context, orderHistory)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Get OrderHistory
//
//	@Summary		Get a specific OrderHistory
//	@Description	Get an instance of OrderHistory
//	@Tags			OrderHistory
//	@Accept			json
//	@Success		200	{object}	responseModel.OrderHistory
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"OrderHistory id"
//	@Router			/order-history/{id} [get]
func (h OrderHistoryHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	orderHistory, errr := h.orderHistorySvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(orderHistory).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active order histories
//
//	@Summary		Get all active order histories
//	@Description	Get all active order histories
//	@Tags			OrderHistory
//	@Accept			json
//	@Success		200		{object}	responseModel.OrderHistory
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/order-history [get]
func (h OrderHistoryHandler) GetAllOrderHistories(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	orderHistories, errr := h.orderHistorySvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(orderHistories).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get order histories by order id
//
//	@Summary		Get order histories by order id
//	@Description	Get order histories by order id
//	@Tags			OrderHistory
//	@Accept			json
//	@Success		200		{object}	responseModel.OrderHistory
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			orderId	path		int	true	"order id"
//	@Router			/order-history/order/{orderId} [get]
func (h OrderHistoryHandler) GetByOrderId(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	orderId, _ := strconv.Atoi(ctx.Param("orderId"))

	orderHistories, errr := h.orderHistorySvc.GetByOrderId(&context, uint(orderId))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(orderHistories).FormatAndSend(&context, ctx, http.StatusOK)
}
