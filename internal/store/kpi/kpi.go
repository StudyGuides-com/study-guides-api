package kpi

import (
	"context"
	"time"
)

// KPIStore defines the interface for KPI database operations
type KPIStore interface {
	// ExecuteTimeStatsProcedure runs the calculate_time_stats_by_group procedure
	// Returns the execution ID of the created job
	ExecuteTimeStatsProcedure(ctx context.Context, group string) (string, error)
	
	// ExecuteUpdateStatsProcedure runs the update_calculated_stats_by_group procedure
	// Returns the execution ID of the created job
	ExecuteUpdateStatsProcedure(ctx context.Context, group string) (string, error)
	
	// GetExecutionStatus retrieves the status of a running procedure (if trackable)
	GetExecutionStatus(ctx context.Context, executionID string) (*ExecutionStatus, error)
	
	// ListRunningExecutions lists all currently running KPI procedures
	ListRunningExecutions(ctx context.Context) ([]*ExecutionStatus, error)
	
	// ListRecentExecutions lists recent executions (running and completed)
	ListRecentExecutions(ctx context.Context) ([]*ExecutionStatus, error)
}

// ExecutionStatus represents the status of a procedure execution
type ExecutionStatus struct {
	ID               string
	Type             string    // job type (e.g., "KPI")
	Description      *string   // human-readable description
	Group            string    // extracted from metadata for KPIs
	Status           string    // "Pending", "Running", "Completed", "Failed", "Cancelled"
	StartedAt        *time.Time
	CompletedAt      *time.Time
	Progress         *int      // 0-100
	DurationSeconds  *int
	Error            *string
	Metadata         map[string]interface{} // JSONB metadata
}