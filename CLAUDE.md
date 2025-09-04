# Study Guides API - Architecture Documentation

## Overview
The Study Guides API is a Go-based microservice providing gRPC and HTTP endpoints for a study guides platform. It implements a clean architecture with clear separation between API layer, business logic, and data access layers.

## Documentation Structure
For detailed component documentation, see:
- `cmd/server/CLAUDE.md` - Server initialization and lifecycle
- `internal/store/CLAUDE.md` - Data access layer with 9 domain stores, including indexing
- `internal/services/CLAUDE.md` - Business logic and service implementations
- `internal/mcp/CLAUDE.md` - Model Context Protocol system for AI-driven operations
- `internal/middleware/CLAUDE.md` - Authentication, rate limiting, and error handling
- `internal/lib/CLAUDE.md` - Shared libraries (AI, routing, formatting)
- `internal/errors/CLAUDE.md` - Error definitions and handling
- `internal/utils/CLAUDE.md` - Utility functions
- `prisma/CLAUDE.md` - Database schema management

## Architecture Principles

### Layered Architecture
The codebase follows a classic layered architecture pattern:
- **API Layer** (`api/v1/`): Protocol Buffers definitions for gRPC services
- **Service Layer** (`internal/services/`): Business logic and service implementations  
- **Data Access Layer** (`internal/store/`): Data persistence abstractions and implementations
- **AI Integration Layer** (`internal/mcp/`): Model Context Protocol for AI-driven operations
- **Middleware Layer** (`internal/middleware/`): Cross-cutting concerns (auth, rate limiting, error handling)
- **Router Layer** (`internal/lib/router/`): Operation routing and handler mapping (legacy)

### Key Design Patterns

#### Repository Pattern
Each domain (user, tag, question, etc.) has its own store interface with SQL implementations:
```go
// Example from internal/store/tag/tag.go
type TagStore interface {
    GetTag(ctx context.Context, id string) (*sharedpb.Tag, error)
    ListTags(ctx context.Context, opts ListTagsOptions) ([]*sharedpb.Tag, error)
    // ... other methods
}
```

#### Dependency Injection
The main store aggregates all domain stores and is injected throughout the application:
```go
// internal/store/store.go
type Store interface {
    SearchStore() search.SearchStore
    TagStore() tag.TagStore
    UserStore() user.UserStore
    // ... other stores
}
```

#### Authentication Middleware Pattern
JWT-based authentication is handled via gRPC interceptors:
```go
// internal/middleware/auth.go
func AuthUnaryInterceptor(secret string) grpc.UnaryServerInterceptor
```

#### Handler Pattern with Base Handlers
Services use base handlers for common authentication patterns:
```go
// internal/services/base.go
func PublicBaseHandler(ctx context.Context, fn func(ctx context.Context) (interface{}, error))
func AuthBaseHandler(ctx context.Context, fn func(ctx context.Context, session *SessionDetails) (interface{}, error))
```

## Directory Structure & Responsibilities

### `/api/v1/`
Protocol Buffer definitions organized by service domain:
- **Service domains**: admin, chat, devops, health, interaction, question, roland, search, tag, user
- **Shared types**: Common protobuf messages in `/shared/` subdirectory
- **Generated code**: `.pb.go` and `_grpc.pb.go` files are auto-generated

### `/cmd/server/`
Application entry point and server management:
- `main.go`: Application bootstrap with environment loading and store initialization
- `servermanager.go`: Graceful shutdown handling with 30-second timeout
- `server.go`: Core server implementation (not examined in detail)

### `/internal/services/`
Business logic layer with one service per domain matching the API structure. Each service implements the corresponding gRPC service interface.

### `/internal/store/`
Data access layer with interface/implementation separation:
- Each domain has its own package with interface definition
- SQL implementations follow naming pattern: `Sql{Domain}Store`
- External service integrations (Algolia for search, DigitalOcean for devops)

### `/internal/middleware/`
Cross-cutting concerns:
- **Authentication**: JWT token validation and context injection
- **Rate limiting**: Request throttling
- **Error handling**: Centralized error processing

