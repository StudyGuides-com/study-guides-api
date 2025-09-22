# Pruning Feature Documentation

## Overview

The pruning feature removes orphaned objects from the Algolia search index that no longer exist in the database. This ensures search results remain accurate and prevents users from finding references to deleted content.

## Architecture

### Design Principles

1. **Resource Efficiency**: Constant memory usage (~1MB) regardless of index size
2. **Stream Processing**: Processes objects one at a time without loading entire datasets
3. **Database Friendly**: Uses prepared statements for efficient repeated queries
4. **Scalable**: Designed to handle millions of records without performance degradation
5. **Observable**: Job-based tracking with progress updates

### Implementation Strategy

The pruning operation follows a streaming pattern:

1. **Create Job Record**: Track operation in database
2. **Stream Algolia Objects**: Use `BrowseObjects()` iterator
3. **Check Existence**: Query database for each object ID
4. **Batch Deletions**: Accumulate up to 1000 IDs before deletion
5. **Update Progress**: Periodic job metadata updates

## API Interface

### gRPC Endpoint

```protobuf
rpc PruneIndex(PruneIndexRequest) returns (PruneIndexResponse);

message PruneIndexRequest {
  string object_type = 1;                      // Default: "Tag"
  repeated TagType tag_types = 2;              // Optional filter
  repeated ContextType context_types = 3;      // Optional filter
}

message PruneIndexResponse {
  string job_id = 1;
  string status = 2;
  string message = 3;
  google.protobuf.Timestamp started_at = 4;
}
```

### Natural Language Interface (MCP)

Supported commands:
- `"prune tags"` - Remove orphaned tag objects
- `"clean algolia index"` - Clean up search index
- `"remove orphaned tags"` - Delete tags not in database
- `"clean search index"` - General index cleanup

## Usage Examples

### gRPC Client

```go
// Simple pruning - all tags
req := &indexingv1.PruneIndexRequest{
    ObjectType: "Tag",
}

// Filtered pruning - specific tag types
req := &indexingv1.PruneIndexRequest{
    ObjectType: "Tag",
    TagTypes: []sharedpb.TagType{
        sharedpb.TagType_Topic,
        sharedpb.TagType_Category,
    },
}

// Filtered by context
req := &indexingv1.PruneIndexRequest{
    ObjectType: "Tag",
    ContextTypes: []sharedpb.ContextType{
        sharedpb.ContextType_Colleges,
    },
}

resp, err := client.PruneIndex(ctx, req)
// Returns job ID for tracking
```

### MCP/Natural Language

```
User: "prune tags"
Response: "Started Tag pruning job (job ID: xyz123)"

User: "clean algolia"
Response: "Started Tag pruning job to remove orphaned objects from search index"

User: "remove orphaned tags from search"
Response: "Started pruning job..."
```

## Implementation Details

### Core Algorithm

```go
func runPruningAsync(jobID, objectType string, filters...) {
    // 1. Prepare existence check query
    stmt := db.Prepare("SELECT 1 FROM Tag WHERE id = $1 AND filters...")

    // 2. Stream through Algolia
    it := algoliaIndex.BrowseObjects()
    deleteBuffer := make([]string, 0, 1000)

    for {
        record := it.Next()

        // 3. Check database
        if !existsInDB(record.ObjectID) {
            deleteBuffer = append(deleteBuffer, record.ObjectID)

            // 4. Batch delete when full
            if len(deleteBuffer) >= 1000 {
                algoliaIndex.DeleteObjects(deleteBuffer)
                updateProgress(jobID, deletedCount)
                deleteBuffer = deleteBuffer[:0]
            }
        }
    }

    // 5. Final cleanup and completion
    algoliaIndex.DeleteObjects(deleteBuffer)
    markJobComplete(jobID)
}
```

### Resource Usage

#### Memory Footprint
- **Algolia Iterator**: ~500 bytes (current record only)
- **Delete Buffer**: Max 50KB (1000 IDs × 50 bytes)
- **Prepared Statement**: Single connection
- **Total**: < 1MB constant regardless of scale

#### Performance Characteristics

For different dataset sizes:

| Records | Memory | Est. Time | DB Queries/sec |
|---------|--------|-----------|----------------|
| 80,000  | ~1MB   | 2-4 min   | 100            |
| 1M      | ~1MB   | 15-20 min | 100            |
| 6M      | ~1MB   | 1-2 hours | 100            |

### Database Impact

- **Query Type**: Primary key lookups only
- **Connection**: Single connection from pool
- **Blocking**: No locks, read-only queries
- **Index Usage**: Efficient PK index lookups

