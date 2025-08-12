package kpi

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lucsky/cuid"
)

// SqlKPIStore implements KPIStore using PostgreSQL
type SqlKPIStore struct {
	db *sql.DB
}

// NewSqlKPIStore creates a new SQL-based KPI store
func NewSqlKPIStore(ctx context.Context, dbURL string) (KPIStore, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &SqlKPIStore{
		db: db,
	}, nil
}

// ExecuteTimeStatsProcedure runs the calculate_time_stats_by_group procedure
func (s *SqlKPIStore) ExecuteTimeStatsProcedure(ctx context.Context, group string) (string, error) {
	executionID := cuid.New()
	now := time.Now()
	
	// Insert job record in database
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "Job" (id, type, status, description, "startedAt", metadata, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, executionID, "KPI", "Running", fmt.Sprintf("Calculate time stats for %s", group), now, fmt.Sprintf(`{"group": "%s", "procedure": "calculate_time_stats_by_group"}`, group), now, now)
	
	if err != nil {
		return "", fmt.Errorf("failed to create job record: %w", err)
	}
	
	// Run in background since these take a long time
	go s.runProcedureAsync(executionID, "SELECT calculate_time_stats_by_group($1)", group)
	
	return executionID, nil
}

// ExecuteUpdateStatsProcedure runs the update_calculated_stats_by_group procedure
func (s *SqlKPIStore) ExecuteUpdateStatsProcedure(ctx context.Context, group string) (string, error) {
	executionID := cuid.New()
	now := time.Now()
	
	// Insert job record in database
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "Job" (id, type, status, description, "startedAt", metadata, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, executionID, "KPI", "Running", fmt.Sprintf("Update calculated stats for %s", group), now, fmt.Sprintf(`{"group": "%s", "procedure": "update_calculated_stats_by_group"}`, group), now, now)
	
	if err != nil {
		return "", fmt.Errorf("failed to create job record: %w", err)
	}
	
	// Run in background since these take a long time
	go s.runProcedureAsync(executionID, "SELECT update_calculated_stats_by_group($1)", group)
	
	return executionID, nil
}

// runProcedureAsync executes a procedure in the background and updates status
func (s *SqlKPIStore) runProcedureAsync(executionID, query string, args ...interface{}) {
	// Use a long timeout context for these heavy procedures
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	
	startTime := time.Now()
	
	// Execute the procedure
	result, err := s.db.ExecContext(ctx, query, args...)
	
	// Update job record in database
	now := time.Now()
	duration := int(now.Sub(startTime).Seconds())
	
	if err != nil {
		// Update as failed
		_, updateErr := s.db.ExecContext(context.Background(), `
			UPDATE "Job" 
			SET status = $1, "completedAt" = $2, "durationSeconds" = $3, "errorMessge" = $4, "updatedAt" = $5
			WHERE id = $6
		`, "Failed", now, duration, err.Error(), now, executionID)
		
		if updateErr != nil {
			fmt.Printf("Failed to update job %s as failed: %v\n", executionID, updateErr)
		}
	} else {
		// Update as completed - preserve original metadata and add completion info
		var completionMetadata string
		if result != nil {
			rowsAffected, _ := result.RowsAffected()
			// Get the original metadata first
			var originalMetadata sql.NullString
			s.db.QueryRowContext(context.Background(), `SELECT metadata FROM "Job" WHERE id = $1`, executionID).Scan(&originalMetadata)
			
			if originalMetadata.Valid && originalMetadata.String != "" {
				// Merge completion data into original metadata (simple approach)
				originalJson := originalMetadata.String
				if originalJson == "{}" {
					completionMetadata = fmt.Sprintf(`{"rowsAffected": %d}`, rowsAffected)
				} else {
					// Insert rowsAffected into existing JSON (simple string manipulation)
					completionMetadata = strings.Replace(originalJson, "}", fmt.Sprintf(`, "rowsAffected": %d}`, rowsAffected), 1)
				}
			} else {
				completionMetadata = fmt.Sprintf(`{"rowsAffected": %d}`, rowsAffected)
			}
		} else {
			// No result, just preserve original metadata
			var originalMetadata sql.NullString
			s.db.QueryRowContext(context.Background(), `SELECT metadata FROM "Job" WHERE id = $1`, executionID).Scan(&originalMetadata)
			if originalMetadata.Valid {
				completionMetadata = originalMetadata.String
			} else {
				completionMetadata = `{}`
			}
		}
		
		_, updateErr := s.db.ExecContext(context.Background(), `
			UPDATE "Job" 
			SET status = $1, "completedAt" = $2, "durationSeconds" = $3, metadata = $4, "updatedAt" = $5
			WHERE id = $6
		`, "Completed", now, duration, completionMetadata, now, executionID)
		
		if updateErr != nil {
			fmt.Printf("Failed to update job %s as completed: %v\n", executionID, updateErr)
		}
	}
}

