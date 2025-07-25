# internal/middleware Package

## Purpose and Responsibility
The `internal/middleware` package provides cross-cutting concerns for the gRPC server through interceptors. It handles authentication, rate limiting, error transformation, and session management in a centralized, reusable way.

## Key Architectural Decisions

### gRPC Interceptor Pattern
All middleware is implemented as gRPC unary server interceptors:
- **Chain execution**: Interceptors are chained in a specific order during server setup
- **Context enrichment**: Middleware adds data to the request context for downstream use
- **Request/response transformation**: Interceptors can modify requests and responses

### Optional Authentication Model
Authentication is extractive rather than enforcing:
- JWT tokens are parsed and validated when present
- User information is stored in context for services to check
- Individual services decide whether authentication is required
- Graceful degradation for unauthenticated requests

### Per-User Rate Limiting
Rate limiting is applied per-user or per-IP:
- Authenticated requests are limited by user ID
- Unauthenticated requests are limited by IP address
- Separate rate limiters maintained for each user/IP

## Implementation Details

### Authentication Interceptor (`auth.go`)
The `AuthUnaryInterceptor` handles JWT token processing:

#### Token Extraction and Validation
- Extracts JWT from `Authorization: Bearer <token>` header
- Validates token signature using provided secret
- Stores user ID and roles in request context if valid
- Continues processing even if token is invalid/missing

#### Role Mapping
Maps string roles from JWT claims to protobuf enums:
```go
"admin" → UserRole_USER_ROLE_ADMIN
"user" → UserRole_USER_ROLE_USER
"freelancer" → UserRole_USER_ROLE_FREELANCER
"tester" → UserRole_USER_ROLE_TESTER
```

#### Context Storage
User information is stored using typed context keys:
- `userIDKey` → User ID string
- `userRoleKey` → Array of UserRole enums

### Rate Limiting Interceptor (`ratelimit.go`)
The `RateLimitUnaryInterceptor` implements token bucket rate limiting:

#### Rate Limiter Storage
- `limiterStore` maintains per-key rate limiters
- Thread-safe with mutex protection
- Creates new limiters on-demand for new keys

#### Key Selection Strategy
1. **Primary**: User ID from JWT (if authenticated)
2. **Fallback**: Client IP address from peer context
3. **Default**: "unknown" if neither available

#### Rate Limit Enforcement
- Uses `golang.org/x/time/rate` package
- Configurable rate (requests per second) and burst capacity
- Returns `ResourceExhausted` status if limit exceeded

### Error Interceptor (`error.go`)
The `ErrorUnaryInterceptor` transforms application errors to user-friendly messages:

#### Custom Error Mapping
Maps specific application errors to friendly messages:
- `ErrToolNotFound` → "I couldn't understand how to help with that request..."
- `ErrNotFound` → "The requested resource was not found."
- `ErrSystemPromptEmpty` → "System configuration error..."

#### Pattern Matching
Performs string matching on error messages for common patterns:
- "AI did not call any tools" → Suggests valid commands
- "AI returned no choices" → Suggests retry
- "failed to parse AI response" → Technical difficulties message

#### Error Preservation
Unmapped errors are returned as-is to preserve detailed error information for debugging.

### Session Management (`utils.go`)
Provides utilities for extracting session information from context:

#### SessionDetails Structure
Encapsulates user session information:
```go
type SessionDetails struct {
    UserID    *string              // User ID from JWT
    UserRoles *[]sharedpb.UserRole // User roles from JWT
    IsAuth    bool                 // Whether user is authenticated
}
```

#### Context Extraction Functions
- `UserIDFromContext()` → Extract user ID
- `UserRolesFromContext()` → Extract user roles  
- `GetSessionDetails()` → Get complete session info

#### Role Checking
`HasRole()` method provides convenient role-based authorization:
```go
if session.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
    // Admin-only logic
}
```

## Common Patterns

### Interceptor Chain Order
Middleware must be applied in correct order:
1. **Error handling** - Outermost to catch all errors
2. **Authentication** - Extract user context
3. **Rate limiting** - Apply per-user limits

### Context Key Pattern
Uses typed context keys to prevent collisions:
```go
type contextKey string
const userIDKey contextKey = "userID"
```

### Graceful Degradation
All middleware gracefully handles missing information:
- Missing JWT → Continue with empty user context
- Missing peer info → Use "unknown" rate limit key
- Invalid roles → Empty role slice

## Key Files

### auth.go
JWT authentication with context enrichment:
- Token extraction from Authorization header
- JWT validation and parsing
- User ID and role context storage
- Non-enforcing authentication model

### ratelimit.go
Per-user/IP rate limiting:
- Token bucket rate limiting implementation
- Thread-safe rate limiter storage
- Flexible key selection (user ID or IP)
- gRPC status code integration

### error.go  
Error message transformation:
- Application error to user-friendly message mapping
- Pattern-based error detection
- Preserves unmapped errors for debugging

### utils.go
Session management utilities:
- SessionDetails structure definition
- Context extraction helpers
- Role-based authorization utilities

## Gotchas & Non-Obvious Behaviors

1. **Authentication is Optional**: The auth middleware extracts JWT info but doesn't enforce authentication - services must check `session.IsAuth`

2. **Rate Limiting Fallback**: If no user ID available, rate limiting falls back to IP address, which may cause issues with proxies/load balancers

3. **Role Case Sensitivity**: Role mapping is case-insensitive (`strings.ToLower()` is used)

4. **Context Key Types**: Uses typed context keys to prevent accidental collisions with other packages

5. **Error Message Patterns**: Error interceptor uses string matching, which may be brittle if error messages change

6. **Empty vs Nil Pointers**: SessionDetails always returns valid pointers, but they may point to empty values

7. **Thread Safety**: Rate limiter store uses mutex, but individual limiters are not thread-safe (relies on rate.Limiter's internal synchronization)

8. **No Rate Limit Cleanup**: Rate limiters are created but never cleaned up, which could lead to memory leaks with many unique users

9. **JWT Secret Validation**: JWT parsing continues even with invalid signatures - only stores context if token is valid

10. **Peer Context Dependency**: Rate limiting fallback depends on gRPC peer context, which may not always be available

## Dependencies

### External Libraries
- `github.com/golang-jwt/jwt/v5` - JWT token parsing and validation
- `golang.org/x/time/rate` - Token bucket rate limiting
- `google.golang.org/grpc` - gRPC interceptor framework and status codes

### Internal Dependencies
- `api/v1/shared` - UserRole protobuf enums
- `internal/errors` - Application-specific error types for transformation

## Configuration
Middleware configuration is handled during server setup:
- **JWT Secret**: Environment variable for token validation
- **Rate Limits**: Environment variables for rate and burst capacity
- **Interceptor Order**: Defined in server setup code

## Security Considerations
- JWT secrets should be kept secure and rotated regularly
- Rate limiting helps prevent abuse but may need tuning based on usage patterns
- Error messages balance user-friendliness with avoiding information disclosure
- Authentication context is trusted once validated - no re-validation in services