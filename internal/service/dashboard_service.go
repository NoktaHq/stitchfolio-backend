package service

import (
	"context"
	"time"

	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type DashboardService interface {
	GetTaskDashboard(ctx *context.Context, assigneeID *uint) (*responseModel.TaskDashboardResponse, *errs.XError)
	GetOrderDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.OrderDashboardResponse, *errs.XError)
	GetStatsDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.StatsDashboardResponse, *errs.XError)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func ProvideDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetTaskDashboard(ctx *context.Context, assigneeID *uint) (*responseModel.TaskDashboardResponse, *errs.XError) {
	return s.dashboardRepo.GetTaskDashboard(ctx, assigneeID)
}

func (s *dashboardService) GetOrderDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.OrderDashboardResponse, *errs.XError) {
	return s.dashboardRepo.GetOrderDashboard(ctx, from, to)
}

func (s *dashboardService) GetStatsDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.StatsDashboardResponse, *errs.XError) {
	return s.dashboardRepo.GetStatsDashboard(ctx, from, to)
}
