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

type PurchaseItemHandler struct {
	purchaseItemSvc service.PurchaseItemService
	resp            response.Response
	dataResp        response.DataResponse
}

func ProvidePurchaseItemHandler(svc service.PurchaseItemService) *PurchaseItemHandler {
	return &PurchaseItemHandler{purchaseItemSvc: svc}
}

// Save creates a purchase item for a purchase.
func (h *PurchaseItemHandler) Save(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	purchaseId, _ := strconv.Atoi(ctx.Param("id"))
	var req requestModel.PurchaseItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	if errr := h.purchaseItemSvc.Save(&c, req, uint(purchaseId)); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&c, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Save success").FormatAndSend(&c, ctx, http.StatusCreated)
}

// Update updates a purchase item.
func (h *PurchaseItemHandler) Update(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	detailIdStr := ctx.Param("detailId")
	if detailIdStr == "" {
		detailIdStr = ctx.Param("id")
	}
	id, _ := strconv.Atoi(detailIdStr)
	var req requestModel.PurchaseItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	if errr := h.purchaseItemSvc.Update(&c, req, uint(id)); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&c, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Update success").FormatAndSend(&c, ctx, http.StatusAccepted)
}

// Get returns a purchase item by id.
func (h *PurchaseItemHandler) Get(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	item, errr := h.purchaseItemSvc.Get(&c, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(item).FormatAndSend(&c, ctx, http.StatusOK)
}

// GetByPurchaseId returns all purchase items for a purchase.
func (h *PurchaseItemHandler) GetByPurchaseId(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	purchaseId, _ := strconv.Atoi(ctx.Param("id"))
	items, errr := h.purchaseItemSvc.GetByPurchaseId(&c, uint(purchaseId))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(items).FormatAndSend(&c, ctx, http.StatusOK)
}

// Delete soft-deletes a purchase item.
func (h *PurchaseItemHandler) Delete(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	detailIdStr := ctx.Param("detailId")
	if detailIdStr == "" {
		detailIdStr = ctx.Param("id")
	}
	id, _ := strconv.Atoi(detailIdStr)
	if errr := h.purchaseItemSvc.Delete(&c, uint(id)); errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.resp.SuccessResponse("Delete success").FormatAndSend(&c, ctx, http.StatusOK)
}
