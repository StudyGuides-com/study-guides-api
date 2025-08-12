package kpi

import (
	"context"
	"time"
)

// KPIStore defines the interface for KPI database operations
type KPIStore interface {
	// ExecuteTimeStatsProcedure runs the calculate_time_stats_by_group procedure
	ExecuteTimeStatsProcedure(ctx context.Context, group string) error
	
	// ExecuteUpdateStatsProcedure runs the update_calculated_stats_by_group procedure
	ExecuteUpdateStatsProcedure(ctx context.Context, group string) error
	
	// GetExecutionStatus retrieves the status of a running procedure (if trackable)
	GetExecutionStatus(ctx context.Context, executionID string) (*ExecutionStatus, error)
	
	// ListRunningExecutions lists all currently running KPI procedures
	ListRunningExecutions(ctx context.Context) ([]*ExecutionStatus, error)
}

// ExecutionStatus represents the status of a procedure execution
type ExecutionStatus struct {
	ID          string
	ProcedureName string
	Group       string
	Status      string // running, complete, failed
	StartedAt   time.Time
	CompletedAt *time.Time
	Error       *string
	RowsAffected *int64
}