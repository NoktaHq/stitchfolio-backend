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

type ChannelHandler struct {
	channelSvc service.ChannelService
	resp       response.Response
	dataResp   response.DataResponse
}

func ProvideChannelHandler(svc service.ChannelService) *ChannelHandler {
	return &ChannelHandler{channelSvc: svc}
}

// Save Channel
//
//	@Summary		Save Channel
//	@Description	Saves an instance of Channel
//	@Tags			Channel
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			channel	body		requestModel.Channel	true	"channel"
//	@Router			/channel [post]
func (h ChannelHandler) SaveChannel(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	var channel requesModel.Channel
	err := ctx.Bind(&channel)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.channelSvc.SaveChannel(&context, channel)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)

}

// Update Channel
//
//	@Summary		Update Channel
//	@Description	Updates an instance of Channel
//	@Tags			Channel
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			channel	body		requestModel.Channel	true	"channel"
//	@Param			id		path		int						true	"Channel id"
//	@Router			/channel/{id} [put]
func (h ChannelHandler) UpdateChannel(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	var channel requesModel.Channel
	err := ctx.Bind(&channel)

	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.channelSvc.UpdateChannel(&context, channel, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)

}

// Get Channel
//
//	@Summary		Get a specific  Channel
//	@Description	Get an instance of Channel
//	@Tags			Channel
//	@Accept			json
//	@Success		200	{object}	responseModel.Channel
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Channel id"
//	@Router			/channel/{id} [get]
func (h ChannelHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	channel, errr := h.channelSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(channel).FormatAndSend(&context, ctx, http.StatusOK)

}

// Get all active channels
//
//	@Summary		Get all active channels
//	@Description	Get all active channels
//	@Tags			Channel
//	@Accept			json
//	@Success		200		{object}	responseModel.Channel
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			name	query		string	false	"name"
//	@Router			/channel [get]
func (h ChannelHandler) GetAllChannels(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	autoCompName := ctx.Query("name")
	autoCompName = util.EncloseWithSingleQuote(autoCompName)

	channels, errr := h.channelSvc.GetAllChannels(&context, autoCompName)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(channels).FormatAndSend(&context, ctx, http.StatusOK)

}

// Autocomplete for channels
//
//	@Summary		Autocomplete for channels
//	@Description	Autocomplete for channels
//	@Tags			Channel
//	@Accept			json
//	@Success		200		{object}	responseModel.ChannelAutoComplete
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			name	query		string	false	"name"
//	@Router			/channel/autocomplete [get]
func (h ChannelHandler) ChannelAutoComplete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	autoCompName := ctx.Query("name")
	autoCompName = util.EncloseWithSingleQuote(autoCompName)

	channels, errr := h.channelSvc.ChannelAutoComplete(&context, autoCompName)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(channels).FormatAndSend(&context, ctx, http.StatusOK)

}

// Delete a Channel
//
//	@Summary		Delete Channel
//	@Description	Deletes an instance of Channel
//	@Tags			Channel
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"Channel id"
//
//	@Router			/channel/{id} [delete]
func (h ChannelHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.channelSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)

}
