# internal/services Package

## Purpose and Responsibility
The `internal/services` package contains the business logic layer for the Study Guides API. Each service implements a corresponding gRPC service interface, handling request validation, authentication checks, and coordinating with the data access layer.

## Key Architectural Decisions

### Service-per-Domain Pattern
Each service corresponds to a specific business domain and gRPC service:
- `AdminService` - Administrative operations (user management, tag creation)
- `ChatService` - AI-powered conversational interface with tool routing
- `DevopsService` - Deployment and infrastructure management
- `HealthService` - Health check endpoints
- `InteractionService` - User interaction tracking
- `QuestionService` - Question and assessment management
- `RolandService` - AI assistant functionality
- `SearchService` - Search operations via Algolia
- `TagService` - Tag hierarchy and metadata management
- `UserService` - User account management

### Base Handler Pattern
Two base handlers provide consistent authentication patterns:
- `PublicBaseHandler` - For requests that may or may not be authenticated
- `AuthBaseHandler` - For requests requiring authentication with session details

### Embedded Unimplemented Servers
All services embed the corresponding `Unimplemented*ServiceServer` from generated protobuf code, ensuring forward compatibility when new RPC methods are added.

## Implementation Details

### Chat Service Architecture
The ChatService is the most complex, implementing both legacy tool routing and the new MCP (Model Context Protocol) system:

#### MCP Integration (Primary System)
- **MCP Processor**: Handles natural language processing via `internal/mcp`
- **Repository Adapters**: Including `IndexingRepositoryAdapter` for index management
- **Dynamic tool generation**: Creates OpenAI tools from registered repositories
- **Indexing triggers**: 
  - "index tags" → incremental indexing (only changed items)
  - "force index tags" → complete rebuild (all items)

#### Legacy Tool Integration
- Dynamic system prompt generation based on available tools
- OpenAI function calling for intent classification
- Tool routing via the internal router package
- Format detection from natural language (csv, json, list)

#### Conversation History Management
- `ConversationHistory` struct manages message history with size limits
- Messages are stored in context metadata as JSON
- Smart truncation keeps only recent messages (max 10 messages, 1000 chars each)
- Response summarization prevents token limit issues

#### Response Flow
1. Extract conversation history from context metadata
2. Add user message to history
3. Generate system prompt with available tools
4. Call AI with history and tools
5. Parse tool call from AI response
6. Route to appropriate handler
7. Add response summary to history
8. Return response with updated context

### Authentication Patterns
Services use base handlers for authentication (see `internal/middleware/CLAUDE.md` for auth implementation details):
```go
resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
    if session.UserID == nil {
        return nil, status.Error(codes.Unauthenticated, "authentication required")
    }
    // Business logic here
})
```

### Error Handling
- Consistent use of gRPC status codes
- Internal errors are wrapped with `codes.Internal`
- Authentication errors use `codes.Unauthenticated`
- Invalid requests use appropriate validation codes

## Common Patterns

### Service Constructor Pattern
All services follow the same constructor pattern:
```go
func NewServiceType(store store.Store) *ServiceType {
    return &ServiceType{
        store: store,
    }
}
```

### Store Delegation
Services delegate data operations to the appropriate store:
```go
result, err := s.store.DomainStore().Operation(ctx, params)
```

### Response Type Assertion
Base handlers return `interface{}`, requiring type assertion:
```go
return resp.(*ServiceTypeResponse), nil
```

## Key Files

### base.go
Defines the base handler patterns for public and authenticated requests.

### chat.go
Most complex service implementing:
- Conversation history management
- AI integration with OpenAI
- Tool routing and intent classification
- Format detection and response summarization

### admin.go
Administrative operations requiring authentication:
- User management (KillUser)
- Tag creation (NewTag - placeholder implementation)

### devops.go
Infrastructure management operations:
- Application deployment
- Rollback functionality  
- Deployment listing and status checking

### Other service files
Standard CRUD operations for their respective domains, following the same patterns.

## Gotchas & Non-Obvious Behaviors

1. **Chat Service Dependencies**: Only the chat service requires router and AI client dependencies - others only need the store

2. **Conversation History Limits**: Chat service truncates conversations to prevent token limits (10 messages max, 1000 chars per message)

3. **Response Summarization**: Long responses are intelligently summarized based on operation type for conversation history

4. **Format Detection**: Chat service detects output format requests from natural language ("as csv", "in json", etc.)

5. **Authentication is Optional at Middleware Level**: Services must explicitly check `session.UserID` - the middleware only extracts auth info

6. **Type Assertion Required**: Base handlers return `interface{}` requiring explicit type assertion

7. **AI Tool Routing**: Chat service uses OpenAI function calling to classify user intent and route to appropriate tools

8. **Extensive Debug Logging**: Chat service has extensive debug output for troubleshooting AI responses

9. **Pointer Conversion**: Some services need to convert slices of structs to slices of pointers for protobuf compatibility

10. **Dynamic System Prompts**: Chat service builds system prompts dynamically based on available tools, not hardcoded

## Dependencies

### External Libraries
- `google.golang.org/grpc` - gRPC framework and status codes
- `github.com/sashabaranov/go-openai` - OpenAI API client (chat service only)

### Internal Dependencies
- `api/v1/*` - Generated protobuf service interfaces and types
- `internal/store` - Data access layer
- `internal/middleware` - Authentication context and session details
- `internal/lib/router` - Tool routing (chat service only)
- `internal/lib/ai` - AI client wrapper (chat service only)
- `internal/lib/tools` - Tool definitions and classification (chat service only)

## Authentication Requirements
- **Public Services**: Health, Search (some operations)
- **Authenticated Services**: Admin, Devops, most User operations
- **Mixed Services**: Chat (public), Tag (mixed), Question (mixed)

Services requiring authentication will return `codes.Unauthenticated` if no valid session is present.