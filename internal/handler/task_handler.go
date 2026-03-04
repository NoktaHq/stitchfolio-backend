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

type TaskHandler struct {
	taskSvc  service.TaskService
	resp     response.Response
	dataResp response.DataResponse
}

func ProvideTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{taskSvc: svc}
}

// SaveTask
//
//	@Summary		Save Task
//	@Description	Saves an instance of Task
//	@Tags			Task
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			task	body		requestModel.Task	true	"task"
//	@Router			/task [post]
func (h TaskHandler) SaveTask(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var task requestModel.Task
	err := ctx.Bind(&task)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.taskSvc.SaveTask(&context, task)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// UpdateTask
//
//	@Summary		Update Task
//	@Description	Updates an instance of Task
//	@Tags			Task
//	@Accept			json
//	@Success		202		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			task	body		requestModel.Task	true	"task"
//	@Param			id		path		int					true	"Task id"
//	@Router			/task/{id} [put]
func (h TaskHandler) UpdateTask(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var task requestModel.Task
	err := ctx.Bind(&task)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.taskSvc.UpdateTask(&context, task, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get Task
//
//	@Summary		Get a specific Task
//	@Description	Get an instance of Task
//	@Tags			Task
//	@Accept			json
//	@Success		200	{object}	responseModel.Task
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Task id"
//	@Router			/task/{id} [get]
func (h TaskHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	task, errr := h.taskSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(task).FormatAndSend(&context, ctx, http.StatusOK)
}

// GetAllTasks
//
//	@Summary		Get all active tasks
//	@Description	Get all active tasks
//	@Tags			Task
//	@Accept			json
//	@Success		200		{object}	[]responseModel.Task
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/task [get]
func (h TaskHandler) GetAllTasks(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	tasks, errr := h.taskSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(tasks).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete Task
//
//	@Summary		Delete Task
//	@Description	Deletes an instance of Task
//	@Tags			Task
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"Task id"
//	@Router			/task/{id} [delete]
func (h TaskHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.taskSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
