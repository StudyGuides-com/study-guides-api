package indexing

import (
	"context"
	"fmt"
	"time"

	"github.com/studyguides-com/study-guides-api/internal/store/indexing"
)

// IndexingRepositoryAdapter adapts the indexing store to implement the MCP repository pattern
type IndexingRepositoryAdapter struct {
	store indexing.IndexingStore
}

// NewIndexingRepositoryAdapter creates a new adapter for indexing operations
func NewIndexingRepositoryAdapter(store indexing.IndexingStore) *IndexingRepositoryAdapter {
	return &IndexingRepositoryAdapter{
		store: store,
	}
}

// Find implements the generic find operation for indexing executions
func (a *IndexingRepositoryAdapter) Find(ctx context.Context, filter IndexingFilter) ([]IndexingExecution, error) {
	var results []IndexingExecution
	
	// Handle trigger reindex request
	if filter.TriggerReindex {
		objectType := "Tag" // Default
		if filter.ObjectType != nil {
			objectType = *filter.ObjectType
		}
		
		force := false
		if filter.Force != nil {
			force = *filter.Force
		}
		
		// Start indexing job
		jobID, err := a.store.StartIndexingJob(ctx, objectType, force)
		if err != nil {
			return nil, fmt.Errorf("failed to start indexing job: %w", err)
		}
		
		// Return immediate response
		now := time.Now()
		execution := IndexingExecution{
			ID:         jobID,
			ObjectType: objectType,
			Status:     string(IndexingStatusRunning),
			StartedAt:  &now,
			Message:    fmt.Sprintf("Started indexing %s (force=%t)", objectType, force),
			Force:      force,
		}
		
		results = append(results, execution)
		return results, nil
	}
	
	// Handle status filter - get running jobs
	if filter.Status != nil && *filter.Status == string(IndexingStatusRunning) {
		jobs, err := a.store.ListRunningJobs(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list running jobs: %w", err)
		}
		
		for _, job := range jobs {
			execution := a.jobToExecution(job)
			results = append(results, *execution)
		}
		
		return results, nil
	}
	
	// Get recent jobs for specific object type
	if filter.ObjectType != nil {
		jobs, err := a.store.ListRecentJobs(ctx, *filter.ObjectType)
		if err != nil {
			return nil, fmt.Errorf("failed to list recent jobs: %w", err)
		}
		
		for _, job := range jobs {
			execution := a.jobToExecution(job)
			results = append(results, *execution)
		}
		
		return results, nil
	}
	
	// Get all recent jobs (for all object types)
	for _, objType := range AllObjectTypes() {
		jobs, err := a.store.ListRecentJobs(ctx, objType)
		if err != nil {
			continue // Skip on error
		}
		
		// Just take the most recent one for each type
		if len(jobs) > 0 {
			execution := a.jobToExecution(jobs[0])
			results = append(results, *execution)
		}
	}
	
	return results, nil
}

// FindByID returns a specific indexing execution by job ID
func (a *IndexingRepositoryAdapter) FindByID(ctx context.Context, id string) (*IndexingExecution, error) {
	job, err := a.store.GetJobStatus(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}
	
	return a.jobToExecution(*job), nil
}

// Create starts a new indexing job
func (a *IndexingRepositoryAdapter) Create(ctx context.Context, entity IndexingExecution) (*IndexingExecution, error) {
	if entity.ObjectType == "" {
		entity.ObjectType = "Tag" // Default
	}
	
	jobID, err := a.store.StartIndexingJob(ctx, entity.ObjectType, entity.Force)
	if err != nil {
		return nil, fmt.Errorf("failed to start indexing job: %w", err)
	}
	
	now := time.Now()
	return &IndexingExecution{
		ID:         jobID,
		ObjectType: entity.ObjectType,
		Status:     string(IndexingStatusRunning),
		StartedAt:  &now,
		Message:    fmt.Sprintf("Started indexing %s", entity.ObjectType),
		Force:      entity.Force,
	}, nil
}

// Update is not supported for indexing executions
func (a *IndexingRepositoryAdapter) Update(ctx context.Context, id string, update IndexingUpdate) (*IndexingExecution, error) {
	return nil, fmt.Errorf("update operation not supported for indexing executions")
}

// Delete is not supported for indexing executions (jobs run to completion)
func (a *IndexingRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("delete operation not supported for indexing executions")
}

// Count returns the number of running indexing jobs
func (a *IndexingRepositoryAdapter) Count(ctx context.Context, filter IndexingFilter) (int, error) {
	if filter.Status != nil && *filter.Status == string(IndexingStatusRunning) {
		jobs, err := a.store.ListRunningJobs(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to list running jobs: %w", err)
		}
		return len(jobs), nil
	}
	
	// Count all recent jobs
	total := 0
	for _, objType := range AllObjectTypes() {
		jobs, err := a.store.ListRecentJobs(ctx, objType)
		if err != nil {
			continue
		}
		total += len(jobs)
	}
	
	return total, nil
}

// Helper function to convert JobStatus to IndexingExecution
func (a *IndexingRepositoryAdapter) jobToExecution(job indexing.JobStatus) *IndexingExecution {
	execution := &IndexingExecution{
		ID:        job.ID,
		StartedAt: job.StartedAt,
		Message:   job.Description,
	}
	
	// Extract object type from metadata
	if job.Metadata != nil {
		if objType, ok := job.Metadata["objectType"].(string); ok {
			execution.ObjectType = objType
		}
		if force, ok := job.Metadata["force"].(bool); ok {
			execution.Force = force
		}
		if items, ok := job.Metadata["itemsProcessed"].(float64); ok {
			execution.ItemsProcessed = int(items)
		}
	}
	
	// Map status
	switch job.Status {
	case "Running":
		execution.Status = string(IndexingStatusRunning)
	case "Completed":
		execution.Status = string(IndexingStatusComplete)
		execution.CompletedAt = job.CompletedAt
		if execution.CompletedAt != nil && execution.StartedAt != nil {
			duration := execution.CompletedAt.Sub(*execution.StartedAt)
			execution.Duration = &duration
		}
	case "Failed":
		execution.Status = string(IndexingStatusFailed)
		if job.ErrorMessage != nil {
			execution.Error = *job.ErrorMessage
		}
	default:
		execution.Status = string(IndexingStatusPending)
	}
	
	return execution
}