package indexing

import (
	"github.com/studyguides-com/study-guides-api/internal/repository"
)

// GetResourceSchema returns the schema for indexing operations
func GetResourceSchema() repository.ResourceSchema {
	return repository.ResourceSchema{
		EntityType: IndexingExecution{},
		FilterType: IndexingFilter{},
		UpdateType: IndexingUpdate{},
	}
}