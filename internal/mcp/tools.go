package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/studyguides-com/study-guides-api/internal/repository"
)

// SimpleToolGenerator generates basic OpenAI tools for testing
type SimpleToolGenerator struct {
	registry *RepositoryRegistry
}

// NewSimpleToolGenerator creates a simple tool generator for prototyping
func NewSimpleToolGenerator(registry *RepositoryRegistry) *SimpleToolGenerator {
	return &SimpleToolGenerator{
		registry: registry,
	}
}

// GenerateTools creates basic OpenAI tools for registered repositories
func (g *SimpleToolGenerator) GenerateTools() []openai.Tool {
	var tools []openai.Tool

	for _, resource := range g.registry.ListResources() {
		// Generate basic CRUD tools for each resource
		tools = append(tools, g.generateFindTool(resource))
		tools = append(tools, g.generateCountTool(resource))
		tools = append(tools, g.generateFindByIDTool(resource))
		
		// Generate resource-specific tools
		resourceTools := g.generateResourceSpecificTools(resource)
		tools = append(tools, resourceTools...)
	}

	return tools
}

// generateFindTool creates a simple find tool
func (g *SimpleToolGenerator) generateFindTool(resource string) openai.Tool {
	baseProperties := map[string]interface{}{
		"public": map[string]interface{}{
			"type":        "boolean",
			"description": "Filter by public/private status",
		},
		"type": map[string]interface{}{
			"type":        "string", 
			"description": "Filter by type (e.g., Category, Topic, etc.)",
		},
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Search by name (partial match)",
		},
		"limit": map[string]interface{}{
			"type":        "integer",
			"description": "Maximum number of results to return",
		},
	}
	
	// Add resource-specific properties
	if resource == "indexing" {
		baseProperties["triggerReindex"] = map[string]interface{}{
			"type":        "boolean",
			"description": "Trigger a new indexing job (set to true for reindexing)",
		}
		baseProperties["objectType"] = map[string]interface{}{
			"type":        "string",
			"description": "Type of object to index (Tag, User, Contact, FAQ)",
			"enum":        []string{"Tag"},
		}
		baseProperties["force"] = map[string]interface{}{
			"type":        "boolean",
			"description": "Force reindex even if content hasn't changed",
		}
		baseProperties["status"] = map[string]interface{}{
			"type":        "string",
			"description": "Filter by job status",
			"enum":        []string{"running", "complete", "failed"},
		}
	}
	
	description := fmt.Sprintf("Find %s entities with optional filters", resource)
	if resource == "indexing" {
		description = "Trigger indexing operations or check status of running indexing jobs. Use triggerReindex:true with objectType to start reindexing."
	}
	
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        fmt.Sprintf("%s_find", resource),
			Description: description,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filter": map[string]interface{}{
						"type":        "object",
						"description": fmt.Sprintf("Filter criteria for %s entities", resource),
						"properties":  baseProperties,
					},
				},
			},
		},
	}
}

// generateCountTool creates a simple count tool
func (g *SimpleToolGenerator) generateCountTool(resource string) openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        fmt.Sprintf("%s_count", resource),
			Description: fmt.Sprintf("Count %s entities with optional filters", resource),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"filter": map[string]interface{}{
						"type":        "object",
						"description": fmt.Sprintf("Filter criteria for counting %s entities", resource),
						"properties": map[string]interface{}{
							"public": map[string]interface{}{
								"type":        "boolean",
								"description": "Filter by public/private status",
							},
							"type": map[string]interface{}{
								"type":        "string",
								"description": "Filter by type",
							},
						},
					},
				},
			},
		},
	}
}

// generateFindByIDTool creates a simple findById tool
func (g *SimpleToolGenerator) generateFindByIDTool(resource string) openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        fmt.Sprintf("%s_findById", resource),
			Description: fmt.Sprintf("Find a specific %s by its ID", resource),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("The ID of the %s to retrieve", resource),
					},
				},
				"required": []string{"id"},
			},
		},
	}
}

// generateResourceSpecificTools creates resource-specific tools
func (g *SimpleToolGenerator) generateResourceSpecificTools(resource string) []openai.Tool {
	var tools []openai.Tool
	
	switch resource {
	case "kpi":
		tools = append(tools, g.generateKPIListGroupsTool())
	}
	
	return tools
}

// generateKPIListGroupsTool creates a tool to list available KPI groups
func (g *SimpleToolGenerator) generateKPIListGroupsTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "kpi_list_groups",
			Description: "List all available KPI groups that can be executed",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

// ParseToolCall parses an OpenAI tool call into a Command
func (g *SimpleToolGenerator) ParseToolCall(toolCall openai.ToolCall) (*repository.Command, error) {
	// Handle special tools first
	if toolCall.Function.Name == "kpi_list_groups" {
		return &repository.Command{
			Resource:  "kpi",
			Operation: repository.CRUDOperation("list_groups"),
		}, nil
	}

	// Parse the function name to extract resource and operation
	parts := strings.Split(toolCall.Function.Name, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid tool name format: %s", toolCall.Function.Name)
	}

	resource := parts[0]
	operation := repository.CRUDOperation(parts[1])

	// Parse the arguments
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
	}

	cmd := &repository.Command{
		Resource:  resource,
		Operation: operation,
	}

	// Extract operation-specific parameters
	switch operation {
	case repository.OperationFindByID:
		if id, ok := args["id"].(string); ok {
			cmd.ID = id
		}
	case repository.OperationFind, repository.OperationCount:
		if filter, ok := args["filter"]; ok {
			cmd.Payload = filter
		}
	}

	return cmd, nil
}