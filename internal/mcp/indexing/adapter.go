package indexing

import (
	"context"
	"fmt"
	"time"

	indexingcore "github.com/studyguides-com/study-guides-api/internal/core/indexing"
	"github.com/studyguides-com/study-guides-api/internal/store/indexing"
)

// IndexingRepositoryAdapter adapts the indexing business service to implement the MCP repository pattern
type IndexingRepositoryAdapter struct {
	business *indexingcore.BusinessService
}

// NewIndexingRepositoryAdapter creates a new adapter for indexing operations
func NewIndexingRepositoryAdapter(business *indexingcore.BusinessService) *IndexingRepositoryAdapter {
	return &IndexingRepositoryAdapter{
		business: business,
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
		
		// Start indexing job using business service
		businessReq := indexingcore.TriggerIndexingRequest{
			ObjectType: objectType,
			Force:      force,
		}
		businessResp, err := a.business.TriggerIndexing(ctx, businessReq)
		if err != nil {
			return nil, fmt.Errorf("failed to start indexing job: %w", err)
		}
		jobID := businessResp.JobID
		
		// Get count of items that will be processed to provide better user feedback
		validItemsMsg := "Processing all valid items"
		if objectType == "Tag" {
			// This gives a rough estimate since we filter at queue time
			validItemsMsg = "Processing all valid tags (invalid records automatically excluded)"
		}
		
		// Return immediate response with more context
		now := time.Now()
		forceMsg := ""
		if force {
			forceMsg = " (full rebuild - all items regardless of changes)"
		} else {
			forceMsg = " (incremental - only changed items)"
		}
		
		execution := IndexingExecution{
			ID:         jobID,
			ObjectType: objectType,
			Status:     string(IndexingStatusRunning),
			StartedAt:  &now,
			Message:    fmt.Sprintf("Started %s indexing job %s%s. %s", objectType, jobID[:8], forceMsg, validItemsMsg),
			Force:      force,
		}
		
		results = append(results, execution)
		return results, nil
	}
	
	// Handle status filter - get running jobs
	if filter.Status != nil && *filter.Status == string(IndexingStatusRunning) {
		jobs, err := a.business.ListRunningJobs(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list running jobs: %w", err)
		}
		
		// Create summary response for running jobs
		if len(jobs) == 0 {
			execution := IndexingExecution{
				ID:      "status",
				Message: "No indexing jobs currently running",
				Status:  string(IndexingStatusComplete),
			}
			results = append(results, execution)
		} else {
			// Add summary first
			summaryMsg := fmt.Sprintf("Indexing Status: %d jobs currently running", len(jobs))
			if len(jobs) == 1 {
				summaryMsg = "Indexing Status: 1 job currently running"
			}
			
			execution := IndexingExecution{
				ID:      "summary",
				Message: summaryMsg,
				Status:  string(IndexingStatusRunning),
			}
			results = append(results, execution)
			
			// Then add individual jobs with enhanced messages
			for _, job := range jobs {
				execution := a.jobToExecution(job)
				// Enhance the message to be more descriptive
				if execution.ItemsProcessed > 0 {
					execution.Message = fmt.Sprintf("%s - %d items processed", execution.Message, execution.ItemsProcessed)
				} else {
					execution.Message = fmt.Sprintf("%s - starting up...", execution.Message)
				}
				results = append(results, *execution)
			}
		}
		
		return results, nil
	}
	
	// Get recent jobs for specific object type
	if filter.ObjectType != nil {
		businessReq := indexingcore.ListRecentJobsRequest{
			ObjectType: *filter.ObjectType,
		}
		jobs, err := a.business.ListRecentJobs(ctx, businessReq)
		if err != nil {
			return nil, fmt.Errorf("failed to list recent jobs: %w", err)
		}
		
		for _, job := range jobs {
			execution := a.jobToExecution(job)
			results = append(results, *execution)
		}
		
		return results, nil
	}
	
	// Get status summary and recent jobs
	runningJobs, _ := a.business.ListRunningJobs(ctx)
	runningCount := len(runningJobs)
	
	// Create a summary response first
	var summaryMsg string
	if runningCount == 0 {
		summaryMsg = "No indexing jobs currently running. Showing recent job history:"
	} else if runningCount == 1 {
		summaryMsg = "1 indexing job currently running. Recent activity:"
	} else {
		summaryMsg = fmt.Sprintf("%d indexing jobs currently running. Recent activity:", runningCount)
	}
	
	summaryExecution := IndexingExecution{
		ID:      "overview",
		Message: summaryMsg,
		Status:  string(IndexingStatusComplete),
	}
	results = append(results, summaryExecution)
	
	// Add running jobs first if any
	for _, job := range runningJobs {
		execution := a.jobToExecution(job)
		if execution.ItemsProcessed > 0 {
			execution.Message = fmt.Sprintf("🟡 %s - %d items processed", execution.Message, execution.ItemsProcessed)
		} else {
			execution.Message = fmt.Sprintf("🟡 %s - initializing...", execution.Message)
		}
		results = append(results, *execution)
	}
	
	// Then get recent completed jobs for context
	for _, objType := range AllObjectTypes() {
		businessReq := indexingcore.ListRecentJobsRequest{
			ObjectType: objType,
		}
		jobs, err := a.business.ListRecentJobs(ctx, businessReq)
		if err != nil {
			continue // Skip on error
		}
		
		// Add the most recent completed job for each type (not running)
		for _, job := range jobs {
			if job.Status != "Running" {
				execution := a.jobToExecution(job)
				statusIcon := "✅"
				if job.Status == "Failed" {
					statusIcon = "❌"
				}
				if execution.ItemsProcessed > 0 {
					execution.Message = fmt.Sprintf("%s %s - completed %d items", statusIcon, execution.Message, execution.ItemsProcessed)
				} else {
					execution.Message = fmt.Sprintf("%s %s", statusIcon, execution.Message)
				}
				results = append(results, *execution)
				break // Only one recent completed job per type
			}
		}
	}
	
	return results, nil
}

// FindByID returns a specific indexing execution by job ID
func (a *IndexingRepositoryAdapter) FindByID(ctx context.Context, id string) (*IndexingExecution, error) {
	job, err := a.business.GetJobStatus(ctx, id)
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
	
	businessReq := indexingcore.TriggerIndexingRequest{
		ObjectType: entity.ObjectType,
		Force:      entity.Force,
	}
	businessResp, err := a.business.TriggerIndexing(ctx, businessReq)
	if err != nil {
		return nil, fmt.Errorf("failed to start indexing job: %w", err)
	}
	jobID := businessResp.JobID
	
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
		jobs, err := a.business.ListRunningJobs(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to list running jobs: %w", err)
		}
		return len(jobs), nil
	}
	
	// Count all recent jobs
	total := 0
	for _, objType := range AllObjectTypes() {
		businessReq := indexingcore.ListRecentJobsRequest{
			ObjectType: objType,
		}
		jobs, err := a.business.ListRecentJobs(ctx, businessReq)
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