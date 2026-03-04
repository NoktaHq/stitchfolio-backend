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

type ExpenseTrackerHandler struct {
	expenseTrackerSvc service.ExpenseTrackerService
	resp              response.Response
	dataResp          response.DataResponse
}

func ProvideExpenseTrackerHandler(svc service.ExpenseTrackerService) *ExpenseTrackerHandler {
	return &ExpenseTrackerHandler{expenseTrackerSvc: svc}
}

// Save ExpenseTracker
//
//	@Summary		Save ExpenseTracker
//	@Description	Saves an instance of ExpenseTracker
//	@Tags			ExpenseTracker
//	@Accept			json
//	@Success		201				{object}	responseModel.Response
//	@Failure		400				{object}	responseModel.Response
//	@Failure		501				{object}	responseModel.Response
//	@Param			expenseTracker	body		requestModel.ExpenseTracker	true	"expenseTracker"
//	@Router			/expense-tracker [post]
func (h ExpenseTrackerHandler) SaveExpenseTracker(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var expenseTracker requesModel.ExpenseTracker
	err := ctx.Bind(&expenseTracker)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.expenseTrackerSvc.SaveExpenseTracker(&context, expenseTracker)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update ExpenseTracker
//
//	@Summary		Update ExpenseTracker
//	@Description	Updates an instance of ExpenseTracker
//	@Tags			ExpenseTracker
//	@Accept			json
//	@Success		201				{object}	responseModel.Response
//	@Failure		400				{object}	responseModel.Response
//	@Failure		501				{object}	responseModel.Response
//	@Param			expenseTracker	body		requestModel.ExpenseTracker	true	"expenseTracker"
//	@Param			id				path		int							true	"ExpenseTracker id"
//	@Router			/expense-tracker/{id} [put]
func (h ExpenseTrackerHandler) UpdateExpenseTracker(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var expenseTracker requesModel.ExpenseTracker
	err := ctx.Bind(&expenseTracker)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.expenseTrackerSvc.UpdateExpenseTracker(&context, expenseTracker, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get ExpenseTracker
//
//	@Summary		Get a specific ExpenseTracker
//	@Description	Get an instance of ExpenseTracker
//	@Tags			ExpenseTracker
//	@Accept			json
//	@Success		200	{object}	responseModel.ExpenseTracker
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"ExpenseTracker id"
//	@Router			/expense-tracker/{id} [get]
func (h ExpenseTrackerHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	expenseTracker, errr := h.expenseTrackerSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(expenseTracker).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active expense trackers
//
//	@Summary		Get all active expense trackers
//	@Description	Get all active expense trackers
//	@Tags			ExpenseTracker
//	@Accept			json
//	@Success		200		{object}	responseModel.ExpenseTracker
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/expense-tracker [get]
func (h ExpenseTrackerHandler) GetAllExpenseTrackers(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	expenseTrackers, errr := h.expenseTrackerSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(expenseTrackers).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete an ExpenseTracker
//
//	@Summary		Delete ExpenseTracker
//	@Description	Deletes an instance of ExpenseTracker
//	@Tags			ExpenseTracker
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"expenseTracker id"
//
//	@Router			/expense-tracker/{id} [delete]
func (h ExpenseTrackerHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.expenseTrackerSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
