package indexing

import (
	"reflect"

	"github.com/studyguides-com/study-guides-api/internal/repository"
)

// GetResourceSchema returns the schema for indexing operations
func GetResourceSchema() repository.ResourceSchema {
	return repository.ResourceSchema{
		Name:        ResourceName,
		Description: "Manage Algolia search indexing operations",
		EntityType:  reflect.TypeOf(IndexingExecution{}),
		FilterType:  reflect.TypeOf(IndexingFilter{}),
		UpdateType:  reflect.TypeOf(IndexingUpdate{}),
		Operations: []string{
			"find",       // List jobs or trigger reindex
			"findById",   // Get specific job status
			"count",      // Count running jobs
		},
		Properties: map[string]repository.PropertySchema{
			"id": {
				Type:        "string",
				Description: "Unique job identifier",
				Required:    false,
			},
			"objectType": {
				Type:        "string",
				Description: "Type of object to index (Tag, User, Contact, FAQ)",
				Required:    false,
				Enum:        AllObjectTypes(),
			},
			"status": {
				Type:        "string",
				Description: "Job execution status",
				Required:    false,
				Enum:        []string{"pending", "running", "complete", "failed"},
			},
			"triggerReindex": {
				Type:        "boolean",
				Description: "Trigger a new reindexing job",
				Required:    false,
			},
			"force": {
				Type:        "boolean",
				Description: "Force reindex even if content hasn't changed",
				Required:    false,
			},
			"itemsProcessed": {
				Type:        "number",
				Description: "Number of items processed",
				Required:    false,
			},
		},
	}
}