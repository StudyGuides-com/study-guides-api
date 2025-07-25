package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/studyguides-com/study-guides-api/internal/repository"
)

// CommandHandler executes commands against registered repositories
type CommandHandler struct {
	registry *RepositoryRegistry
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(registry *RepositoryRegistry) *CommandHandler {
	return &CommandHandler{
		registry: registry,
	}
}

// Handle executes a command and returns a response
func (h *CommandHandler) Handle(ctx context.Context, cmd repository.Command) (*repository.Response, error) {
	repo, exists := h.registry.GetRepository(cmd.Resource)
	if !exists {
		return &repository.Response{
			Success: false,
			Error:   fmt.Sprintf("unknown resource: %s", cmd.Resource),
		}, nil
	}

	schema, exists := h.registry.GetSchema(cmd.Resource)
	if !exists {
		return &repository.Response{
			Success: false,
			Error:   fmt.Sprintf("no schema found for resource: %s", cmd.Resource),
		}, nil
	}

	switch cmd.Operation {
	case OperationFind:
		return h.handleFind(ctx, repo, schema, cmd.Payload)
	case OperationFindByID:
		return h.handleFindByID(ctx, repo, cmd.ID)
	case OperationCreate:
		return h.handleCreate(ctx, repo, schema, cmd.Payload)
	case OperationUpdate:
		return h.handleUpdate(ctx, repo, schema, cmd.ID, cmd.Payload)
	case OperationDelete:
		return h.handleDelete(ctx, repo, cmd.ID)
	case OperationCount:
		return h.handleCount(ctx, repo, schema, cmd.Payload)
	default:
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("unknown operation: %s", cmd.Operation),
		}, nil
	}
}

// handleFind executes a find operation using reflection
func (h *CommandHandler) handleFind(ctx context.Context, repo interface{}, schema repository.ResourceSchema, payload interface{}) (*Response, error) {
	// Convert payload to the correct filter type
	filter, err := h.convertPayload(payload, schema.FilterType)
	if err != nil {
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("invalid filter payload: %v", err),
		}, nil
	}

	// Call the Find method using reflection
	repoValue := reflect.ValueOf(repo)
	method := repoValue.MethodByName("Find")
	if !method.IsValid() {
		return &Response{
			Success: false,
			Error:   "repository does not implement Find method",
		}, nil
	}

	// Call Find(ctx, filter)
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(filter).Elem(), // Dereference the pointer
	})

	// Check for errors
	if len(results) != 2 {
		return &Response{
			Success: false,
			Error:   "Find method should return ([]T, error)",
		}, nil
	}

	errValue := results[1]
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		return &Response{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Extract the data
	data := results[0].Interface()
	dataSlice := reflect.ValueOf(data)
	count := dataSlice.Len()

	return &Response{
		Success: true,
		Data:    data,
		Count:   &count,
		Message: fmt.Sprintf("Found %d items", count),
	}, nil
}

// handleFindByID executes a findByID operation using reflection
func (h *CommandHandler) handleFindByID(ctx context.Context, repo interface{}, id string) (*Response, error) {
	if id == "" {
		return &Response{
			Success: false,
			Error:   "id is required for findById operation",
		}, nil
	}

	repoValue := reflect.ValueOf(repo)
	method := repoValue.MethodByName("FindByID")
	if !method.IsValid() {
		return &Response{
			Success: false,
			Error:   "repository does not implement FindByID method",
		}, nil
	}

	// Call FindByID(ctx, id)
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
	})

	if len(results) != 2 {
		return &Response{
			Success: false,
			Error:   "FindByID method should return (*T, error)",
		}, nil
	}

	errValue := results[1]
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		return &Response{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	data := results[0].Interface()
	return &Response{
		Success: true,
		Data:    data,
		Message: fmt.Sprintf("Found item with id: %s", id),
	}, nil
}

