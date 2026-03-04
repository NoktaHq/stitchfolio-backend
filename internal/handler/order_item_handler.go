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

type OrderItemHandler struct {
	orderItemSvc service.OrderItemService
	resp         response.Response
	dataResp     response.DataResponse
}

func ProvideOrderItemHandler(svc service.OrderItemService) *OrderItemHandler {
	return &OrderItemHandler{orderItemSvc: svc}
}

// Save OrderItem
//
//	@Summary		Save OrderItem
//	@Description	Saves an instance of OrderItem
//	@Tags			OrderItem
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			orderItem	body		requestModel.OrderItem	true	"orderItem"
//	@Router			/order-item [post]
func (h OrderItemHandler) SaveOrderItem(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var orderItem requesModel.OrderItem
	err := ctx.Bind(&orderItem)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.orderItemSvc.SaveOrderItem(&context, orderItem)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update OrderItem
//
//	@Summary		Update OrderItem
//	@Description	Updates an instance of OrderItem
//	@Tags			OrderItem
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			orderItem	body		requestModel.OrderItem	true	"orderItem"
//	@Param			id			path		int						true	"OrderItem id"
//	@Router			/order-item/{id} [put]
func (h OrderItemHandler) UpdateOrderItem(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var orderItem requesModel.OrderItem
	err := ctx.Bind(&orderItem)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.orderItemSvc.UpdateOrderItem(&context, orderItem, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get OrderItem
//
//	@Summary		Get a specific OrderItem
//	@Description	Get an instance of OrderItem
//	@Tags			OrderItem
//	@Accept			json
//	@Success		200	{object}	responseModel.OrderItem
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"OrderItem id"
//	@Router			/order-item/{id} [get]
func (h OrderItemHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	orderItem, errr := h.orderItemSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(orderItem).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active order items
//
//	@Summary		Get all active order items
//	@Description	Get all active order items
//	@Tags			OrderItem
//	@Accept			json
//	@Success		200		{object}	responseModel.OrderItem
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/order-item [get]
func (h OrderItemHandler) GetAllOrderItems(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	orderItems, errr := h.orderItemSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(orderItems).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete an OrderItem
//
//	@Summary		Delete OrderItem
//	@Description	Deletes an instance of OrderItem
//	@Tags			OrderItem
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"orderItem id"
//
//	@Router			/order-item/{id} [delete]
func (h OrderItemHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.orderItemSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
