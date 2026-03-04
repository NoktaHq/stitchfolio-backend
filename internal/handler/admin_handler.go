package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	requesModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type AdminHandler struct {
	adminSvc service.AdminService
	resp     response.Response
	dataResp response.DataResponse
}

func ProvideAdminHandler(svc service.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: svc}
}

var _ = (*responseModel.Response)(nil) // used by swagger comments

// Get All Time Dashboard Stats
//
//	@Summary		Switch a record from one channel to another
//	@Description	Switch a record from one channel to another
//	@Tags			Admin
//	@Accept			json
//	@Success		200		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		501		{object}	responseModel.Response
//	@Param			batch	body		requesModel.SwitchBranch	true	"batch"
//
//	@Router			/admin/switch-branch [POST]
func (h AdminHandler) SwitchBranch(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var branch requesModel.SwitchBranch
	err := ctx.Bind(&branch)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.adminSvc.SwitchBranch(&context, &branch)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.DefaultSuccessResponse().FormatAndSend(&context, ctx, http.StatusOK)
}
