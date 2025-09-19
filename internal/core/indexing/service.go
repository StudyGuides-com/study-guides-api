package indexing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lucsky/cuid"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"github.com/studyguides-com/study-guides-api/internal/store/indexing"
)

// BusinessService provides core indexing business logic
// This is shared between gRPC and MCP interfaces
type BusinessService struct {
	store store.Store
}

// NewBusinessService creates a new business service instance
func NewBusinessService(store store.Store) *BusinessService {
	return &BusinessService{
		store: store,
	}
}

// TriggerIndexingRequest represents an indexing job request
type TriggerIndexingRequest struct {
	ObjectType string
	Force      bool
}

// TriggerIndexingResponse represents the response from triggering an indexing job
type TriggerIndexingResponse struct {
	JobID     string
	Status    string
	Message   string
	StartedAt time.Time
}

// TriggerIndexing starts a new indexing job
func (bs *BusinessService) TriggerIndexing(ctx context.Context, req TriggerIndexingRequest) (*TriggerIndexingResponse, error) {
	// Default object type to "Tag" if not specified
	objectType := req.ObjectType
	if objectType == "" {
		objectType = "Tag"
	}

	// Start indexing job
	jobID, err := bs.store.IndexingStore().StartIndexingJob(ctx, objectType, req.Force)
	if err != nil {
		return nil, fmt.Errorf("failed to start indexing job: %w", err)
	}

	// Create response message
	forceMsg := "incremental"
	if req.Force {
		forceMsg = "force rebuild"
	}
	message := fmt.Sprintf("Started %s indexing job (%s mode)", objectType, forceMsg)

	return &TriggerIndexingResponse{
		JobID:     jobID,
		Status:    "Running",
		Message:   message,
		StartedAt: time.Now(),
	}, nil
}

// GetJobStatus returns the status of a specific indexing job
func (bs *BusinessService) GetJobStatus(ctx context.Context, jobID string) (*indexing.JobStatus, error) {
	job, err := bs.store.IndexingStore().GetJobStatus(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return job, nil
}

// ListRunningJobs returns all currently running indexing jobs
func (bs *BusinessService) ListRunningJobs(ctx context.Context) ([]indexing.JobStatus, error) {
	jobs, err := bs.store.IndexingStore().ListRunningJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list running jobs: %w", err)
	}
	return jobs, nil
}

// ListRecentJobsRequest represents a request for recent jobs
type ListRecentJobsRequest struct {
	ObjectType string
	Limit      int
}

// ListRecentJobs returns recent indexing jobs, optionally filtered by object type
func (bs *BusinessService) ListRecentJobs(ctx context.Context, req ListRecentJobsRequest) ([]indexing.JobStatus, error) {
	// Default object type to "Tag" if not specified
	objectType := req.ObjectType
	if objectType == "" {
		objectType = "Tag"
	}

	// Get recent jobs from store
	jobs, err := bs.store.IndexingStore().ListRecentJobs(ctx, objectType)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent jobs: %w", err)
	}

	// Apply limit if specified
	if req.Limit > 0 && len(jobs) > req.Limit {
		jobs = jobs[:req.Limit]
	}

	return jobs, nil
}

// TriggerTagIndexingRequest represents a tag indexing request with filters
type TriggerTagIndexingRequest struct {
	Force        bool
	TagTypes     []sharedpb.TagType
	ContextTypes []sharedpb.ContextType
}

// TriggerTagIndexing starts a new tag indexing job with filtering
func (bs *BusinessService) TriggerTagIndexing(ctx context.Context, req TriggerTagIndexingRequest) (*TriggerIndexingResponse, error) {
	// Start indexing job with filters
	jobID, err := bs.store.IndexingStore().StartIndexingJobWithFilters(ctx, "Tag", req.Force, req.TagTypes, req.ContextTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to start tag indexing job: %w", err)
	}

	// Create descriptive message about filters
	var filterParts []string

	if len(req.TagTypes) > 0 {
		filterParts = append(filterParts, fmt.Sprintf("%d tag types", len(req.TagTypes)))
	}

	if len(req.ContextTypes) > 0 {
		filterParts = append(filterParts, fmt.Sprintf("%d context types", len(req.ContextTypes)))
	}

	var filterMsg string
	if len(filterParts) > 0 {
		filterMsg = fmt.Sprintf(" with filters: %v", filterParts)
	} else {
		filterMsg = " (all tags)"
	}

	forceMsg := "incremental"
	if req.Force {
		forceMsg = "force rebuild"
	}

	message := fmt.Sprintf("Started tag indexing job (%s mode)%s", forceMsg, filterMsg)

	return &TriggerIndexingResponse{
		JobID:     jobID,
		Status:    "Running",
		Message:   message,
		StartedAt: time.Now(),
	}, nil
}

// TriggerSingleIndexingRequest represents a single item indexing request
type TriggerSingleIndexingRequest struct {
	ObjectType string
	ID         string
	Force      bool
}

// TriggerSingleIndexing starts a new indexing job for a single specific item
func (bs *BusinessService) TriggerSingleIndexing(ctx context.Context, req TriggerSingleIndexingRequest) (*TriggerIndexingResponse, error) {
	// Default object type to "Tag" if not specified
	objectType := req.ObjectType
	if objectType == "" {
		objectType = "Tag"
	}

	// Validate that ID is provided
	if req.ID == "" {
		return nil, fmt.Errorf("ID is required for single item indexing")
	}

	// Start single item indexing job
	jobID, err := bs.store.IndexingStore().StartSingleIndexingJob(ctx, objectType, req.ID, req.Force)
	if err != nil {
		return nil, fmt.Errorf("failed to start single indexing job: %w", err)
	}

	// Create response message
	forceMsg := "incremental"
	if req.Force {
		forceMsg = "force rebuild"
	}
	message := fmt.Sprintf("Started single %s indexing job for ID %s (%s mode)", objectType, req.ID, forceMsg)

	return &TriggerIndexingResponse{
		JobID:     jobID,
		Status:    "Running",
		Message:   message,
		StartedAt: time.Now(),
	}, nil
}