// handleCreate executes a create operation using reflection
func (h *CommandHandler) handleCreate(ctx context.Context, repo interface{}, schema ResourceSchema, payload interface{}) (*Response, error) {
	// Convert payload to entity type
	entity, err := h.convertPayload(payload, schema.EntityType)
	if err != nil {
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("invalid entity payload: %v", err),
		}, nil
	}

	repoValue := reflect.ValueOf(repo)
	method := repoValue.MethodByName("Create")
	if !method.IsValid() {
		return &Response{
			Success: false,
			Error:   "repository does not implement Create method",
		}, nil
	}

	// Call Create(ctx, entity)
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(entity).Elem(),
	})

	if len(results) != 2 {
		return &Response{
			Success: false,
			Error:   "Create method should return (*T, error)",
		}, nil
	}

	errValue := results[1]
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		return &Response{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	data := results[0].Interface()
	return &Response{
		Success: true,
		Data:    data,
		Message: "Entity created successfully",
	}, nil
}

// handleUpdate executes an update operation using reflection
func (h *CommandHandler) handleUpdate(ctx context.Context, repo interface{}, schema ResourceSchema, id string, payload interface{}) (*Response, error) {
	if id == "" {
		return &Response{
			Success: false,
			Error:   "id is required for update operation",
		}, nil
	}

	// Convert payload to update type
	update, err := h.convertPayload(payload, schema.UpdateType)
	if err != nil {
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("invalid update payload: %v", err),
		}, nil
	}

	repoValue := reflect.ValueOf(repo)
	method := repoValue.MethodByName("Update")
	if !method.IsValid() {
		return &Response{
			Success: false,
			Error:   "repository does not implement Update method",
		}, nil
	}

	// Call Update(ctx, id, update)
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
		reflect.ValueOf(update).Elem(),
	})

	if len(results) != 2 {
		return &Response{
			Success: false,
			Error:   "Update method should return (*T, error)",
		}, nil
	}

	errValue := results[1]
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		return &Response{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	data := results[0].Interface()
	return &Response{
		Success: true,
		Data:    data,
		Message: fmt.Sprintf("Entity with id %s updated successfully", id),
	}, nil
}

// handleDelete executes a delete operation using reflection
func (h *CommandHandler) handleDelete(ctx context.Context, repo interface{}, id string) (*Response, error) {
	if id == "" {
		return &Response{
			Success: false,
			Error:   "id is required for delete operation",
		}, nil
	}

	repoValue := reflect.ValueOf(repo)
	method := repoValue.MethodByName("Delete")
	if !method.IsValid() {
		return &Response{
			Success: false,
			Error:   "repository does not implement Delete method",
		}, nil
	}

	// Call Delete(ctx, id)
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
	})

	if len(results) != 1 {
		return &Response{
			Success: false,
			Error:   "Delete method should return error",
		}, nil
	}

	errValue := results[0]
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		return &Response{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &Response{
		Success: true,
		Message: fmt.Sprintf("Entity with id %s deleted successfully", id),
	}, nil
}

// handleCount executes a count operation using reflection
func (h *CommandHandler) handleCount(ctx context.Context, repo interface{}, schema ResourceSchema, payload interface{}) (*Response, error) {
	// Convert payload to filter type
	filter, err := h.convertPayload(payload, schema.FilterType)
	if err != nil {
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("invalid filter payload: %v", err),
		}, nil
	}

	repoValue := reflect.ValueOf(repo)
	method := repoValue.MethodByName("Count")
	if !method.IsValid() {
		return &Response{
			Success: false,
			Error:   "repository does not implement Count method",
		}, nil
	}

	// Call Count(ctx, filter)
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(filter).Elem(),
	})

	if len(results) != 2 {
		return &Response{
			Success: false,
			Error:   "Count method should return (int, error)",
		}, nil
	}

	errValue := results[1]
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		return &Response{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	count := int(results[0].Int())
	return &Response{
		Success: true,
		Count:   &count,
		Message: fmt.Sprintf("Found %d items", count),
	}, nil
}

// convertPayload converts a payload to the target type using JSON marshaling/unmarshaling
func (h *CommandHandler) convertPayload(payload interface{}, targetType interface{}) (interface{}, error) {
	if payload == nil {
		// Create a new instance of the target type
		targetValue := reflect.New(reflect.TypeOf(targetType))
		return targetValue.Interface(), nil
	}

	// First, marshal the payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create a new instance of the target type
	targetValue := reflect.New(reflect.TypeOf(targetType))
	targetPtr := targetValue.Interface()

	// Unmarshal into the target type
	if err := json.Unmarshal(payloadJSON, targetPtr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload to target type: %w", err)
	}

	return targetPtr, nil
}