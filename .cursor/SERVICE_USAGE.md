# Service Usage Guide (Branch Changes)

This document describes how to use the services and APIs that changed on this branch compared to main. It is intended for developers and for Cursor AI when working in this codebase.

---

## 1. Task & Dashboard

### Task status (replacing IsCompleted)

Tasks now use a **status** field instead of **isCompleted**.

- **Entity:** `internal/entities/task.go`  
  - `TaskStatus`: `PENDING` | `IN_PROGRESS` | `COMPLETED` | `CANCELLED`  
  - Default: `PENDING`

- **Request body** (create/update): send `status` (string), not `isCompleted`.  
  Example: `{ "title": "Review design", "status": "IN_PROGRESS", ... }`

- **Response:** Task and TaskSummary include `status` (string); `isCompleted` is no longer returned.

- **Filtering:** List tasks with query filters using `status` (e.g. `status=PENDING` or `status=COMPLETED`). See `internal/repository/scopes/task_scopes.go` — filter map uses `"Status": "status"`.

### Task API (HTTP)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/task` | Create task — body: `requestModel.Task` (include `status`) |
| PUT | `/task/:id` | Update task — body: `requestModel.Task` |
| GET | `/task/:id` | Get one task — response: `responseModel.Task` (has `status`) |
| GET | `/task` | List tasks — query: `search` (optional); response: `[]responseModel.Task` |
| DELETE | `/task/:id` | Soft-delete task |

### Task dashboard API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/dashboard/task` | Task dashboard — query: `assigneeId` (optional, filter by user) |

**Response (`TaskDashboardResponse`):**

- **tasksByStatus** — Count per status: `[{ "status": "PENDING", "count": 5 }, ...]`. Use for status breakdown (e.g. cards or charts).
- **overdueTasks**, **dueToday**, **dueNext7Days** — Tasks with `status != COMPLETED` in the respective date ranges.
- **incompleteByAssignee** — Count of incomplete tasks per assignee.
- **highPriorityIncomplete** — Incomplete tasks with priority set.
- **upcomingReminders** — Tasks with reminder in next 24–48h.
- **completionRate** — `last7Days` / `last30Days`: completed vs created in window.
- **recentCompletions** — Last 10 completed tasks.

All task dashboard queries respect channel and (when provided) `assigneeId`.  
See `docs/dashboard-response-fields.md` for exact field calculations.

---

## 2. File Store (Temp Upload → Confirm → Entity Document)

The **file store** supports **temp upload** of files, then **confirm on save/update** of an entity (OrderItem, Expense). Full structure, tables, and flows are documented in **`.cursor/FILE_STORE_STRUCTURE.md`**. Summary below.

### Flow

1. **Temp upload:** Client calls `POST /file-store/temp` (or `/file-store/temp/bulk`) with `multipart/form-data` (field `file` or `files`). Response: `{ data: { id, tempFileKey, fileName } }` per file.
2. **Confirm:** Client includes in the entity payload a **`files`** array of **ConfirmFile**: `{ id, fileKey, kind, description }` (id/fileKey from step 1). On **Save** or **Update** of that entity, the backend:
   - Inserts **entity_document** (type, documentType = kind, entityName, entityId, description).
   - Moves file from temp S3 key to final key `{entityName}/{entityId}/{kind}`.
   - Updates **file_store_metadata** (entityId, entityType, fileKey, fileUrl).
3. **Get:** When getting the entity, backend loads **entity_document** for that entity and resolves file URLs from **file_store_metadata**; response includes **`files`** (`[]EntityDocument` with `document.fileUrl`).

### Entities with files

- **OrderItem:** `Order.OrderItems[].Files` (request: `[]ConfirmFile`, response: `[]EntityDocument`). Confirm on SaveOrder/UpdateOrder; fill on Get/GetAll Order.
- **Expense (ExpenseTracker):** `ExpenseTracker.Files`. Confirm on SaveExpenseTracker/UpdateExpenseTracker; fill on Get/GetAll.

### Key services

