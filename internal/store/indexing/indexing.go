package indexing

import (
	"context"
	"time"
)

// IndexOperation represents a pending indexing operation from the outbox
type IndexOperation struct {
	ID         int64
	ObjectType string
	ObjectID   string
	Action     string // "upsert" or "delete"
	QueuedAt   time.Time
}

// IndexState represents the current indexing state of an object
type IndexState struct {
	ObjectType      string
	ObjectID        string
	LastIndexedAt   *time.Time
	LastIndexedHash []byte
	LastAttemptAt   *time.Time
	AttemptCount    int
	LastError       *string
}

// JobStatus represents the status of an indexing job
type JobStatus struct {
	ID             string
	Type           string
	Status         string // Running, Completed, Failed
	Description    string
	StartedAt      *time.Time
	CompletedAt    *time.Time
	Progress       int
	DurationSeconds *int
	ErrorMessage   *string
	Metadata       map[string]interface{}
}

// IndexingStore defines the interface for indexing operations
type IndexingStore interface {
	// Job management (like KPIs)
	StartIndexingJob(ctx context.Context, objectType string, force bool) (string, error)
	GetJobStatus(ctx context.Context, jobID string) (*JobStatus, error)
	ListRecentJobs(ctx context.Context, objectType string) ([]JobStatus, error)
	ListRunningJobs(ctx context.Context) ([]JobStatus, error)
	
	// Outbox operations
	QueueIndexOperation(ctx context.Context, objectType, objectID, action string) error
	GetPendingOperations(ctx context.Context, objectType string, limit int) ([]IndexOperation, error)
	RemoveFromOutbox(ctx context.Context, id int64) error
	
	// State tracking
	GetIndexState(ctx context.Context, objectType, objectID string) (*IndexState, error)
	UpdateIndexState(ctx context.Context, objectType, objectID string, hash []byte) error
	UpdateIndexError(ctx context.Context, objectType, objectID string, err error) error
	
	// Batch operations
	QueueBatchForReindex(ctx context.Context, objectType string) error
}