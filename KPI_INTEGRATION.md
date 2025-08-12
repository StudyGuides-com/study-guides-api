# KPI Integration with MCP

## Overview
The KPI (Key Performance Indicator) system has been integrated with the MCP (Model Context Protocol) to allow natural language execution of long-running SQL procedures through the chat interface. Since these procedures can take a long time to run, they execute asynchronously in the background.

## Architecture

### Components
1. **KPI Store** (`internal/store/kpi/`)
   - `SqlKPIStore`: Executes SQL procedures asynchronously
   - Tracks execution status in memory
   - 30-minute timeout for long-running procedures

2. **MCP Adapter** (`internal/mcp/kpi/`)
   - `KPIRepositoryAdapter`: Bridges KPI store to MCP interface
   - Monitors execution status in background
   - Provides status tracking and cancellation

3. **Chat Integration**
   - Registered with ChatService
   - AI can identify KPI requests through natural language
   - Returns immediate response with execution ID

## Available KPI Groups

### Time-based Statistics
- **MonthlyInteractions**: Uses `calculate_time_stats_by_group` procedure
  - Calculates time-based interaction statistics
  - Natural language: "calculate monthly interactions", "monthly stats"

### Update Statistics Groups
All use `update_calculated_stats_by_group` procedure:

- **Tags**: Tag-related statistics
- **TagTypes**: Tag type statistics  
- **Reports**: Report statistics
- **Topics**: Topic statistics
- **MissingData**: Missing data identification
- **Ratings**: Rating statistics
- **Questions**: Question statistics
- **Users**: User statistics
- **UserContent**: User content statistics
- **Contacts**: Contact statistics

## Usage Examples

### Through Chat Interface

```
User: "run all KPIs"
Bot: Starts all 11 KPI calculations in background

User: "calculate monthly interactions"
Bot: Starts MonthlyInteractions calculation

User: "update user stats"
Bot: Starts Users statistics update

User: "check running KPIs"
Bot: Returns list of currently running KPI executions

User: "how many KPIs are running?"
Bot: Returns count of running executions
```

### Testing
Run the test script:
```bash
go run cmd/test-kpi-operations/main.go
```

## Implementation Details

### Asynchronous Execution
- Procedures run in goroutines with 30-minute timeout
- Status tracked in memory (could be moved to Redis/DB)
- Background monitoring updates status every 5 seconds

### Status Tracking
Each execution tracks:
- ID (CUID)
- Group name
- Status (pending/running/complete/failed)
- Start time
- Completion time
- Duration
- Error message (if failed)

### Error Handling
- Database connection errors
- Procedure execution failures
- Timeout after 30 minutes
- User cancellation support

## Future Enhancements

1. **Persistent Status Storage**
   - Move execution tracking to database or Redis
   - Survive server restarts

2. **Progress Tracking**
   - Add progress percentage for long operations
   - Stream updates to client

3. **Scheduling**
   - Cron-based automatic execution
   - Dependency management between KPIs

4. **Notifications**
   - Slack/email notifications on completion
   - Webhook support for external systems

5. **Result Storage**
   - Store KPI results for historical analysis
   - Trend visualization

## Troubleshooting

### KPIs Not Running
- Check database connectivity
- Verify procedures exist in database
- Check server logs for goroutine errors

### Status Not Updating
- Background monitor runs every 5 seconds
- Check if execution ID exists in tracking map
- Verify no panic in goroutine

### Timeouts
- Default timeout is 30 minutes
- Can be adjusted in `sqlkpistore.go`
- Consider breaking large procedures into smaller chunks