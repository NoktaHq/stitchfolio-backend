# Cursor AI Context Files

This directory contains comprehensive documentation to help Cursor AI understand and maintain the Stitchfolio backend codebase consistently.

## 📁 Files Overview

### 1. `ARCHITECTURE_GUIDE.md` (Comprehensive Reference)
**Purpose:** Complete architectural documentation with detailed examples  
**Use When:** 
- Learning the codebase architecture
- Understanding layer responsibilities
- Looking for detailed implementation patterns
- Need full code examples

**Contents:**
- Architecture overview and flow
- Layer-by-layer detailed guidelines
- Complete code examples for each layer
- Standard workflow walkthrough
- Common patterns and best practices

### 2. `QUICK_REFERENCE.md` (Implementation Checklist)
**Purpose:** Quick checklist for adding new features  
**Use When:**
- Adding a new resource/entity
- Need a step-by-step checklist
- Want code snippets for common patterns
- Testing new implementations

**Contents:**
- 12-step implementation checklist
- Standard method signatures
- Common code snippets
- File templates
- Testing checklist

### 3. `rules/backend-standards.md` (AI Rules)
**Purpose:** Cursor AI-specific coding standards and rules  
**Use When:**
- Cursor needs to understand project conventions
- Automatically enforcing coding standards
- Quick reference for mandatory patterns

**Contents:**
- Mandatory implementation patterns
- Naming conventions
- Code style rules
- File organization
- Import paths

### 4. `FILE_STORE_STRUCTURE.md` (File store: temp upload, confirm, entity_document)
**Purpose:** Single reference for file store flow, tables, and adding files to entities  
**Use When:**
- Implementing or changing file upload (temp → confirm on save/update)
- Working with OrderItem or ExpenseTracker files
- Adding files to a new entity

**Contents:**
- Overview (temp upload → confirm → entity_document + file_store_metadata)
- Tables: FileStoreMetadata, EntityDocuments
- Models: ConfirmFile, TempFileUpload, EntityDocument
- HTTP: POST /file-store/temp, /file-store/temp/bulk
- Entities with files: OrderItem, Expense (ExpenseTracker)
- Generic confirm flow and get flow
- Checklist: adding files to a new entity
- File locations quick reference

See also **`rules/file-store.mdc`** (Cursor rule when editing file-store/entity_document code).

### 5. `SERVICE_USAGE.md` (Branch / service usage)
**Purpose:** How to use services and APIs that changed on the current branch  
**Use When:**
- Using Task/Dashboard (status, tasksByStatus)
- Using File Store (temp upload, confirm, OrderItem/Expense files — see FILE_STORE_STRUCTURE.md)
- Using Expense/ExpenseDetail with Balance
- Working with Measurement JSON types

**Contents:**
- Task status and dashboard API
- File Store (temp upload → confirm; entity_document + file_store_metadata; see FILE_STORE_STRUCTURE.md)
- Expense Balance and RecalculateAndUpdateBalance
- Measurement/MeasurementHistory JSON type
- Branch-specific notes and quick reference

## 🚀 Quick Start for AI

When asked to add a new feature:
1. **First**, review `QUICK_REFERENCE.md` for the implementation checklist
2. **Then**, refer to `ARCHITECTURE_GUIDE.md` for detailed patterns
3. **Always**, follow the rules in `rules/backend-standards.md`

## 📋 Workflow Example

**Task:** "Add a Product feature"

```
Step 1: Open QUICK_REFERENCE.md
        → Follow the 12-step checklist

Step 2: For each step, refer to ARCHITECTURE_GUIDE.md
        → See detailed examples for that layer

Step 3: Implement following rules/backend-standards.md
        → Ensure naming, patterns, and structure comply

Step 4: Test using the testing checklist
        → Verify all functionality works
```

## 🎯 Document Usage Guide

### For New Features (Full CRUD)
Use in order:
1. `QUICK_REFERENCE.md` - Get the checklist
2. `ARCHITECTURE_GUIDE.md` - Get detailed examples
3. `rules/backend-standards.md` - Verify compliance

### For Modifying Existing Code
1. Identify the layer being modified
2. Check `ARCHITECTURE_GUIDE.md` for that layer's patterns
3. Verify changes follow `rules/backend-standards.md`

### For Understanding the Codebase
1. Start with `ARCHITECTURE_GUIDE.md` - Overview section
2. Read through layer-by-layer guidelines
3. Reference `rules/backend-standards.md` for specifics

## 📊 Layer Flow Quick Reference

```
┌──────────────────────────────────────────────────┐
│                    HTTP Request                  │
└────────────────────┬─────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────────┐
│  Router (internal/router/router.go)               │
│  • Route definitions                               │
│  • Middleware application                          │
│  • Group organization                              │
└────────────────────┬───────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────────┐
│  Handler (internal/handler/)                       │
│  • Request binding & validation                    │
│  • Extract params (path, query, body)              │
│  • Call service layer                              │
│  • Format & send response                          │
└────────────────────┬───────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────────┐
│  Service (internal/service/)                       │
│  • Business logic orchestration                    │
│  • Request → Entity (via Mapper)                   │
│  • Call repository                                 │
│  • Entity → Response (via ResponseMapper)          │
└────────────────────┬───────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────────┐
│  Repository (internal/repository/)                 │
│  • Database operations (GORM)                      │
│  • Apply scopes (Channel, IsActive, Search)        │
│  • Handle preloading                               │
│  • Return entities                                 │
└────────────────────┬───────────────────────────────┘
                     ↓
┌────────────────────────────────────────────────────┐
│  Database (PostgreSQL)                             │
└────────────────────────────────────────────────────┘
```

