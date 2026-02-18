# Quick Reference - Adding a New Feature

Use this checklist when adding a new resource/feature to the Stitchfolio backend.

## Checklist for Adding New Feature (e.g., "Product")

- [ ] **1. Entity** (`internal/entities/product.go`)
  - [ ] Define struct with `*Model` embedded
  - [ ] Add fields with JSON tags
  - [ ] Add GORM relation tags
  - [ ] Implement `TableNameForQuery()` method

- [ ] **2. Request Model** (`internal/model/request/product.go`)
  - [ ] Define struct with `ID` and `IsActive`
  - [ ] Add fields with `json:"fieldName,omitempty"`
  - [ ] Include only input-relevant fields

- [ ] **3. Response Model** (`internal/model/response/product.go`)
  - [ ] Define struct with `ID` and `IsActive`
  - [ ] Add fields with `json:"fieldName,omitempty"`
  - [ ] Embed `AuditFields`
  - [ ] Include related entities for nesting
  - [ ] Add autocomplete struct if needed

- [ ] **4. Request Mapper** (`internal/mapper/mapper.go`)
  - [ ] Add method to `Mapper` interface
  - [ ] Implement mapper method (requestModel → entity)
  - [ ] Handle date/time conversions if needed
  - [ ] Handle nested objects if needed

- [ ] **5. Response Mapper** (`internal/mapper/response_mapper.go`)
  - [ ] Add methods to `ResponseMapper` interface (singular + plural)
  - [ ] Implement singular mapper (entity → responseModel)
  - [ ] Check for nil at start
  - [ ] Map nested relations recursively
  - [ ] Implement plural mapper using loop

- [ ] **6. Repository** (`internal/repository/product_repository.go`)
  - [ ] Define interface with CRUD methods
  - [ ] Define struct embedding `GormDAL`
  - [ ] Implement `ProvideProductRepository` function
  - [ ] Implement `Create` method
  - [ ] Implement `Update` method
  - [ ] Implement `Get` method with Preloads
  - [ ] Implement `GetAll` method with Scopes
  - [ ] Implement `Delete` method (soft delete)
  - [ ] Add custom methods as needed

- [ ] **7. Service** (`internal/service/product_service.go`)
  - [ ] Define interface with business methods
  - [ ] Define struct with repo, mapper dependencies
  - [ ] Implement `ProvideProductService` function
  - [ ] Implement `SaveProduct` (request → entity → repo)
  - [ ] Implement `UpdateProduct` (request → entity → repo)
  - [ ] Implement `Get` (repo → entity → response)
  - [ ] Implement `GetAll` (repo → entities → responses)
  - [ ] Implement `Delete` (repo)
  - [ ] Add business logic methods as needed

- [ ] **8. Handler** (`internal/handler/product_handler.go`)
  - [ ] Define struct with service dependency
  - [ ] Implement `ProvideProductHandler` function
  - [ ] Implement `SaveProduct` with swagger annotations
  - [ ] Implement `UpdateProduct` with swagger annotations
  - [ ] Implement `Get` with swagger annotations
  - [ ] Implement `GetAllProducts` with swagger annotations
  - [ ] Implement `Delete` with swagger annotations
  - [ ] Add custom handlers as needed

- [ ] **9. Wire DI** (`internal/di/wire.go`)
  - [ ] Add `handler.ProvideProductHandler` to `handlerSet`
  - [ ] Add `service.ProvideProductService` to `svcSet`
  - [ ] Add `repository.ProvideProductRepository` to `repoSet`

- [ ] **10. Base Handler** (`internal/handler/base/base_handler.go`)
  - [ ] Add `ProductHandler *handler.ProductHandler` to struct
  - [ ] Add parameter to `ProvideBaseHandler` function
  - [ ] Initialize in return statement

- [ ] **11. Router** (`internal/router/router.go`)
  - [ ] Create endpoint group with middleware
  - [ ] Add POST route for create
  - [ ] Add PUT route for update
  - [ ] Add GET routes (single, list, custom)
  - [ ] Add DELETE route

- [ ] **12. Generate Wire** (Terminal)
  - [ ] Run `cd internal/di && wire`
  - [ ] Verify `wire_gen.go` updated without errors

## Standard Method Signatures

### Repository
```go
Create(*context.Context, *entities.Product) *errs.XError
Update(*context.Context, *entities.Product) *errs.XError
Get(*context.Context, uint) (*entities.Product, *errs.XError)
GetAll(*context.Context, string) ([]entities.Product, *errs.XError)
Delete(*context.Context, uint) *errs.XError
```

### Service
```go
SaveProduct(*context.Context, requestModel.Product) *errs.XError
UpdateProduct(*context.Context, requestModel.Product, uint) *errs.XError
Get(*context.Context, uint) (*responseModel.Product, *errs.XError)
GetAll(*context.Context, string) ([]responseModel.Product, *errs.XError)
Delete(*context.Context, uint) *errs.XError
```

