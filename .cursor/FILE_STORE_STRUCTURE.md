# File Store Structure

This document describes the **file store** flow: temp upload → confirm on entity Save/Update → storage in **entity_document** and **file_store_metadata**, and how to add files to an entity.

---

## 1. Overview

- **Temp upload:** Client uploads files to **temporary** S3 storage via `POST /file-store/temp` (or bulk). Backend returns `id`, `tempFileKey`, `fileName` per file.
- **Confirm on save/update:** When the client saves or updates an entity (e.g. Order, Expense) that has a `files` array, each item is a **ConfirmFile** (`id`, `fileKey`, `kind`, `description`). The backend:
  1. Inserts a row into **entity_document** (type, documentType, entityName, entityId, description).
  2. Moves the file from the temp S3 key to the final key `{entityName}/{entityId}/{kind}` (copy + delete source).
  3. Updates the **file_store_metadata** row (by `id`) with `entityId`, `entityType`, `fileKey`, and new presigned `fileUrl`.
- **Get:** When getting the entity, the backend loads **entity_document** rows for that entity and resolves each file URL from **file_store_metadata** (key `entityName/entityId/type`), and returns them as `files` on the response.

---

## 2. Tables

### 2.1 FileStoreMetadata (`stich.FileStoreMetadata`)

- **Purpose:** One row per file; stores S3 key, bucket, presigned URL, size, type, and link to entity.
- **Key fields:** `file_key`, `file_url`, `file_bucket`, `entity_id`, `entity_type`, `file_name`, `file_size`, `file_type`.
- **Key format:** Final key is `{entityType}/{entityId}/{kind}` (e.g. `OrderItem/123/design`, `Expense/456/receipt`). Temp key is `temp/{uuid}`.
- **Entity:** `internal/entities/file_store_metadata.go`
- **Request/response:** `internal/model/request/file_store_metadata.go`, `internal/model/response/` (FileStoreMetadata, TempFileUpload).

### 2.2 EntityDocuments (entity_document table)

- **Purpose:** One row per “document” attached to an entity; holds type, description, and reference (entityName, entityId). File URL is resolved via file_store_metadata using `entityName/entityId/type`.
- **Key fields:** `type`, `document_type`, `entity_name`, `entity_id`, `description`.
- **Entity:** `internal/entities/entity_documents.go`
- **Request/response:** `internal/model/request/entity_documents.go`, `internal/model/response/entity_documents.go` (EntityDocument with optional `document` = FileResponse).

---

## 3. Models

### 3.1 ConfirmFile (request – when confirming a temp file)

**Location:** `internal/model/models/file.go`

```go
type ConfirmFile struct {
    Id          uint   `json:"id,omitempty"`      // file_store_metadata id from UploadTemp response
    FileKey     string `json:"fileKey,omitempty"` // temp key from UploadTemp response
    Kind        string `json:"kind,omitempty"`    // e.g. "design", "receipt"
    Description string `json:"description,omitempty"`
}
```

Used in request bodies (e.g. `OrderItem.Files`, `ExpenseTracker.Files`).

### 3.2 TempFileUpload (response – after UploadTemp)

**Location:** `internal/model/response/` (TempFileUpload)

- `Id` – file_store_metadata id (use as `ConfirmFile.Id`).
- `FileKey` – temp S3 key (use as `ConfirmFile.FileKey`).
- `FileName` – original filename.

### 3.3 EntityDocument (response – file with URL when getting entity)

- `Id`, `Type`, `DocumentType`, `Description`, `EntityName`, `EntityId`
- `Document` – `*FileResponse` with `FileUrl`, `FileName` (resolved from file_store_metadata).

---

## 4. HTTP Endpoints (FileStore)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/file-store/temp` | Upload one file to temp storage. Body: `multipart/form-data`, field **`file`**. Returns `{ data: { id, tempFileKey, fileName } }`. |
| POST | `/file-store/temp/bulk` | Upload multiple files. Body: `multipart/form-data`, field **`files`** (multiple). Returns `{ data: [ { id, tempFileKey, fileName }, ... ] }`. |

Both require **BearerAuth**. Max form size 32 MB.

---

## 5. Service Layer

### 5.1 FileStoreService (`internal/service/file_store_service.go`)

- **UploadTemp(ctx, file)** – Upload to S3 `temp/{uuid}`, insert FileStoreMetadata (entityType=temp, entityId=0), return id, tempKey, fileName.
- **ConfirmTempUpload(ctx, existingTempKey, newKey)** – Copy S3 object from temp key to new key, delete source. Returns **newUrl** (presigned for new key).
- **UpdateEntityIdAndKey(ctx, id, entityId, entityType, fileKey, fileUrl)** – Update FileStoreMetadata row by `id` with entity id/type/key/url.
- **GetFileKey(ctx, entityName, entityId, fileType)** – Returns `"{entityName}/{entityId}/{fileType}"`.

### 5.2 EntityDocumentService (`internal/service/entity_document_service.go`)

- **SaveEntityDocument(ctx, req)** – Insert entity_document; if `Document` has content, also upload to S3 and create/update file_store_metadata (used for non–temp flows).
- **GetEntityDocumentsByEntity(ctx, entityId, entityName, typ)** – Get all entity_document rows for that entity (optional filter by `typ`). Resolves each file URL from **file_store_metadata** using `(entityDocument.EntityName, entityDocument.EntityId, entityDocument.Type)` so it works for both “EntityDocuments” and “OrderItem”/“Expense” style keys.

---

## 6. Entities That Have Files

### 6.1 OrderItem

