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

type DressTypeHandler struct {
	dressTypeSvc service.DressTypeService
	resp         response.Response
	dataResp     response.DataResponse
}

func ProvideDressTypeHandler(svc service.DressTypeService) *DressTypeHandler {
	return &DressTypeHandler{dressTypeSvc: svc}
}

// Save DressType
//
//	@Summary		Save DressType
//	@Description	Saves an instance of DressType
//	@Tags			DressType
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			dressType	body		requestModel.DressType	true	"dressType"
//	@Router			/dress-type [post]
func (h DressTypeHandler) SaveDressType(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var dressType requesModel.DressType
	err := ctx.Bind(&dressType)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.dressTypeSvc.SaveDressType(&context, dressType)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update DressType
//
//	@Summary		Update DressType
//	@Description	Updates an instance of DressType
//	@Tags			DressType
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			dressType	body		requestModel.DressType	true	"dressType"
//	@Param			id			path		int						true	"DressType id"
//	@Router			/dress-type/{id} [put]
func (h DressTypeHandler) UpdateDressType(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var dressType requesModel.DressType
	err := ctx.Bind(&dressType)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.dressTypeSvc.UpdateDressType(&context, dressType, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get DressType
//
//	@Summary		Get a specific DressType
//	@Description	Get an instance of DressType
//	@Tags			DressType
//	@Accept			json
//	@Success		200	{object}	responseModel.DressType
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"DressType id"
//	@Router			/dress-type/{id} [get]
func (h DressTypeHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	dressType, errr := h.dressTypeSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(dressType).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active dress types
//
//	@Summary		Get all active dress types
//	@Description	Get all active dress types
//	@Tags			DressType
//	@Accept			json
//	@Success		200		{object}	responseModel.DressType
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/dress-type [get]
func (h DressTypeHandler) GetAllDressTypes(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	dressTypes, errr := h.dressTypeSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(dressTypes).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete a DressType
//
//	@Summary		Delete DressType
//	@Description	Deletes an instance of DressType
//	@Tags			DressType
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"dressType id"
//
//	@Router			/dress-type/{id} [delete]
func (h DressTypeHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.dressTypeSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
