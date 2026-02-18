# Stitchfolio Backend - Cursor AI Rules

## Project Overview
This is a Go backend application using Gin framework with clean architecture pattern. The application follows a strict layered architecture with dependency injection using Google Wire.

## Architecture Layers (must follow in order)
1. **Router** → HTTP endpoint definitions
2. **Handler** → HTTP request/response handling
3. **Service** → Business logic
4. **Repository** → Data access
5. **Entities** → Database models (GORM)
6. **Mapper** → DTO transformations

## Code Generation Rules

### When Adding a New Resource/Feature
Always implement in this exact order:
1. Entity definition (`internal/entities/`)
2. Request model (`internal/model/request/`)
3. Response model (`internal/model/response/`)
4. Request mapper (`internal/mapper/mapper.go`)
5. Response mapper (`internal/mapper/response_mapper.go`)
6. Repository interface and implementation (`internal/repository/`)
7. Service interface and implementation (`internal/service/`)
8. Handler implementation (`internal/handler/`)
9. Wire DI registration (`internal/di/wire.go`)
10. Base handler registration (`internal/handler/base/base_handler.go`)
11. Router registration (`internal/router/router.go`)

### Mandatory Patterns

#### All Entities Must:
- Embed `*Model` for common fields
- Implement `TableNameForQuery()` returning `"\"stich\".\"TableName\" E"`
- Use GORM tags for relations: `gorm:"foreignKey:FieldId;constraint:OnDelete:CASCADE"`

#### All Request Models Must:
- Include `ID uint` and `IsActive bool` fields
- Use `json:"fieldName,omitempty"` tags

#### All Response Models Must:
- Include `ID uint` and `IsActive bool` fields
- Embed `AuditFields` for CreatedAt/UpdatedAt/CreatedBy/UpdatedBy
- Use `json:"fieldName,omitempty"` tags
- Include nested related entities

#### All Repositories Must:
- Embed `GormDAL`
- Implement interface with standard methods: Create, Update, Get, GetAll, Delete
- Use `WithDB(ctx)` for all database operations
- Apply `scopes.Channel()` and `scopes.IsActive()` in GetAll
- Use `scopes.ILike(search, "field1", "field2")` for search
- Apply `db.Paginate(ctx)` for list endpoints
- Return `*errs.XError` for errors
- Use `Preload()` for eager loading relations in Get method

#### All Services Must:
- Define interface first, then implementation struct
- Accept mapper dependencies: `mapper.Mapper` and `mapper.ResponseMapper`
- Use mapper to convert request → entity before repository call
- Use response mapper to convert entity → response after repository call
- Return response models, never entities
- Return `*errs.XError` for errors

#### All Handlers Must:
- Copy Gin context: `context := util.CopyContextFromGin(ctx)`
- Use `ctx.Bind()` for request body
- Use `ctx.Param()` for path parameters
- Use `ctx.Query()` for query parameters
- Use `util.EncloseWithSingleQuote()` for search queries
- Use `strconv.Atoi()` for ID conversion
- Return errors with `h.resp.DefaultFailureResponse(err).FormatAndSend()`
- Return success with `h.resp.SuccessResponse()` or `h.dataResp.DefaultSuccessResponse()`
- Include Swagger annotations: @Summary, @Description, @Tags, @Accept, @Success, @Failure, @Param, @Router

#### All Mappers Must:
- Implement both singular and plural methods for response mappers
- Check for nil at start of response mapper: `if e == nil { return nil, nil }`
- Map nested relations recursively
- Return error if mapping fails

### Provider Functions
All components must have a `Provide*` function for Wire:
```go
func ProvideCustomerHandler(svc service.CustomerService) *CustomerHandler
func ProvideCustomerService(repo repository.CustomerRepository, ...) CustomerService
func ProvideCustomerRepository(customDB GormDAL) CustomerRepository
```

### Naming Conventions
- Files: lowercase_with_underscores (e.g., `customer_handler.go`)
- Types: PascalCase with suffix (e.g., `CustomerHandler`, `CustomerService`)
- Interfaces: Same as struct but without suffix (e.g., `CustomerService` interface, `customerService` struct)
- Methods: PascalCase for exported, camelCase for internal
- Context parameter: always `ctx *context.Context`
- Gin context: always `ctx *gin.Context`
- Copied context: always `context`

### Standard CRUD Methods

#### Repository
- `Create(ctx *context.Context, entity *entities.T) *errs.XError`
- `Update(ctx *context.Context, entity *entities.T) *errs.XError`
- `Get(ctx *context.Context, id uint) (*entities.T, *errs.XError)`
- `GetAll(ctx *context.Context, search string) ([]entities.T, *errs.XError)`
- `Delete(ctx *context.Context, id uint) *errs.XError` (soft delete)

