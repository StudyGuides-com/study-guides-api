# cmd/server Package

## Purpose and Responsibility
The `cmd/server` package is the main entry point for the Study Guides API server. It handles application bootstrap, server lifecycle management, and provides a unified HTTP/gRPC server implementation.

## Key Architectural Decisions

### Unified HTTP/gRPC Server
- Single server handles both HTTP and gRPC traffic on the same port
- Uses HTTP/2 with h2c (HTTP/2 over cleartext) to support gRPC over HTTP/1.1
- Content-Type header routing: `application/grpc*` → gRPC, everything else → HTTP web router
- Enables serving both API endpoints and static web content from one server

### Middleware Chain Architecture
gRPC requests pass through a chain of interceptors in order:
1. `ErrorUnaryInterceptor()` - Error handling and logging
2. `AuthUnaryInterceptor(JWT_SECRET)` - JWT authentication (optional)
3. `RateLimitUnaryInterceptor(rate, burst)` - Rate limiting per user

### Environment-Based Configuration
Configuration is loaded from environment variables with sensible defaults:
- `PORT` (default: 8080) - Server listening port
- `JWT_SECRET` - JWT token validation key
- `RATE_LIMIT_USER_PER_SECOND` (default: 1.0) - Rate limit per user
- `RATE_LIMIT_USER_BURST` (default: 5) - Rate limit burst capacity
- `OPENAI_API_KEY` & `OPENAI_MODEL` - AI service configuration

## Implementation Details

### Server Lifecycle Management
- **ServerManager**: Orchestrates server startup and graceful shutdown
- **Graceful Shutdown**: 30-second timeout for cleanup on SIGINT/SIGTERM
- **Service Registration**: All gRPC services are registered in `registerServices()`

### Service Registration Pattern
Services are instantiated with their dependencies and registered with the gRPC server:
```go
// Example pattern
serviceInstance := services.NewServiceType(appStore)
servicepb.RegisterServiceTypeServer(s.grpcServer, serviceInstance)
```

All gRPC services including IndexingService are registered in the `registerServices()` method:
```go
indexingpb.RegisterIndexingServiceServer(s.grpcServer, indexingservice.NewIndexingService(appStore))
```

### Special Case: Chat Service
The chat service has additional dependencies beyond the standard store:
- Router for operation handling
- AI client for OpenAI integration

### Detailed Request Logging
Comprehensive request logging includes:
- HTTP method, path, protocol version
- Content-Type and User-Agent headers
- Routing decisions (gRPC vs web router)
- Request lifecycle markers

## Common Patterns

### Environment Variable Parsing
Utility functions provide safe environment variable parsing with fallbacks:
- `parseEnvAsInt(key, fallback)` - Integer parsing
- `parseEnvAsRate(key, fallback)` - Rate limit parsing  
- `getPort()` - Port configuration

### Dependency Injection
- Store is initialized in main and injected into all services
- Services receive their dependencies through constructors
- No global state or singleton patterns

### Error Handling
- Fail-fast initialization: server won't start if store initialization fails
- Graceful error handling during shutdown
- Force stop capability for emergency shutdown

## Key Files

### main.go
- Application entry point
- Environment loading with godotenv
- Store initialization
- Server manager startup

### servermanager.go  
- Server lifecycle management
- Signal handling for graceful shutdown
- Shutdown coordination with timeout

### server.go
- Unified HTTP/gRPC server implementation
- Service registration
- Request routing logic
- Middleware chain setup

### utils.go
- Environment variable parsing utilities
- Configuration helpers with fallback values

## Gotchas & Non-Obvious Behaviors

1. **Single Port Architecture**: Both HTTP and gRPC traffic use the same port, differentiated by Content-Type header

2. **h2c Requirement**: HTTP/2 cleartext (h2c) is essential for gRPC-over-HTTP/1.1 support

3. **Optional Authentication**: JWT middleware extracts auth info but doesn't enforce it - individual services must check authentication

4. **Rate Limiting Scope**: Rate limiting is applied per-user, not per-request

5. **Chat Service Dependencies**: Only the chat service requires the router and AI client - other services only need the store

6. **Force Stop**: Emergency shutdown (`ForceStop()`) immediately terminates gRPC server without waiting for connections to close

7. **Verbose Logging**: Development logging is extremely detailed, logging every request with full headers

8. **Service Order**: Service registration order doesn't matter - gRPC handles routing by service name

## Dependencies

### External Libraries
- `golang.org/x/net/http2` - HTTP/2 support
- `golang.org/x/net/http2/h2c` - HTTP/2 cleartext for gRPC
- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/grpc/reflection` - gRPC reflection for debugging
- `golang.org/x/time/rate` - Rate limiting
- `github.com/joho/godotenv` - Environment file loading

### Internal Dependencies
- `internal/store` - Data access layer
- `internal/services` - Business logic layer including IndexingService
- `internal/core` - Shared business logic services
- `internal/middleware` - Cross-cutting concerns
- `internal/lib/router` - Operation routing (chat service only)
- `internal/lib/ai` - AI integration (chat service only)
- `internal/lib/webrouter` - HTTP web routing
- `api/v1/*` - Generated protobuf services including indexing service