// GetExecutionStatus retrieves the status of a running procedure
func (s *SqlKPIStore) GetExecutionStatus(ctx context.Context, executionID string) (*ExecutionStatus, error) {
	var execution ExecutionStatus
	var metadata sql.NullString
	var description sql.NullString
	var startedAt sql.NullTime
	var completedAt sql.NullTime
	var progress sql.NullInt64
	var durationSeconds sql.NullInt64
	var errorMessage sql.NullString
	
	err := s.db.QueryRowContext(ctx, `
		SELECT id, type, description, status, "startedAt", "completedAt", 
		       progress, "durationSeconds", "errorMessge", metadata
		FROM "Job" 
		WHERE id = $1 AND type = 'KPI'
	`, executionID).Scan(
		&execution.ID, &execution.Type, &description, &execution.Status,
		&startedAt, &completedAt, &progress, &durationSeconds, &errorMessage, &metadata,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("execution %s not found", executionID)
		}
		return nil, fmt.Errorf("failed to get execution status: %w", err)
	}
	
	// Extract group from metadata if available
	if metadata.Valid && metadata.String != "" {
		// Simple JSON parsing for group field
		if len(metadata.String) > 0 {
			// Extract group from JSON metadata (basic parsing)
			metadataStr := metadata.String
			if idx := strings.Index(metadataStr, `"group": "`); idx >= 0 {
				start := idx + len(`"group": "`)
				if end := strings.Index(metadataStr[start:], `"`); end >= 0 {
					execution.Group = metadataStr[start : start+end]
				}
			}
		}
	}
	
	// Set optional fields
	if description.Valid {
		execution.Description = &description.String
	}
	if startedAt.Valid {
		execution.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		execution.CompletedAt = &completedAt.Time
	}
	if progress.Valid {
		progressInt := int(progress.Int64)
		execution.Progress = &progressInt
	}
	if durationSeconds.Valid {
		durationInt := int(durationSeconds.Int64)
		execution.DurationSeconds = &durationInt
	}
	if errorMessage.Valid {
		execution.Error = &errorMessage.String
	}
	
	return &execution, nil
}

// ListRunningExecutions lists all currently running KPI procedures
func (s *SqlKPIStore) ListRunningExecutions(ctx context.Context) ([]*ExecutionStatus, error) {
	return s.listExecutionsByStatus(ctx, "Running")
}

// ListRecentExecutions lists recent executions (running and completed)
func (s *SqlKPIStore) ListRecentExecutions(ctx context.Context) ([]*ExecutionStatus, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, description, status, "startedAt", "completedAt", 
		       progress, "durationSeconds", "errorMessge", metadata
		FROM "Job" 
		WHERE type = 'KPI'
		ORDER BY "createdAt" DESC 
		LIMIT 50
	`)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list recent executions: %w", err)
	}
	defer rows.Close()
	
	return s.scanExecutions(rows)
}

// listExecutionsByStatus is a helper to list executions by status
func (s *SqlKPIStore) listExecutionsByStatus(ctx context.Context, status string) ([]*ExecutionStatus, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, description, status, "startedAt", "completedAt", 
		       progress, "durationSeconds", "errorMessge", metadata
		FROM "Job" 
		WHERE type = 'KPI' AND status = $1
		ORDER BY "createdAt" DESC
	`, status)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list executions by status: %w", err)
	}
	defer rows.Close()
	
	return s.scanExecutions(rows)
}

// scanExecutions is a helper to scan rows into ExecutionStatus structs
func (s *SqlKPIStore) scanExecutions(rows *sql.Rows) ([]*ExecutionStatus, error) {
	var executions []*ExecutionStatus
	
	for rows.Next() {
		var execution ExecutionStatus
		var metadata sql.NullString
		var description sql.NullString
		var startedAt sql.NullTime
		var completedAt sql.NullTime
		var progress sql.NullInt64
		var durationSeconds sql.NullInt64
		var errorMessage sql.NullString
		
		err := rows.Scan(
			&execution.ID, &execution.Type, &description, &execution.Status,
			&startedAt, &completedAt, &progress, &durationSeconds, &errorMessage, &metadata,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan execution: %w", err)
		}
		
		// Extract group from metadata if available
		if metadata.Valid && metadata.String != "" {
			// Simple JSON parsing for group field
			metadataStr := metadata.String
			if idx := strings.Index(metadataStr, `"group": "`); idx >= 0 {
				start := idx + len(`"group": "`)
				if end := strings.Index(metadataStr[start:], `"`); end >= 0 {
					execution.Group = metadataStr[start : start+end]
				}
			}
		}
		
		// Set optional fields
		if description.Valid {
			execution.Description = &description.String
		}
		if startedAt.Valid {
			execution.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			execution.CompletedAt = &completedAt.Time
		}
		if progress.Valid {
			progressInt := int(progress.Int64)
			execution.Progress = &progressInt
		}
		if durationSeconds.Valid {
			durationInt := int(durationSeconds.Int64)
			execution.DurationSeconds = &durationInt
		}
		if errorMessage.Valid {
			execution.Error = &errorMessage.String
		}
		
		executions = append(executions, &execution)
	}
	
	return executions, rows.Err()
}

// CleanupOldExecutions removes completed executions older than the specified duration
func (s *SqlKPIStore) CleanupOldExecutions(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM "Job" 
		WHERE type = 'KPI' 
		  AND status IN ('Completed', 'Failed', 'Cancelled')
		  AND "completedAt" < $1
	`, cutoff)
	
	return err
}