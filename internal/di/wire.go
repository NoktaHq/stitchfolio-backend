//go:build wireinject
// +build wireinject

// go:build wireinject

package di

import (
	"context"

	"github.com/google/wire"
	"github.com/imkarthi24/sf-backend/internal/app"
	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/cron"
	"github.com/imkarthi24/sf-backend/internal/handler"
	baseHandler "github.com/imkarthi24/sf-backend/internal/handler/base"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/internal/router"
	"github.com/imkarthi24/sf-backend/internal/service"
	baseService "github.com/imkarthi24/sf-backend/internal/service/base"
	"github.com/loop-kar/pixie/db"
	pkgservice "github.com/loop-kar/pixie/service"
)

var appConfigSet = wire.NewSet(
	config.ProvideAppConfig,
	wire.FieldsOf(new(config.AppConfig), "SMTP", "Server", "Database"),
	// wire.FieldsOf(new(config.AppConfig), "Server"),
	// wire.FieldsOf(new(config.AppConfig), "Database"),
)

var pkgServiceSet = wire.NewSet(
	ProvideServiceContainer,
	wire.FieldsOf(new(*pkgservice.Service), "EmailService"),
)

var handlerSet = wire.NewSet(
	baseHandler.ProvideHealthHandler,
	baseHandler.ProvideBaseHandler,
	handler.ProvideUserHandler,
	handler.ProvideChannelHandler,
	handler.ProvideMasterConfigHandler,
	handler.ProvideAdminHandler,
	handler.ProvideCustomerHandler,
	handler.ProvideEnquiryHandler,
	handler.ProvideOrderHandler,
	handler.ProvideOrderItemHandler,
	handler.ProvideMeasurementHandler,
	handler.ProvidePersonHandler,
	handler.ProvideDressTypeHandler,
	handler.ProvideOrderHistoryHandler,
	handler.ProvideMeasurementHistoryHandler,
	handler.ProvideEnquiryHistoryHandler,
	handler.ProvideExpenseTrackerHandler,
	handler.ProvideExpenseDetailHandler,
	handler.ProvideTaskHandler,
	handler.ProvideCategoryHandler,
	handler.ProvideProductHandler,
	handler.ProvideInventoryHandler,
	handler.ProvideInventoryLogHandler,
	handler.ProvideDashboardHandler,
)
var logSet = wire.NewSet(
	ProvideNewRelic,
)

var routerSet = wire.NewSet(
	router.InitRouter,
)

var dbSet = wire.NewSet(
	ProvideDatabaseConnectionParams,
	db.ProvideDatabase,
	db.ProvideDBTransactionManager,
)

var mapperSet = wire.NewSet(
	mapper.ProvideMapper,
	mapper.ProvideResponseMapper,
)

var svcSet = wire.NewSet(
	service.ProvideUserService,
	service.ProvideNotificationService,
	service.ProvideChannelService,
	service.ProvideMasterConfigService,
	service.ProvideAdminService,
	service.ProvideCustomerService,
	service.ProvideEnquiryService,
	service.ProvideOrderService,
	service.ProvideOrderItemService,
	service.ProvideMeasurementService,
	service.ProvidePersonService,
	service.ProvideDressTypeService,
	service.ProvideOrderHistoryService,
	service.ProvideMeasurementHistoryService,
	service.ProvideEnquiryHistoryService,
	service.ProvideExpenseTrackerService,
	service.ProvideExpenseDetailService,
	service.ProvideTaskService,
	service.ProvideCategoryService,
	service.ProvideProductService,
	service.ProvideInventoryService,
	service.ProvideInventoryLogService,
	service.ProvideDashboardService,
)

var baseSvc = wire.NewSet(
	baseService.ProvideBaseService,
)

var repoSet = wire.NewSet(
	repository.ProvideGormDAL,
	repository.ProvideUserRepository,
	repository.ProvideNotificationRepository,
	repository.ProvideChannelRepository,
	repository.ProvideMasterConfigRepository,
	repository.ProvideAdminRepository,
	repository.ProvideCustomerRepository,
	repository.ProvideEnquiryRepository,
	repository.ProvideOrderRepository,
	repository.ProvideOrderItemRepository,
	repository.ProvideMeasurementRepository,
	repository.ProvidePersonRepository,
	repository.ProvideDressTypeRepository,
	repository.ProvideOrderHistoryRepository,
	repository.ProvideMeasurementHistoryRepository,
	repository.ProvideEnquiryHistoryRepository,
	repository.ProvideExpenseTrackerRepository,
	repository.ProvideExpenseDetailRepository,
	repository.ProvideTaskRepository,
	repository.ProvideCategoryRepository,
	repository.ProvideProductRepository,
	repository.ProvideInventoryRepository,
	repository.ProvideInventoryLogRepository,
	repository.ProvideDashboardRepository,
)

var cronSet = wire.NewSet(
	cron.ProvideCron,
)

func InitApp(ctx *context.Context) (*app.App, error) {
	wire.Build(
		appConfigSet,
		pkgServiceSet,
		logSet,
		mapperSet,
		routerSet,
		dbSet,
		repoSet,
		svcSet,
		handlerSet,
		wire.Struct(new(app.App), "*"),
	)
	return &app.App{}, nil

}

func InitJobService(ctx *context.Context) (*app.Task, error) {
	wire.Build(
		appConfigSet,
		pkgServiceSet,
		logSet,
		mapperSet,
		dbSet,
		repoSet,
		baseSvc,
		svcSet,
		cronSet,
		wire.Struct(new(app.Task), "*"),
	)
	return &app.Task{}, nil
}
