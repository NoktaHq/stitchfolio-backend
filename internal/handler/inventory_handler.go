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

type InventoryHandler struct {
	inventorySvc service.InventoryService
	resp         response.Response
	dataResp     response.DataResponse
}

func ProvideInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventorySvc: svc}
}

//	@Summary		Get a specific Inventory
//	@Description	Get an instance of Inventory
//	@Tags			Inventory
//	@Accept			json
//	@Success		200	{object}	responseModel.Inventory
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Inventory id"
//	@Router			/inventory/{id} [get]
func (h InventoryHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	inventory, errr := h.inventorySvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(inventory).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get all inventories
//	@Description	Get all inventory records
//	@Tags			Inventory
//	@Accept			json
//	@Success		200		{object}	responseModel.Inventory
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/inventory [get]
func (h InventoryHandler) GetAllInventories(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	inventories, errr := h.inventorySvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(inventories).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get inventory by product ID
//	@Description	Get inventory for a specific product
//	@Tags			Inventory
//	@Accept			json
//	@Success		200			{object}	responseModel.Inventory
//	@Failure		400			{object}	responseModel.DataResponse
//	@Param			productId	path		int	true	"Product ID"
//	@Router			/inventory/product/{productId} [get]
func (h InventoryHandler) GetByProductId(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	productId, _ := strconv.Atoi(ctx.Param("productId"))

	inventory, errr := h.inventorySvc.GetByProductId(&context, uint(productId))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(inventory).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Update low stock threshold
//	@Description	Update the low stock threshold for a product
//	@Tags			Inventory
//	@Accept			json
//	@Success		202			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		500			{object}	responseModel.Response
//	@Param			inventory	body		requestModel.Inventory	true	"inventory"
//	@Param			id			path		int						true	"Inventory id"
//	@Router			/inventory/{id}/threshold [put]
func (h InventoryHandler) UpdateThreshold(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var inventory requestModel.Inventory
	err := ctx.Bind(&inventory)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.inventorySvc.UpdateThreshold(&context, inventory, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Threshold updated successfully").FormatAndSend(&context, ctx, http.StatusAccepted)
}

//	@Summary		Get low stock items
//	@Description	Get all items with stock below threshold
//	@Tags			Inventory
//	@Accept			json
//	@Success		200	{object}	responseModel.LowStockItem
//	@Failure		400	{object}	responseModel.DataResponse
//	@Router			/inventory/low-stock [get]
func (h InventoryHandler) GetLowStockItems(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	items, errr := h.inventorySvc.GetLowStockItems(&context)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(items).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Record stock movement
//	@Description	Record a stock IN, OUT, or ADJUST movement
//	@Tags			Inventory
//	@Accept			json
//	@Success		201			{object}	responseModel.StockMovementResponse
//	@Failure		400			{object}	responseModel.DataResponse
//	@Failure		500			{object}	responseModel.DataResponse
//	@Param			movement	body		requestModel.StockMovementRequest	true	"stock movement"
//	@Router			/inventory/movement [post]
func (h InventoryHandler) RecordStockMovement(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var movement requestModel.StockMovementRequest
	err := ctx.Bind(&movement)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	response, errr := h.inventorySvc.RecordStockMovement(&context, movement)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.dataResp.DefaultSuccessResponse(response).FormatAndSend(&context, ctx, http.StatusCreated)
}
