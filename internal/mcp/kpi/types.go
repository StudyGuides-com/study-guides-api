package kpi

import (
	"time"
)

// KPIGroup represents a group of statistics that can be calculated
type KPIGroup string

const (
	KPIGroupMonthlyInteractions KPIGroup = "MonthlyInteractions"
	KPIGroupTags                KPIGroup = "Tags"
	KPIGroupTagTypes            KPIGroup = "TagTypes"
	KPIGroupReports             KPIGroup = "Reports"
	KPIGroupTopics              KPIGroup = "Topics"
	KPIGroupMissingData         KPIGroup = "MissingData"
	KPIGroupRatings             KPIGroup = "Ratings"
	KPIGroupQuestions           KPIGroup = "Questions"
	KPIGroupUsers               KPIGroup = "Users"
	KPIGroupUserContent         KPIGroup = "UserContent"
	KPIGroupContacts            KPIGroup = "Contacts"
)

// AllKPIGroups returns all available KPI groups
func AllKPIGroups() []KPIGroup {
	return []KPIGroup{
		KPIGroupMonthlyInteractions,
		KPIGroupTags,
		KPIGroupTagTypes,
		KPIGroupReports,
		KPIGroupTopics,
		KPIGroupMissingData,
		KPIGroupRatings,
		KPIGroupQuestions,
		KPIGroupUsers,
		KPIGroupUserContent,
		KPIGroupContacts,
	}
}

// KPIStatus represents the status of a KPI calculation
type KPIStatus string

const (
	KPIStatusPending  KPIStatus = "pending"
	KPIStatusRunning  KPIStatus = "running"
	KPIStatusComplete KPIStatus = "complete"
	KPIStatusFailed   KPIStatus = "failed"
)

// KPIExecution represents a KPI calculation execution
type KPIExecution struct {
	ID          string        `json:"id"`
	Group       KPIGroup      `json:"group"`
	Status      KPIStatus     `json:"status"`
	StartedAt   *time.Time    `json:"started_at,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	Duration    *time.Duration `json:"duration,omitempty"`
	Error       string        `json:"error,omitempty"`
	Result      interface{}   `json:"result,omitempty"`
}

// KPIFilter for finding executions
type KPIFilter struct {
	Group    *KPIGroup  `json:"group,omitempty"`
	Status   *KPIStatus `json:"status,omitempty"`
	RunAll   bool       `json:"run_all,omitempty"`
	Since    *time.Time `json:"since,omitempty"`
}

// KPIUpdate for updating execution status
type KPIUpdate struct {
	Status      *KPIStatus  `json:"status,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	Error       *string     `json:"error,omitempty"`
	Result      interface{} `json:"result,omitempty"`
}

// KPIRepository defines the interface for KPI operations
type KPIRepository interface {
	// Execute starts a KPI calculation (returns immediately with execution ID)
	Execute(group KPIGroup) (*KPIExecution, error)
	
	// ExecuteAll starts all KPI calculations
	ExecuteAll() ([]*KPIExecution, error)
	
	// GetStatus returns the status of a KPI execution
	GetStatus(id string) (*KPIExecution, error)
	
	// GetLatestStatus returns the most recent execution for a group
	GetLatestStatus(group KPIGroup) (*KPIExecution, error)
	
	// GetRunningExecutions returns all currently running executions
	GetRunningExecutions() ([]*KPIExecution, error)
	
	// CancelExecution attempts to cancel a running execution
	CancelExecution(id string) error
}

// ResourceName is the MCP resource name for KPI operations
const ResourceName = "kpi"