- **FileStoreService:** `UploadTemp`, `ConfirmTempUpload` (returns newUrl), `UpdateEntityIdAndKey`, `GetFileKey`.
- **EntityDocumentService:** `SaveEntityDocument`, `GetEntityDocumentsByEntity` (resolves file URL per row using entityName/entityId/type).

### Reference

- **Full structure, tables, models, and “add files to new entity”:** `.cursor/FILE_STORE_STRUCTURE.md`
- **Handler:** `internal/handler/file_store_handler.go` (UploadTemp, UploadTempBulk)
- **Models:** `internal/model/models/file.go` (ConfirmFile), response EntityDocument in `internal/model/response/entity_documents.go`

---

## 3. Expense Tracker & Expense Detail (Balance)

### Balance on Expense

- **Entity:** `Expense` has a **Balance** field: `Price - Sum(ExpenseDetail.Price)` (remaining amount).
- **Request/Response:** `ExpenseTracker` request and response models include **balance** (optional in request; set by backend).

### Recalculation

- **RecalculateAndUpdateBalance(ctx, expenseId)** (on `ExpenseTrackerRepository`) sets  
  `Expense.Balance = Expense.Price - Sum(ExpenseDetail.Price)` for that expense.
- **ExpenseDetailService** calls it after:
  - **Save** (add detail) → recalculates that expense.
  - **Update** (edit detail) → recalculates that expense.
  - **Delete** → recalculates the expense the detail belonged to.
- **ExpenseTrackerService** calls it after **Save** and **Update** of the expense (so when `Price` changes, balance is updated).

### Using Expense Detail from another service

If you add/update/delete expense details programmatically, prefer using **ExpenseDetailService** so balance stays in sync. If you bypass it and write to the repository directly, you must call **RecalculateAndUpdateBalance** on the expense (via `ExpenseTrackerRepository`) after the change.

### Migration

- Migration **008_add_expense_balance** adds the `balance` column to the Expense table.  
- Ensure migrations are run when deploying this branch (e.g. `make migrate-dev` or app startup with `--migrate true`).

---

## 4. Measurement & MeasurementHistory (JSON type)

- **Request:** `requestModel.Measurement` and `BulkMeasurementItem` use **`entitiy_types.JSON`** for **Values** (replacing `json.RawMessage`).
- **Response:** `responseModel.Measurement` uses **`entitiy_types.JSON`** for **Values**; **MeasurementHistory** uses **`entitiy_types.JSON`** for **OldValues** (replacing `RawMessage`).
- **Type:** `internal/entities/types/json.go` — use this JSON type for request/response and entity so marshalling and DB storage stay consistent.

---

## 5. Other branch notes

- **Enquiry / DressType response:** Audit fields (`createdAt`, `updatedAt`, etc.) are now included in the response mapper for Enquiry and DressType.
- **Wire:** `ExpenseDetailService` now depends on **ExpenseTrackerRepository** (in addition to ExpenseDetailRepository) for balance recalculation. Wire is already updated in `internal/di/wire_gen.go`.
- **App migrate:** `internal/app/app.go` currently enables migration for **Expense** (e.g. `008_add_expense_balance`); Task/Inventory/Product/Category migrations are commented out. Toggle as needed for your environment.

---

## Quick reference

| Area | Main change | Where to look |
|------|-------------|----------------|
| Tasks | `status` instead of `isCompleted` | `entities/task.go`, request/response task models, task_scopes, dashboard_repository |
| Task dashboard | New `tasksByStatus` | `TaskDashboardResponse`, `GetTaskDashboard` |
| File Store | Temp upload → confirm on save; entity_document + file_store_metadata | `.cursor/FILE_STORE_STRUCTURE.md`, `file_store_service.go`, `order_service.go`, `expense_tracker_service.go` |
| Expense | Balance + recalc on detail changes | `Expense.Balance`, `ExpenseDetailService`, `RecalculateAndUpdateBalance` |
| Measurement | JSON type for Values/OldValues | `entities/types/json.go`, request/response measurement models |

For dashboard field formulas, see **`docs/dashboard-response-fields.md`**.
