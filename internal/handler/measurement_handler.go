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

type MeasurementHandler struct {
	measurementSvc service.MeasurementService
	resp           response.Response
	dataResp       response.DataResponse
}

func ProvideMeasurementHandler(svc service.MeasurementService) *MeasurementHandler {
	return &MeasurementHandler{measurementSvc: svc}
}

// Save Measurement
//
//	@Summary		Save Measurement
//	@Description	Saves an instance of Measurement
//	@Tags			Measurement
//	@Accept			json
//	@Success		201			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			measurement	body		requestModel.Measurement	true	"measurement"
//	@Router			/measurement [post]
func (h MeasurementHandler) SaveMeasurement(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var measurement requesModel.Measurement
	err := ctx.Bind(&measurement)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.measurementSvc.SaveMeasurement(&context, measurement)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Save Bulk Measurements
//
//	@Summary		Save Bulk Measurements
//	@Description	Saves multiple measurements for multiple persons in bulk
//	@Tags			Measurement
//	@Accept			json
//	@Success		201				{object}	responseModel.Response
//	@Failure		400				{object}	responseModel.Response
//	@Failure		501				{object}	responseModel.Response
//	@Param			measurements	body		[]requestModel.BulkMeasurementRequest	true	"Array of bulk measurement requests"
//	@Router			/measurement/bulk [post]
func (h MeasurementHandler) SaveBulkMeasurements(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var bulkRequests []requesModel.BulkMeasurementRequest
	err := ctx.Bind(&bulkRequests)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.measurementSvc.SaveBulkMeasurements(&context, bulkRequests)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Bulk save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update Measurement
//
//	@Summary		Update Measurement
//	@Description	Updates a single measurement by its ID
//	@Tags			Measurement
//	@Accept			json
//	@Success		200			{object}	responseModel.Response
//	@Failure		400			{object}	responseModel.Response
//	@Failure		501			{object}	responseModel.Response
//	@Param			id			path		int							true	"Measurement id"
//	@Param			measurement	body		requestModel.Measurement	true	"Measurement to update"
//	@Router			/measurement/{id} [put]
func (h MeasurementHandler) UpdateMeasurement(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var measurement requesModel.Measurement
	err := ctx.Bind(&measurement)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.measurementSvc.UpdateMeasurement(&context, measurement, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusOK)
}

// Bulk Update Measurements
//
//	@Summary		Bulk Update Measurements
//	@Description	Updates an array of measurements by their IDs
//	@Tags			Measurement
//	@Accept			json
//	@Success		200				{object}	responseModel.Response
//	@Failure		400				{object}	responseModel.Response
//	@Failure		501				{object}	responseModel.Response
//	@Param			measurements	body		[]requestModel.Measurement	true	"Array of measurements to update"
//	@Router			/measurement/bulk [put]
func (h MeasurementHandler) BulkUpdateMeasurements(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var measurements []requesModel.Measurement
	err := ctx.Bind(&measurements)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.measurementSvc.BulkUpdateMeasurements(&context, measurements)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Bulk update success").FormatAndSend(&context, ctx, http.StatusOK)
}

// Get Measurement
//
//	@Summary		Get a specific Measurement
//	@Description	Get an instance of Measurement
//	@Tags			Measurement
//	@Accept			json
//	@Success		200	{object}	responseModel.Measurement
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Measurement id"
//	@Router			/measurement/{id} [get]
func (h MeasurementHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	measurement, errr := h.measurementSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(measurement).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active measurements
//
//	@Summary		Get all active measurements grouped by PersonId
//	@Description	Get all active measurements grouped by PersonId with DressTypes as CSV. Searches by Customer Name.
//	@Tags			Measurement
//	@Accept			json
//	@Success		200		{object}	responseModel.MeasurementBrowse
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search by Customer Name (returns all Persons of that customer)"
//	@Router			/measurement [get]
func (h MeasurementHandler) GetAllMeasurements(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	//search = util.EncloseWithSingleQuote(search)

	measurements, errr := h.measurementSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(measurements).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete a Measurement
//
//	@Summary		Delete Measurement
//	@Description	Deletes an instance of Measurement
//	@Tags			Measurement
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"measurement id"
//
//	@Router			/measurement/{id} [delete]
func (h MeasurementHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.measurementSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