## 🔄 Data Flow

```
Request Body → Request Model → [Mapper] → Entity → Repository → Database
                                                                     ↓
Response Body ← Response Model ← [ResponseMapper] ← Entity ← Repository
```

## 🛠️ Cross-Cutting Concerns

### Dependency Injection
- **Location:** `internal/di/`
- **Tool:** Google Wire
- **Pattern:** Provider functions in wire sets
- **Reference:** See "Dependency Injection" section in `ARCHITECTURE_GUIDE.md`

### Models (DTOs)
- **Request:** `internal/model/request/`
- **Response:** `internal/model/response/`
- **Purpose:** Separate API contracts from database models

### Mappers
- **Location:** `internal/mapper/`
- **Types:** Request Mapper (input), Response Mapper (output)
- **Pattern:** Interface + implementation struct

### Entities
- **Location:** `internal/entities/`
- **Purpose:** Database table definitions (GORM models)
- **Base:** All embed `*Model` for common fields

### Scopes
- **Location:** `internal/repository/scopes/`
- **Purpose:** Reusable query filters
- **Common:** Channel(), IsActive(), ILike(), Paginate()

## 📝 Key Patterns

### Multi-Tenancy
- All data filtered by `channel_id`
- Automatic via `scopes.Channel()`
- System admin (channel_id=0) sees all data

### Soft Delete
- Sets `is_active = false`
- Filtered via `scopes.IsActive()`
- Never hard delete records

### Audit Trail
- CreatedAt, UpdatedAt, CreatedBy, UpdatedBy
- Automatic via GORM hooks
- Populated from context

### Error Handling
- Custom error type: `*errs.XError`
- Consistent across all layers
- HTTP status codes standardized

## 🎨 Naming Conventions Summary

| Item | Convention | Example |
|------|-----------|---------|
| Files | lowercase_underscore | `customer_handler.go` |
| Types | PascalCase + Suffix | `CustomerHandler` |
| Interfaces | Same as struct (no suffix) | `CustomerService` |
| Methods | PascalCase (exported) | `SaveCustomer` |
| Variables | camelCase | `customerRepo` |
| Providers | Provide + Name | `ProvideCustomerHandler` |

## 🚦 Implementation Order (Critical!)

When adding a new feature, implement in this exact order:

1. Entity
2. Request Model
3. Response Model
4. Request Mapper
5. Response Mapper
6. Repository
7. Service
8. Handler
9. Wire DI
10. Base Handler
11. Router
12. Generate Wire (`cd internal/di && wire`)

**Why?** Each layer depends on the previous ones.

## 🔍 Finding Information

| Need | Document | Section |
|------|----------|---------|
| File store (temp upload, confirm, entity_document) | FILE_STORE_STRUCTURE.md | Full doc |
| Task/Dashboard/FileStore/Expense usage | SERVICE_USAGE.md | Per-service sections |
| Complete entity example | ARCHITECTURE_GUIDE.md | Entities Layer |
| Repository patterns | ARCHITECTURE_GUIDE.md | Repository Layer |
| Service patterns | ARCHITECTURE_GUIDE.md | Service Layer |
| Handler patterns | ARCHITECTURE_GUIDE.md | Handler Layer |
| Implementation checklist | QUICK_REFERENCE.md | Checklist |
| Code snippets | QUICK_REFERENCE.md | Code Snippets |
| Mandatory rules | rules/backend-standards.md | All sections |
| Scopes usage | ARCHITECTURE_GUIDE.md | Scopes subsection |
| Wire DI setup | ARCHITECTURE_GUIDE.md | Dependency Injection |
| Dashboard field formulas | docs/dashboard-response-fields.md | Task/Order/Stats tables |

## ⚠️ Common Mistakes to Avoid

1. ❌ Not following implementation order
2. ❌ Returning entities from service (should return response models)
3. ❌ Forgetting to run `wire` after DI changes
4. ❌ Not applying Channel() and IsActive() scopes in GetAll
5. ❌ Hard deleting instead of soft delete
6. ❌ Not copying Gin context in handlers
7. ❌ Placing generic routes before specific ones
8. ❌ Forgetting to add handler to BaseHandler
9. ❌ Not implementing plural response mapper
10. ❌ Missing nil check in response mapper

## 🧪 Verification Steps

After implementing a new feature:

```bash
# 1. Generate Wire DI
cd internal/di && wire

# 2. Run the application
go run main.go

# 3. Test endpoints
# - POST /api/v1/resource
# - GET /api/v1/resource/:id
# - GET /api/v1/resource?search=query
# - PUT /api/v1/resource/:id
# - DELETE /api/v1/resource/:id

# 4. Check Swagger docs
# Visit: http://localhost:8080/api/v1/swagger/index.html
```

## 📚 Additional Resources

- **Wire Documentation:** https://github.com/google/wire
- **GORM Documentation:** https://gorm.io/docs/
- **Gin Documentation:** https://gin-gonic.com/docs/

## 🔄 Updating These Docs

When architectural patterns change:
1. Update `ARCHITECTURE_GUIDE.md` with detailed explanation
2. Update `QUICK_REFERENCE.md` with checklist/snippet changes
3. Update `rules/backend-standards.md` with new mandatory rules
4. Update this README if structure changes

---

**Last Updated:** 2026-02-15  
**Version:** 1.0  
**Maintained By:** Development Team
