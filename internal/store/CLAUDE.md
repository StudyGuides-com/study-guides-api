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
    QuestionStore() question.QuestionStore
    InteractionStore() interaction.InteractionStore
    RolandStore() roland.RolandStore
    DevopsStore() devops.DevopsStore
    KPIStore() kpi.KPIStore
    IndexingStore() indexing.IndexingStore
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

### KPI Store
The `KPIStore` manages performance metrics and analytics:
- **Stored procedures execution**: Runs `calculate_time_stats_by_group` and `update_calculated_stats_by_group` procedures
- **Job tracking**: Monitors execution status of long-running KPI calculations
- **Group-based metrics**: Calculates statistics per user group (e.g., class, school)
- **Execution management**: Tracks running, completed, and failed KPI jobs with metadata

### Indexing Store
The `IndexingStore` manages search index synchronization:
- **Connection pooling**: Uses dedicated `pgxpool.Pool` for concurrent operations
- **Outbox pattern**: Queue-based approach for reliable index updates
- **State tracking**: Maintains index state with hashing for change detection
- **Batch operations**: Supports bulk reindexing of entire object types
- **Filtering support**: SQL-level filtering by TagType and ContextType via `StartIndexingJobWithFilters()`
- **Job management**: Similar to KPI store, tracks indexing job progress
- **Error handling**: Retry limits and automatic removal of failed/non-existent items
- **Algolia integration**: Direct integration with Algolia for index updates
- **Dependency injection**: Requires TagStore for tag hierarchy operations

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
- Main `Store` interface definition with 9 domain stores (AdminStore exists separately)
- Store aggregation implementation
- Initialization coordination for all stores
- Environment variable dependency management
- Connection pool creation for IndexingStore
- Special initialization for stores with dependencies (IndexingStore needs TagStore)

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

### kpi/sqlkpistore.go
KPI metrics and analytics management:
- Execution of stored procedures for time-based statistics
- Job status tracking with metadata storage
- Group-based metric calculation
- Execution history and monitoring

### indexing/sqlindexingstore.go
Search index synchronization management:
- Outbox pattern for incremental change tracking
- Dual-mode indexing:
  - **Incremental mode** (`force=false`): Only queues changed items via `QueueChangedForIndex`
  - **Force rebuild** (`force=true`): Queues all items via `QueueBatchForReindex`
- Change detection includes:
  - Tag updates (`updatedAt` comparison)
  - Access permission changes (TagAccess updates)
  - Ancestry changes (parent tag modifications)
- Index state tracking with SHA-256 content hashing
- Job-based async processing with 30-minute timeout
- Direct Algolia API integration
- Batch processing (100 items at a time)

### Domain-specific stores
Each domain (tag, user, question, interaction, roland) follows the same pattern:
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

11. **KPI Store Procedures**: KPI calculations run as database stored procedures, not application code

12. **Indexing Store Pool**: IndexingStore creates its own connection pool separate from other stores for concurrent operations

13. **Indexing Outbox Pattern**:
    - Outbox tracks incremental changes from CRUD operations
    - `force=false`: Processes only changed items (incremental sync)
    - `force=true`: Rebuilds entire index regardless of changes
    - Failed operations remain in outbox for retry

14. **Indexing Filtering Methods**:
    - `StartIndexingJobWithFilters()`: Creates jobs with TagType/ContextType filters
    - `QueueBatchForReindexWithFilters()`: Queues all matching items for force rebuild
    - `QueueChangedForIndexWithFilters()`: Queues only changed matching items
    - SQL-level filtering using `WHERE type = ANY($1) AND context = ANY($2)`

15. **Indexing Change Detection**: The incremental mode detects changes via:
    - Direct tag updates (compares `updatedAt` timestamps)
    - TagAccess permission changes
    - Parent tag modifications (affects ancestry chain)

16. **Indexing Error Handling**: Recent improvements prevent infinite retry loops:
    - Retry limits: Failed items removed after 10 attempts
    - Error detection: Special handling for "no rows in result set" (non-existent items)
    - Job completion: Fixed metadata query failures that prevented completion

17. **Store Initialization Order**: IndexingStore must be initialized after TagStore due to dependency

18. **Indexing Job Timeout**: Jobs have a 30-minute timeout - large datasets (80k+ tags) may timeout if processing is too slow

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
- **Known schema issue**: Job table has misspelled column `"errorMessge"` (missing 'a') - code correctly uses this spelling

## Performance Considerations
- Connection pooling for database operations
- Recursive query depth limits (maxTagDepth = 5)
- Batch operations for bulk imports
- Index caching for search performance
- Lazy initialization of external service clients