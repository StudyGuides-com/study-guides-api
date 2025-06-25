package tools

import (
	"github.com/sashabaranov/go-openai"
)

type ToolNames string

const (
	ToolNameTagCount     ToolNames = "TagCount"
	ToolNameListTags     ToolNames = "ListTags"
	ToolNameUniqueTagTypes ToolNames = "UniqueTagTypes"
	ToolNameUnknown      ToolNames = "Unknown"
)

// ClassificationToolDefinitions contains all available tool definitions for classification
var ClassificationToolDefinitions = []ToolDefinition{
	NewToolDefinition(
		string(ToolNameTagCount),
		"Returns the number of tags. Optional filters: type and contextType.",
	).WithParameters(NoRequiredParams, typeProperty, contextProperty),
	
	NewToolDefinition(
		string(ToolNameListTags),
		"Returns a list of tags. Optional filters: type and contextType.",
	).WithParameters(NoRequiredParams, typeProperty, contextProperty),
	
	NewToolDefinition(
		string(ToolNameUniqueTagTypes),
		"Returns a list of all unique tag types available in the system.",
	).WithParameters(NoRequiredParams),
	
	NewToolDefinition(
		string(ToolNameUnknown),
		"Use when the user's request doesn't match any other available operations.",
	).WithParameters(NoRequiredParams),
}

// ClassificationToolMap provides efficient access to both tools and names
var ClassificationToolMap = func() map[string]openai.Tool {
	toolMap := make(map[string]openai.Tool)
	for _, toolDef := range ClassificationToolDefinitions {
		toolMap[toolDef.Name] = toolDef.AsTool()
	}
	return toolMap
}()

// GetClassificationData returns both the tool definitions and the tools map
func GetClassificationDefinitions() ([]ToolDefinition) {
	return ClassificationToolDefinitions
}

// GetClassificationTools returns the tools as a slice
func GetClassificationTools() []openai.Tool {
	tools := make([]openai.Tool, 0, len(ClassificationToolMap))
	for _, tool := range ClassificationToolMap {
		tools = append(tools, tool)
	}
	return tools
}
