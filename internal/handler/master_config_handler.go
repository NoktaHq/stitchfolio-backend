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

type MasterConfigHandler struct {
	masterConfigSvc service.MasterConfigService
	resp            response.Response
	dataResp        response.DataResponse
}

func ProvideMasterConfigHandler(svc service.MasterConfigService) *MasterConfigHandler {
	return &MasterConfigHandler{
		masterConfigSvc: svc,
	}
}

// Create Master Config
//
//	@Summary		Create Master Config
//	@Description	Creates a new master config
//	@Tags			MasterConfig
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			config	body		requestModel.MasterConfig	true	"Master Config"
//	@Router			/masterConfig [post]
func (h MasterConfigHandler) Create(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var config requestModel.MasterConfig
	err := ctx.Bind(&config)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.masterConfigSvc.Save(&context, config)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Master config created successfully").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update Master Config
//
//	@Summary		Update Master Config
//	@Description	Updates an existing master config
//	@Tags			MasterConfig
//	@Accept			json
//	@Success		200		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			id		path		int							true	"Master Config ID"
//	@Param			config	body		requestModel.MasterConfig	true	"Master Config"
//	@Router			/masterConfig/{id} [put]
func (h MasterConfigHandler) Update(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var config requestModel.MasterConfig
	err := ctx.Bind(&config)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.masterConfigSvc.Update(&context, config, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Master config updated successfully").FormatAndSend(&context, ctx, http.StatusOK)
}

// Get Master Config
//
//	@Summary		Get Master Config
//	@Description	Get a specific master config by ID
//	@Tags			MasterConfig
//	@Accept			json
//	@Success		200	{object}	responseModel.DataResponse
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"Master Config ID"
//	@Router			/masterConfig/{id} [get]
func (h MasterConfigHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))

	config, errr := h.masterConfigSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(config).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get Master Config by Name
//
//	@Summary		Get Master Config by Name
//	@Description	Get a specific master config by name
//	@Tags			MasterConfig
//	@Accept			json
//	@Success		200		{object}	responseModel.DataResponse
//	@Failure		400		{object}	responseModel.Response
//	@Param			name	query		string	true	"Master Config Name"
//	@Router			/masterConfig/value [get]
func (h MasterConfigHandler) GetValue(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	name := ctx.Query("name")

	config, errr := h.masterConfigSvc.GetByName(&context, name)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(config).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get Master Config for browse
//
//	@Summary		Get Master Config for browse
//	@Description	Get Master Config for browse
//	@Tags			MasterConfig
//	@Accept			json
//	@Success		200		{object}	responseModel.DataResponse
//	@Failure		400		{object}	responseModel.Response
//	@Param			search	query		string	false	"search"
//	@Router			/masterConfig/browse [get]
func (h MasterConfigHandler) Browse(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	query := ctx.Query("search")
	query = util.EncloseWithSingleQuote(query)

	config, errr := h.masterConfigSvc.Browse(&context, query)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(config).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get Multiple Master Config Values
//
//	@Summary		Get Multiple Master Config Values
//	@Description	Get multiple master config values by their names (format: "Type.Name")
//	@Tags			MasterConfig
//	@Accept			json
//	@Success		200		{object}	responseModel.DataResponse
//	@Failure		400		{object}	responseModel.Response
//	@Param			names	body		[]string	true	"Array of config names (e.g., [\"a.b\", \"a.c\", \"x.a\"])"
//	@Router			/masterConfig/values [post]
func (h MasterConfigHandler) GetMultipleValues(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var names []string
	err := ctx.BindJSON(&names)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	values, errr := h.masterConfigSvc.GetMultipleValues(&context, names)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(values).FormatAndSend(&context, ctx, http.StatusOK)
}