#### Service
- `SaveT(ctx *context.Context, req requestModel.T) *errs.XError`
- `UpdateT(ctx *context.Context, req requestModel.T, id uint) *errs.XError`
- `Get(ctx *context.Context, id uint) (*responseModel.T, *errs.XError)`
- `GetAll(ctx *context.Context, search string) ([]responseModel.T, *errs.XError)`
- `Delete(ctx *context.Context, id uint) *errs.XError`

#### Handler
- `SaveT(ctx *gin.Context)`
- `UpdateT(ctx *gin.Context)`
- `Get(ctx *gin.Context)`
- `GetAllTs(ctx *gin.Context)`
- `Delete(ctx *gin.Context)`

### HTTP Routes
- POST `/resource` - Create
- PUT `/resource/:id` - Update
- GET `/resource/:id` - Get single
- GET `/resource` - Get list with search
- DELETE `/resource/:id` - Delete
- GET `/resource/autocomplete` - Autocomplete (if needed)

**Important:** Place specific routes before generic ones (e.g., `/autocomplete` before `/:id`)

### Scopes Usage
Always apply these scopes in GetAll:
```go
Scopes(scopes.Channel(), scopes.IsActive())
Scopes(scopes.ILike(search, "field1", "field2"))
Scopes(db.Paginate(ctx))
```

### Error Handling
- Repository: `errs.NewXError(errs.DATABASE, "message", err)`
- Service: `errs.NewXError(errs.INVALID_REQUEST, "message", err)` or `errs.NewXError(errs.MAPPING_ERROR, "message", err)`
- Handler: Return appropriate HTTP status (400 for bad request, 500 for internal errors)

### HTTP Status Codes
- 200 OK - GET success
- 201 Created - POST success
- 202 Accepted - PUT success
- 400 Bad Request - Validation/binding errors
- 500 Internal Server Error - Service/repository errors

### Multi-tenancy & Audit
- All entities automatically include `ChannelId` via embedded `Model`
- All entities automatically include audit fields via embedded `Model`
- Channel filtering applied via `scopes.Channel()`
- Audit fields populated automatically via GORM hooks

### Import Paths
```go
"github.com/imkarthi24/sf-backend/internal/entities"
"github.com/imkarthi24/sf-backend/internal/handler"
"github.com/imkarthi24/sf-backend/internal/mapper"
requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
"github.com/imkarthi24/sf-backend/internal/repository"
"github.com/imkarthi24/sf-backend/internal/service"
"github.com/loop-kar/pixie/errs"
"github.com/loop-kar/pixie/response"
"github.com/loop-kar/pixie/util"
"github.com/loop-kar/pixie/db"
```

## Code Style

### Always:
- Use pointer receivers for methods
- Use `omitempty` in JSON tags
- Copy Gin context before using
- Apply Channel and IsActive scopes in queries
- Implement Provide functions for Wire
- Add Swagger annotations to handlers
- Run `wire` after modifying DI

### Never:
- Return entities from service layer (use response models)
- Skip Channel or IsActive scopes in GetAll
- Forget to add handler to BaseHandler
- Use hard delete (always soft delete with IsActive = false)
- Mix up Gin context (`ctx *gin.Context`) with standard context (`ctx *context.Context`)

## File Organization
```
internal/
├── entities/          # Database models
├── model/
│   ├── request/      # Request DTOs
│   └── response/     # Response DTOs
├── mapper/           # DTO transformations
├── repository/       # Data access
│   └── scopes/      # Query filters
├── service/          # Business logic
├── handler/          # HTTP handlers
│   └── base/        # Handler aggregation
├── router/           # Route definitions
└── di/              # Dependency injection
```

## When Modifying Existing Code

### Adding New Field to Entity:
1. Add to entity struct
2. Add to request model
3. Add to response model
4. Update request mapper
5. Update response mapper
6. Update any custom repository methods if needed

### Adding New Endpoint:
1. Add method to handler
2. Add route to router
3. Add method to service (if new)
4. Add method to repository (if new)

## Wire DI Pattern
After adding any new component:
1. Add provider to appropriate wire set in `internal/di/wire.go`
2. Add to BaseHandler if it's a handler
3. Run `cd internal/di && wire`
4. Verify no errors in `wire_gen.go`

## Documentation
- Add Swagger annotations to all handler methods
- Use clear, concise summary and description
- Specify all parameters (path, query, body)
- Document success and failure responses

## References
See `.cursor/ARCHITECTURE_GUIDE.md` for detailed patterns and examples.
See `.cursor/QUICK_REFERENCE.md` for implementation checklist.
