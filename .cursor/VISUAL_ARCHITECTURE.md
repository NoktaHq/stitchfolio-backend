# Stitchfolio Backend - Visual Architecture

This document provides visual representations of the architecture to complement the detailed guides.

## 📊 Complete Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT REQUEST                           │
│                      (HTTP/JSON/REST API)                        │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                        MIDDLEWARE LAYER                          │
│  ┌────────────┐ ┌───────────┐ ┌─────────┐ ┌──────────────┐    │
│  │ NewRelic   │→│   Logging │→│  CORS   │→│ JWT Verify   │    │
│  └────────────┘ └───────────┘ └─────────┘ └──────────────┘    │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                     ROUTER LAYER                                 │
│  internal/router/router.go                                       │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Route Groups:                                            │  │
│  │  • /api/v1/customer                                       │  │
│  │  • /api/v1/order                                          │  │
│  │  • /api/v1/user                                           │  │
│  │  • etc.                                                   │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                     HANDLER LAYER                                │
│  internal/handler/                                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Responsibilities:                                        │  │
│  │  1. Extract request data (params, query, body)           │  │
│  │  2. Validate & bind to request model                     │  │
│  │  3. Call service layer                                   │  │
│  │  4. Format response                                      │  │
│  │  5. Set HTTP status code                                 │  │
│  └──────────────────────────────────────────────────────────┘  │
│  Components:                                                     │
│  • CustomerHandler, OrderHandler, UserHandler, etc.              │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                   REQUEST MODEL LAYER                            │
│  internal/model/request/                                         │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Purpose: Define API input contracts                     │  │
│  │  • Customer, Order, User, etc.                           │  │
│  │  • Validation rules                                      │  │
│  │  • JSON tags for deserialization                         │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                      MAPPER LAYER (IN)                           │
│  internal/mapper/mapper.go                                       │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Transform: Request Model → Entity                       │  │
│  │  • Handle type conversions                               │  │
│  │  • Parse dates/times                                     │  │
│  │  • Map nested objects                                    │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                     SERVICE LAYER                                │
│  internal/service/                                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Responsibilities:                                        │  │
│  │  1. Business logic execution                             │  │
│  │  2. Orchestrate multiple repositories                    │  │
│  │  3. Transaction management                               │  │
│  │  4. Data transformation (via mappers)                    │  │
│  └──────────────────────────────────────────────────────────┘  │
│  Components:                                                     │
│  • CustomerService, OrderService, UserService, etc.              │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER                              │
│  internal/repository/                                            │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Responsibilities:                                        │  │
│  │  1. Database CRUD operations                             │  │
│  │  2. Query construction with GORM                         │  │
│  │  3. Apply scopes (channel, active, search)               │  │
│  │  4. Handle relationships (preload)                       │  │
│  └──────────────────────────────────────────────────────────┘  │
│  Components:                                                     │
│  • CustomerRepository, OrderRepository, UserRepository, etc.     │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                       SCOPES LAYER                               │
│  internal/repository/scopes/                                     │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Reusable Query Filters:                                 │  │
│  │  • Channel() - Multi-tenancy filtering                   │  │
│  │  • IsActive() - Soft delete filtering                    │  │
│  │  • ILike() - Case-insensitive search                     │  │
│  │  • Paginate() - Pagination                               │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                      ENTITY LAYER                                │
│  internal/entities/                                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  GORM Models (Database Tables):                          │  │
│  │  • Customer, Order, User, etc.                           │  │
│  │  • Embedded *Model (ID, timestamps, audit, channel)      │  │
│  │  • Relations & constraints                               │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                    DATABASE (PostgreSQL)                         │
│  Schema: "stich"                                                 │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
                        (Data Stored)
                             ↑
