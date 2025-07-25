# internal/lib Package

## Purpose and Responsibility
The `internal/lib` package contains shared library code providing common functionality across the application. It includes AI integration, tool definitions, routing logic, web serving, and response formatting utilities.

## Key Architectural Decisions

### Modular Library Design
The lib package is organized into focused sub-packages:
- **ai**: OpenAI API integration and client abstraction
- **tools**: Tool definition system for AI function calling
- **router**: Operation routing and handler management
- **webrouter**: HTTP web server for static content
- **formatting**: Response formatting (JSON, CSV, list) utilities

### AI-First Architecture
The system is designed around AI interactions:
- Tools define available operations for AI function calling
- Router maps tool names to actual handler implementations
- Conversation history management for stateful interactions
- Dynamic system prompt generation based on available tools

### Multi-Format Response System
Handlers support multiple output formats:
- **List**: Human-readable text format
- **JSON**: Structured data for APIs
- **CSV**: Spreadsheet-compatible format
- Format selection via AI natural language detection

## Implementation Details

### AI Package (`ai/`)
Provides OpenAI API integration with three main methods:

#### ChatCompletion
Basic chat completion for simple interactions:
- System and user prompt validation
- Configurable temperature and token limits
- Error handling for empty responses

#### ChatCompletionWithTools
Advanced completion with function calling:
- Tool definitions for AI function selection
- JSON response format enforcement
- Tool choice configuration support

#### ChatCompletionWithHistory
Stateful conversations with history:
- Message history management
- System prompt prepending
- Full JSON response for parsing tool calls

### Tools Package (`tools/`)
Defines the available operations for AI function calling:

#### Tool Definition System
Each tool includes:
- Name and description for AI understanding
- Parameter definitions with types and descriptions
- Required vs optional parameter specifications
- Usage examples and format guidance

#### Available Tools
Comprehensive set of operations:
- **Tag Operations**: ListTags, TagCount, GetTag, ListRootTags
- **Metadata**: UniqueTagTypes, UniqueContextTypes  
- **User Management**: UserCount, GetUser
- **DevOps**: Deploy, Rollback, ListDeployments, GetDeploymentStatus
- **Fallback**: Unknown for unmatched requests

#### Parameter System
Rich parameter definitions with:
- Type validation (string, boolean, number)
- Format detection (list, json, csv)
- Time-based filtering (days, months, years)
- Friendly app name mapping

### Router Package (`router/`)
Maps tool names to actual handler implementations:

#### Operation Router
- Tool name to handler function mapping
- Context and parameter passing
- Error handling and fallback to unknown handler
- Store dependency injection

#### Handler Pattern
All handlers follow consistent signature:
```go
func(ctx context.Context, store store.Store, params map[string]string) (string, error)
```

#### Response Processing
Handlers return formatted JSON responses with:
- Data payload in requested format
- Descriptive messages
- Applied filters information
- Content type metadata

### Formatting Package (`formatting/`)
Standardizes response formats across all handlers:

#### Format Detection
Automatic format selection from parameters:
- "csv", "spreadsheet", "excel" → CSV format
- "json", "data", "api" → JSON format  
- Default or "list" → Human-readable list

#### Formatter System
Type-specific formatters for different data types:
- TagFormatter for tag data
- UserFormatter for user data
- Generic formatters for counts and status

#### Response Structure
Standardized API response format:
```go
type APIResponse struct {
    Type        string
    Data        interface{}
    Message     string
    ContentType string
    Filters     map[string]string
}
```

### Web Router Package (`webrouter/`)
Serves static web content and health endpoints:
- Static file serving (CSS, images, favicon)
- Health check endpoints
- 404 handling with custom templates
- Environment-based configuration

## Common Patterns

### Error Propagation
Consistent error handling across all packages:
- Custom error types in internal/errors
- gRPC status code mapping
- User-friendly error messages via middleware

### Context Usage
All operations accept and propagate context:
- Request cancellation support
- Timeout handling
- Trace information passing

### Parameter Handling
Standardized parameter processing:
- String map interface for all handlers
- Type conversion utilities
- Validation and default value handling

### Interface Abstraction
Clean interface definitions:
- AiClient interface for testability
- Router interface for flexibility
- Store interfaces for data access

## Key Files

### ai/ai.go
OpenAI integration with three completion methods and request building utilities.

### tools/tools.go
Comprehensive tool definitions with parameters and AI guidance.

### router/router.go
Operation routing with handler mapping and error handling.

### router/handlers/
Individual handler implementations for each tool:
- handlelisttags.go, handletagcount.go, etc.
- Consistent parameter processing and response formatting
- Store integration and error handling

### router/formatting/
Response formatting utilities:
- Format detection and conversion
- Type-specific formatters
- Standardized response structures

### webrouter/
HTTP web server for static content and health checks.

## Gotchas & Non-Obvious Behaviors

1. **Tool Choice Configuration**: AI function calling requires specific tool choice setup to force tool selection

2. **Format Detection Priority**: Format parameter takes precedence over natural language format requests

3. **Parameter Case Sensitivity**: Tool parameters are case-sensitive and must match exactly

4. **Context Propagation**: All handlers must accept context even if not used for cancellation support

5. **JSON Response Wrapping**: All responses are JSON-wrapped even for plain text to maintain API consistency

6. **Tool Name Mapping**: Router uses exact string matching for tool names - typos cause fallback to Unknown

7. **History Management**: Conversation history is managed in the chat service, not in the AI client

8. **Static Content Serving**: Web router serves both API endpoints and static files from the same server

9. **Error Message Transformation**: Error messages are transformed by middleware for user-friendliness

10. **App Name Resolution**: DevOps tools support both DigitalOcean app IDs and friendly names

## Dependencies

### External Libraries
- `github.com/sashabaranov/go-openai` - OpenAI API client
- Standard library packages for HTTP, JSON, and string processing

### Internal Dependencies
- `internal/store` - Data access layer for all handlers
- `internal/errors` - Custom error types
- `api/v1/shared` - Protobuf types for data structures

## Configuration
- **OpenAI API Key**: Required for AI functionality
- **OpenAI Model**: Configurable model selection (GPT-4, etc.)
- **Tool Definitions**: Centralized in tools package
- **Handler Mapping**: Defined in router setup

## Performance Considerations
- AI requests have configurable timeouts
- Response caching not implemented (each request hits AI)
- Tool parameter validation happens before AI calls
- Static file serving uses standard HTTP file server