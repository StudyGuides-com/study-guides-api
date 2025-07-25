package repository

import (
	"context"
)

// Repository defines the generic CRUD interface for all domain entities
type Repository[T any, F any, U any] interface {
	// Find returns entities matching the filter
	Find(ctx context.Context, filter F) ([]T, error)

	// FindByID returns a single entity by its ID
	FindByID(ctx context.Context, id string) (*T, error)

	// Create creates a new entity and returns it
	Create(ctx context.Context, entity T) (*T, error)

	// Update updates an entity by ID with the provided update data
	Update(ctx context.Context, id string, update U) (*T, error)

	// Delete removes an entity by ID
	Delete(ctx context.Context, id string) error

	// Count returns the number of entities matching the filter
	Count(ctx context.Context, filter F) (int, error)
}

// CRUDOperation represents the type of operation being performed
type CRUDOperation string

const (
	OperationFind     CRUDOperation = "find"
	OperationFindByID CRUDOperation = "findById"
	OperationCreate   CRUDOperation = "create"
	OperationUpdate   CRUDOperation = "update"
	OperationDelete   CRUDOperation = "delete"
	OperationCount    CRUDOperation = "count"
)

// Command represents a request to perform an operation on a resource
type Command struct {
	Resource  string        `json:"resource"`  // e.g., "tag", "user", "question"
	Operation CRUDOperation `json:"operation"` // CRUD operation
	ID        string        `json:"id,omitempty"`        // For operations requiring an ID
	Payload   interface{}   `json:"payload,omitempty"`   // Filter, entity, or update data
}

// Response represents the result of executing a command
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Count   *int        `json:"count,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ResourceSchema defines the types used for a resource
type ResourceSchema struct {
	EntityType interface{} // The main entity type
	FilterType interface{} // The filter type for queries
	UpdateType interface{} // The update type for modifications
}