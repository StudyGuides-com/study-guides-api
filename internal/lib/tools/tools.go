package tools

import (
	"github.com/sashabaranov/go-openai"
)

func GetClassificationTools() []openai.Tool {
	tools := []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "GetTagCount",
				Description: "Count tags. Optional filters: type and contextType.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Tag type filter, e.g. 'Course'",
						},
						"contextType": map[string]interface{}{
							"type":        "string",
							"description": "Context type filter, e.g. 'College'",
						},
					},
					"required": []string{},
				},
			},
		},
	}
	return tools
}
