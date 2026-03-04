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

type MeasurementHistoryHandler struct {
	measurementHistorySvc service.MeasurementHistoryService
	resp                  response.Response
	dataResp              response.DataResponse
}

func ProvideMeasurementHistoryHandler(svc service.MeasurementHistoryService) *MeasurementHistoryHandler {
	return &MeasurementHistoryHandler{measurementHistorySvc: svc}
}

// Save MeasurementHistory
//
//	@Summary		Save MeasurementHistory
//	@Description	Saves an instance of MeasurementHistory
//	@Tags			MeasurementHistory
//	@Accept			json
//	@Success		201					{object}	responseModel.Response
//	@Failure		400					{object}	responseModel.Response
//	@Failure		501					{object}	responseModel.Response
//	@Param			measurementHistory	body		requestModel.MeasurementHistory	true	"measurementHistory"
//	@Router			/measurement-history [post]
func (h MeasurementHistoryHandler) SaveMeasurementHistory(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var measurementHistory requesModel.MeasurementHistory
	err := ctx.Bind(&measurementHistory)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.measurementHistorySvc.SaveMeasurementHistory(&context, measurementHistory)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Get MeasurementHistory
//
//	@Summary		Get a specific MeasurementHistory
//	@Description	Get an instance of MeasurementHistory
//	@Tags			MeasurementHistory
//	@Accept			json
//	@Success		200	{object}	responseModel.MeasurementHistory
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"MeasurementHistory id"
//	@Router			/measurement-history/{id} [get]
func (h MeasurementHistoryHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	measurementHistory, errr := h.measurementHistorySvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(measurementHistory).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active measurement histories
//
//	@Summary		Get all active measurement histories
//	@Description	Get all active measurement histories
//	@Tags			MeasurementHistory
//	@Accept			json
//	@Success		200		{object}	responseModel.MeasurementHistory
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/measurement-history [get]
func (h MeasurementHistoryHandler) GetAllMeasurementHistories(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	measurementHistories, errr := h.measurementHistorySvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(measurementHistories).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get measurement histories by measurement id
//
//	@Summary		Get measurement histories by measurement id
//	@Description	Get measurement histories by measurement id
//	@Tags			MeasurementHistory
//	@Accept			json
//	@Success		200				{object}	responseModel.MeasurementHistory
//	@Failure		400				{object}	responseModel.DataResponse
//	@Param			measurementId	path		int	true	"measurement id"
//	@Router			/measurement-history/measurement/{measurementId} [get]
func (h MeasurementHistoryHandler) GetByMeasurementId(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	measurementId, _ := strconv.Atoi(ctx.Param("measurementId"))

	measurementHistories, errr := h.measurementHistorySvc.GetByMeasurementId(&context, uint(measurementId))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(measurementHistories).FormatAndSend(&context, ctx, http.StatusOK)
}