### `/internal/lib/`
Shared library code:
- **router**: Operation routing system mapping tool names to handlers
- **tools**: Tool definitions and mapping utilities
- **webrouter**: HTTP web server for static content and health checks
- **ai**: AI integration utilities

## Key Architectural Decisions

### Protocol Buffers + gRPC
The API is defined using Protocol Buffers, enabling:
- Strong typing across service boundaries
- Code generation for multiple languages
- Efficient binary serialization
- Built-in versioning support

### Multi-Database Architecture
- Main PostgreSQL database for core entities
- Separate Roland database (likely for AI/chat functionality)
- Algolia for search functionality
- DigitalOcean for deployment operations

### Environment-Based Configuration
Configuration is handled via environment variables:
- `DATABASE_URL`: Main database connection
- `ROLAND_DATABASE_URL`: Roland-specific database
- `ALGOLIA_APP_ID` / `ALGOLIA_ADMIN_API_KEY`: Search service credentials

### Graceful Shutdown
The server implements proper graceful shutdown with:
- Signal handling for SIGINT/SIGTERM
- 30-second timeout for cleanup
- Coordinated shutdown across all services

## Common Patterns

### Error Handling
The application uses gRPC status codes consistently across all layers:
```go
return nil, status.Error(codes.Internal, err.Error())
return nil, status.Error(codes.FailedPrecondition, "missing required environment variables")
```
For detailed error handling implementation, see:
- `internal/errors/CLAUDE.md` - Custom error definitions
- `internal/middleware/CLAUDE.md` - Error middleware and transformation

### Context Propagation
Context is passed through all layers for:
- Request cancellation
- Authentication details (see `internal/middleware/CLAUDE.md`)
- Distributed tracing support

### Store Initialization
All stores are initialized at startup with fail-fast behavior - if any store fails to initialize, the application terminates. For implementation details, see `internal/store/CLAUDE.md`.

## Naming Conventions

### Files
- Service files: `{domain}.go` in `/internal/services/`
- Store interfaces: `{domain}.go` in `/internal/store/{domain}/`
- Store implementations: `sql{domain}store.go`
- Protocol buffers: `{domain}.proto`

### Packages
- Service packages match API structure: `api/v1/{domain}` → `internal/services/{domain}`
- Store packages: `internal/store/{domain}`
- Generated code packages: `{domain}v1` for imports

## Gotchas & Non-Obvious Behaviors

1. **Authentication is Optional**: The auth middleware extracts JWT claims but doesn't enforce authentication - individual handlers must check if user is authenticated

2. **Store Aggregation**: The main `Store` interface aggregates all domain stores, but each domain store is initialized independently

3. **Environment Variable Dependencies**: Missing critical environment variables cause immediate application startup failure

4. **Roland Separation**: The Roland functionality (likely AI/chat) uses a separate database connection, suggesting it may have different scaling or data requirements

5. **Dual AI Systems**: 
   - **MCP (Primary)**: Model Context Protocol for modern AI operations via ChatService
   - **Legacy Router**: Tool-based routing for older integrations

6. **Indexing Modes**: Search index synchronization has two distinct modes:
   - **"index tags"**: Incremental sync (only changed items)
   - **"force index tags"**: Complete rebuild (all items)

7. **Static Content Serving**: Despite being primarily a gRPC API, the service also serves static web content (favicon, images, CSS) through the webrouter

## Dependencies

### Core Go Packages
- `google.golang.org/grpc`: gRPC framework
- `google.golang.org/protobuf`: Protocol buffer support
- `github.com/jackc/pgx/v5`: PostgreSQL driver
- `github.com/golang-jwt/jwt/v5`: JWT authentication

### External Services
- **Algolia**: Search functionality
- **DigitalOcean**: DevOps/deployment operations
- **OpenAI**: AI functionality (via `github.com/sashabaranov/go-openai`)

### Development Tools
- `github.com/joho/godotenv`: Environment file loading
- `github.com/dustinkirkland/golang-petname`: Name generation utilities
- `github.com/lucsky/cuid`: Unique ID generation