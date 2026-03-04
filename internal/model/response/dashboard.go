package responseModel

import "time"

// StatusCountStat holds status and count (orders by status, enquiries by status, etc.).
type StatusCountStat struct {
	Status string `json:"status,omitempty"`
	Count  int    `json:"count,omitempty"`
}

// SourceCountStat holds source and count (e.g. enquiries by source).
type SourceCountStat struct {
	Source string `json:"source,omitempty"`
	Count  int    `json:"count,omitempty"`
}

// ReferrerCountStat holds referrer info and enquiry count (top referrers).
type ReferrerCountStat struct {
	Referrer       string `json:"referrer,omitempty"`
	ReferrerPhone  string `json:"referrerPhone,omitempty"`
	EnquiriesCount int    `json:"enquiriesCount,omitempty"`
}

// TaskDashboardResponse is the API response for the task dashboard.
// Filter by ChannelId (and optionally assignee). Uses DueDate, ReminderDate, IsCompleted, CompletedAt, AssignedToId, Priority.
type TaskDashboardResponse struct {
	OverdueTasks          TaskDashboardTaskList `json:"overdueTasks"`                    // DueDate < today, IsCompleted = false
	DueToday              TaskDashboardTaskList `json:"dueToday"`                         // due today
	DueNext7Days          TaskDashboardTaskList `json:"dueNext7Days"`                     // due in next 7 days
	IncompleteByAssignee  []AssigneeTaskCount   `json:"incompleteByAssignee"`             // count per user
	HighPriorityIncomplete TaskDashboardTaskList `json:"highPriorityIncomplete"`         // Priority set, not done
	UpcomingReminders    TaskDashboardTaskList `json:"upcomingReminders"`                // ReminderDate in next 24–48h
	CompletionRate       CompletionRateStat    `json:"completionRate"`                   // % completed in last 7/30 days
	RecentCompletions    []TaskSummary         `json:"recentCompletions"`                // last N tasks with CompletedAt
}

type TaskDashboardTaskList struct {
	Count int          `json:"count"`
	Tasks []TaskSummary `json:"tasks,omitempty"`
}

type TaskSummary struct {
	ID           uint       `json:"id,omitempty"`
	Title        string     `json:"title,omitempty"`
	DueDate      *time.Time `json:"dueDate,omitempty"`
	ReminderDate *time.Time `json:"reminderDate,omitempty"`
	Priority     *int       `json:"priority,omitempty"`
	IsCompleted  bool       `json:"isCompleted"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	AssignedToId *uint      `json:"assignedToId,omitempty"`
	AssignedTo   string     `json:"assignedTo,omitempty"` // display name
}

type AssigneeTaskCount struct {
	UserID     uint   `json:"userId"`
	UserName   string `json:"userName"`
	Incomplete int    `json:"incomplete"`
}

type CompletionRateStat struct {
	Last7Days  CompletionRateWindow `json:"last7Days,omitempty"`
	Last30Days CompletionRateWindow `json:"last30Days,omitempty"`
}

type CompletionRateWindow struct {
	Completed int     `json:"completed"`
	Total     int     `json:"total"`
	Percent   float64 `json:"percent"`
}

// OrderDashboardResponse is the API response for the order dashboard.
// Filter by ChannelId, date range, status. Uses Order.Status, ExpectedDeliveryDate, DeliveredDate, OrderValue, AdditionalCharges, OrderTakenById; OrderHistory for recent activity.
type OrderDashboardResponse struct {
	OrdersByStatus       []StatusCountStat   `json:"ordersByStatus"`       // count per status DRAFT → DELIVERED, CANCELLED
	OverdueAtRiskOrders  OrderDashboardList `json:"overdueAtRiskOrders"`  // ExpectedDeliveryDate passed or soon, not DELIVERED
	RevenueInPeriod      float64            `json:"revenueInPeriod"`       // sum OrderValue + AdditionalCharges in period
	DeliveriesDueThisWeek OrderDashboardList `json:"deliveriesDueThisWeek"` // by ExpectedDeliveryDate
	RecentDeliveries     OrderDashboardList `json:"recentDeliveries"`      // DeliveredDate last 7/30 days
	OrdersByTakenBy      []UserOrderCount    `json:"ordersByTakenBy"`      // count per OrderTakenById
	OrderCountInPeriod   int                `json:"orderCountInPeriod"`   // last 7/30 days
	RecentOrderActivity  []OrderActivityItem `json:"recentOrderActivity"`  // from OrderHistory
}

type OrderDashboardList struct {
	Count  int             `json:"count"`
	Orders []OrderSummary  `json:"orders,omitempty"`
}

type OrderSummary struct {
	ID                    uint       `json:"id,omitempty"`
	Status                string     `json:"status,omitempty"`
	OrderValue            float64    `json:"orderValue,omitempty"`
	AdditionalCharges     float64    `json:"additionalCharges,omitempty"`
	ExpectedDeliveryDate  *time.Time `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *time.Time `json:"deliveredDate,omitempty"`
	OrderTakenById       *uint      `json:"orderTakenById,omitempty"`
	OrderTakenBy         string    `json:"orderTakenBy,omitempty"`
	CustomerName         string    `json:"customerName,omitempty"`
}

type UserOrderCount struct {
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
}

type OrderActivityItem struct {
	ID            uint      `json:"id,omitempty"`
	OrderId       uint      `json:"orderId,omitempty"`
	Action        string    `json:"action,omitempty"`
	ChangedFields string    `json:"changedFields,omitempty"`
	PerformedAt   time.Time `json:"performedAt,omitempty"`
	PerformedBy   string    `json:"performedBy,omitempty"`
}

// StatsDashboardResponse is the API response for the stats dashboard.
// Aggregates; support date range and ChannelId for revenue, expenses, new customers, task completion.
type StatsDashboardResponse struct {
	RevenueInPeriod       float64             `json:"revenueInPeriod"`       // delivered orders in period
	OrderPipelineValue   float64             `json:"orderPipelineValue"`   // sum value for orders not CANCELLED/DELIVERED
	EnquiriesByStatus    []StatusCountStat   `json:"enquiriesByStatus"`    // new / accepted / callback / closed
	EnquiryOrderConversion *EnquiryConversionStat `json:"enquiryOrderConversion,omitempty"` // orders linked to customers who had enquiries in period
	ExpenseTotalInPeriod float64             `json:"expenseTotalInPeriod"` // Expense.PurchaseDate + Price
	NewCustomersInPeriod int                 `json:"newCustomersInPeriod"`  // Customer.CreatedAt in range
	TaskCompletionInPeriod *CompletionRateStat `json:"taskCompletionInPeriod,omitempty"`   // completed vs created in period
	LowStockItems        []LowStockItem      `json:"lowStockItems"`        // Inventory.Quantity <= LowStockThreshold
	EnquiriesBySource   []SourceCountStat   `json:"enquiriesBySource"`   // Enquiry.Source
	TopReferrers         []ReferrerCountStat `json:"topReferrers"`         // Enquiry.ReferredBy + count
}

type EnquiryConversionStat struct {
	EnquiriesInPeriod int `json:"enquiriesInPeriod"`
	OrdersFromEnquiry int `json:"ordersFromEnquiry"`
}
