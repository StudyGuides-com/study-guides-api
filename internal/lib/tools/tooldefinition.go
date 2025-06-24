package tools

import (
	"github.com/sashabaranov/go-openai"
)

// ToolDefinition represents a tool definition with name, description, and parameters
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []Property
	Required    []string
}

// NewToolDefinition creates a new ToolDefinition
func NewToolDefinition(name, description string) ToolDefinition {
	return ToolDefinition{
		Name:        name,
		Description: description,
		Parameters:  []Property{},
		Required:    []string{},
	}
}

// WithParameters adds parameters to the tool definition
func (td ToolDefinition) WithParameters(required []string, params ...Property) ToolDefinition {
	td.Parameters = params
	td.Required = required
	return td
}

// AsTool converts the ToolDefinition to an openai.Tool
func (td ToolDefinition) AsTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        td.Name,
			Description: td.Description,
			Parameters:  BuildParameterSchemaFromProps(td.Required, td.Parameters...),
		},
	}
} 