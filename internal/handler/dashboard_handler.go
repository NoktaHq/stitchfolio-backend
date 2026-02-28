package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type DashboardHandler struct {
	dashboardSvc service.DashboardService
	resp         response.Response
	dataResp     response.DataResponse
}

func ProvideDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: svc}
}

// GetTaskDashboard
//
//	@Summary		Task dashboard
//	@Description	Returns task dashboard: overdue, due today/next 7 days, by assignee, high priority, reminders, completion rate, recent completions.
//	@Tags			Dashboard
//	@Accept			json
//	@Success		200			{object}	responseModel.TaskDashboardResponse
//	@Failure		400			{object}	responseModel.Response
//	@Param			assigneeId	query		int	false	"Filter by assignee user ID"
//	@Router			/dashboard/task [get]
func (h *DashboardHandler) GetTaskDashboard(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	var assigneeID *uint
	if idStr := ctx.Query("assigneeId"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			uid := uint(id)
			assigneeID = &uid
		}
	}
	data, err := h.dashboardSvc.GetTaskDashboard(&c, assigneeID)
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(data).FormatAndSend(&c, ctx, http.StatusOK)
}

// GetOrderDashboard
//
//	@Summary		Order dashboard
//	@Description	Returns order dashboard: by status, overdue/at-risk, revenue, deliveries due, recent deliveries, by taken-by, count in period, recent activity.
//	@Tags			Dashboard
//	@Accept			json
//	@Success		200		{object}	responseModel.OrderDashboardResponse
//	@Failure		400		{object}	responseModel.Response
//	@Param			from	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to		query		string	false	"End date (YYYY-MM-DD)"
//	@Router			/dashboard/order [get]
func (h *DashboardHandler) GetOrderDashboard(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	from, to := parseDateRange(ctx, "from", "to")
	data, err := h.dashboardSvc.GetOrderDashboard(&c, from, to)
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(data).FormatAndSend(&c, ctx, http.StatusOK)
}

// GetStatsDashboard
//
//	@Summary		Stats dashboard
//	@Description	Returns stats dashboard: revenue, pipeline value, enquiries by status/source, conversion, expenses, new customers, task completion, low stock, top referrers.
//	@Tags			Dashboard
//	@Accept			json
//	@Success		200		{object}	responseModel.StatsDashboardResponse
//	@Failure		400		{object}	responseModel.Response
//	@Param			from	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to		query		string	false	"End date (YYYY-MM-DD)"
//	@Router			/dashboard/stats [get]
func (h *DashboardHandler) GetStatsDashboard(ctx *gin.Context) {
	c := util.CopyContextFromGin(ctx)
	from, to := parseDateRange(ctx, "from", "to")
	data, err := h.dashboardSvc.GetStatsDashboard(&c, from, to)
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&c, ctx, http.StatusBadRequest)
		return
	}
	h.dataResp.DefaultSuccessResponse(data).FormatAndSend(&c, ctx, http.StatusOK)
}

func parseDateRange(ctx *gin.Context, fromKey, toKey string) (from, to *time.Time) {
	if s := ctx.Query(fromKey); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			t = t.Truncate(24 * time.Hour)
			from = &t
		}
	}
	if s := ctx.Query(toKey); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			t = t.Truncate(24 * time.Hour)
			to = &t
		}
	}
	return from, to
}
