package tools

import (
	"github.com/sashabaranov/go-openai"
)

type ToolNames string

const (
	ToolNameTagCount     ToolNames = "TagCount"
	ToolNameListTags     ToolNames = "ListTags"
	ToolNameListRootTags ToolNames = "ListRootTags"
	ToolNameUniqueTagTypes ToolNames = "UniqueTagTypes"
	ToolNameUniqueContextTypes ToolNames = "UniqueContextTypes"
	ToolNameUnknown      ToolNames = "Unknown"
)

// ClassificationToolDefinitions contains all available tool definitions for classification
var ClassificationToolDefinitions = []ToolDefinition{
	NewToolDefinition(
		string(ToolNameTagCount),
		"Returns the number of tags. Use 'type' for tag categories (Course, Subject, etc.) and 'contextType' for organizational contexts (College, DoD, etc.). Optional filters: type, contextType, and public status.",
	).WithParameters(NoRequiredParams, typeProperty, contextProperty, publicProperty, formatProperty),
	
	NewToolDefinition(
		string(ToolNameListTags),
		"Returns a list of tags. Use 'type' for tag categories (Course, Subject, etc.) and 'contextType' for organizational contexts (College, DoD, etc.). Optional filters: type, contextType, and public status. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters(NoRequiredParams, typeProperty, contextProperty, publicProperty, formatProperty),
	
	NewToolDefinition(
		string(ToolNameListRootTags),
		"Returns a list of root tags (tags with no parent). Optional filter: public status. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters(NoRequiredParams, publicProperty, formatProperty),
	
	NewToolDefinition(
		string(ToolNameUniqueTagTypes),
		"Returns a list of all unique tag types available in the system.",
	).WithParameters(NoRequiredParams),
	
	NewToolDefinition(
		string(ToolNameUniqueContextTypes),
		"Returns a list of all unique context types (organizational contexts) available in the system.",
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
