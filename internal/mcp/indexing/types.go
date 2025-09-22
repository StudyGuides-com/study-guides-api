package indexing

import "time"

// ResourceName is the MCP resource name for indexing operations
const ResourceName = "indexing"

// IndexingFilter defines the filter criteria for indexing operations
type IndexingFilter struct {
	TriggerReindex bool    `json:"triggerReindex,omitempty"`
	TriggerPruning bool    `json:"triggerPruning,omitempty"` // Remove orphaned objects from index
	ObjectType     *string `json:"objectType,omitempty"`     // "Tag", "User", "Contact", "FAQ"
	Force          *bool   `json:"force,omitempty"`          // Bypass hash comparison
	Status         *string `json:"status,omitempty"`          // Filter by job status
	JobID          *string `json:"jobId,omitempty"`           // Get specific job
}

// IndexingExecution represents an indexing job execution
type IndexingExecution struct {
	ID             string         `json:"id"`
	ObjectType     string         `json:"objectType"`
	Status         string         `json:"status"` // running/complete/failed
	StartedAt      *time.Time     `json:"startedAt,omitempty"`
	CompletedAt    *time.Time     `json:"completedAt,omitempty"`
	Duration       *time.Duration `json:"duration,omitempty"`
	ItemsProcessed int            `json:"itemsProcessed"`
	Error          string         `json:"error,omitempty"`
	Message        string         `json:"message,omitempty"`
	Force          bool           `json:"force"`
}

// IndexingUpdate is used for update operations (not applicable for indexing)
type IndexingUpdate struct{}

// IndexingStatus represents the status of indexing operations
type IndexingStatus string

const (
	IndexingStatusPending  IndexingStatus = "pending"
	IndexingStatusRunning  IndexingStatus = "running"
	IndexingStatusComplete IndexingStatus = "complete"
	IndexingStatusFailed   IndexingStatus = "failed"
)

// AllObjectTypes returns all supported object types for indexing
func AllObjectTypes() []string {
	return []string{
		"Tag",
		// Future: "User", "Contact", "FAQ"
	}
}