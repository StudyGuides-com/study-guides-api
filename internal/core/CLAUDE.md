# internal/core Package

## Purpose and Responsibility
The `internal/core` package contains shared business logic that can be used by multiple interface layers (gRPC services, MCP adapters, etc.). This layer implements the core domain operations while remaining independent of specific interface concerns.

## Architectural Purpose

### Shared Business Logic Layer
The core layer serves as a central location for business logic that would otherwise be duplicated across multiple interface implementations:

```
MCP Adapter ──┐
              ├── Core Business Service → Store Layer
gRPC Service ─┘
```

This pattern ensures:
- **Single Source of Truth**: Business rules exist in one place
- **DRY Principle**: No duplication of logic between interfaces
- **Consistency**: All interfaces behave identically
- **Testability**: Business logic can be tested independently of interface concerns

## Implementation Details

### `/internal/core/indexing/`

#### BusinessService (`service.go`)
The core indexing business service provides fundamental indexing operations:

**Core Operations:**
- `TriggerIndexing(req)` - Generic object indexing
- `TriggerTagIndexing(req)` - Tag-specific indexing with filtering
- `GetJobStatus(jobID)` - Job status retrieval
- `ListRunningJobs()` - Active job monitoring
- `ListRecentJobs(req)` - Job history queries

**Key Features:**
```go
type BusinessService struct {
    store store.Store
}

// Generic indexing for any object type
func (bs *BusinessService) TriggerIndexing(ctx context.Context, req TriggerIndexingRequest) (*TriggerIndexingResponse, error)

// Tag-specific indexing with filtering capabilities
func (bs *BusinessService) TriggerTagIndexing(ctx context.Context, req TriggerTagIndexingRequest) (*TriggerIndexingResponse, error)
```

#### Flexible Tag Filtering
The `TriggerTagIndexing` method supports flexible filtering combinations:

**Filter Types:**
```go
type TriggerTagIndexingRequest struct {
    Force        bool                     // Force complete reindex vs incremental
    TagTypes     []sharedpb.TagType      // Filter by tag types (optional)
    ContextTypes []sharedpb.ContextType  // Filter by context types (optional)
}
```

**Filter Combinations:**
- **No filters**: `TagTypes=[]` AND `ContextTypes=[]` → Index all tags
- **TagType only**: `TagTypes=[Topic, Category]` → Index specific tag types
- **ContextType only**: `ContextTypes=[DoD, Colleges]` → Index specific contexts
- **Combined**: `TagTypes=[Topic]` AND `ContextTypes=[DoD]` → Tag type AND context

#### Type Definitions
The core service defines framework-agnostic types:

```go
// Request/Response types using native Go types and protobuf enums
type TriggerIndexingRequest struct {
    ObjectType string
    Force      bool
}

type TriggerIndexingResponse struct {
    JobID     string
    Status    string
    Message   string
    StartedAt time.Time
}
```

## Integration Patterns

### gRPC Service Integration
gRPC services use the core service as a dependency and handle type conversion between protobuf and business types. See `internal/services/CLAUDE.md` for implementation details.

### MCP Adapter Integration
MCP adapters also use the core service for natural language triggers. See `internal/mcp/CLAUDE.md` for implementation details.

## Key Design Decisions

### Framework Independence
Core services avoid dependencies on specific frameworks:
- **No gRPC types**: Uses native Go types and standard protobuf enums
- **No MCP types**: Converts from interface-specific types to business types
- **Store delegation**: Delegates data operations to store layer

### Type Conversion Strategy
Each interface layer is responsible for converting between its types and core business types:
- **gRPC**: `protobuf types ↔ business types`
- **MCP**: `MCP filter types ↔ business types`
- **Core**: Only works with business-domain types

### Future Extensibility
The core layer is designed to support additional domains:
- `internal/core/user/` - User management business logic
- `internal/core/question/` - Question domain business logic
- `internal/core/analytics/` - Analytics business logic

## Common Patterns

### Constructor Pattern
All core services follow the same initialization pattern:
```go
func NewBusinessService(store store.Store) *BusinessService {
    return &BusinessService{
        store: store,
    }
}
```

### Error Handling
Core services return domain errors that are later converted by interface layers:
```go
if err != nil {
    return nil, fmt.Errorf("failed to start indexing job: %w", err)
}
```

### Store Delegation
Core services delegate data operations to the store layer:
```go
jobID, err := bs.store.IndexingStore().StartIndexingJob(ctx, objectType, req.Force)
```

## Benefits

### Code Organization
- **Separation of Concerns**: Interface logic vs business logic
- **Reusability**: Same business logic across multiple interfaces
- **Maintainability**: Changes in one place affect all interfaces

### Testing Strategy
- **Unit Testing**: Business logic can be tested independently
- **Interface Testing**: Interface layers can mock business services
- **Integration Testing**: Full stack testing through any interface

### Development Efficiency
- **Parallel Development**: Interface and business logic can be developed separately
- **Consistency**: All interfaces automatically stay in sync
- **Debugging**: Business logic issues are centralized

## Future Enhancements

### Recent Improvements
- **Enhanced Filtering**: ✅ Implemented database-level filtering for TagType and ContextType via `StartIndexingJobWithFilters()`
- **Error Handling**: ✅ Added retry limits and graceful handling of non-existent items

### Planned Improvements
- **Batch Operations**: Support for multiple object processing
- **Validation Layer**: Input validation at business service level
- **Metrics Integration**: Business-level operation metrics

The core business service layer provides a solid foundation for scalable, maintainable business logic that can support multiple interface patterns while maintaining consistency and avoiding code duplication.