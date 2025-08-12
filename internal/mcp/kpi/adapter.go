package kpi

import (
	"context"
	"fmt"
	"time"

	"github.com/lucsky/cuid"
	"github.com/studyguides-com/study-guides-api/internal/store/kpi"
)

// KPIRepositoryAdapter adapts the KPI store to implement the MCP repository pattern
type KPIRepositoryAdapter struct {
	store kpi.KPIStore
}

// NewKPIRepositoryAdapter creates a new adapter for KPI operations
func NewKPIRepositoryAdapter(store kpi.KPIStore) *KPIRepositoryAdapter {
	return &KPIRepositoryAdapter{
		store: store,
	}
}

// Find implements the generic find operation for KPI executions
func (a *KPIRepositoryAdapter) Find(ctx context.Context, filter KPIFilter) ([]KPIExecution, error) {
	var results []KPIExecution
	
	// If RunAll is set, execute all KPI groups
	if filter.RunAll {
		executions, err := a.ExecuteAll()
		if err != nil {
			return nil, err
		}
		for _, exec := range executions {
			results = append(results, *exec)
		}
		return results, nil
	}
	
	// If a specific group is requested, execute it
	if filter.Group != nil {
		exec, err := a.Execute(*filter.Group)
		if err != nil {
			return nil, err
		}
		results = append(results, *exec)
		return results, nil
	}
	
	// Otherwise return status of latest execution for each group
	latestExecutions, err := a.GetLatestExecutionsPerGroup()
	if err != nil {
		return nil, err
	}
	
	for _, exec := range latestExecutions {
		results = append(results, *exec)
	}
	
	return results, nil
}

// FindByID returns a specific KPI execution by ID
func (a *KPIRepositoryAdapter) FindByID(ctx context.Context, id string) (*KPIExecution, error) {
	return a.GetStatus(id)
}

// Create starts a new KPI execution
func (a *KPIRepositoryAdapter) Create(ctx context.Context, entity KPIExecution) (*KPIExecution, error) {
	if entity.Group == "" {
		return nil, fmt.Errorf("group is required")
	}
	
	return a.Execute(entity.Group)
}

// Update is not supported for KPI executions
func (a *KPIRepositoryAdapter) Update(ctx context.Context, id string, update KPIUpdate) (*KPIExecution, error) {
	return nil, fmt.Errorf("update operation not supported for KPI executions")
}

// Delete cancels a KPI execution
func (a *KPIRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return a.CancelExecution(id)
}

// Count returns the number of running executions
func (a *KPIRepositoryAdapter) Count(ctx context.Context, filter KPIFilter) (int, error) {
	if filter.Status != nil && *filter.Status == KPIStatusRunning {
		running, err := a.GetRunningExecutions()
		if err != nil {
			return 0, err
		}
		return len(running), nil
	}
	
	// Count all recent executions
	recent, err := a.GetRecentExecutions()
	if err != nil {
		return 0, err
	}
	return len(recent), nil
}

// Execute starts a KPI calculation (returns immediately with execution ID)
func (a *KPIRepositoryAdapter) Execute(group KPIGroup) (*KPIExecution, error) {
	now := time.Now()
	
	// Determine which procedure to run (this creates the Job record and returns the ID)
	ctx := context.Background()
	var executionID string
	var err error
	
	switch group {
	case KPIGroupMonthlyInteractions:
		executionID, err = a.store.ExecuteTimeStatsProcedure(ctx, string(group))
	default:
		executionID, err = a.store.ExecuteUpdateStatsProcedure(ctx, string(group))
	}
	
	if err != nil {
		return &KPIExecution{
			ID:        "",
			Group:     group,
			Status:    KPIStatusFailed,
			StartedAt: &now,
			Error:     err.Error(),
		}, err
	}
	
	// Return execution info (the actual job is running in the background)
	return &KPIExecution{
		ID:        executionID,
		Group:     group,
		Status:    KPIStatusRunning,
		StartedAt: &now,
	}, nil
}

// ExecuteAll starts all KPI calculations
func (a *KPIRepositoryAdapter) ExecuteAll() ([]*KPIExecution, error) {
	groups := AllKPIGroups()
	executions := make([]*KPIExecution, 0, len(groups))
	
	for _, group := range groups {
		exec, err := a.Execute(group)
		if err != nil {
			// Continue with other groups even if one fails
			exec = &KPIExecution{
				ID:     cuid.New(),
				Group:  group,
				Status: KPIStatusFailed,
				Error:  err.Error(),
			}
		}
		executions = append(executions, exec)
	}
	
	return executions, nil
}

// GetStatus returns the status of a KPI execution
func (a *KPIRepositoryAdapter) GetStatus(id string) (*KPIExecution, error) {
	// Check with the store
	storeStatus, err := a.store.GetExecutionStatus(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("execution %s not found", id)
	}
	
	// Convert store status to KPI execution
	execution := &KPIExecution{
		ID:        storeStatus.ID,
		Group:     KPIGroup(storeStatus.Group),
		StartedAt: storeStatus.StartedAt,
	}
	
	switch storeStatus.Status {
	case "Running":
		execution.Status = KPIStatusRunning
	case "Completed":
		execution.Status = KPIStatusComplete
		execution.CompletedAt = storeStatus.CompletedAt
		if execution.CompletedAt != nil && execution.StartedAt != nil {
			duration := execution.CompletedAt.Sub(*execution.StartedAt)
			execution.Duration = &duration
		}
	case "Failed":
		execution.Status = KPIStatusFailed
		if storeStatus.Error != nil {
			execution.Error = *storeStatus.Error
		}
	}
	
	return execution, nil
}

