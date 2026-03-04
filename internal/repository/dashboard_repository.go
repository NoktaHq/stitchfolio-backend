package repository

import (
	"context"
	"time"

	"github.com/imkarthi24/sf-backend/internal/entities"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/errs"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetTaskDashboard(ctx *context.Context, assigneeID *uint) (*responseModel.TaskDashboardResponse, *errs.XError)
	GetOrderDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.OrderDashboardResponse, *errs.XError)
	GetStatsDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.StatsDashboardResponse, *errs.XError)
}

type dashboardRepository struct {
	GormDAL
}

func ProvideDashboardRepository(dal GormDAL) DashboardRepository {
	return &dashboardRepository{GormDAL: dal}
}

func (dr *dashboardRepository) GetTaskDashboard(ctx *context.Context, assigneeID *uint) (*responseModel.TaskDashboardResponse, *errs.XError) {
	now := time.Now().Truncate(24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)
	sevenDaysLater := now.Add(7 * 24 * time.Hour)
	reminderStart := now
	reminderEnd := now.Add(48 * time.Hour)
	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)
	thirtyDaysAgo := now.Add(-30 * 24 * time.Hour)

	baseTask := func() *gorm.DB {
		q := dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive())
		if assigneeID != nil && *assigneeID != 0 {
			q = q.Where("assigned_to_id = ?", *assigneeID)
		}
		return q
	}

	resp := &responseModel.TaskDashboardResponse{}

	// 0. Tasks by status (count per status)
	var statusCounts []struct {
		Status string
		Count  int64
	}
	res := baseTask().Select("status", "count(*) as count").Group("status").Scan(&statusCounts)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard tasks by status", res.Error)
	}
	resp.TasksByStatus = make([]responseModel.StatusCountStat, 0, len(statusCounts))
	for _, s := range statusCounts {
		resp.TasksByStatus = append(resp.TasksByStatus, responseModel.StatusCountStat{Status: s.Status, Count: int(s.Count)})
	}

	// 1. Overdue tasks
	var overdue []entities.Task
	tx := baseTask().Where("status != ?", entities.TaskStatusCompleted).Where("due_date < ?", now)
	tx = tx.Preload("AssignedTo", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&overdue)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard overdue tasks", res.Error)
	}
	resp.OverdueTasks = buildTaskList(overdue)

	// 2. Due today
	var dueToday []entities.Task
	tx = baseTask().Where("status != ?", entities.TaskStatusCompleted).Where("due_date >= ? AND due_date < ?", now, tomorrow)
	tx = tx.Preload("AssignedTo", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&dueToday)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard due today", res.Error)
	}
	resp.DueToday = buildTaskList(dueToday)

	// 3. Due next 7 days (excluding today)
	var dueNext7 []entities.Task
	tx = baseTask().Where("status != ?", entities.TaskStatusCompleted).Where("due_date >= ? AND due_date < ?", tomorrow, sevenDaysLater)
	tx = tx.Preload("AssignedTo", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&dueNext7)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard due next 7 days", res.Error)
	}
	resp.DueNext7Days = buildTaskList(dueNext7)

	// 4. Incomplete by assignee
	type assigneeCount struct {
		AssignedToID *uint
		Count       int64
	}
	var byAssignee []assigneeCount
	res = baseTask().Where("status != ?", entities.TaskStatusCompleted).Select("assigned_to_id, count(*) as count").Group("assigned_to_id").Scan(&byAssignee)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard incomplete by assignee", res.Error)
	}
	resp.IncompleteByAssignee = make([]responseModel.AssigneeTaskCount, 0, len(byAssignee))
	for _, r := range byAssignee {
		name := ""
		if r.AssignedToID != nil && *r.AssignedToID != 0 {
			var u entities.User
			res := dr.WithDB(ctx).Model(&entities.User{}).Select("id", "first_name", "last_name").First(&u, *r.AssignedToID)
			if res.Error == nil {
				name = u.FirstName + " " + u.LastName
			}
		}
		resp.IncompleteByAssignee = append(resp.IncompleteByAssignee, responseModel.AssigneeTaskCount{
			UserID:     uintPtrToUint(r.AssignedToID),
			UserName:   name,
			Incomplete: int(r.Count),
		})
	}

	// 5. High-priority incomplete (Priority set and not done)
	var highPrio []entities.Task
	tx = baseTask().Where("status != ?", entities.TaskStatusCompleted).Where("priority IS NOT NULL AND priority > 0")
	tx = tx.Preload("AssignedTo", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&highPrio)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard high priority", res.Error)
	}
	resp.HighPriorityIncomplete = buildTaskList(highPrio)

	// 6. Upcoming reminders (24-48h)
	var reminders []entities.Task
	tx = baseTask().Where("reminder_date IS NOT NULL AND reminder_date >= ? AND reminder_date <= ?", reminderStart, reminderEnd)
	tx = tx.Preload("AssignedTo", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&reminders)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard reminders", res.Error)
	}
	resp.UpcomingReminders = buildTaskList(reminders)

	// 7. Completion rate (last 7 and 30 days)
	var completed7, total7, completed30, total30 int64
	res = dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("completed_at >= ?", sevenDaysAgo).Where("status = ?", entities.TaskStatusCompleted).Count(&completed7)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard completion rate 7d", res.Error)
	}
	res = dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("created_at >= ?", sevenDaysAgo).Count(&total7)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard total 7d", res.Error)
	}
	res = dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("completed_at >= ?", thirtyDaysAgo).Where("status = ?", entities.TaskStatusCompleted).Count(&completed30)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard completion rate 30d", res.Error)
	}
	res = dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("created_at >= ?", thirtyDaysAgo).Count(&total30)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard total 30d", res.Error)
	}

	resp.CompletionRate = responseModel.CompletionRateStat{
		Last7Days:  responseModel.CompletionRateWindow{Completed: int(completed7), Total: int(total7), Percent: percent(int(completed7), int(total7))},
		Last30Days: responseModel.CompletionRateWindow{Completed: int(completed30), Total: int(total30), Percent: percent(int(completed30), int(total30))},
	}

	// 8. Recent completions (last 10)
	var recent []entities.Task
	tx = dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("status = ?", entities.TaskStatusCompleted).Where("completed_at IS NOT NULL").Order("completed_at DESC").Limit(10)
	tx = tx.Preload("AssignedTo", scopes.SelectFields("first_name", "last_name"))
	if assigneeID != nil && *assigneeID != 0 {
		tx = tx.Where("assigned_to_id = ?", *assigneeID)
	}
	res = tx.Find(&recent)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard recent completions", res.Error)
	}
	resp.RecentCompletions = taskSummaries(recent)

	return resp, nil
}