- **Entity name:** `OrderItem` (`entities.Entity_OrderItem`).
- **Request:** `Order.OrderItems[].Files` = `[]ConfirmFile`.
- **Response:** `Order.OrderItems[].Files` = `[]EntityDocument` (with `document.fileUrl`, `document.fileName`).
- **Flow:** On **SaveOrder** / **UpdateOrder**, after persisting order/order items, for each order item with `Files`, for each ConfirmFile: save entity_document (type=kind, documentType=kind, entityName=OrderItem, entityId=orderItemId, description), then ConfirmTempUpload + UpdateEntityIdAndKey.
- **Get:** **Get** / **GetAll** Order: after mapping, **fillOrderItemFiles** calls `GetEntityDocumentsByEntity(ctx, orderItem.ID, Entity_OrderItem, "")` and sets `orderItem.Files`.
- **Code:** `internal/service/order_service.go` (confirmTempFileForOrderItem, fillOrderItemFiles).

### 6.2 Expense (ExpenseTracker)

- **Entity name:** `Expense` (`entities.Entity_Expense`).
- **Request:** `ExpenseTracker.Files` = `[]ConfirmFile`.
- **Response:** `ExpenseTracker.Files` = `[]EntityDocument`.
- **Flow:** On **SaveExpenseTracker** / **UpdateExpenseTracker**, after create/update, for each `expenseTracker.Files` call confirmTempFileForExpense (same pattern: entity_document + ConfirmTempUpload + UpdateEntityIdAndKey). Then RecalculateAndUpdateBalance.
- **Get:** **Get** / **GetAll** ExpenseTracker: after mapping, **fillExpenseFiles** calls `GetEntityDocumentsByEntity(ctx, expense.ID, Entity_Expense, "")` and sets `expense.Files`.
- **Code:** `internal/service/expense_tracker_service.go` (confirmTempFileForExpense, fillExpenseFiles).

---

## 7. Confirm Flow (Generic)

For any entity that supports files via temp upload:

1. **Persist the entity** (so you have its ID).
2. For each **ConfirmFile** in the request:
   - **Entity document:** Save one row: type = kind, documentType = kind, entityName = `<EntityName>`, entityId = `<persistedId>`, description = ConfirmFile.Description.
   - **Final key:** `newKey = GetFileKey(ctx, entityName, persistedId, confirmFile.Kind)`.
   - **ConfirmTempUpload(ctx, confirmFile.FileKey, newKey)** → get `newUrl`.
   - **UpdateEntityIdAndKey(ctx, confirmFile.Id, persistedId, entityName, newKey, newUrl)**.

Use the same **entityName** as in `entities.Entity_*` (e.g. `"OrderItem"`, `"Expense"`).

---

## 8. Get Flow (Filling Files on Response)

After mapping entity → response model:

- Call **GetEntityDocumentsByEntity(ctx, entityId, entities.Entity_<Name>, entities.EntityDocumentsType(""))**.
- Set **response.Files = result**.

File URLs are resolved inside **GetEntityDocumentsByEntity** using each document’s `EntityName` and `EntityId` (so OrderItem/Expense keys work without a separate “WithFileFromEntity” method).

---

## 9. Adding Files to a New Entity

1. **Entity name:** Add constant in `internal/entities/const.go` if new (e.g. `Entity_MyEntity EntityName = "MyEntity"`).
2. **Request model:** Add `Files []models.ConfirmFile` to the request struct for the resource that “has” files (e.g. MyEntityRequest).
3. **Response model:** Add `Files []EntityDocument` to the response struct.
4. **Service:** Inject **FileStoreService** and **EntityDocumentService**.
5. **Save/Create:** After creating the entity, loop over `request.Files` and for each call:
   - Save entity_document (entityName, entityId = new id, type/documentType = kind, description).
   - newKey = GetFileKey(ctx, entityName, id, f.Kind).
   - newUrl, xerr := ConfirmTempUpload(ctx, f.FileKey, newKey).
   - UpdateEntityIdAndKey(ctx, f.Id, id, entityName, newKey, newUrl).
6. **Update:** Same loop after update, using the entity’s id.
7. **Get/GetAll:** After mapping to response, call GetEntityDocumentsByEntity(ctx, id, Entity_MyEntity, "") and set response.Files.
8. **Wire:** Ensure the service provider receives FileStoreService and EntityDocumentService (wire will inject them).

---

## 10. File Locations Quick Reference

| What | Location |
|------|----------|
| ConfirmFile, FileUpload | `internal/model/models/file.go` |
| FileStoreMetadata entity | `internal/entities/file_store_metadata.go` |
| EntityDocuments entity | `internal/entities/entity_documents.go` |
| FileStoreRepository | `internal/repository/file_store_repository.go` |
| FileStoreService | `internal/service/file_store_service.go` |
| EntityDocumentService | `internal/service/entity_document_service.go` |
| FileStoreHandler (UploadTemp, UploadTempBulk) | `internal/handler/file_store_handler.go` |
| OrderItem files flow | `internal/service/order_service.go` |
| Expense files flow | `internal/service/expense_tracker_service.go` |

---

## 11. Key Conventions

- **S3 key:** Final key is always `{entityName}/{entityId}/{kind}`. Temp key is `temp/{uuid}`.
- **entity_document:** Stores metadata only; file URL comes from file_store_metadata lookup by (entityName, entityId, type).
- **ConfirmFile.Id** is the **file_store_metadata** id from the UploadTemp response; do not use entity_document id here.
- Use **GetEntityDocumentsByEntity** for both “classic” EntityDocuments-backed docs and OrderItem/Expense; file resolution uses each row’s EntityName/EntityId so it works for all.