┌─────────────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER (READ)                       │
│  Returns entities from database                                  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                     MAPPER LAYER (OUT)                           │
│  internal/mapper/response_mapper.go                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Transform: Entity → Response Model                      │  │
│  │  • Format dates for API                                  │  │
│  │  • Map nested relations                                  │  │
│  │  • Add computed fields                                   │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                   RESPONSE MODEL LAYER                           │
│  internal/model/response/                                        │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Purpose: Define API output contracts                    │  │
│  │  • Customer, Order, User, etc.                           │  │
│  │  • Includes audit fields                                 │  │
│  │  • Nested related objects                                │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                     HANDLER LAYER (RESPONSE)                     │
│  Formats response and sends to client                            │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT RESPONSE                          │
│                      (HTTP/JSON/REST API)                        │
└─────────────────────────────────────────────────────────────────┘
```

## 🔄 Data Flow Diagram

### Write Operation (POST/PUT)

```
Client Request
    │
    │ JSON Body
    ↓
┌─────────────────┐
│  Handler        │  ctx.Bind() → Request Model
└────────┬────────┘
         │
         │ Request Model
         ↓
┌─────────────────┐
│  Mapper (IN)    │  Request Model → Entity
└────────┬────────┘
         │
         │ Entity
         ↓
┌─────────────────┐
│  Service        │  Business Logic
└────────┬────────┘
         │
         │ Entity
         ↓
┌─────────────────┐
│  Repository     │  db.Create() / db.Update()
└────────┬────────┘
         │
         │ SQL
         ↓
┌─────────────────┐
│  Database       │  INSERT / UPDATE
└─────────────────┘
```

### Read Operation (GET)

```
Client Request
    │
    │ Query Params
    ↓
┌─────────────────┐
│  Handler        │  Extract params (id, search)
└────────┬────────┘
         │
         │ Params
         ↓
┌─────────────────┐
│  Service        │  Call repository
└────────┬────────┘
         │
         │ ID / Search
         ↓
┌─────────────────┐
│  Repository     │  db.Find() with Scopes & Preloads
└────────┬────────┘
         │
         │ SQL Query
         ↓
┌─────────────────┐
│  Database       │  SELECT with JOINs
└────────┬────────┘
         │
         │ Rows
         ↓
┌─────────────────┐
│  Repository     │  Map to Entities
└────────┬────────┘
         │
         │ Entities
         ↓
┌─────────────────┐
│  Service        │  Pass to mapper
└────────┬────────┘
         │
         │ Entities
         ↓
┌─────────────────┐
│  Mapper (OUT)   │  Entity → Response Model
└────────┬────────┘
         │
         │ Response Model
         ↓
┌─────────────────┐
│  Handler        │  Format & Send JSON
└────────┬────────┘
         │
         │ JSON Response
         ↓
Client Response
```

## 🏗️ Dependency Injection Structure

```
                    ┌─────────────────┐
                    │   Wire Config   │
                    │  (wire.go)      │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              ↓              ↓              ↓
        ┌──────────┐  ┌──────────┐  ┌──────────┐
        │ Handlers │  │ Services │  │  Repos   │
        └────┬─────┘  └────┬─────┘  └────┬─────┘
             │             │             │
             └─────────────┼─────────────┘
                           ↓
                   ┌───────────────┐
                   │ BaseHandler   │
                   │ (aggregation) │
                   └───────┬───────┘
                           ↓
                   ┌───────────────┐
                   │   Router      │
                   │ (init routes) │
                   └───────┬───────┘
                           ↓
                   ┌───────────────┐
                   │   Gin Engine  │
                   │ (HTTP server) │
                   └───────────────┘

Wire Sets:
┌─────────────────────────────────────┐
│ handlerSet                          │
│  • ProvideCustomerHandler           │
│  • ProvideOrderHandler              │
│  • ProvideUserHandler               │
│  • ...                              │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ svcSet                              │
│  • ProvideCustomerService           │
│  • ProvideOrderService              │
│  • ProvideUserService               │
│  • ...                              │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ repoSet                             │
│  • ProvideCustomerRepository        │
│  • ProvideOrderRepository           │
│  • ProvideUserRepository            │
│  • ...                              │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ mapperSet                           │
│  • ProvideMapper                    │
│  • ProvideResponseMapper            │
└─────────────────────────────────────┘
```

## 🔐 Multi-Tenancy Architecture

```
┌────────────────────────────────────────────────────────┐
│                    JWT Token                            │
│  Contains: user_id, channel_id, role                    │
└───────────────────────┬────────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────────────┐
│              JWT Verify Middleware                      │
│  Extracts: channel_id → Context                         │
└───────────────────────┬────────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────────────┐
│                   Handler Layer                         │
│  Copies context: util.CopyContextFromGin(ctx)           │
└───────────────────────┬────────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────────────┐
│                  Service Layer                          │
│  Passes context to repository                           │
└───────────────────────┬────────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────────────┐
│                Repository Layer                         │
│  Applies scopes.Channel() filter                        │
└───────────────────────┬────────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────────────┐
│                Channel Scope Logic                      │
│  if channel_id == 0:                                    │
│      return all data (System Admin)                     │
│  else:                                                  │
│      WHERE channel_id = {extracted_channel_id}          │
└───────────────────────┬────────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────────────┐
│                   Database Query                        │
│  SELECT * FROM table WHERE channel_id = X               │
└────────────────────────────────────────────────────────┘

