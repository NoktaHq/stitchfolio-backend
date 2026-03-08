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

type PurchaseHandler struct {
	purchaseSvc service.PurchaseService
	resp        response.Response
	dataResp    response.DataResponse
}

func ProvidePurchaseHandler(svc service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchaseSvc: svc}
}

// SavePurchase creates a new purchase.
func (h PurchaseHandler) SavePurchase(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var purchase requestModel.Purchase
	if err := ctx.Bind(&purchase); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	if errr := h.purchaseSvc.SavePurchase(&context, purchase); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// UpdatePurchase updates an existing purchase.
func (h PurchaseHandler) UpdatePurchase(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var purchase requestModel.Purchase
	if err := ctx.Bind(&purchase); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(ctx.Param("id"))
	if errr := h.purchaseSvc.UpdatePurchase(&context, purchase, uint(id)); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get returns a purchase by id with vendor and items.
func (h PurchaseHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	purchase, errr := h.purchaseSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(purchase).FormatAndSend(&context, ctx, http.StatusOK)
}

// GetAllPurchases returns all purchases with optional search.
func (h PurchaseHandler) GetAllPurchases(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)
	purchases, errr := h.purchaseSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(purchases).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete soft-deletes a purchase.
func (h PurchaseHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	if errr := h.purchaseSvc.Delete(&context, uint(id)); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.resp.SuccessResponse("Delete success").FormatAndSend(&context, ctx, http.StatusOK)
}

// ReceivePurchase receives all pending quantities for the purchase and updates inventory.
func (h PurchaseHandler) ReceivePurchase(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	result, errr := h.purchaseSvc.ReceivePurchase(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(result).FormatAndSend(&context, ctx, http.StatusOK)
}
