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

type PersonHandler struct {
	personSvc service.PersonService
	resp      response.Response
	dataResp  response.DataResponse
}

func ProvidePersonHandler(svc service.PersonService) *PersonHandler {
	return &PersonHandler{personSvc: svc}
}

// Save Person
//
//	@Summary		Save Person
//	@Description	Saves an instance of Person
//	@Tags			Person
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			person	body		requestModel.Person	true	"person"
//	@Router			/person [post]
func (h PersonHandler) SavePerson(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var person requesModel.Person
	err := ctx.Bind(&person)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.personSvc.SavePerson(&context, person)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update Person
//
//	@Summary		Update Person
//	@Description	Updates an instance of Person
//	@Tags			Person
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			person	body		requestModel.Person	true	"person"
//	@Param			id		path		int					true	"Person id"
//	@Router			/person/{id} [put]
func (h PersonHandler) UpdatePerson(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var person requesModel.Person
	err := ctx.Bind(&person)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.personSvc.UpdatePerson(&context, person, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get Person
//
//	@Summary		Get a specific Person
//	@Description	Get an instance of Person
//	@Tags			Person
//	@Accept			json
//	@Success		200	{object}	responseModel.Person
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Person id"
//	@Router			/person/{id} [get]
func (h PersonHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	person, errr := h.personSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(person).FormatAndSend(&context, ctx, http.StatusOK)
}

// Get all active persons
//
//	@Summary		Get all active persons
//	@Description	Get all active persons
//	@Tags			Person
//	@Accept			json
//	@Success		200		{object}	responseModel.Person
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/person [get]
func (h PersonHandler) GetAllPersons(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	persons, errr := h.personSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(persons).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete a Person
//
//	@Summary		Delete Person
//	@Description	Deletes an instance of Person
//	@Tags			Person
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"person id"
//
//	@Router			/person/{id} [delete]
func (h PersonHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.personSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}

// Get persons by customer id
//
//	@Summary		Get persons by customer id
//	@Description	Get persons by customer id
//	@Tags			Person
//	@Accept			json
//	@Success		200			{object}	responseModel.Person
//	@Failure		400			{object}	responseModel.DataResponse
//	@Param			customerId	path		int	true	"customer id"
//	@Router			/person/customer/{customerId} [get]
func (h PersonHandler) GetByCustomerId(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	customerId, _ := strconv.Atoi(ctx.Param("customerId"))

	persons, errr := h.personSvc.GetByCustomerId(&context, uint(customerId))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(persons).FormatAndSend(&context, ctx, http.StatusOK)
}