Result: Each tenant sees only their data
```

## 📝 Soft Delete Flow

```
Delete Request
      │
      ↓
┌──────────────┐
│  Handler     │  DELETE /customer/:id
└──────┬───────┘
       │
       ↓
┌──────────────┐
│  Service     │  Delete(ctx, id)
└──────┬───────┘
       │
       ↓
┌──────────────┐
│  Repository  │  customer.IsActive = false
└──────┬───────┘
       │
       ↓
┌──────────────┐
│  Database    │  UPDATE Customers SET is_active=false WHERE id=X
└──────────────┘

Future Queries:
      │
      ↓
┌──────────────┐
│  Repository  │  Scopes(scopes.IsActive())
└──────┬───────┘
       │
       ↓
┌──────────────┐
│  Database    │  SELECT * WHERE is_active = true
└──────────────┘

Result: Deleted records hidden but recoverable
```

## 🔍 Query Scopes Flow

```
Repository GetAll Request
         │
         ↓
┌────────────────────────────────────────┐
│  Apply Scopes Chain:                   │
│  .Model(Customer{})                    │
│  (Use .Table(TableNameForQuery()) only  │
│   for raw SQL; then use Channel("E"),  │
│   IsActive("E").)                      │
└──────────────┬─────────────────────────┘
               │
               ↓
┌────────────────────────────────────────┐
│  Scope: Channel()                      │
│  WHERE channel_id = X                  │
└──────────────┬─────────────────────────┘
               │
               ↓
┌────────────────────────────────────────┐
│  Scope: IsActive()                     │
│  AND is_active = true                  │
└──────────────┬─────────────────────────┘
               │
               ↓
┌────────────────────────────────────────┐
│  Scope: ILike(search, fields...)       │
│  AND (field1 ILIKE '%search%' OR ...)  │
└──────────────┬─────────────────────────┘
               │
               ↓
┌────────────────────────────────────────┐
│  Scope: Paginate(ctx)                  │
│  LIMIT X OFFSET Y                      │
└──────────────┬─────────────────────────┘
               │
               ↓
         Final SQL Query
```

## 🎯 Component Relationship Diagram

```
┌─────────────────────────────────────────────────────────┐
│                  BaseHandler                            │
│  ┌───────────┐ ┌───────────┐ ┌────────────┐           │
│  │  Customer │ │   Order   │ │    User    │  ...       │
│  │  Handler  │ │  Handler  │ │  Handler   │            │
│  └─────┬─────┘ └─────┬─────┘ └──────┬─────┘           │
└────────┼─────────────┼───────────────┼─────────────────┘
         │             │               │
         ↓             ↓               ↓
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│  Customer   │ │   Order     │ │    User     │
│  Service    │ │  Service    │ │  Service    │
└──────┬──────┘ └──────┬──────┘ └──────┬──────┘
       │               │               │
       │      ┌────────┴────────┐      │
       │      ↓                 ↓      │
       │  ┌─────────┐      ┌─────────┐│
       │  │ Request │      │Response ││
       │  │ Mapper  │      │ Mapper  ││
       │  └────┬────┘      └────┬────┘│
       │       │                │     │
       ↓       ↓                ↓     ↓