func buildTaskList(tasks []entities.Task) responseModel.TaskDashboardTaskList {
	return responseModel.TaskDashboardTaskList{
		Count: len(tasks),
		Tasks: taskSummaries(tasks),
	}
}

func taskSummaries(tasks []entities.Task) []responseModel.TaskSummary {
	out := make([]responseModel.TaskSummary, 0, len(tasks))
	for _, t := range tasks {
		name := ""
		if t.AssignedTo != nil {
			name = t.AssignedTo.FirstName + " " + t.AssignedTo.LastName
		}
		out = append(out, responseModel.TaskSummary{
			ID:           t.ID,
			Title:        t.Title,
			DueDate:      t.DueDate,
			ReminderDate: t.ReminderDate,
			Priority:     t.Priority,
			Status:       string(t.Status),
			CompletedAt:  t.CompletedAt,
			AssignedToId: t.AssignedToId,
			AssignedTo:   name,
		})
	}
	return out
}

func (dr *dashboardRepository) GetOrderDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.OrderDashboardResponse, *errs.XError) {
	db := dr.WithDB(ctx)
	now := time.Now().Truncate(24 * time.Hour)
	weekEnd := now.Add(7 * 24 * time.Hour)
	thirtyDaysAgo := now.Add(-30 * 24 * time.Hour)

	if from == nil {
		t := thirtyDaysAgo
		from = &t
	}
	if to == nil {
		to = &now
	}

	// Base query for Order: new session each call so aggregation (Group) never leaks into list queries.
	orderBase := func() *gorm.DB {
		return dr.WithDB(ctx).Model(&entities.Order{}).Scopes(scopes.Channel(), scopes.IsActive())
	}
	// Base with order_value/order_quantity subqueries for list queries that need OrderValue.
	orderBaseWithValue := func() *gorm.DB {
		return orderBase().
			Select(`"stich"."Orders".*,
				(SELECT COALESCE(SUM(quantity), 0) FROM "stich"."OrderItems" WHERE "stich"."OrderItems".order_id = "stich"."Orders".id) as order_quantity,
				(SELECT COALESCE(SUM(total), 0) FROM "stich"."OrderItems" WHERE "stich"."OrderItems".order_id = "stich"."Orders".id) as order_value`)
	}

	resp := &responseModel.OrderDashboardResponse{}

	// 1. Orders by status (aggregation: no need for order_value subquery)
	var statusCounts []struct {
		Status string
		Count  int64
	}
	res := orderBase().Select("status", "count(*) as count").Group("status").Scan(&statusCounts)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard orders by status", res.Error)
	}
	resp.OrdersByStatus = make([]responseModel.StatusCountStat, 0, len(statusCounts))
	for _, s := range statusCounts {
		resp.OrdersByStatus = append(resp.OrdersByStatus, responseModel.StatusCountStat{Status: s.Status, Count: int(s.Count)})
	}

	// 2. Overdue / at-risk (ExpectedDeliveryDate passed or soon, status not DELIVERED)
	var atRisk []entities.Order
	tx := orderBaseWithValue().Where("status != ?", entities.DELIVERED).Where("expected_delivery_date IS NOT NULL AND expected_delivery_date <= ?", weekEnd)
	tx = tx.Preload("Customer", scopes.SelectFields("first_name", "last_name")).
		Preload("OrderTakenBy", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&atRisk)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard at-risk orders", res.Error)
	}
	resp.OverdueAtRiskOrders = orderListFromEntities(atRisk)

	// 3. Revenue in period (OrderValue + AdditionalCharges, by CreatedAt)
	var ordersInPeriod []entities.Order
	res = orderBaseWithValue().Where("created_at >= ? AND created_at <= ?", from, to).Find(&ordersInPeriod)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard revenue", res.Error)
	}
	for _, o := range ordersInPeriod {
		resp.RevenueInPeriod += o.OrderValue + o.AdditionalCharges
	}

	// 4. Deliveries due this week
	var dueThisWeek []entities.Order
	tx = orderBaseWithValue().Where("expected_delivery_date >= ? AND expected_delivery_date < ?", now, weekEnd)
	tx = tx.Preload("Customer", scopes.SelectFields("first_name", "last_name")).
		Preload("OrderTakenBy", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&dueThisWeek)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard deliveries due", res.Error)
	}
	resp.DeliveriesDueThisWeek = orderListFromEntities(dueThisWeek)

	// 5. Recent deliveries (last 30 days)
	var recentDel []entities.Order
	tx = orderBaseWithValue().Where("delivered_date IS NOT NULL AND delivered_date >= ?", thirtyDaysAgo)
	tx = tx.Preload("Customer", scopes.SelectFields("first_name", "last_name")).
		Preload("OrderTakenBy", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&recentDel)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard recent deliveries", res.Error)
	}
	resp.RecentDeliveries = orderListFromEntities(recentDel)

	// 6. Orders taken by user (aggregation)
	var byUser []struct {
		OrderTakenById *uint
		Count         int64
	}
	res = orderBase().Select("order_taken_by_id", "count(*) as count").Group("order_taken_by_id").Scan(&byUser)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard orders by user", res.Error)
	}
	resp.OrdersByTakenBy = make([]responseModel.UserOrderCount, 0, len(byUser))
	for _, r := range byUser {
		name := ""
		if r.OrderTakenById != nil && *r.OrderTakenById != 0 {
			var u entities.User
			res := dr.WithDB(ctx).Model(&entities.User{}).Select("id", "first_name", "last_name").First(&u, *r.OrderTakenById)
			if res.Error == nil {
				name = u.FirstName + " " + u.LastName
			}
		}
		resp.OrdersByTakenBy = append(resp.OrdersByTakenBy, responseModel.UserOrderCount{
			UserID: uintPtrToUint(r.OrderTakenById),
			Name:   name,
			Count:  int(r.Count),
		})
	}

	// 7. Order count in period
	var countInPeriod int64
	res = orderBase().Where("created_at >= ? AND created_at <= ?", from, to).Count(&countInPeriod)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard order count", res.Error)
	}
	resp.OrderCountInPeriod = int(countInPeriod)

	// 8. Recent order activity (OrderHistory)
	var histories []entities.OrderHistory
	tx = db.Model(&entities.OrderHistory{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Order("performed_at DESC").Limit(20).
		Preload("PerformedBy", scopes.SelectFields("first_name", "last_name"))
	res = tx.Find(&histories)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "dashboard order activity", res.Error)
	}
	resp.RecentOrderActivity = make([]responseModel.OrderActivityItem, 0, len(histories))
	for _, h := range histories {
		name := ""
		if h.PerformedBy != nil {
			name = h.PerformedBy.FirstName + " " + h.PerformedBy.LastName
		}
		resp.RecentOrderActivity = append(resp.RecentOrderActivity, responseModel.OrderActivityItem{
			ID:            h.ID,
			OrderId:       h.OrderId,
			Action:        string(h.Action),
			ChangedFields: h.ChangedFields,
			PerformedAt:   h.PerformedAt,
			PerformedBy:   name,
		})
	}

	return resp, nil
}

