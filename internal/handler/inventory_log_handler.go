package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type InventoryLogHandler struct {
	inventoryLogSvc service.InventoryLogService
	resp            response.Response
	dataResp        response.DataResponse
}

func ProvideInventoryLogHandler(svc service.InventoryLogService) *InventoryLogHandler {
	return &InventoryLogHandler{inventoryLogSvc: svc}
}

//	@Summary		Get a specific Inventory Log
//	@Description	Get an instance of Inventory Log
//	@Tags			InventoryLog
//	@Accept			json
//	@Success		200	{object}	responseModel.InventoryLog
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Inventory Log id"
//	@Router			/inventory-log/{id} [get]
func (h InventoryLogHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	log, errr := h.inventoryLogSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(log).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get all inventory logs
//	@Description	Get all inventory log records
//	@Tags			InventoryLog
//	@Accept			json
//	@Success		200		{object}	responseModel.InventoryLog
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/inventory-log [get]
func (h InventoryLogHandler) GetAllInventoryLogs(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	logs, errr := h.inventoryLogSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(logs).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get inventory logs by product ID
//	@Description	Get all inventory logs for a specific product
//	@Tags			InventoryLog
//	@Accept			json
//	@Success		200			{object}	responseModel.InventoryLog
//	@Failure		400			{object}	responseModel.DataResponse
//	@Param			productId	path		int	true	"Product ID"
//	@Router			/inventory-log/product/{productId} [get]
func (h InventoryLogHandler) GetByProductId(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	productId, _ := strconv.Atoi(ctx.Param("productId"))

	logs, errr := h.inventoryLogSvc.GetByProductId(&context, uint(productId))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(logs).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get inventory logs by change type
//	@Description	Get all inventory logs filtered by change type (IN, OUT, ADJUST)
//	@Tags			InventoryLog
//	@Accept			json
//	@Success		200			{object}	responseModel.InventoryLog
//	@Failure		400			{object}	responseModel.DataResponse
//	@Param			changeType	query		string	true	"Change Type"	Enums(IN, OUT, ADJUST)
//	@Router			/inventory-log/change-type [get]
func (h InventoryLogHandler) GetByChangeType(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	changeType := ctx.Query("changeType")

	logs, errr := h.inventoryLogSvc.GetByChangeType(&context, changeType)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(logs).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get inventory logs by date range
//	@Description	Get all inventory logs within a date range
//	@Tags			InventoryLog
//	@Accept			json
//	@Success		200			{object}	responseModel.InventoryLog
//	@Failure		400			{object}	responseModel.DataResponse
//	@Param			startDate	query		string	false	"Start Date (YYYY-MM-DD)"
//	@Param			endDate		query		string	false	"End Date (YYYY-MM-DD)"
//	@Router			/inventory-log/date-range [get]
func (h InventoryLogHandler) GetByDateRange(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	logs, errr := h.inventoryLogSvc.GetByDateRange(&context, startDate, endDate)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(logs).FormatAndSend(&context, ctx, http.StatusOK)
}
