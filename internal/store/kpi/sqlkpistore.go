package kpi

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lucsky/cuid"
)

// SqlKPIStore implements KPIStore using PostgreSQL
type SqlKPIStore struct {
	db *sql.DB
	
	// In-memory tracking of executions (could be moved to Redis/DB table)
	mu         sync.RWMutex
	executions map[string]*ExecutionStatus
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
		db:         db,
		executions: make(map[string]*ExecutionStatus),
	}, nil
}

// ExecuteTimeStatsProcedure runs the calculate_time_stats_by_group procedure
func (s *SqlKPIStore) ExecuteTimeStatsProcedure(ctx context.Context, group string) error {
	executionID := cuid.New()
	
	// Track execution
	s.mu.Lock()
	s.executions[executionID] = &ExecutionStatus{
		ID:            executionID,
		ProcedureName: "calculate_time_stats_by_group",
		Group:         group,
		Status:        "running",
		StartedAt:     time.Now(),
	}
	s.mu.Unlock()
	
	// Run in background since these take a long time
	go s.runProcedureAsync(executionID, "SELECT calculate_time_stats_by_group($1)", group)
	
	return nil
}

// ExecuteUpdateStatsProcedure runs the update_calculated_stats_by_group procedure
func (s *SqlKPIStore) ExecuteUpdateStatsProcedure(ctx context.Context, group string) error {
	executionID := cuid.New()
	
	// Track execution
	s.mu.Lock()
	s.executions[executionID] = &ExecutionStatus{
		ID:            executionID,
		ProcedureName: "update_calculated_stats_by_group",
		Group:         group,
		Status:        "running",
		StartedAt:     time.Now(),
	}
	s.mu.Unlock()
	
	// Run in background since these take a long time
	go s.runProcedureAsync(executionID, "SELECT update_calculated_stats_by_group($1)", group)
	
	return nil
}

// runProcedureAsync executes a procedure in the background and updates status
func (s *SqlKPIStore) runProcedureAsync(executionID, query string, args ...interface{}) {
	// Use a long timeout context for these heavy procedures
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	
	// Execute the procedure
	result, err := s.db.ExecContext(ctx, query, args...)
	
	// Update execution status
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if execution, exists := s.executions[executionID]; exists {
		now := time.Now()
		execution.CompletedAt = &now
		
		if err != nil {
			execution.Status = "failed"
			errStr := err.Error()
			execution.Error = &errStr
		} else {
			execution.Status = "complete"
			if result != nil {
				rowsAffected, _ := result.RowsAffected()
				execution.RowsAffected = &rowsAffected
			}
		}
	}
}

// GetExecutionStatus retrieves the status of a running procedure
func (s *SqlKPIStore) GetExecutionStatus(ctx context.Context, executionID string) (*ExecutionStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if execution, exists := s.executions[executionID]; exists {
		// Return a copy to avoid race conditions
		status := *execution
		return &status, nil
	}
	
	return nil, fmt.Errorf("execution %s not found", executionID)
}

// ListRunningExecutions lists all currently running KPI procedures
func (s *SqlKPIStore) ListRunningExecutions(ctx context.Context) ([]*ExecutionStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var running []*ExecutionStatus
	for _, execution := range s.executions {
		if execution.Status == "running" {
			// Return copies to avoid race conditions
			status := *execution
			running = append(running, &status)
		}
	}
	
	return running, nil
}

// CleanupOldExecutions removes completed executions older than the specified duration
func (s *SqlKPIStore) CleanupOldExecutions(olderThan time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	cutoff := time.Now().Add(-olderThan)
	for id, execution := range s.executions {
		if execution.Status != "running" && execution.CompletedAt != nil && execution.CompletedAt.Before(cutoff) {
			delete(s.executions, id)
		}
	}
}