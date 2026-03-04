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

## 2. File Store Service

The **File Store Service** stores file metadata in the DB and files in S3. It is used **in-process** by other services (e.g. when uploading a file for an entity). There is no dedicated HTTP handler; wire injects `FileStoreService` where needed.

**Location:** `internal/service/file_store_service.go`  
**Repository:** `internal/repository/file_store_repository.go`  
**Request/Response:** `internal/model/request/file_store_metadata.go`, `internal/model/response/file_store_metadata.go`

### Interface

```go
type FileStoreService interface {
    SaveFileStoreMetadata(*context.Context, requestModel.FileStoreMetadata) (uint, *errs.XError)
    UpdateFileStoreMetadata(*context.Context, requestModel.FileStoreMetadata, uint) *errs.XError
    GetFileStoreMetadata(*context.Context, uint) (*responseModel.FileStoreMetadata, *errs.XError)
    DeleteFileStoreMetadata(*context.Context, uint) *errs.XError

    GetFileStoreMetadataByKey(ctx *context.Context, entityName string, entityId uint, kind string) (*responseModel.FileStoreMetadata, *errs.XError)
    GetFileKey(ctx *context.Context, entityName string, entityId uint, fileType string) string
    Upload(ctx *context.Context, file models.FileUpload) *errs.XError
    GetFileMetadataIfExists(ctx *context.Context, entityType string, id uint, kind string) (bool, *responseModel.FileStoreMetadata, *errs.XError)
}
```

### How to use

**1. Upload a file (recommended)**  
- Call **`Upload(ctx, file)`** with `models.FileUpload` containing:
  - `EntityType` (e.g. `"Expense"`, `"Order"`),
  - `EntityId`,
  - `Kind` (file type/key suffix),
  - `Content` (bytes),
  - `Metadata` (e.g. `Filename`, `Size`, `Header` for Content-Type).
- The service will:
  - Upload to S3 using key `{entityType}/{entityId}/{kind}`.
  - Create or update `FileStoreMetadata` with presigned URL and metadata.
- Use this when the client uploads a file for a known entity.

**2. Metadata-only (no S3 upload)**  
- **Save:** `SaveFileStoreMetadata(ctx, req)` — returns new ID.  
- **Update:** `UpdateFileStoreMetadata(ctx, req, id)`.  
- **Get by ID:** `GetFileStoreMetadata(ctx, id)`.  
- **Get by entity + kind:** `GetFileStoreMetadataByKey(ctx, entityType, entityId, kind)`.  
- **Delete:** `DeleteFileStoreMetadata(ctx, id)` (soft delete).

**3. Check if file exists for an entity**  
- **GetFileMetadataIfExists(ctx, entityType, id, kind)** returns `(exists bool, metadata *FileStoreMetadata, err)`.

**4. Build S3 key for your own logic**  
- **GetFileKey(ctx, entityName, entityId, fileType)** returns `"{entityName}/{entityId}/{fileType}"`.

### Request model (`requestModel.FileStoreMetadata`)

- `ID`, `IsActive`
- `FileName`, `FileSize`, `FileType`, `FileUrl`, `FileKey`, `FileBucket`
- `EntityId`, `EntityType` (e.g. order, expense)

Metadata is scoped by **channel** in the repository. Ensure context has channel set (e.g. from auth middleware).

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
| File Store | Metadata CRUD + Upload, key by entity | `file_store_service.go`, `file_store_repository.go`, request/response FileStoreMetadata |
| Expense | Balance + recalc on detail changes | `Expense.Balance`, `ExpenseDetailService`, `RecalculateAndUpdateBalance` |
| Measurement | JSON type for Values/OldValues | `entities/types/json.go`, request/response measurement models |

For dashboard field formulas, see **`docs/dashboard-response-fields.md`**.