// GetLatestStatus returns the most recent execution for a group
func (a *KPIRepositoryAdapter) GetLatestStatus(group KPIGroup) (*KPIExecution, error) {
	// Get recent executions and find the latest for this group
	recent, err := a.GetRecentExecutions()
	if err != nil {
		return nil, err
	}
	
	var latest *KPIExecution
	for _, exec := range recent {
		if exec.Group == group {
			if latest == nil || (exec.StartedAt != nil && latest.StartedAt != nil && exec.StartedAt.After(*latest.StartedAt)) {
				latest = exec
			}
		}
	}
	
	if latest == nil {
		return nil, fmt.Errorf("no executions found for group %s", group)
	}
	
	return latest, nil
}

// GetRunningExecutions returns all currently running executions
func (a *KPIRepositoryAdapter) GetRunningExecutions() ([]*KPIExecution, error) {
	var running []*KPIExecution
	
	// Get running executions from the store
	storeRunning, err := a.store.ListRunningExecutions(context.Background())
	if err != nil {
		return nil, err
	}
	
	for _, storeExec := range storeRunning {
		execution := &KPIExecution{
			ID:        storeExec.ID,
			Group:     KPIGroup(storeExec.Group),
			Status:    KPIStatusRunning,
			StartedAt: storeExec.StartedAt,
		}
		running = append(running, execution)
	}
	
	return running, nil
}

// CancelExecution attempts to cancel a running execution
func (a *KPIRepositoryAdapter) CancelExecution(id string) error {
	// Check if execution exists and is running
	status, err := a.GetStatus(id)
	if err != nil {
		return fmt.Errorf("execution %s not found", id)
	}
	
	if status.Status != KPIStatusRunning {
		return fmt.Errorf("execution %s is not running", id)
	}
	
	// TODO: Update the Job record to mark as cancelled
	// For now, just return an error since we can't easily cancel background procedures
	return fmt.Errorf("cancellation not yet implemented")
}

// GetLatestExecutionsPerGroup returns the latest execution for each KPI group
func (a *KPIRepositoryAdapter) GetLatestExecutionsPerGroup() ([]*KPIExecution, error) {
	var latest []*KPIExecution
	
	// Get all recent executions from the store
	storeExecutions, err := a.store.ListRecentExecutions(context.Background())
	if err != nil {
		return nil, err
	}
	
	// Group by group name and keep only the latest for each
	groupLatest := make(map[string]*KPIExecution)
	
	for _, storeExec := range storeExecutions {
		execution := &KPIExecution{
			ID:        storeExec.ID,
			Group:     KPIGroup(storeExec.Group),
			StartedAt: storeExec.StartedAt,
		}
		
		// Map status from store
		switch storeExec.Status {
		case "Running":
			execution.Status = KPIStatusRunning
		case "Completed":
			execution.Status = KPIStatusComplete
			execution.CompletedAt = storeExec.CompletedAt
			if execution.CompletedAt != nil && execution.StartedAt != nil {
				duration := execution.CompletedAt.Sub(*execution.StartedAt)
				execution.Duration = &duration
			}
		case "Failed":
			execution.Status = KPIStatusFailed
			if storeExec.Error != nil {
				execution.Error = *storeExec.Error
			}
		}
		
		// Keep only the latest for each group (by StartedAt)
		groupName := string(execution.Group)
		if existing, exists := groupLatest[groupName]; !exists || 
			(execution.StartedAt != nil && existing.StartedAt != nil && execution.StartedAt.After(*existing.StartedAt)) {
			groupLatest[groupName] = execution
		}
	}
	
	// Convert map to slice
	for _, exec := range groupLatest {
		latest = append(latest, exec)
	}
	
	return latest, nil
}

// GetRecentExecutions returns all recent executions (running and completed)
func (a *KPIRepositoryAdapter) GetRecentExecutions() ([]*KPIExecution, error) {
	var recent []*KPIExecution
	
	// Get recent executions from the store
	storeExecutions, err := a.store.ListRecentExecutions(context.Background())
	if err != nil {
		return nil, err
	}
	
	for _, storeExec := range storeExecutions {
		execution := &KPIExecution{
			ID:        storeExec.ID,
			Group:     KPIGroup(storeExec.Group),
			StartedAt: storeExec.StartedAt,
		}
		
		// Map status from store
		switch storeExec.Status {
		case "Running":
			execution.Status = KPIStatusRunning
		case "Completed":
			execution.Status = KPIStatusComplete
			execution.CompletedAt = storeExec.CompletedAt
			if execution.CompletedAt != nil && execution.StartedAt != nil {
				duration := execution.CompletedAt.Sub(*execution.StartedAt)
				execution.Duration = &duration
			}
		case "Failed":
			execution.Status = KPIStatusFailed
			if storeExec.Error != nil {
				execution.Error = *storeExec.Error
			}
		}
		
		recent = append(recent, execution)
	}
	
	return recent, nil
}