### Handler
```go
SaveProduct(ctx *gin.Context)
UpdateProduct(ctx *gin.Context)
Get(ctx *gin.Context)
GetAllProducts(ctx *gin.Context)
Delete(ctx *gin.Context)
```

## Common Code Snippets

### Entity with Relations
```go
type Product struct {
    *Model `mapstructure:",squash"`
    Name   string `json:"name"`
    
    CategoryId uint     `json:"categoryId"`
    Category   Category `gorm:"foreignKey:CategoryId"`
}

func (Product) TableNameForQuery() string {
    return "\"stich\".\"Products\" E"
}
```

### Repository GetAll Pattern
```go
var products []entities.Product
res := pr.WithDB(ctx).Table(entities.Product{}.TableNameForQuery()).
    Scopes(scopes.Channel(), scopes.IsActive()).
    Scopes(scopes.ILike(search, "name", "description")).
    Scopes(db.Paginate(ctx)).
    Find(&products)
```

### Service Save Pattern
```go
dbProduct, err := svc.mapper.Product(product)
if err != nil {
    return errs.NewXError(errs.INVALID_REQUEST, "Unable to save product", err)
}

errr := svc.productRepo.Create(ctx, dbProduct)
if errr != nil {
    return errr
}
return nil
```

### Handler Pattern
```go
context := util.CopyContextFromGin(ctx)
var product requesModel.Product
err := ctx.Bind(&product)
if err != nil {
    x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
    h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
    return
}

errr := h.productSvc.SaveProduct(&context, product)
if errr != nil {
    h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
    return
}

h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
```

### Router Pattern
```go
productEndpoints := appRouter.Group("product", router.VerifyJWT(srvConfig.JwtSecretKey))
{
    productEndpoints.POST("", handler.ProductHandler.SaveProduct)
    productEndpoints.PUT(":id", handler.ProductHandler.UpdateProduct)
    productEndpoints.GET("autocomplete", handler.ProductHandler.AutocompleteProduct)
    productEndpoints.GET(":id", handler.ProductHandler.Get)
    productEndpoints.GET("", handler.ProductHandler.GetAllProducts)
    productEndpoints.DELETE(":id", handler.ProductHandler.Delete)
}
```

## File Templates

### Minimal Entity
```go
package entities

type Product struct {
    *Model `mapstructure:",squash"`
    Name   string `json:"name"`
}

func (Product) TableNameForQuery() string {
    return "\"stich\".\"Products\" E"
}
```

### Minimal Request Model
```go
package requestModel

type Product struct {
    ID       uint   `json:"id,omitempty"`
    IsActive bool   `json:"isActive,omitempty"`
    Name     string `json:"name,omitempty"`
}
```

### Minimal Response Model
```go
package responseModel

type Product struct {
    ID       uint   `json:"id,omitempty"`
    IsActive bool   `json:"isActive,omitempty"`
    Name     string `json:"name,omitempty"`
    AuditFields
}
```

## Common Scopes

```go
// Always apply for multi-tenant filtering
scopes.Channel()

// Always apply for soft delete filtering
scopes.IsActive()

// Search across multiple fields
scopes.ILike(search, "field1", "field2", "field3")

// Pagination
db.Paginate(ctx)

// Audit info (created_by, updated_by names)
scopes.WithAuditInfo()
```

## HTTP Status Codes

- **200 OK** - GET success (single or list)
- **201 Created** - POST success
- **202 Accepted** - PUT success
- **400 Bad Request** - Validation error, malformed request
- **500 Internal Server Error** - Service/repository error

## Common Mistakes to Avoid

- ❌ Forgetting to add `omitempty` to JSON tags
- ❌ Not copying Gin context: use `util.CopyContextFromGin(ctx)`
- ❌ Returning entity from service (should return response model)
- ❌ Forgetting to add handler to BaseHandler
- ❌ Not running `wire` after adding to wire sets
- ❌ Placing specific routes after generic ones (`:id` should be last)
- ❌ Not applying `Channel()` and `IsActive()` scopes in GetAll
- ❌ Forgetting nil check in response mapper
- ❌ Not implementing both singular and plural response mappers

## Testing Checklist

After implementation, verify:

- [ ] Wire generation succeeds (`cd internal/di && wire`)
- [ ] Application starts without errors
- [ ] POST /resource creates new record
- [ ] GET /resource/:id retrieves record with relations
- [ ] GET /resource lists records with search
- [ ] PUT /resource/:id updates record
- [ ] DELETE /resource/:id soft deletes record
- [ ] Swagger docs generated correctly
- [ ] Channel filtering works (data isolated by channel)
- [ ] Audit fields populated correctly

---

**Quick Tip:** Keep this checklist open while implementing. Work through it sequentially for consistency.
