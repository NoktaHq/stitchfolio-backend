package base

import "github.com/imkarthi24/sf-backend/internal/handler"

type BaseHandler struct {
	HealthHandler             Health
	UserHandler               *handler.UserHandler
	ChannelHandler            *handler.ChannelHandler
	MasterConfigHandler       *handler.MasterConfigHandler
	AdminHandler              *handler.AdminHandler
	CustomerHandler           *handler.CustomerHandler
	EnquiryHandler            *handler.EnquiryHandler
	OrderHandler              *handler.OrderHandler
	OrderItemHandler          *handler.OrderItemHandler
	MeasurementHandler        *handler.MeasurementHandler
	PersonHandler             *handler.PersonHandler
	DressTypeHandler          *handler.DressTypeHandler
	OrderHistoryHandler       *handler.OrderHistoryHandler
	MeasurementHistoryHandler *handler.MeasurementHistoryHandler
	EnquiryHistoryHandler     *handler.EnquiryHistoryHandler
	ExpenseTrackerHandler     *handler.ExpenseTrackerHandler
	ExpenseDetailHandler     *handler.ExpenseDetailHandler
	TaskHandler               *handler.TaskHandler
	CategoryHandler           *handler.CategoryHandler
	ProductHandler            *handler.ProductHandler
	InventoryHandler          *handler.InventoryHandler
	InventoryLogHandler       *handler.InventoryLogHandler
	VendorHandler             *handler.VendorHandler
	PurchaseHandler            *handler.PurchaseHandler
	PurchaseItemHandler        *handler.PurchaseItemHandler
	DashboardHandler          *handler.DashboardHandler
}

func ProvideBaseHandler(health Health,
	user *handler.UserHandler,
	channelHandler *handler.ChannelHandler,
	masterConfigHandler *handler.MasterConfigHandler,
	adminHandler *handler.AdminHandler,
	customerHandler *handler.CustomerHandler,
	enquiryHandler *handler.EnquiryHandler,
	orderHandler *handler.OrderHandler,
	orderItemHandler *handler.OrderItemHandler,
	measurementHandler *handler.MeasurementHandler,
	personHandler *handler.PersonHandler,
	dressTypeHandler *handler.DressTypeHandler,
	orderHistoryHandler *handler.OrderHistoryHandler,
	measurementHistoryHandler *handler.MeasurementHistoryHandler,
	enquiryHistoryHandler *handler.EnquiryHistoryHandler,
	expenseTrackerHandler *handler.ExpenseTrackerHandler,
	expenseDetailHandler *handler.ExpenseDetailHandler,
	taskHandler *handler.TaskHandler,
	categoryHandler *handler.CategoryHandler,
	productHandler *handler.ProductHandler,
	inventoryHandler *handler.InventoryHandler,
	inventoryLogHandler *handler.InventoryLogHandler,
	vendorHandler *handler.VendorHandler,
	purchaseHandler *handler.PurchaseHandler,
	purchaseItemHandler *handler.PurchaseItemHandler,
	dashboardHandler *handler.DashboardHandler,
) BaseHandler {
	return BaseHandler{
		HealthHandler:             health,
		UserHandler:               user,
		ChannelHandler:            channelHandler,
		MasterConfigHandler:       masterConfigHandler,
		AdminHandler:              adminHandler,
		CustomerHandler:           customerHandler,
		EnquiryHandler:            enquiryHandler,
		OrderHandler:              orderHandler,
		OrderItemHandler:          orderItemHandler,
		MeasurementHandler:        measurementHandler,
		PersonHandler:             personHandler,
		DressTypeHandler:          dressTypeHandler,
		OrderHistoryHandler:       orderHistoryHandler,
		MeasurementHistoryHandler: measurementHistoryHandler,
		EnquiryHistoryHandler:     enquiryHistoryHandler,
		ExpenseTrackerHandler:     expenseTrackerHandler,
		ExpenseDetailHandler:     expenseDetailHandler,
		TaskHandler:               taskHandler,
		CategoryHandler:           categoryHandler,
		ProductHandler:            productHandler,
		InventoryHandler:          inventoryHandler,
		InventoryLogHandler:       inventoryLogHandler,
		VendorHandler:             vendorHandler,
		PurchaseHandler:           purchaseHandler,
		PurchaseItemHandler:       purchaseItemHandler,
		DashboardHandler:          dashboardHandler,
	}
}
