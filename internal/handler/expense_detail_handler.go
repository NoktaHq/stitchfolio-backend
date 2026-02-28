package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type ExpenseDetailHandler struct {
	expenseDetailSvc service.ExpenseDetailService
	resp             response.Response
	dataResp         response.DataResponse
}

func ProvideExpenseDetailHandler(svc service.ExpenseDetailService) *ExpenseDetailHandler {
	return &ExpenseDetailHandler{expenseDetailSvc: svc}
}

var _ = (*responseModel.Response)(nil) // used by swagger comments

// Save ExpenseDetail
//
//	@Summary		Save ExpenseDetail
//	@Description	Saves an expense detail for an expense
//	@Tags			ExpenseDetail
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Param			expenseId	path		int							true	"Expense id"
//	@Param			body		body		requestModel.ExpenseDetail	true	"expense detail"
//	@Router			/expense-tracker/{expenseId}/expense-detail [post]
func (h *ExpenseDetailHandler) Save(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	expenseId, _ := strconv.Atoi(ctx.Param("id"))
	var req requestModel.ExpenseDetail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	if err := h.expenseDetailSvc.Save(&c, req, uint(expenseId)); err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Save success").FormatAndSend(&c, ctx, http.StatusCreated)
}

// Update ExpenseDetail
//
//	@Summary		Update ExpenseDetail
//	@Description	Updates an expense detail
//	@Tags			ExpenseDetail
//	@Accept			json
//	@Success		202		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Param			id		path		int							true	"ExpenseDetail id"
//	@Param			body	body		requestModel.ExpenseDetail	true	"expense detail"
//	@Router			/expense-detail/{id} [put]
func (h *ExpenseDetailHandler) Update(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	detailIdStr := ctx.Param("detailId")
	if detailIdStr == "" {
		detailIdStr = ctx.Param("id")
	}
	id, _ := strconv.Atoi(detailIdStr)
	var req requestModel.ExpenseDetail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.resp.DefaultFailureResponse(errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	if err := h.expenseDetailSvc.Update(&c, req, uint(id)); err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusInternalServerError)
		return
	}
	h.resp.SuccessResponse("Update success").FormatAndSend(&c, ctx, http.StatusAccepted)
}

// Get ExpenseDetail
//
//	@Summary		Get ExpenseDetail
//	@Description	Get an expense detail by id
//	@Tags			ExpenseDetail
//	@Accept			json
//	@Success		200	{object}	responseModel.ExpenseDetail
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"ExpenseDetail id"
//	@Router			/expense-detail/{id} [get]
func (h *ExpenseDetailHandler) Get(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	detail, err := h.expenseDetailSvc.Get(&c, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(detail).FormatAndSend(&c, ctx, http.StatusOK)
}

// GetByExpenseId returns all expense details for an expense
//
//	@Summary		Get expense details by expense id
//	@Description	Get all expense details for an expense
//	@Tags			ExpenseDetail
//	@Accept			json
//	@Success		200			{object}	[]responseModel.ExpenseDetail
//	@Failure		400			{object}	responseModel.DataResponse
//	@Param			expenseId	path		int	true	"Expense id"
//	@Router			/expense-tracker/{expenseId}/expense-detail [get]
func (h *ExpenseDetailHandler) GetByExpenseId(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	expenseId, _ := strconv.Atoi(ctx.Param("id"))
	details, err := h.expenseDetailSvc.GetByExpenseId(&c, uint(expenseId))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(details).FormatAndSend(&c, ctx, http.StatusOK)
}

// Delete ExpenseDetail
//
//	@Summary		Delete ExpenseDetail
//	@Description	Deletes an expense detail
//	@Tags			ExpenseDetail
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"ExpenseDetail id"
//	@Router			/expense-detail/{id} [delete]
func (h *ExpenseDetailHandler) Delete(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	detailIdStr := ctx.Param("detailId")
	if detailIdStr == "" {
		detailIdStr = ctx.Param("id")
	}
	id, _ := strconv.Atoi(detailIdStr)
	if err := h.expenseDetailSvc.Delete(&c, uint(id)); err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.resp.SuccessResponse("Delete success").FormatAndSend(&c, ctx, http.StatusOK)
}
