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
	
	// Track MCP-level executions
	executions map[string]*KPIExecution
}

// NewKPIRepositoryAdapter creates a new adapter for KPI operations
func NewKPIRepositoryAdapter(store kpi.KPIStore) *KPIRepositoryAdapter {
	return &KPIRepositoryAdapter{
		store:      store,
		executions: make(map[string]*KPIExecution),
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
	
	// Otherwise return status of running executions
	running, err := a.GetRunningExecutions()
	if err != nil {
		return nil, err
	}
	
	for _, exec := range running {
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
	
	// Count all tracked executions
	return len(a.executions), nil
}

// Execute starts a KPI calculation (returns immediately with execution ID)
func (a *KPIRepositoryAdapter) Execute(group KPIGroup) (*KPIExecution, error) {
	executionID := cuid.New()
	now := time.Now()
	
	execution := &KPIExecution{
		ID:        executionID,
		Group:     group,
		Status:    KPIStatusRunning,
		StartedAt: &now,
	}
	
	// Store execution
	a.executions[executionID] = execution
	
	// Determine which procedure to run
	ctx := context.Background()
	var err error
	
	switch group {
	case KPIGroupMonthlyInteractions:
		err = a.store.ExecuteTimeStatsProcedure(ctx, string(group))
	default:
		err = a.store.ExecuteUpdateStatsProcedure(ctx, string(group))
	}
	
	if err != nil {
		execution.Status = KPIStatusFailed
		execution.Error = err.Error()
		return execution, err
	}
	
	// Start background status monitoring
	go a.monitorExecution(executionID)
	
	return execution, nil
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
	if execution, exists := a.executions[id]; exists {
		return execution, nil
	}
	
	// Check with the store
	storeStatus, err := a.store.GetExecutionStatus(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("execution %s not found", id)
	}
	
	// Convert store status to KPI execution
	execution := &KPIExecution{
		ID:        storeStatus.ID,
		Group:     KPIGroup(storeStatus.Group),
		StartedAt: &storeStatus.StartedAt,
	}
	
	switch storeStatus.Status {
	case "running":
		execution.Status = KPIStatusRunning
	case "complete":
		execution.Status = KPIStatusComplete
		execution.CompletedAt = storeStatus.CompletedAt
		if execution.CompletedAt != nil && execution.StartedAt != nil {
			duration := execution.CompletedAt.Sub(*execution.StartedAt)
			execution.Duration = &duration
		}
	case "failed":
		execution.Status = KPIStatusFailed
		if storeStatus.Error != nil {
			execution.Error = *storeStatus.Error
		}
	}
	
	return execution, nil
}

// GetLatestStatus returns the most recent execution for a group
func (a *KPIRepositoryAdapter) GetLatestStatus(group KPIGroup) (*KPIExecution, error) {
	var latest *KPIExecution
	
	for _, exec := range a.executions {
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
	
	// Check our tracked executions
	for _, exec := range a.executions {
		if exec.Status == KPIStatusRunning {
			running = append(running, exec)
		}
	}
	
	// Also check with the store
	storeRunning, err := a.store.ListRunningExecutions(context.Background())
	if err == nil {
		for _, storeExec := range storeRunning {
			// Check if we already have this execution
			found := false
			for _, exec := range running {
				if exec.ID == storeExec.ID {
					found = true
					break
				}
			}
			
			if !found {
				execution := &KPIExecution{
					ID:        storeExec.ID,
					Group:     KPIGroup(storeExec.Group),
					Status:    KPIStatusRunning,
					StartedAt: &storeExec.StartedAt,
				}
				running = append(running, execution)
			}
		}
	}
	
	return running, nil
}

// CancelExecution attempts to cancel a running execution
func (a *KPIRepositoryAdapter) CancelExecution(id string) error {
	if execution, exists := a.executions[id]; exists {
		if execution.Status == KPIStatusRunning {
			execution.Status = KPIStatusFailed
			execution.Error = "Cancelled by user"
			now := time.Now()
			execution.CompletedAt = &now
			if execution.StartedAt != nil {
				duration := now.Sub(*execution.StartedAt)
				execution.Duration = &duration
			}
			return nil
		}
		return fmt.Errorf("execution %s is not running", id)
	}
	
	return fmt.Errorf("execution %s not found", id)
}

// monitorExecution monitors the status of an execution in the background
func (a *KPIRepositoryAdapter) monitorExecution(executionID string) {
	// Poll every 5 seconds for status updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		storeStatus, err := a.store.GetExecutionStatus(context.Background(), executionID)
		if err != nil {
			// Execution might not be tracked by store yet
			continue
		}
		
		if execution, exists := a.executions[executionID]; exists {
			// Update status based on store
			if storeStatus.Status == "complete" {
				execution.Status = KPIStatusComplete
				execution.CompletedAt = storeStatus.CompletedAt
				if execution.StartedAt != nil && execution.CompletedAt != nil {
					duration := execution.CompletedAt.Sub(*execution.StartedAt)
					execution.Duration = &duration
				}
				return // Stop monitoring
			} else if storeStatus.Status == "failed" {
				execution.Status = KPIStatusFailed
				if storeStatus.Error != nil {
					execution.Error = *storeStatus.Error
				}
				now := time.Now()
				execution.CompletedAt = &now
				if execution.StartedAt != nil {
					duration := now.Sub(*execution.StartedAt)
					execution.Duration = &duration
				}
				return // Stop monitoring
			}
		}
	}
}