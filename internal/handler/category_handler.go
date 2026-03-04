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

type CategoryHandler struct {
	categorySvc service.CategoryService
	resp        response.Response
	dataResp    response.DataResponse
}

func ProvideCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc: svc}
}

//	@Summary		Save Category
//	@Description	Saves an instance of Category
//	@Tags			Category
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		500			{object}	responseModel.Response
//	@Param			category	body		requestModel.Category	true	"category"
//	@Router			/category [post]
func (h CategoryHandler) SaveCategory(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var category requestModel.Category
	err := ctx.Bind(&category)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.categorySvc.SaveCategory(&context, category)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

//	@Summary		Update Category
//	@Description	Updates an instance of Category
//	@Tags			Category
//	@Accept			json
//	@Success		202			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		500			{object}	responseModel.Response
//	@Param			category	body		requestModel.Category	true	"category"
//	@Param			id			path		int						true	"Category id"
//	@Router			/category/{id} [put]
func (h CategoryHandler) UpdateCategory(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var category requestModel.Category
	err := ctx.Bind(&category)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.categorySvc.UpdateCategory(&context, category, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

//	@Summary		Get a specific Category
//	@Description	Get an instance of Category
//	@Tags			Category
//	@Accept			json
//	@Success		200	{object}	responseModel.Category
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Category id"
//	@Router			/category/{id} [get]
func (h CategoryHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	category, errr := h.categorySvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(category).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get all active categories
//	@Description	Get all active categories
//	@Tags			Category
//	@Accept			json
//	@Success		200		{object}	responseModel.Category
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/category [get]
func (h CategoryHandler) GetAllCategories(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	categories, errr := h.categorySvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(categories).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Delete Category
//	@Description	Deletes an instance of Category
//	@Tags			Category
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"category id"
//	@Router			/category/{id} [delete]
func (h CategoryHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.categorySvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Autocomplete for categories
//	@Description	Autocomplete for categories
//	@Tags			Category
//	@Accept			json
//	@Success		200		{object}	responseModel.CategoryAutoComplete
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/category/autocomplete [get]
func (h CategoryHandler) AutocompleteCategory(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	categories, errr := h.categorySvc.AutocompleteCategory(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(categories).FormatAndSend(&context, ctx, http.StatusOK)
}
