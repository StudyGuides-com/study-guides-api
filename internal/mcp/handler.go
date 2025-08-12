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
		return h.handleFind(ctx, repo, schema, cmd.Resource, cmd.Payload)
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
	case OperationListGroups:
		return h.handleListGroups(ctx, repo, cmd.Resource)
	default:
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("unknown operation: %s", cmd.Operation),
		}, nil
	}
}

// handleFind executes a find operation using reflection
func (h *CommandHandler) handleFind(ctx context.Context, repo interface{}, schema repository.ResourceSchema, resource string, payload interface{}) (*Response, error) {
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

	// Generate a more informative message based on resource type and operation
	message := h.generateFindMessage(resource, data, count)

	return &Response{
		Success: true,
		Data:    data,
		Count:   &count,
		Message: message,
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

// generateFindMessage creates context-appropriate messages for find operations
func (h *CommandHandler) generateFindMessage(resource string, data interface{}, count int) string {
	if resource == "kpi" {
		return h.generateKPIMessage(data, count)
	}
	
	// Default message for other resources
	return fmt.Sprintf("Found %d items", count)
}

// generateKPIMessage creates specific messages for KPI operations
func (h *CommandHandler) generateKPIMessage(data interface{}, count int) string {
	// Try to extract execution details from the data
	dataSlice := reflect.ValueOf(data)
	if dataSlice.Kind() != reflect.Slice {
		return fmt.Sprintf("Found %d items", count)
	}

	if count == 0 {
		return "No KPI executions found"
	}

	if count == 1 {
		// Single execution - show details
		if dataSlice.Len() > 0 {
			item := dataSlice.Index(0).Interface()
			
			// Handle both map[string]interface{} and direct struct types
			var group, status interface{}
			var duration interface{}
			var errorMsg interface{}
			
			if execution, ok := item.(map[string]interface{}); ok {
				group = execution["group"]
				status = execution["status"]
				duration = execution["duration"]
				errorMsg = execution["error"]
			} else {
				// Try to access as struct fields using reflection
				itemValue := reflect.ValueOf(item)
				if itemValue.Kind() == reflect.Ptr {
					itemValue = itemValue.Elem()
				}
				if itemValue.Kind() == reflect.Struct {
					if groupField := itemValue.FieldByName("Group"); groupField.IsValid() {
						group = groupField.Interface()
					}
					if statusField := itemValue.FieldByName("Status"); statusField.IsValid() {
						status = statusField.Interface()
					}
					if durationField := itemValue.FieldByName("Duration"); durationField.IsValid() && !durationField.IsNil() {
						duration = durationField.Interface()
					}
					if errorField := itemValue.FieldByName("Error"); errorField.IsValid() {
						errorMsg = errorField.Interface()
					}
				}
			}
			
			if group != nil && status != nil {
				statusStr := fmt.Sprintf("%v", status)
				groupStr := fmt.Sprintf("%v", group)
				
				switch statusStr {
				case "running":
					return fmt.Sprintf("KPI execution for %s is running", groupStr)
				case "complete":
					// Try to get duration info
					if duration != nil {
						return fmt.Sprintf("KPI execution for %s completed in %v", groupStr, duration)
					}
					return fmt.Sprintf("KPI execution for %s completed", groupStr)
				case "failed":
					if errorMsg != nil {
						return fmt.Sprintf("KPI execution for %s failed: %v", groupStr, errorMsg)
					}
					return fmt.Sprintf("KPI execution for %s failed", groupStr)
				default:
					return fmt.Sprintf("KPI execution for %s is %s", groupStr, statusStr)
				}
			} else if group != nil {
				return fmt.Sprintf("Started KPI execution for %v", group)
			}
		}
		return "Found 1 KPI execution"
	}

	// Multiple executions - show summary
	var groups []string
	runningCount := 0
	completedCount := 0
	failedCount := 0
	
	for i := 0; i < dataSlice.Len() && i < 5; i++ { // Show first 5 groups
		item := dataSlice.Index(i).Interface()
		
		// Handle both map[string]interface{} and direct struct types
		var group, status interface{}
		
		if execution, ok := item.(map[string]interface{}); ok {
			group = execution["group"]
			status = execution["status"]
		} else {
			// Try to access as struct fields using reflection
			itemValue := reflect.ValueOf(item)
			if itemValue.Kind() == reflect.Ptr {
				itemValue = itemValue.Elem()
			}
			if itemValue.Kind() == reflect.Struct {
				if groupField := itemValue.FieldByName("Group"); groupField.IsValid() {
					group = groupField.Interface()
				}
				if statusField := itemValue.FieldByName("Status"); statusField.IsValid() {
					status = statusField.Interface()
				}
			}
		}
		
		if group != nil {
			groups = append(groups, fmt.Sprintf("%v", group))
		}
		if status != nil {
			statusStr := fmt.Sprintf("%v", status)
			switch statusStr {
			case "running":
				runningCount++
			case "complete":
				completedCount++
			case "failed":
				failedCount++
			}
		}
	}

	// Generate a smart message based on the mix of statuses
	message := fmt.Sprintf("Found %d KPI executions", count)
	if len(groups) > 0 {
		if len(groups) <= 3 {
			message += fmt.Sprintf(" (%s)", joinStrings(groups, ", "))
		} else {
			message += fmt.Sprintf(" (%s and %d more)", joinStrings(groups[:3], ", "), len(groups)-3)
		}
	}
	
	// Add status summary
	var statusParts []string
	if runningCount > 0 {
		statusParts = append(statusParts, fmt.Sprintf("%d running", runningCount))
	}
	if completedCount > 0 {
		statusParts = append(statusParts, fmt.Sprintf("%d completed", completedCount))
	}
	if failedCount > 0 {
		statusParts = append(statusParts, fmt.Sprintf("%d failed", failedCount))
	}
	
	if len(statusParts) > 0 {
		message += ": " + joinStrings(statusParts, ", ")
	}

	return message
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// handleListGroups executes a list_groups operation for KPI resource
func (h *CommandHandler) handleListGroups(ctx context.Context, repo interface{}, resource string) (*Response, error) {
	if resource != "kpi" {
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("list_groups operation not supported for resource: %s", resource),
		}, nil
	}

	// For KPI, return the available groups
	groups := []map[string]interface{}{
		{"name": "MonthlyInteractions", "description": "Monthly user interaction statistics"},
		{"name": "Tags", "description": "Tag usage and hierarchy stats"},
		{"name": "TagTypes", "description": "Tag type distribution stats"},
		{"name": "Reports", "description": "Report generation statistics"},
		{"name": "Topics", "description": "Topic usage and engagement stats"},
		{"name": "MissingData", "description": "Data quality and completeness metrics"},
		{"name": "Ratings", "description": "Rating and feedback statistics"},
		{"name": "Questions", "description": "Question performance and usage stats"},
		{"name": "Users", "description": "User registration and activity metrics"},
		{"name": "UserContent", "description": "User-generated content statistics"},
		{"name": "Contacts", "description": "Contact and communication statistics"},
	}

	return &Response{
		Success: true,
		Data:    groups,
		Count:   &[]int{len(groups)}[0],
		Message: fmt.Sprintf("Available KPI groups: %s", joinGroupNames(groups)),
	}, nil
}

// joinGroupNames extracts and joins group names for the message
func joinGroupNames(groups []map[string]interface{}) string {
	var names []string
	for _, group := range groups {
		if name, ok := group["name"].(string); ok {
			names = append(names, name)
		}
	}
	return joinStrings(names, ", ")
}