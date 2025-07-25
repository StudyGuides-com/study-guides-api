# internal/store Package

## Purpose and Responsibility
The `internal/store` package implements the data access layer for the Study Guides API. It provides a unified interface for all data operations while abstracting away the specifics of different storage backends (PostgreSQL, Algolia, DigitalOcean).

## Key Architectural Decisions

### Repository Pattern with Interface Segregation
Each domain has its own store interface and implementation:
- **Domain-specific interfaces**: `TagStore`, `UserStore`, `QuestionStore`, etc.
- **Implementation separation**: Interfaces in `{domain}.go`, implementations in `sql{domain}store.go`
- **External service integration**: Special implementations for Algolia (`algoliasearch.go`) and DigitalOcean (`digitaloceandevops.go`)

### Aggregating Store Pattern
The main `Store` interface aggregates all domain stores:
```go
type Store interface {
    SearchStore() search.SearchStore
    TagStore() tag.TagStore
    UserStore() user.UserStore
    // ... other stores
}
```

### Multi-Backend Architecture
Different operations use different storage backends:
- **PostgreSQL**: Primary data storage (tags, users, questions, admin operations)
- **Algolia**: Search functionality with full-text search and filtering
- **DigitalOcean**: DevOps operations (deployment, rollback)
- **Separate Roland database**: AI/chat functionality isolation

## Implementation Details

### SQL Implementation Pattern
Most stores follow consistent SQL patterns:
- **Connection pooling**: Use `pgxpool.Pool` for PostgreSQL connections
- **Structured scanning**: Use `github.com/georgysavva/scany/v2/pgxscan` for result mapping
- **Error handling**: Convert database errors to gRPC status codes
- **Upsert operations**: Use `ON CONFLICT` for data deduplication

### Admin Store Complexity
The `SqlAdminStore` is the most complex implementation, providing:
- **Complex hierarchical queries**: Recursive CTEs for tag trees and ancestry
- **Bulk operations**: Import/export functionality with GOB serialization
- **Tree operations**: Tag hierarchy traversal with cycle detection
- **Content management**: Tag, passage, and question CRUD with metadata

### Search Store Integration
The `AlgoliaStore` provides sophisticated search capabilities:
- **Role-based filtering**: Admin users see all results, regular users filtered by ownership
- **Dynamic filter building**: Context, type, and permission-based filtering
- **Result transformation**: Convert Algolia hits to protobuf messages
- **Index management**: Multi-index support (tags, users)

### DevOps Store Abstraction
The `DigitalOceanDevopsStore` wraps the DigitalOcean API:
- **Lazy client initialization**: Client created on first use
- **Status mapping**: Convert DigitalOcean deployment phases to internal enums
- **Error wrapping**: Consistent error handling across API calls

## Common Patterns

### Constructor Pattern
All stores follow the same initialization pattern:
```go
func NewSqlDomainStore(ctx context.Context, dbURL string) (DomainStore, error) {
    // Connection setup
    // Interface implementation return
}
```

### Context Propagation
All store operations accept `context.Context` for:
- Request cancellation
- Database transaction management
- Timeout handling

### Error Handling Strategy
Consistent error handling across all stores:
- Database errors → `codes.Internal`
- Not found errors → `codes.NotFound`
- Permission errors → `codes.PermissionDenied`
- Validation errors → `codes.InvalidArgument`

### Metadata Management
Complex metadata handling in admin operations:
- JSON serialization/deserialization
- Metadata merging for updates
- Timestamp management with protobuf types

## Key Files

### store.go
- Main `Store` interface definition
- Store aggregation implementation
- Initialization coordination for all stores
- Environment variable dependency management

### admin/sqladminstore.go
Most complex store implementation featuring:
- Recursive tag hierarchy operations
- Complex import/export with GOB serialization
- Tree traversal with cycle detection
- Bulk operations for content management

### search/algoliasearch.go
Search functionality with:
- Role-based access control
- Dynamic filter construction
- Multi-index support
- Result transformation from Algolia format

### devops/digitaloceandevops.go
DevOps operations abstraction:
- DigitalOcean API integration
- Deployment lifecycle management
- Status mapping and error handling

### Domain-specific stores
Each domain (tag, user, question, etc.) follows the same pattern:
- Interface definition in `{domain}.go`
- SQL implementation in `sql{domain}store.go`
- Standard CRUD operations with domain-specific extensions

## Gotchas & Non-Obvious Behaviors

1. **Multiple Database Connections**: Main database and separate Roland database require different connection strings

2. **Admin Store Complexity**: The admin store handles the most complex operations including recursive tree queries and bulk imports

3. **Search Permission Model**: Search results are filtered based on user roles and ownership, not just authentication

4. **Lazy Client Initialization**: External service clients (DigitalOcean) are created on first use, not during store initialization

5. **Metadata Merging**: Tag metadata updates merge with existing metadata rather than replacing it

6. **Cycle Detection**: Tag hierarchy operations include cycle detection to prevent infinite loops

7. **GOB Serialization**: Import functionality uses Go's GOB encoding for complex data structures

8. **Upsert Strategy**: Most operations use PostgreSQL's `ON CONFLICT` for deduplication based on hash or ID

9. **Index Cache Management**: Admin store manages an index cache for performance optimization

10. **Tree Deletion Order**: Tag tree deletion works from leaf nodes up to prevent foreign key violations

## Dependencies

### External Libraries
- `github.com/jackc/pgx/v5` - PostgreSQL driver and connection pooling
- `github.com/georgysavva/scany/v2/pgxscan` - Structured result scanning
- `github.com/algolia/algoliasearch-client-go/v3` - Algolia search client
- `github.com/digitalocean/godo` - DigitalOcean API client
- `google.golang.org/grpc` - gRPC status codes and error handling
- `google.golang.org/protobuf` - Protocol buffer types and timestamps

### Internal Dependencies
- `api/v1/shared` - Shared protobuf types
- `api/v1/{domain}` - Domain-specific protobuf types
- `internal/utils` - Utility functions (CUID generation, type conversion)

## Database Schema Assumptions
The store implementations assume specific PostgreSQL schema:
- Tables use quoted identifiers (e.g., `"Tag"`, `"User"`)
- Complex relationships between tags, questions, and users
- Metadata stored as JSONB
- Recursive functions for user data deletion
- Index cache table for performance optimization

## Performance Considerations
- Connection pooling for database operations
- Recursive query depth limits (maxTagDepth = 5)
- Batch operations for bulk imports
- Index caching for search performance
- Lazy initialization of external service clients