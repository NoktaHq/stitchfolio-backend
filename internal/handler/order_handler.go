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

type OrderHandler struct {
	orderSvc service.OrderService
	resp     response.Response
	dataResp response.DataResponse
}

func ProvideOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: svc}
}

// Save Order
//
//	@Summary		Save Order
//	@Description	Saves an instance of Order
//	@Tags			Order
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			order	body		requestModel.Order	true	"order"
//	@Router			/order [post]
func (h OrderHandler) SaveOrder(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var order requesModel.Order
	err := ctx.Bind(&order)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.orderSvc.SaveOrder(&context, order)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update Order
//
//	@Summary		Update Order
//	@Description	Updates an instance of Order
//	@Tags			Order
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			order	body		requestModel.Order	true	"order"
//	@Param			id		path		int					true	"Order id"
//	@Router			/order/{id} [put]
func (h OrderHandler) UpdateOrder(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var order requesModel.Order
	err := ctx.Bind(&order)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.orderSvc.UpdateOrder(&context, order, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get Order
//
//	@Summary		Get a specific Order
//	@Description	Get an instance of Order
//	@Tags			Order
//	@Accept			json
//	@Success		200	{object}	responseModel.Order
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Order id"
//	@Router			/order/{id} [get]
func (h OrderHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	order, errr := h.orderSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(order).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active orders
//
//	@Summary		Get all active orders
//	@Description	Get all active orders
//	@Tags			Order
//	@Accept			json
//	@Success		200		{object}	responseModel.Order
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Param			filters	query		string	false	"filters (e.g., name eq 'Shirt', status eq 'Active')"
//	@Router			/order [get]
func (h OrderHandler) GetAllOrders(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	orders, errr := h.orderSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(orders).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete an Order
//
//	@Summary		Delete Order
//	@Description	Deletes an instance of Order
//	@Tags			Order
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"order id"
//
//	@Router			/order/{id} [delete]
func (h OrderHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.orderSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
