package mcp

import (
	"github.com/studyguides-com/study-guides-api/internal/repository"
)

// Re-export common types for convenience
type Command = repository.Command
type Response = repository.Response
type CRUDOperation = repository.CRUDOperation
type ResourceSchema = repository.ResourceSchema

// Re-export constants
const (
	OperationFind     = repository.OperationFind
	OperationFindByID = repository.OperationFindByID
	OperationCreate   = repository.OperationCreate
	OperationUpdate   = repository.OperationUpdate
	OperationDelete   = repository.OperationDelete
	OperationCount    = repository.OperationCount
)

// RepositoryRegistry manages the mapping of resource names to repositories
type RepositoryRegistry struct {
	repositories map[string]interface{}
	schemas      map[string]repository.ResourceSchema
}

// NewRepositoryRegistry creates a new repository registry
func NewRepositoryRegistry() *RepositoryRegistry {
	return &RepositoryRegistry{
		repositories: make(map[string]interface{}),
		schemas:      make(map[string]repository.ResourceSchema),
	}
}

// Register adds a repository for a specific resource
func (r *RepositoryRegistry) Register(resourceName string, repo interface{}, schema repository.ResourceSchema) {
	r.repositories[resourceName] = repo
	r.schemas[resourceName] = schema
}

// GetRepository returns the repository for a resource
func (r *RepositoryRegistry) GetRepository(resourceName string) (interface{}, bool) {
	repo, exists := r.repositories[resourceName]
	return repo, exists
}

// GetSchema returns the schema for a resource
func (r *RepositoryRegistry) GetSchema(resourceName string) (repository.ResourceSchema, bool) {
	schema, exists := r.schemas[resourceName]
	return schema, exists
}

// ListResources returns all registered resource names
func (r *RepositoryRegistry) ListResources() []string {
	resources := make([]string, 0, len(r.repositories))
	for name := range r.repositories {
		resources = append(resources, name)
	}
	return resources
}