┌─────────────┐           ┌─────────────┐
│  Customer   │           │    User     │
│  Repository │←──────────│ Repository  │
└──────┬──────┘           └──────┬──────┘
       │                         │
       │    ┌────────────────────┘
       │    │
       ↓    ↓
┌──────────────┐
│   GormDAL    │
│   (Base DB)  │
└──────┬───────┘
       │
       ↓
┌──────────────┐
│  Database    │
│ (PostgreSQL) │
└──────────────┘
```

## 📦 File Organization Map

```
internal/
│
├── entities/                      # Database models (GORM)
│   ├── base_model.go             # Common Model struct
│   ├── customer.go               # Customer entity
│   ├── order.go                  # Order entity
│   └── ...
│
├── model/
│   ├── request/                  # API input DTOs
│   │   ├── customer.go
│   │   ├── order.go
│   │   └── ...
│   │
│   └── response/                 # API output DTOs
│       ├── customer.go
│       ├── order.go
│       ├── audit.go              # Common audit fields
│       └── ...
│
├── mapper/                       # DTO ↔ Entity transformations
│   ├── mapper.go                # Request → Entity
│   └── response_mapper.go       # Entity → Response
│
├── repository/                   # Data access layer
│   ├── db.go                    # GormDAL base
│   ├── customer_repository.go
│   ├── order_repository.go
│   └── scopes/                  # Reusable query filters
│       ├── scope.go            # Common scopes
│       ├── channel_scopes.go
│       └── ...
│
├── service/                      # Business logic layer
│   ├── base/
│   │   └── base_service.go
│   ├── customer_service.go
│   ├── order_service.go
│   └── ...
│
├── handler/                      # HTTP handlers
│   ├── base/
│   │   ├── base_handler.go     # Handler aggregation
│   │   └── health_handler.go
│   ├── customer_handler.go
│   ├── order_handler.go
│   └── ...
│
├── router/                       # Route definitions
│   ├── router.go                # Main router
│   └── middleware/
│       └── jwt.go
│
├── di/                          # Dependency injection (Wire)
│   ├── wire.go                 # Wire config
│   ├── wire_gen.go             # Generated (don't edit)
│   └── provider.go             # Custom providers
│
├── config/                      # Configuration
├── constants/                   # Constants
└── app/                        # Application entry point
```

## 🚀 Request Lifecycle Example

**Scenario:** Create a new customer

```
1. CLIENT
   POST /api/v1/customer
   Body: {"firstName": "John", "lastName": "Doe", ...}
   Header: Authorization: Bearer <JWT>
   
   ↓

2. MIDDLEWARE
   • JWT Verify → Extract channel_id, user_id
   • CORS → Allow origin
   • Logging → Log request
   
   ↓

3. ROUTER
   Match route: POST /api/v1/customer
   → handler.CustomerHandler.SaveCustomer
   
   ↓

4. HANDLER (customer_handler.go)
   • Copy context from Gin
   • Bind JSON to requestModel.Customer
   • Validate request
   • Call service: customerSvc.SaveCustomer(&context, customer)
   
   ↓

5. SERVICE (customer_service.go)
   • Map request → entity: mapper.Customer(customer)
   • Call repo: customerRepo.Create(ctx, dbCustomer)
   • Business logic: Create related Person record
   • Return error or success
   
   ↓

6. REPOSITORY (customer_repository.go)
   • Get DB with context: WithDB(ctx)
   • Execute: db.Create(&customer)
   • Return error or nil
   
   ↓

7. GORM (ORM)
   • Before hooks: BeforeCreate() → Set audit fields
   • Generate SQL: INSERT INTO Customers (...)
   • Execute query
   
   ↓

8. DATABASE (PostgreSQL)
   • Execute INSERT
   • Return generated ID
   • Commit transaction
   
   ↓

9. RETURN PATH
   Repository → Service → Handler
   
   ↓

10. HANDLER (response)
    • Format success response
    • Set HTTP status: 201 Created
    • Send JSON: {"status": "success", "message": "Save success"}
    
    ↓

11. CLIENT
    Receives: 201 Created
    Body: {"status": "success", ...}
```

---

**Note:** These diagrams are textual representations. For actual visual diagrams, consider using tools like draw.io, Lucidchart, or Mermaid.

