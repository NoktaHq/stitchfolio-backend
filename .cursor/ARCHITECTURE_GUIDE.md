# Stitchfolio Backend - Architecture Guide

This guide documents the standard workflow and patterns used throughout the Stitchfolio backend application. Follow these patterns when adding new features or modifying existing ones.

## Table of Contents
- [Architecture Overview](#architecture-overview)
- [Layer Structure](#layer-structure)
- [Standard Workflow for Adding a New Feature](#standard-workflow-for-adding-a-new-feature)
- [Layer-by-Layer Guidelines](#layer-by-layer-guidelines)
- [Naming Conventions](#naming-conventions)
- [Common Patterns](#common-patterns)

## Architecture Overview

The application follows a clean architecture pattern with clear separation of concerns across 6 main layers:

```
Router → Handler → Service → Repository → Database
              ↓         ↓          ↓
           Models   Mapper    Entities
```

**Flow Direction:**
1. **Router** (`internal/router/router.go`) - Defines HTTP endpoints and routes them to handlers
2. **Handler** (`internal/handler/`) - HTTP layer, handles request/response, validation
3. **Service** (`internal/service/`) - Business logic layer
4. **Repository** (`internal/repository/`) - Data access layer
5. **Entities** (`internal/entities/`) - Database table definitions (GORM models)
6. **Mapper** (`internal/mapper/`) - Transforms between request/response models and entities

**Cross-Cutting Concerns:**
- **DI** (`internal/di/`) - Dependency injection using Google Wire
- **Models** (`internal/model/`) - Request and response DTOs
- **Scopes** (`internal/repository/scopes/`) - Reusable query filters

## Layer Structure

### 1. Entities Layer (`internal/entities/`)
Database table definitions using GORM.

**Structure:**
```go
type Customer struct {
    *Model `mapstructure:",squash"`
    
    // Fields
    FirstName      string `json:"firstName"`
    LastName       string `json:"lastName"`
    Email          string `json:"email"`
    PhoneNumber    string `json:"phoneNumber"`
    WhatsappNumber string `json:"whatsappNumber"`
    Address        string `json:"address"`
    
    // Relations
    Persons   []Person  `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE"`
    Enquiries []Enquiry `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE"`
    Orders    []Order   `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE"`
}

func (Customer) TableNameForQuery() string {
    return "\"stich\".\"Customers\" E"
}
```

**Key Points:**
- Embed `*Model` for common fields (ID, timestamps, audit fields, channel_id)
- Use GORM tags for relations and constraints
- Implement `TableNameForQuery()` for custom table names with schema prefix

### 2. Models Layer (`internal/model/`)

#### Request Models (`internal/model/request/`)
DTOs for incoming HTTP requests.

```go
type Customer struct {
    ID       uint   `json:"id,omitempty"`
    IsActive bool   `json:"isActive,omitempty"`
    
    FirstName      string `json:"firstName,omitempty"`
    LastName       string `json:"lastName,omitempty"`
    Email          string `json:"email,omitempty"`
    PhoneNumber    string `json:"phoneNumber,omitempty"`
    WhatsappNumber string `json:"whatsappNumber,omitempty"`
    Address        string `json:"address,omitempty"`
    Age            int    `json:"age,omitempty"`
    Gender         string `json:"gender,omitempty"`
}
```

#### Response Models (`internal/model/response/`)
DTOs for outgoing HTTP responses.

```go
type Customer struct {
    ID       uint `json:"id,omitempty"`
    IsActive bool `json:"isActive,omitempty"`
    
    FirstName      string `json:"firstName,omitempty"`
    LastName       string `json:"lastName,omitempty"`
    Email          string `json:"email,omitempty"`
    PhoneNumber    string `json:"phoneNumber,omitempty"`
    WhatsappNumber string `json:"whatsappNumber,omitempty"`
    Address        string `json:"address,omitempty"`
    
    AuditFields
    
    Persons   []Person  `json:"persons,omitempty"`
    Enquiries []Enquiry `json:"enquiries,omitempty"`
    Orders    []Order   `json:"orders,omitempty"`
}
```

**Key Points:**
- Response models include `AuditFields` (CreatedAt, UpdatedAt, CreatedBy, UpdatedBy)
- Response models include related entities for nested responses
- Use `omitempty` for optional fields

### 3. Mapper Layer (`internal/mapper/`)

Two mapper types:

#### Request Mapper (`mapper.go`)
Converts request models to entities.

```go
type Mapper interface {
    Customer(e requestModel.Customer) (*entities.Customer, error)
    // ... other methods
}

func (m *mapper) Customer(e requestModel.Customer) (*entities.Customer, error) {
    return &entities.Customer{
        Model:          &entities.Model{ID: e.ID, IsActive: e.IsActive},
        FirstName:      e.FirstName,
        LastName:       e.LastName,
        Email:          e.Email,
        PhoneNumber:    e.PhoneNumber,
        WhatsappNumber: e.WhatsappNumber,
        Address:        e.Address,
    }, nil
}
```

#### Response Mapper (`response_mapper.go`)
Converts entities to response models.

```go
type ResponseMapper interface {
    Customer(e *entities.Customer) (*responseModel.Customer, error)
    Customers(items []entities.Customer) ([]responseModel.Customer, error)
    // ... other methods
}

func (m *responseMapper) Customer(e *entities.Customer) (*responseModel.Customer, error) {
    if e == nil {
        return nil, nil
    }
    
    persons, err := m.Persons(e.Persons)
    if err != nil {
        return nil, err
    }
    
    enquiries, err := m.Enquiries(e.Enquiries)
    if err != nil {
        return nil, err
    }
    
    return &responseModel.Customer{
        ID:             e.ID,
        IsActive:       e.IsActive,
        FirstName:      e.FirstName,
        LastName:       e.LastName,
        Email:          e.Email,
        PhoneNumber:    e.PhoneNumber,
        WhatsappNumber: e.WhatsappNumber,
        Address:        e.Address,
        Persons:        persons,
        Enquiries:      enquiries,
        Orders:         orders,
    }, nil
}
```

**Key Points:**
- Always implement both singular and plural mapper methods
- Handle nested relations recursively
- Check for nil before mapping
- Return error if mapping fails

### 4. Repository Layer (`internal/repository/`)

Data access layer using GORM.

```go
type CustomerRepository interface {
    Create(*context.Context, *entities.Customer) *errs.XError
    Update(*context.Context, *entities.Customer) *errs.XError
    Get(*context.Context, uint) (*entities.Customer, *errs.XError)
    GetAll(*context.Context, string) ([]entities.Customer, *errs.XError)
    Delete(*context.Context, uint) *errs.XError
    // Additional domain-specific methods
}

type customerRepository struct {
    GormDAL
}

func ProvideCustomerRepository(customDB GormDAL) CustomerRepository {
    return &customerRepository{GormDAL: customDB}
}
```

**Standard CRUD Methods:**

```go
// Create
func (cr *customerRepository) Create(ctx *context.Context, customer *entities.Customer) *errs.XError {
    res := cr.WithDB(ctx).Create(&customer)
    if res.Error != nil {
        return errs.NewXError(errs.DATABASE, "Unable to save customer", res.Error)
    }
    return nil
}

// Update
func (cr *customerRepository) Update(ctx *context.Context, customer *entities.Customer) *errs.XError {
    return cr.GormDAL.Update(ctx, *customer)
}

// Get (with preloads)
func (cr *customerRepository) Get(ctx *context.Context, id uint) (*entities.Customer, *errs.XError) {
    customer := entities.Customer{}
    res := cr.WithDB(ctx).
        Preload("Persons").
        Preload("Persons.Measurements").
        Preload("Persons.Measurements.DressType").
        Preload("Enquiries").
        Preload("Orders").
        Find(&customer, id)
    if res.Error != nil {
        return nil, errs.NewXError(errs.DATABASE, "Unable to find customer", res.Error)
    }
    return &customer, nil
}

// GetAll (with scopes)
func (cr *customerRepository) GetAll(ctx *context.Context, search string) ([]entities.Customer, *errs.XError) {
    var customers []entities.Customer
    res := cr.WithDB(ctx).Table(entities.Customer{}.TableNameForQuery()).
        Scopes(scopes.Channel(), scopes.IsActive()).
        Scopes(scopes.ILike(search, "first_name", "last_name", "email", "phone_number")).
        Scopes(db.Paginate(ctx)).
        Find(&customers)
    if res.Error != nil {
        return nil, errs.NewXError(errs.DATABASE, "Unable to find customers", res.Error)
    }
    return customers, nil
}

// Delete (soft delete)
func (cr *customerRepository) Delete(ctx *context.Context, id uint) *errs.XError {
    customer := &entities.Customer{Model: &entities.Model{ID: id, IsActive: false}}
    err := cr.GormDAL.Delete(ctx, customer)
    if err != nil {
        return err
    }
    return nil
}
```

**Key Points:**
- Embed `GormDAL` for base database operations
- Use `WithDB(ctx)` to get context-aware DB instance
- Apply scopes for common filters (Channel, IsActive, search)
- Use `Preload` for eager loading relations
- Return `*errs.XError` for consistent error handling
- Delete is soft delete (sets IsActive = false)

#### Scopes (`internal/repository/scopes/`)

Reusable query filters.

**Common Scopes:**

```go
// IsActive - Filter by is_active = true
func IsActive(params ...string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        stmt := &gorm.Statement{DB: db}
        stmt.Parse(db.Statement.Model)
        tableName := stmt.Schema.Table
        return db.Where(fmt.Sprintf("%s.is_active", tableName), true)
    }
}

// Channel - Filter by channel_id from context
func Channel(params ...string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        var channelId uint
        if id, ok := db.Get(constants.CHANNEL_ID); ok {
            channelId = id.(uint)
        }
        
        if channelId == 0 {
            return db // System Admin - access all data
        }
        
        stmt := &gorm.Statement{DB: db}
        stmt.Parse(db.Statement.Model)
        tableName := stmt.Schema.Table
        
        return db.Where(fmt.Sprintf("%s.channel_id", tableName), channelId)
    }
}

// ILike - Case-insensitive search across multiple fields
func ILike(query string, params ...string) func(db *gorm.DB) *gorm.DB {
    if len(params) == 0 || util.IsNilOrEmptyString(&query) {
        return func(db *gorm.DB) *gorm.DB { return db }
    }
    
    return func(db *gorm.DB) *gorm.DB {
        whereClause := ""
        query = util.EncloseWithPercentageOperator(query)
        funk.ForEach(params, func(param string) {
            whereClause = whereClause + fmt.Sprintf(`%s ILIKE %s`, param, query) + OR
        })
        whereClause = strings.Trim(whereClause, OR)
        return db.Where(whereClause)
    }
}
```

**Usage in Repository:**
```go
res := cr.WithDB(ctx).Table(entities.Customer{}.TableNameForQuery()).
    Scopes(scopes.Channel(), scopes.IsActive()).
    Scopes(scopes.ILike(search, "first_name", "last_name", "email")).
    Find(&customers)
```

### 5. Service Layer (`internal/service/`)

Business logic layer.

```go
type CustomerService interface {
    SaveCustomer(*context.Context, requestModel.Customer) *errs.XError
    UpdateCustomer(*context.Context, requestModel.Customer, uint) *errs.XError
    Get(*context.Context, uint) (*responseModel.Customer, *errs.XError)
    GetAll(*context.Context, string) ([]responseModel.Customer, *errs.XError)
    Delete(*context.Context, uint) *errs.XError
}

type customerService struct {
    customerRepo repository.CustomerRepository
    personRepo   repository.PersonRepository // Additional repos as needed
    mapper       mapper.Mapper
    respMapper   mapper.ResponseMapper
}

func ProvideCustomerService(
    repo repository.CustomerRepository,
    personRepo repository.PersonRepository,
    mapper mapper.Mapper,
    respMapper mapper.ResponseMapper,
) CustomerService {
    return customerService{
        customerRepo: repo,
        personRepo:   personRepo,
        mapper:       mapper,
        respMapper:   respMapper,
    }
}
```

**Standard Service Methods:**

```go
// Save
func (svc customerService) SaveCustomer(ctx *context.Context, customer requestModel.Customer) *errs.XError {
    dbCustomer, err := svc.mapper.Customer(customer)
    if err != nil {
        return errs.NewXError(errs.INVALID_REQUEST, "Unable to save customer", err)
    }
    
    errr := svc.customerRepo.Create(ctx, dbCustomer)
    if errr != nil {
        return errr
    }
    
    // Additional business logic (e.g., create related person)
    person := &entities.Person{
        Model:      &entities.Model{IsActive: true},
        FirstName:  customer.FirstName,
        LastName:   customer.LastName,
        CustomerId: dbCustomer.ID,
        Age:        &customer.Age,
        Gender:     entities.Gender(customer.Gender),
    }
    
    errr = svc.personRepo.Create(ctx, person)
    if errr != nil {
        return errr
    }
    
    return nil
}

// Update
func (svc customerService) UpdateCustomer(ctx *context.Context, customer requestModel.Customer, id uint) *errs.XError {
    dbCustomer, err := svc.mapper.Customer(customer)
    if err != nil {
        return errs.NewXError(errs.INVALID_REQUEST, "Unable to update customer", err)
    }
    
    dbCustomer.ID = id
    errr := svc.customerRepo.Update(ctx, dbCustomer)
    if errr != nil {
        return errr
    }
    return nil
}

// Get
func (svc customerService) Get(ctx *context.Context, id uint) (*responseModel.Customer, *errs.XError) {
    customer, err := svc.customerRepo.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    
    mappedCustomer, mapErr := svc.respMapper.Customer(customer)
    if mapErr != nil {
        return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Customer data", mapErr)
    }
    
    return mappedCustomer, nil
}

// GetAll
func (svc customerService) GetAll(ctx *context.Context, search string) ([]responseModel.Customer, *errs.XError) {
    customers, err := svc.customerRepo.GetAll(ctx, search)
    if err != nil {
        return nil, err
    }
    
    mappedCustomers, mapErr := svc.respMapper.Customers(customers)
    if mapErr != nil {
        return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Customer data", mapErr)
    }
    
    return mappedCustomers, nil
}

// Delete
func (svc customerService) Delete(ctx *context.Context, id uint) *errs.XError {
    err := svc.customerRepo.Delete(ctx, id)
    if err != nil {
        return err
    }
    return nil
}
```

**Key Points:**
- Service orchestrates business logic
- Uses mapper to convert request → entity
- Uses repository for data access
- Uses response mapper to convert entity → response
- Handles cross-repository transactions
- Returns response models, not entities

### 6. Handler Layer (`internal/handler/`)

HTTP request/response handling.

```go
type CustomerHandler struct {
    customerSvc service.CustomerService
    resp        response.Response
    dataResp    response.DataResponse
}

func ProvideCustomerHandler(svc service.CustomerService) *CustomerHandler {
    return &CustomerHandler{customerSvc: svc}
}
```

**Standard Handler Methods:**

```go
// Save
// @Summary     Save Customer
// @Description Saves an instance of Customer
// @Tags        Customer
// @Accept      json
// @Success     201         {object} response.Response
// @Failure     400         {object} response.Response
// @Failure     501         {object} response.Response
// @Param       customer    body     requestModel.Customer true "customer"
// @Router      /customer [post]
func (h CustomerHandler) SaveCustomer(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    var customer requesModel.Customer
    err := ctx.Bind(&customer)
    if err != nil {
        x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
        h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    errr := h.customerSvc.SaveCustomer(&context, customer)
    if errr != nil {
        h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
        return
    }
    
    h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

// Update
// @Summary     Update Customer
// @Description Updates an instance of Customer
// @Tags        Customer
// @Accept      json
// @Success     201         {object} response.Response
// @Failure     400         {object} response.Response
// @Param       customer    body     requestModel.Customer true "customer"
// @Param       id          path     int                   true "Customer id"
// @Router      /customer/{id} [put]
func (h CustomerHandler) UpdateCustomer(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    var customer requesModel.Customer
    err := ctx.Bind(&customer)
    if err != nil {
        x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
        h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    id, _ := strconv.Atoi(ctx.Param("id"))
    errr := h.customerSvc.UpdateCustomer(&context, customer, uint(id))
    if errr != nil {
        h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
        return
    }
    
    h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// Get
// @Summary     Get a specific Customer
// @Description Get an instance of Customer
// @Tags        Customer
// @Accept      json
// @Success     200 {object} responseModel.Customer
// @Failure     400 {object} response.DataResponse
// @Param       id  path     int true "Customer id"
// @Router      /customer/{id} [get]
func (h CustomerHandler) Get(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    
    id, _ := strconv.Atoi(ctx.Param("id"))
    
    customer, errr := h.customerSvc.Get(&context, uint(id))
    if errr != nil {
        h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    h.dataResp.DefaultSuccessResponse(customer).FormatAndSend(&context, ctx, http.StatusOK)
}

// GetAll
// @Summary     Get all active customers
// @Description Get all active customers
// @Tags        Customer
// @Accept      json
// @Success     200    {object} responseModel.Customer
// @Failure     400    {object} response.DataResponse
// @Param       search query    string false "search"
// @Router      /customer [get]
func (h CustomerHandler) GetAllCustomers(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    
    search := ctx.Query("search")
    search = util.EncloseWithSingleQuote(search)
    
    customers, errr := h.customerSvc.GetAll(&context, search)
    if errr != nil {
        h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    h.dataResp.DefaultSuccessResponse(customers).FormatAndSend(&context, ctx, http.StatusOK)
}

// Delete
// @Summary     Delete Customer
// @Description Deletes an instance of Customer
// @Tags        Customer
// @Accept      json
// @Success     200 {object} response.Response
// @Failure     400 {object} response.Response
// @Param       id  path     int true "customer id"
// @Router      /customer/{id} [delete]
func (h CustomerHandler) Delete(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    
    id, _ := strconv.Atoi(ctx.Param("id"))
    err := h.customerSvc.Delete(&context, uint(id))
    if err != nil {
        h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
```

**Key Points:**
- Always copy context: `context := util.CopyContextFromGin(ctx)`
- Use `ctx.Bind()` for request body parsing
- Extract path params: `ctx.Param("id")`
- Extract query params: `ctx.Query("search")`
- Use `resp.DefaultFailureResponse()` for errors
- Use `resp.SuccessResponse()` for simple success
- Use `dataResp.DefaultSuccessResponse()` for data responses
- Include Swagger annotations

### 7. Router Layer (`internal/router/router.go`)

HTTP route definitions.

```go
func InitRouter(handler baseHandler.BaseHandler, newRelic *newrelic.Application, srvConfig config.ServerConfig) *gin.Engine {
    g := gin.Default()
    g.Use(gin.Recovery())
    
    // Middlewares
    g.Use(middleware.NewRelicMiddleWare(newRelic))
    g.Use(middleware.LogMiddleware())
    g.Use(middleware.Security())
    g.Use(middleware.CORS())
    g.Use(middleware.RequestParser())
    g.Use(gzip.Gzip(gzip.DefaultCompression))
    
    appRouter := g.Group(constants.API_PREFIX_V1)
    {
        // Non-JWT endpoints
        nonJwtEndpoints := appRouter.Group("user")
        {
            nonJwtEndpoints.POST("login", handler.UserHandler.Login)
            nonJwtEndpoints.POST("forgot-password", handler.UserHandler.ForgotPassword)
        }
        
        // JWT-protected endpoints
        customerEndpoints := appRouter.Group("customer", router.VerifyJWT(srvConfig.JwtSecretKey))
        {
            customerEndpoints.POST("", handler.CustomerHandler.SaveCustomer)
            customerEndpoints.PUT(":id", handler.CustomerHandler.UpdateCustomer)
            customerEndpoints.GET("autocomplete", handler.CustomerHandler.AutocompleteCustomer)
            customerEndpoints.GET(":id", handler.CustomerHandler.Get)
            customerEndpoints.GET("", handler.CustomerHandler.GetAllCustomers)
            customerEndpoints.DELETE(":id", handler.CustomerHandler.Delete)
        }
    }
    
    return g
}
```

**Key Points:**
- Use route groups for logical organization
- Apply middlewares at appropriate levels
- JWT endpoints use `router.VerifyJWT()` middleware
- Place more specific routes before generic ones (e.g., `/autocomplete` before `/:id`)
- Use semantic HTTP methods (POST for create, PUT for update, GET for read, DELETE for delete)

### 8. Dependency Injection (`internal/di/`)

Using Google Wire for dependency injection.

**wire.go:**
```go
var handlerSet = wire.NewSet(
    handler.ProvideCustomerHandler,
    // ... other handlers
)

var svcSet = wire.NewSet(
    service.ProvideCustomerService,
    // ... other services
)

var repoSet = wire.NewSet(
    repository.ProvideCustomerRepository,
    // ... other repositories
)

var mapperSet = wire.NewSet(
    mapper.ProvideMapper,
    mapper.ProvideResponseMapper,
)

func InitApp(ctx *context.Context) (*app.App, error) {
    wire.Build(
        appConfigSet,
        pkgServiceSet,
        logSet,
        mapperSet,
        routerSet,
        dbSet,
        repoSet,
        svcSet,
        handlerSet,
        wire.Struct(new(app.App), "*"),
    )
    return &app.App{}, nil
}
```

**BaseHandler Registration:**
```go
// internal/handler/base/base_handler.go
type BaseHandler struct {
    HealthHandler      Health
    UserHandler        *handler.UserHandler
    CustomerHandler    *handler.CustomerHandler
    // ... other handlers
}

func ProvideBaseHandler(
    health Health,
    user *handler.UserHandler,
    customerHandler *handler.CustomerHandler,
    // ... other handlers
) BaseHandler {
    return BaseHandler{
        HealthHandler:   health,
        UserHandler:     user,
        CustomerHandler: customerHandler,
        // ... other handlers
    }
}
```

## Standard Workflow for Adding a New Feature

Follow these steps in order when adding a new feature (e.g., "Product"):

### Step 1: Entity Definition
**File:** `internal/entities/product.go`

```go
package entities

type Product struct {
    *Model `mapstructure:",squash"`
    
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    SKU         string  `json:"sku"`
    
    // Relations
    CategoryId uint     `json:"categoryId"`
    Category   Category `gorm:"foreignKey:CategoryId"`
}

func (Product) TableNameForQuery() string {
    return "\"stich\".\"Products\" E"
}
```

### Step 2: Request Model
**File:** `internal/model/request/product.go`

```go
package requestModel

type Product struct {
    ID       uint   `json:"id,omitempty"`
    IsActive bool   `json:"isActive,omitempty"`
    
    Name        string  `json:"name,omitempty"`
    Description string  `json:"description,omitempty"`
    Price       float64 `json:"price,omitempty"`
    SKU         string  `json:"sku,omitempty"`
    CategoryId  uint    `json:"categoryId,omitempty"`
}
```

### Step 3: Response Model
**File:** `internal/model/response/product.go`

```go
package responseModel

type Product struct {
    ID       uint `json:"id,omitempty"`
    IsActive bool `json:"isActive,omitempty"`
    
    Name        string   `json:"name,omitempty"`
    Description string   `json:"description,omitempty"`
    Price       float64  `json:"price,omitempty"`
    SKU         string   `json:"sku,omitempty"`
    CategoryId  uint     `json:"categoryId,omitempty"`
    Category    *Category `json:"category,omitempty"`
    
    AuditFields
}
```

### Step 4: Mapper Methods
**File:** `internal/mapper/mapper.go` (add to interface and implementation)

```go
// Interface
type Mapper interface {
    // ... existing methods
    Product(e requestModel.Product) (*entities.Product, error)
}

// Implementation
func (m *mapper) Product(e requestModel.Product) (*entities.Product, error) {
    return &entities.Product{
        Model:       &entities.Model{ID: e.ID, IsActive: e.IsActive},
        Name:        e.Name,
        Description: e.Description,
        Price:       e.Price,
        SKU:         e.SKU,
        CategoryId:  e.CategoryId,
    }, nil
}
```

**File:** `internal/mapper/response_mapper.go` (add to interface and implementation)

```go
// Interface
type ResponseMapper interface {
    // ... existing methods
    Product(e *entities.Product) (*responseModel.Product, error)
    Products(items []entities.Product) ([]responseModel.Product, error)
}

// Implementation
func (m *responseMapper) Product(e *entities.Product) (*responseModel.Product, error) {
    if e == nil {
        return nil, nil
    }
    
    var category *responseModel.Category
    if e.Category != nil {
        cat, err := m.Category(&e.Category)
        if err != nil {
            return nil, err
        }
        category = cat
    }
    
    return &responseModel.Product{
        ID:          e.ID,
        IsActive:    e.IsActive,
        Name:        e.Name,
        Description: e.Description,
        Price:       e.Price,
        SKU:         e.SKU,
        CategoryId:  e.CategoryId,
        Category:    category,
        AuditFields: responseModel.AuditFields{
            CreatedAt: e.CreatedAt,
            UpdatedAt: e.UpdatedAt,
            CreatedBy: e.CreatedBy,
            UpdatedBy: e.UpdatedBy,
        },
    }, nil
}

func (m *responseMapper) Products(items []entities.Product) ([]responseModel.Product, error) {
    result := make([]responseModel.Product, 0)
    for _, item := range items {
        mappedItem, err := m.Product(&item)
        if err != nil {
            return nil, err
        }
        result = append(result, *mappedItem)
    }
    return result, nil
}
```

### Step 5: Repository
**File:** `internal/repository/product_repository.go`

```go
package repository

import (
    "context"
    
    "github.com/imkarthi24/sf-backend/internal/entities"
    "github.com/imkarthi24/sf-backend/internal/repository/scopes"
    "github.com/loop-kar/pixie/db"
    "github.com/loop-kar/pixie/errs"
)

type ProductRepository interface {
    Create(*context.Context, *entities.Product) *errs.XError
    Update(*context.Context, *entities.Product) *errs.XError
    Get(*context.Context, uint) (*entities.Product, *errs.XError)
    GetAll(*context.Context, string) ([]entities.Product, *errs.XError)
    Delete(*context.Context, uint) *errs.XError
}

type productRepository struct {
    GormDAL
}

func ProvideProductRepository(customDB GormDAL) ProductRepository {
    return &productRepository{GormDAL: customDB}
}

func (pr *productRepository) Create(ctx *context.Context, product *entities.Product) *errs.XError {
    res := pr.WithDB(ctx).Create(&product)
    if res.Error != nil {
        return errs.NewXError(errs.DATABASE, "Unable to save product", res.Error)
    }
    return nil
}

func (pr *productRepository) Update(ctx *context.Context, product *entities.Product) *errs.XError {
    return pr.GormDAL.Update(ctx, *product)
}

func (pr *productRepository) Get(ctx *context.Context, id uint) (*entities.Product, *errs.XError) {
    product := entities.Product{}
    res := pr.WithDB(ctx).
        Preload("Category").
        Find(&product, id)
    if res.Error != nil {
        return nil, errs.NewXError(errs.DATABASE, "Unable to find product", res.Error)
    }
    return &product, nil
}

func (pr *productRepository) GetAll(ctx *context.Context, search string) ([]entities.Product, *errs.XError) {
    var products []entities.Product
    res := pr.WithDB(ctx).Table(entities.Product{}.TableNameForQuery()).
        Scopes(scopes.Channel(), scopes.IsActive()).
        Scopes(scopes.ILike(search, "name", "description", "sku")).
        Scopes(db.Paginate(ctx)).
        Find(&products)
    if res.Error != nil {
        return nil, errs.NewXError(errs.DATABASE, "Unable to find products", res.Error)
    }
    return products, nil
}

func (pr *productRepository) Delete(ctx *context.Context, id uint) *errs.XError {
    product := &entities.Product{Model: &entities.Model{ID: id, IsActive: false}}
    err := pr.GormDAL.Delete(ctx, product)
    if err != nil {
        return err
    }
    return nil
}
```

### Step 6: Service
**File:** `internal/service/product_service.go`

```go
package service

import (
    "context"
    
    "github.com/imkarthi24/sf-backend/internal/entities"
    "github.com/imkarthi24/sf-backend/internal/mapper"
    requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
    responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
    "github.com/imkarthi24/sf-backend/internal/repository"
    "github.com/loop-kar/pixie/errs"
)

type ProductService interface {
    SaveProduct(*context.Context, requestModel.Product) *errs.XError
    UpdateProduct(*context.Context, requestModel.Product, uint) *errs.XError
    Get(*context.Context, uint) (*responseModel.Product, *errs.XError)
    GetAll(*context.Context, string) ([]responseModel.Product, *errs.XError)
    Delete(*context.Context, uint) *errs.XError
}

type productService struct {
    productRepo repository.ProductRepository
    mapper      mapper.Mapper
    respMapper  mapper.ResponseMapper
}

func ProvideProductService(
    repo repository.ProductRepository,
    mapper mapper.Mapper,
    respMapper mapper.ResponseMapper,
) ProductService {
    return productService{
        productRepo: repo,
        mapper:      mapper,
        respMapper:  respMapper,
    }
}

func (svc productService) SaveProduct(ctx *context.Context, product requestModel.Product) *errs.XError {
    dbProduct, err := svc.mapper.Product(product)
    if err != nil {
        return errs.NewXError(errs.INVALID_REQUEST, "Unable to save product", err)
    }
    
    errr := svc.productRepo.Create(ctx, dbProduct)
    if errr != nil {
        return errr
    }
    
    return nil
}

func (svc productService) UpdateProduct(ctx *context.Context, product requestModel.Product, id uint) *errs.XError {
    dbProduct, err := svc.mapper.Product(product)
    if err != nil {
        return errs.NewXError(errs.INVALID_REQUEST, "Unable to update product", err)
    }
    
    dbProduct.ID = id
    errr := svc.productRepo.Update(ctx, dbProduct)
    if errr != nil {
        return errr
    }
    return nil
}

func (svc productService) Get(ctx *context.Context, id uint) (*responseModel.Product, *errs.XError) {
    product, err := svc.productRepo.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    
    mappedProduct, mapErr := svc.respMapper.Product(product)
    if mapErr != nil {
        return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Product data", mapErr)
    }
    
    return mappedProduct, nil
}

func (svc productService) GetAll(ctx *context.Context, search string) ([]responseModel.Product, *errs.XError) {
    products, err := svc.productRepo.GetAll(ctx, search)
    if err != nil {
        return nil, err
    }
    
    mappedProducts, mapErr := svc.respMapper.Products(products)
    if mapErr != nil {
        return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Product data", mapErr)
    }
    
    return mappedProducts, nil
}

func (svc productService) Delete(ctx *context.Context, id uint) *errs.XError {
    err := svc.productRepo.Delete(ctx, id)
    if err != nil {
        return err
    }
    return nil
}
```

### Step 7: Handler
**File:** `internal/handler/product_handler.go`

```go
package handler

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    requesModel "github.com/imkarthi24/sf-backend/internal/model/request"
    "github.com/imkarthi24/sf-backend/internal/service"
    "github.com/loop-kar/pixie/errs"
    "github.com/loop-kar/pixie/response"
    "github.com/loop-kar/pixie/util"
)

type ProductHandler struct {
    productSvc service.ProductService
    resp       response.Response
    dataResp   response.DataResponse
}

func ProvideProductHandler(svc service.ProductService) *ProductHandler {
    return &ProductHandler{productSvc: svc}
}

// @Summary     Save Product
// @Description Saves an instance of Product
// @Tags        Product
// @Accept      json
// @Success     201      {object} response.Response
// @Failure     400      {object} response.Response
// @Failure     501      {object} response.Response
// @Param       product  body     requestModel.Product true "product"
// @Router      /product [post]
func (h ProductHandler) SaveProduct(ctx *gin.Context) {
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
}

// @Summary     Update Product
// @Description Updates an instance of Product
// @Tags        Product
// @Accept      json
// @Success     201      {object} response.Response
// @Failure     400      {object} response.Response
// @Param       product  body     requestModel.Product true "product"
// @Param       id       path     int                  true "Product id"
// @Router      /product/{id} [put]
func (h ProductHandler) UpdateProduct(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    var product requesModel.Product
    err := ctx.Bind(&product)
    if err != nil {
        x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
        h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    id, _ := strconv.Atoi(ctx.Param("id"))
    errr := h.productSvc.UpdateProduct(&context, product, uint(id))
    if errr != nil {
        h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
        return
    }
    
    h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

// @Summary     Get a specific Product
// @Description Get an instance of Product
// @Tags        Product
// @Accept      json
// @Success     200 {object} responseModel.Product
// @Failure     400 {object} response.DataResponse
// @Param       id  path     int true "Product id"
// @Router      /product/{id} [get]
func (h ProductHandler) Get(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    
    id, _ := strconv.Atoi(ctx.Param("id"))
    
    product, errr := h.productSvc.Get(&context, uint(id))
    if errr != nil {
        h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    h.dataResp.DefaultSuccessResponse(product).FormatAndSend(&context, ctx, http.StatusOK)
}

// @Summary     Get all active products
// @Description Get all active products
// @Tags        Product
// @Accept      json
// @Success     200    {object} responseModel.Product
// @Failure     400    {object} response.DataResponse
// @Param       search query    string false "search"
// @Router      /product [get]
func (h ProductHandler) GetAllProducts(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    
    search := ctx.Query("search")
    search = util.EncloseWithSingleQuote(search)
    
    products, errr := h.productSvc.GetAll(&context, search)
    if errr != nil {
        h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    h.dataResp.DefaultSuccessResponse(products).FormatAndSend(&context, ctx, http.StatusOK)
}

// @Summary     Delete Product
// @Description Deletes an instance of Product
// @Tags        Product
// @Accept      json
// @Success     200 {object} response.Response
// @Failure     400 {object} response.Response
// @Param       id  path     int true "product id"
// @Router      /product/{id} [delete]
func (h ProductHandler) Delete(ctx *gin.Context) {
    context := util.CopyContextFromGin(ctx)
    
    id, _ := strconv.Atoi(ctx.Param("id"))
    err := h.productSvc.Delete(&context, uint(id))
    if err != nil {
        h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
        return
    }
    
    h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}
```

### Step 8: Wire Registration
**File:** `internal/di/wire.go`

Add to respective wire sets:

```go
var handlerSet = wire.NewSet(
    // ... existing
    handler.ProvideProductHandler,
)

var svcSet = wire.NewSet(
    // ... existing
    service.ProvideProductService,
)

var repoSet = wire.NewSet(
    // ... existing
    repository.ProvideProductRepository,
)
```

### Step 9: BaseHandler Registration
**File:** `internal/handler/base/base_handler.go`

```go
type BaseHandler struct {
    // ... existing handlers
    ProductHandler *handler.ProductHandler
}

func ProvideBaseHandler(
    // ... existing parameters
    productHandler *handler.ProductHandler,
) BaseHandler {
    return BaseHandler{
        // ... existing handlers
        ProductHandler: productHandler,
    }
}
```

### Step 10: Router Registration
**File:** `internal/router/router.go`

```go
productEndpoints := appRouter.Group("product", router.VerifyJWT(srvConfig.JwtSecretKey))
{
    productEndpoints.POST("", handler.ProductHandler.SaveProduct)
    productEndpoints.PUT(":id", handler.ProductHandler.UpdateProduct)
    productEndpoints.GET(":id", handler.ProductHandler.Get)
    productEndpoints.GET("", handler.ProductHandler.GetAllProducts)
    productEndpoints.DELETE(":id", handler.ProductHandler.Delete)
}
```

### Step 11: Generate Wire
Run wire to generate dependency injection code:

```bash
cd internal/di
wire
```

## Naming Conventions

### File Names
- Lowercase with underscores: `customer_handler.go`, `product_service.go`
- Match the primary type: `CustomerHandler` → `customer_handler.go`

### Package Names
- Single word, lowercase: `handler`, `service`, `repository`
- Sub-packages use descriptive names: `request`, `response`, `base`

### Type Names
- PascalCase: `CustomerHandler`, `ProductService`
- Suffix with layer name: `CustomerHandler`, `CustomerService`, `CustomerRepository`

### Interface Names
- Same as implementation without suffix: `CustomerService` (interface), `customerService` (struct)

### Method Names
- PascalCase for exported: `SaveCustomer`, `GetAll`
- camelCase for internal: `validateInput`, `processData`
- Standard CRUD names: `Create`, `Update`, `Get`, `GetAll`, `Delete` (repository)
- Standard CRUD names: `Save`, `Update`, `Get`, `GetAll`, `Delete` (service/handler)

### Variable Names
- camelCase: `customerRepo`, `dbCustomer`
- Context always named `ctx` in parameters
- Gin context always named `ctx` in handler methods
- Copied context always named `context`: `context := util.CopyContextFromGin(ctx)`

### Provider Functions
- Prefix with `Provide`: `ProvideCustomerHandler`, `ProvideCustomerService`

## Common Patterns

### Error Handling
```go
// In Repository
if res.Error != nil {
    return errs.NewXError(errs.DATABASE, "Unable to find customer", res.Error)
}

// In Service
if err != nil {
    return errs.NewXError(errs.INVALID_REQUEST, "Unable to save customer", err)
}

// In Handler
if errr != nil {
    h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
    return
}
```

### Context Handling
```go
// Handler - always copy context
context := util.CopyContextFromGin(ctx)

// Pass context as pointer
h.customerSvc.SaveCustomer(&context, customer)

// Repository - use WithDB to get context-aware DB
res := cr.WithDB(ctx).Create(&customer)
```

### Search/Filter Pattern
```go
// In Handler
search := ctx.Query("search")
search = util.EncloseWithSingleQuote(search)

// In Repository
res := cr.WithDB(ctx).Table(entities.Customer{}.TableNameForQuery()).
    Scopes(scopes.Channel(), scopes.IsActive()).
    Scopes(scopes.ILike(search, "first_name", "last_name", "email")).
    Scopes(db.Paginate(ctx)).
    Find(&customers)
```

### Preloading Relations
```go
res := cr.WithDB(ctx).
    Preload("Category").
    Preload("Variants").
    Preload("Variants.Options").
    Find(&product, id)
```

### Soft Delete Pattern
```go
func (cr *customerRepository) Delete(ctx *context.Context, id uint) *errs.XError {
    customer := &entities.Customer{Model: &entities.Model{ID: id, IsActive: false}}
    err := cr.GormDAL.Delete(ctx, customer)
    if err != nil {
        return err
    }
    return nil
}
```

### Mapper Nil Check
```go
func (m *responseMapper) Customer(e *entities.Customer) (*responseModel.Customer, error) {
    if e == nil {
        return nil, nil
    }
    // ... mapping logic
}
```

### Swagger Annotations
```go
// @Summary     Get a specific Customer
// @Description Get an instance of Customer
// @Tags        Customer
// @Accept      json
// @Success     200 {object} responseModel.Customer
// @Failure     400 {object} response.DataResponse
// @Param       id  path     int true "Customer id"
// @Router      /customer/{id} [get]
```

## Additional Notes

### Multi-tenancy (Channel Support)
- All entities include `ChannelId` via embedded `Model`
- Repository queries automatically filter by channel using `scopes.Channel()`
- Channel ID extracted from JWT and stored in context
- System admin (channel_id = 0) can access all data

### Audit Trail
- All entities include audit fields via embedded `Model`:
  - `CreatedAt`, `UpdatedAt` (timestamps)
  - `CreatedById`, `UpdatedById` (user IDs)
  - `CreatedBy`, `UpdatedBy` (user names - transient)
- Automatically populated via GORM hooks (`BeforeCreate`, `BeforeUpdate`)
- User ID extracted from context

### Pagination
- Applied via `db.Paginate(ctx)` scope in repository
- Page number and size passed via context from request headers
- Middleware extracts pagination params from request

### Transaction Management
- Use `db.DBTransactionManager` for multi-repository transactions
- Service layer orchestrates transactions when needed
- Repository methods accept context for transaction support

### Scope Reusability
- Common scopes defined in `internal/repository/scopes/scope.go`
- Entity-specific scopes in respective scope files (e.g., `channel_scopes.go`)
- Apply multiple scopes in chain: `Scopes(scope1(), scope2(), scope3())`

---

**Last Updated:** 2026-02-15  
**Version:** 1.0
