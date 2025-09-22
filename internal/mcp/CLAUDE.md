# internal/mcp Package

## Purpose and Responsibility
The `internal/mcp` package implements the Model Context Protocol (MCP) system for the Study Guides API. MCP provides a unified interface for AI-driven operations, replacing the legacy tool routing system with a more flexible repository-based approach.

## Key Architectural Decisions

### Repository Pattern for AI Operations
MCP uses repository adapters to expose domain operations to AI:
- **IndexingRepositoryAdapter**: Manages search index synchronization
- **KPIRepositoryAdapter**: Handles performance metrics operations
- **TagRepositoryAdapter**: Provides tag management operations

### Dynamic Tool Generation
Tools are automatically generated from registered repositories:
- Schema-driven tool definitions
- Parameter validation and type checking
- Consistent tool naming (e.g., `indexing_find`, `tag_find`)
- OpenAI function calling integration

### Natural Language Processing Flow
1. User sends natural language query
2. MCP Processor generates tools from repositories
3. AI classifies intent and selects appropriate tool
4. Tool arguments are parsed and validated
5. Repository method is executed
6. Results are formatted and returned

## Implementation Details

### MCP Processor (`server.go`)
Central orchestrator that manages:
- Repository registration and schema management
- Tool generation for OpenAI function calling
- AI request processing and response handling
- Error handling and logging

### Repository Adapters
Each domain has an adapter implementing standard CRUD operations:
```go
type Repository interface {
    Find(ctx context.Context, filter FilterType) ([]EntityType, error)
    FindByID(ctx context.Context, id string) (*EntityType, error)
    Create(ctx context.Context, entity EntityType) (*EntityType, error)
    Update(ctx context.Context, id string, update UpdateType) (*EntityType, error)
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context, filter FilterType) (int, error)
}
```

### Shared Business Service Integration
The IndexingRepositoryAdapter uses the shared business service for consistency with gRPC interface:
- **Core Integration**: Delegates to `internal/core/indexing.BusinessService`
- **Incremental indexing**: `"index tags"` → `triggerReindex: true, force: false`
- **Force rebuild**: `"force index tags"` → `triggerReindex: true, force: true`
- **Pruning operations**: `"prune index"` → Remove orphaned Algolia objects
- **Status monitoring**: `"check indexing status"` → `status: "running"`
- **Job tracking**: Progress monitoring and error reporting via business service

## Natural Language Triggers

### Indexing Operations
- **"index tags", "reindex tags", "sync tags"** → Incremental indexing
- **"force reindex", "rebuild index"** → Complete rebuild
- **"prune index", "clean index", "remove orphaned objects"** → Pruning operations
- **"prune tags with filters"** → Filtered pruning by TagType/ContextType
- **"check indexing", "indexing status"** → Status check

### KPI Operations  
- **"calculate stats", "run kpi"** → Execute KPI procedures
- **"check kpi status"** → Monitor running jobs

### Tag Operations
- **"list tags", "show tags"** → Query tag repository
- **"count tags"** → Get tag counts with filters

## Common Patterns

### Repository Registration
Each repository is registered with the MCP processor:
```go
mcpProcessor := mcp.NewMCPProcessor(aiClient)
indexingRepo := indexing.NewIndexingRepositoryAdapter(store)
mcpProcessor.Register(indexing.ResourceName, indexingRepo, indexing.GetResourceSchema())
```

### Filter-Based Queries
All repository operations use filter objects for consistency:
```go
type IndexingFilter struct {
    TriggerReindex bool    `json:"triggerReindex,omitempty"`
    ObjectType     *string `json:"objectType,omitempty"`
    Force          *bool   `json:"force,omitempty"`
    Status         *string `json:"status,omitempty"`
}
```

### Schema Definitions
Each repository provides schema information for tool generation:
- Operation descriptions and examples
- Parameter definitions with types
- Tool-specific guidance for AI

## Key Files

### server.go
- MCP Processor implementation
- Repository management and tool generation
- AI request processing pipeline

### tools.go
- Dynamic tool generation from repository schemas
- OpenAI function definition creation
- Parameter validation and type conversion

### handler.go
- Command execution and response formatting
- Error handling and logging
- Repository method dispatch

### indexing/adapter.go
- IndexingRepositoryAdapter implementation
- Indexing job management and status tracking
- Integration with shared business service (`internal/core/indexing.BusinessService`)

### indexing/schema.go
- Indexing operation schema definitions
- Natural language trigger examples
- Parameter descriptions and validation rules

## Gotchas & Non-Obvious Behaviors

1. **Chat Service Integration**: Only the ChatService uses MCP - other services use legacy patterns

2. **Repository State**: Repository adapters are stateless - they delegate to stores for persistence

3. **Tool Name Mapping**: Tool names follow `{resource}_{operation}` pattern (e.g., `indexing_find`)

4. **Filter Parsing**: AI arguments are parsed into filter objects with strict type validation

5. **Async Operations**: Indexing and KPI operations return immediately with job IDs, actual work happens asynchronously

6. **Error Handling**: Repository errors are wrapped and returned as structured responses

7. **Schema Evolution**: Adding new repositories or operations requires schema updates and tool regeneration

## Dependencies

### External Libraries
- `github.com/sashabaranov/go-openai` - OpenAI API client for function calling
- Standard library for JSON parsing and validation

### Internal Dependencies
- `internal/core/indexing` - Shared business service for indexing operations
- `internal/store/indexing` - Indexing store for search operations
- `internal/store/kpi` - KPI store for metrics operations
- `internal/store/tag` - Tag store for tag operations
- `internal/lib/ai` - AI client wrapper

## Performance Considerations
- Tool generation happens once per request (not cached)
- Repository operations are synchronous (except indexing/KPI jobs)
- Large result sets are handled by the underlying stores
- Error responses are structured for AI consumption