### Safety Features

1. **28-minute timeout**: Prevents runaway jobs
2. **Progress tracking**: Updates job metadata periodically
3. **Error recovery**: Graceful handling of Algolia API failures
4. **Idempotent**: Can be run multiple times safely
5. **Admin-only**: Requires USER_ROLE_ADMIN

## Filtering Options

### Tag Type Filtering

Filter pruning to specific tag types:

```go
TagTypes: []sharedpb.TagType{
    TagType_Topic,
    TagType_Category,
    TagType_Course,
}
```

### Context Type Filtering

Filter by context:

```go
ContextTypes: []sharedpb.ContextType{
    ContextType_Colleges,
    ContextType_DoD,
    ContextType_Encyclopedia,
}
```

### Combined Filtering

Both filters can be used together for precise control:

```go
// Only prune Topic tags in DoD context
TagTypes: []TagType{TagType_Topic}
ContextTypes: []ContextType{ContextType_DoD}
```

## Job Monitoring

### Job Metadata

Track pruning progress via job metadata:

```json
{
  "objectType": "Tag",
  "pruneCount": 1523,
  "tagTypes": [5, 6],
  "contextTypes": [1]
}
```

### Status Checking

```go
// Check specific job
status, err := client.GetJobStatus(ctx, &GetJobStatusRequest{
    JobId: "xyz123",
})

// List running jobs
jobs, err := client.ListRunningJobs(ctx, &ListRunningJobsRequest{})
```

## Error Handling

### Failure Scenarios

1. **Database Connection Loss**: Job marked as failed
2. **Algolia API Errors**: Retries with logged failures
3. **Timeout (28 min)**: Partial completion, can resume
4. **Invalid Filters**: Immediate job failure

### Recovery Strategy

- Jobs are idempotent - safe to retry
- Partial completions are valid (orphans removed)
- Failed deletions logged but don't stop processing

## Future Enhancements

### Planned Improvements

1. **Batch Existence Checking**: Check 100 IDs at once for faster processing
2. **Parallel Processing**: Multiple workers for large datasets
3. **Resume Capability**: Continue from last processed ID after failure
4. **Additional Object Types**: Support User, Question, FAQ objects
5. **Metrics Collection**: Track pruning statistics over time

### Optimization Opportunities

```go
// Future: Batch existence checking
SELECT id FROM Tag WHERE id = ANY($1)
// Process results to find missing IDs

// Future: Parallel workers
for i := 0; i < workers; i++ {
    go processPruningBatch(startIdx, endIdx)
}
```

## Security Considerations

- **Admin Only**: All pruning operations require admin role
- **No User Data Exposure**: Only processes object IDs
- **Audit Trail**: Job records provide operation history
- **Rate Limiting**: Natural throttling via sequential processing

## Operational Notes

### When to Run Pruning

- **After bulk deletions**: Clean up after removing content
- **Regular maintenance**: Weekly/monthly cleanup schedule
- **Before reindexing**: Ensure clean slate for full rebuild
- **Data inconsistencies**: When search returns deleted items

### Best Practices

1. **Monitor job progress**: Check status during long operations
2. **Run during low traffic**: Minimize impact on search performance
3. **Verify results**: Check pruneCount in job metadata
4. **Test with filters first**: Use specific filters before full pruning
5. **Keep job history**: Track pruning patterns over time

## Troubleshooting

### Common Issues

**Job times out after 28 minutes**
- Normal for very large datasets
- Check pruneCount to see progress
- Run again to continue (idempotent)

**High pruneCount unexpected**
- Verify recent deletion operations
- Check for data sync issues
- Review TagAccess changes

**Zero items pruned**
- Index may be in sync
- Check filters are correct
- Verify object type exists

### Debug Commands

```bash
# Check job status
grpcurl -d '{"job_id": "xyz123"}' \
  localhost:50051 indexing.v1.IndexingService/GetJobStatus

# View job metadata
SELECT metadata FROM "Job" WHERE id = 'xyz123';

# Count Algolia vs DB records
SELECT COUNT(*) FROM "Tag"; -- Database count
# Compare with Algolia dashboard count
```

## Code Locations

- **Proto Definition**: `/api/v1/indexing/indexing.proto`
- **Business Logic**: `/internal/core/indexing/service.go`
- **Store Implementation**: `/internal/store/indexing/sqlindexingstore.go`
- **gRPC Handler**: `/internal/services/indexing/indexing.go`
- **MCP Adapter**: `/internal/mcp/indexing/adapter.go`
- **MCP Schema**: `/internal/mcp/indexing/schema.go`