func orderListFromEntities(orders []entities.Order) responseModel.OrderDashboardList {
	summaries := make([]responseModel.OrderSummary, 0, len(orders))
	for _, o := range orders {
		customerName := ""
		if o.Customer != nil {
			customerName = o.Customer.FirstName + " " + o.Customer.LastName
		}
		takenBy := ""
		if o.OrderTakenBy != nil {
			takenBy = o.OrderTakenBy.FirstName + " " + o.OrderTakenBy.LastName
		}
		summaries = append(summaries, responseModel.OrderSummary{
			ID:                   o.ID,
			Status:               string(o.Status),
			OrderValue:           o.OrderValue,
			AdditionalCharges:    o.AdditionalCharges,
			ExpectedDeliveryDate: o.ExpectedDeliveryDate,
			DeliveredDate:        o.DeliveredDate,
			OrderTakenById:       o.OrderTakenById,
			OrderTakenBy:         takenBy,
			CustomerName:         customerName,
		})
	}
	return responseModel.OrderDashboardList{Count: len(orders), Orders: summaries}
}

func (dr *dashboardRepository) GetStatsDashboard(ctx *context.Context, from, to *time.Time) (*responseModel.StatsDashboardResponse, *errs.XError) {
	now := time.Now().Truncate(24 * time.Hour)
	thirtyDaysAgo := now.Add(-30 * 24 * time.Hour)
	if from == nil {
		from = &thirtyDaysAgo
	}
	if to == nil {
		to = &now
	}

	resp := &responseModel.StatsDashboardResponse{}

	// 1. Revenue (delivered) in period
	var deliveredOrders []entities.Order
	tx := dr.WithDB(ctx).Model(&entities.Order{}).
		Select(`"stich"."Orders".*,
			(SELECT COALESCE(SUM(quantity), 0) FROM "stich"."OrderItems" WHERE "stich"."OrderItems".order_id = "stich"."Orders".id) as order_quantity,
			(SELECT COALESCE(SUM(total), 0) FROM "stich"."OrderItems" WHERE "stich"."OrderItems".order_id = "stich"."Orders".id) as order_value`).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("status = ?", entities.DELIVERED).
		Where("delivered_date >= ? AND delivered_date <= ?", from, to)
	res := tx.Find(&deliveredOrders)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats revenue", res.Error)
	}
	for _, o := range deliveredOrders {
		resp.RevenueInPeriod += o.OrderValue + o.AdditionalCharges
	}

	// 2. Order pipeline value (not CANCELLED/DELIVERED)
	var pipelineOrders []entities.Order
	res = dr.WithDB(ctx).Model(&entities.Order{}).
		Select(`"stich"."Orders".*,
			(SELECT COALESCE(SUM(quantity), 0) FROM "stich"."OrderItems" WHERE "stich"."OrderItems".order_id = "stich"."Orders".id) as order_quantity,
			(SELECT COALESCE(SUM(total), 0) FROM "stich"."OrderItems" WHERE "stich"."OrderItems".order_id = "stich"."Orders".id) as order_value`).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("status NOT IN ?", []entities.OrderStatus{entities.DELIVERED, entities.CANCELLED}).
		Find(&pipelineOrders)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats pipeline", res.Error)
	}
	for _, o := range pipelineOrders {
		resp.OrderPipelineValue += o.OrderValue + o.AdditionalCharges
	}

	// 3. Enquiries by status
	var enqStatus []struct {
		Status string
		Count  int64
	}
	res = dr.WithDB(ctx).Table(`"stich"."Enquiries"`).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Select("status", "count(*) as count").
		Group("status").
		Scan(&enqStatus)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats enquiries by status", res.Error)
	}
	resp.EnquiriesByStatus = make([]responseModel.StatusCountStat, 0, len(enqStatus))
	for _, s := range enqStatus {
		resp.EnquiriesByStatus = append(resp.EnquiriesByStatus, responseModel.StatusCountStat{Status: s.Status, Count: int(s.Count)})
	}

	// 4. Enquiry → order conversion (customers who had enquiries in period and have orders in period)
	var enquiryConv struct {
		Enquiries int64
		Orders    int64
	}
	res = dr.WithDB(ctx).Table(`"stich"."Enquiries"`).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("created_at >= ? AND created_at <= ?", from, to).
		Select("COUNT(DISTINCT customer_id) as enquiries").
		Scan(&enquiryConv)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats enquiry count", res.Error)
	}
	sub := dr.WithDB(ctx).Table(`"stich"."Enquiries"`).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("created_at >= ? AND created_at <= ?", from, to).
		Select("customer_id")
	res = dr.WithDB(ctx).Model(&entities.Order{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("customer_id IN (?)", sub).Where("created_at >= ? AND created_at <= ?", from, to).
		Count(&enquiryConv.Orders)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats enquiry conversion", res.Error)
	}
	resp.EnquiryOrderConversion = &responseModel.EnquiryConversionStat{
		EnquiriesInPeriod: int(enquiryConv.Enquiries),
		OrdersFromEnquiry: int(enquiryConv.Orders),
	}

	// 5. Expense total in period
	var expenseTotal float64
	res = dr.WithDB(ctx).Model(&entities.Expense{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("purchase_date >= ? AND purchase_date <= ?", from, to).
		Select("COALESCE(SUM(price), 0)").Scan(&expenseTotal)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats expenses", res.Error)
	}
	resp.ExpenseTotalInPeriod = expenseTotal

	// 6. New customers in period
	var newCustomers int64
	res = dr.WithDB(ctx).Model(&entities.Customer{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("created_at >= ? AND created_at <= ?", from, to).Count(&newCustomers)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats new customers", res.Error)
	}
	resp.NewCustomersInPeriod = int(newCustomers)

	// 7. Task completion in period
	var taskCompleted, taskTotal int64
	res = dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("completed_at >= ? AND completed_at <= ?", from, to).Where("status = ?", entities.TaskStatusCompleted).Count(&taskCompleted)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats task completed", res.Error)
	}
	res = dr.WithDB(ctx).Model(&entities.Task{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("created_at >= ? AND created_at <= ?", from, to).Count(&taskTotal)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats task total", res.Error)
	}
	resp.TaskCompletionInPeriod = &responseModel.CompletionRateStat{
		Last7Days:  responseModel.CompletionRateWindow{},
		Last30Days: responseModel.CompletionRateWindow{Completed: int(taskCompleted), Total: int(taskTotal), Percent: percent(int(taskCompleted), int(taskTotal))},
	}

	// 8. Low-stock items
	var lowStock []entities.Inventory
	res = dr.WithDB(ctx).Model(&entities.Inventory{}).Scopes(scopes.Channel(), scopes.IsActive()).
		Where("quantity <= low_stock_threshold").
		Preload("Product").Preload("Product.Category").
		Find(&lowStock)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats low stock", res.Error)
	}
	resp.LowStockItems = make([]responseModel.LowStockItem, 0, len(lowStock))
	for _, i := range lowStock {
		name := ""
		sku := ""
		categoryName := ""
		if i.Product != nil {
			name = i.Product.Name
			sku = i.Product.SKU
			if i.Product.Category != nil {
				categoryName = i.Product.Category.Name
			}
		}
		resp.LowStockItems = append(resp.LowStockItems, responseModel.LowStockItem{
			ProductId:         i.ProductId,
			ProductName:       name,
			ProductSKU:        sku,
			CurrentStock:      i.Quantity,
			LowStockThreshold: i.LowStockThreshold,
			CategoryName:      categoryName,
		})
	}

	// 9. Enquiries by source
	var bySource []struct {
		Source string
		Count  int64
	}
	res = dr.WithDB(ctx).Table(`"stich"."Enquiries"`).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("source != ''").
		Select("source", "count(*) as count").
		Group("source").
		Scan(&bySource)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats enquiries by source", res.Error)
	}
	resp.EnquiriesBySource = make([]responseModel.SourceCountStat, 0, len(bySource))
	for _, s := range bySource {
		resp.EnquiriesBySource = append(resp.EnquiriesBySource, responseModel.SourceCountStat{Source: s.Source, Count: int(s.Count)})
	}

	// 10. Top referrers (ReferredBy + count)
	var referrers []struct {
		ReferredBy string
		Count     int64
	}
	res = dr.WithDB(ctx).Table(`"stich"."Enquiries"`).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Where("referred_by != ''").
		Select("referred_by", "count(*) as count").
		Group("referred_by").
		Order("count DESC").
		Limit(10).
		Scan(&referrers)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "stats top referrers", res.Error)
	}
	resp.TopReferrers = make([]responseModel.ReferrerCountStat, 0, len(referrers))
	for _, r := range referrers {
		resp.TopReferrers = append(resp.TopReferrers, responseModel.ReferrerCountStat{
			Referrer:       r.ReferredBy,
			EnquiriesCount: int(r.Count),
		})
	}

	return resp, nil
}

func percent(completed, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(completed) / float64(total) * 100
}

func uintPtrToUint(p *uint) uint {
	if p == nil {
		return 0
	}
	return *p
}
