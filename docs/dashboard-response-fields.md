# Dashboard Response Fields — How They Are Calculated

This document describes how each field in the three dashboard API responses is computed in `internal/repository/dashboard_repository.go`.

---

## Completion rate (summary)

**Task Dashboard** (`completionRate`):

- **Last 7 days**: `completed` = count of tasks with `status = 'COMPLETED'` and `completed_at >= sevenDaysAgo`. `total` = count of all tasks with `created_at >= sevenDaysAgo`. **Percent** = `(completed / total) * 100` (0 if total is 0).
- **Last 30 days**: Same logic over the last 30 days (`completed_at` / `created_at >= thirtyDaysAgo`).

So completion rate = **tasks completed in the window** ÷ **tasks created in the same window** × 100. All counts use channel + active scopes; Task Dashboard can be filtered by `assigneeID`.

**Stats Dashboard** (`taskCompletionInPeriod`):

- **Last 7 days**: Not populated (empty).
- **Last 30 days**: `completed` = tasks with `completed_at` in the requested `from`–`to` range and `status = 'COMPLETED'`. `total` = tasks with `created_at` in the same range. **Percent** = `(completed / total) * 100`.

So it’s the same formula, but over the **requested date range** instead of fixed 7/30-day windows.

---

## 1. Task Dashboard (`TaskDashboardResponse`)

| Field | Calculation |
|-------|-------------|
| **tasksByStatus** | `GROUP BY status`, count per status (PENDING, IN_PROGRESS, COMPLETED, CANCELLED). Channel + active; optional assignee filter. |
| **overdueTasks** | Tasks where `status != 'COMPLETED'`, `due_date < today` (midnight). Channel + active; optional assignee filter. |
| **dueToday** | Tasks where `status != 'COMPLETED'`, `due_date` in [today 00:00, tomorrow 00:00). |
| **dueNext7Days** | Tasks where `status != 'COMPLETED'`, `due_date` in [tomorrow 00:00, today + 7 days). |
| **incompleteByAssignee** | For each `assigned_to_id`, count of tasks with `status != 'COMPLETED'`. User names resolved from `User`. |
| **highPriorityIncomplete** | Tasks with `status != 'COMPLETED'` and `priority IS NOT NULL AND priority > 0`. |
| **upcomingReminders** | Tasks with `reminder_date` in the next 24–48 hours from now (inclusive). |
| **completionRate** | **Last 7 days**: `completed` = tasks with `completed_at >= sevenDaysAgo` and `status = 'COMPLETED'`; `total` = tasks with `created_at >= sevenDaysAgo`. **Last 30 days**: same with `thirtyDaysAgo`. **Percent** = `(completed / total) * 100`, 0 if total is 0. |
| **recentCompletions** | Last 10 tasks with `status = 'COMPLETED'` and `completed_at IS NOT NULL`, ordered by `completed_at DESC`. Optional assignee filter. |

All task queries use channel and active scopes; `assigneeID` filters where applicable.

---

## 2. Order Dashboard (`OrderDashboardResponse`)

Date range: if `from` is nil, defaults to 30 days ago; if `to` is nil, defaults to “today” (midnight). All order queries use channel + active.

| Field | Calculation |
|-------|-------------|
| **ordersByStatus** | `GROUP BY status`, count per status (DRAFT, DELIVERED, CANCELLED, etc.). |
| **overdueAtRiskOrders** | Orders where `status != DELIVERED`, `expected_delivery_date IS NOT NULL` and `expected_delivery_date <= end of this week` (today + 7 days). Order value from OrderItems subquery. |
| **revenueInPeriod** | Sum of `(OrderValue + AdditionalCharges)` for orders with `created_at` in `[from, to]`. OrderValue = sum of OrderItems `total`. |
| **deliveriesDueThisWeek** | Orders with `expected_delivery_date` in [today 00:00, today + 7 days). |
| **recentDeliveries** | Orders with `delivered_date IS NOT NULL` and `delivered_date >= thirtyDaysAgo`. |
| **ordersByTakenBy** | `GROUP BY order_taken_by_id`, count per user; user names from `User`. |
| **orderCountInPeriod** | Count of orders with `created_at` in `[from, to]`. |
| **recentOrderActivity** | Last 20 `OrderHistory` rows (channel + active), ordered by `performed_at DESC`, with performer name. |

---

## 3. Stats Dashboard (`StatsDashboardResponse`)

Date range: if `from`/`to` are nil, defaults to last 30 days to “today”. All queries use channel + active unless noted.

| Field | Calculation |
|-------|-------------|
| **revenueInPeriod** | Sum of `(OrderValue + AdditionalCharges)` for orders with `status = DELIVERED` and `delivered_date` in `[from, to]`. OrderValue from OrderItems. |
| **orderPipelineValue** | Sum of `(OrderValue + AdditionalCharges)` for orders where `status` is not `DELIVERED` or `CANCELLED`. |
| **enquiriesByStatus** | Enquiries table: `GROUP BY status`, count per status. |
| **enquiryOrderConversion** | **EnquiriesInPeriod**: count of distinct `customer_id` in Enquiries with `created_at` in `[from, to]`. **OrdersFromEnquiry**: count of Orders (same channel/active) where `customer_id` is in that enquiry customer set and order `created_at` is in `[from, to]`. |
| **expenseTotalInPeriod** | Sum of `Expense.price` where `purchase_date` is in `[from, to]`. |
| **newCustomersInPeriod** | Count of Customers with `created_at` in `[from, to]`. |
| **taskCompletionInPeriod** | **Last 7 days**: empty. **Last 30 days**: `completed` = tasks with `completed_at` in `[from, to]` and `status = 'COMPLETED'`; `total` = tasks with `created_at` in `[from, to]`; **Percent** = `(completed / total) * 100`. |
| **lowStockItems** | Inventory rows where `quantity <= low_stock_threshold`, with Product and Category preloaded. |
| **enquiriesBySource** | Enquiries where `source != ''`: `GROUP BY source`, count per source. |
| **topReferrers** | Enquiries where `referred_by != ''`: `GROUP BY referred_by`, count, ordered by count DESC, limit 10. |

---

## Helper used for percentages

- **percent(completed, total int) float64**: returns `(completed / total) * 100`; returns `0` if `total == 0`.

---

## Summary: Completion rate in one place

- **Task Dashboard `completionRate`**: For last 7 and last 30 days, **percent = (tasks completed in that window) / (tasks created in that window) × 100**.
- **Stats Dashboard `taskCompletionInPeriod`**: Same formula over the **requested `from`–`to` period**; only the “Last 30 days” slot is filled (Last 7 days is empty).

So “completion rate” everywhere is **completed-in-window ÷ created-in-window**, not “completed ÷ all tasks”.
