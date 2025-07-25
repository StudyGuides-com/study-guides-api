package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/studyguides-com/study-guides-api/internal/errors"
)

type AiClient interface {
	// ChatCompletion generates a chat completion using the AI service
	ChatCompletion(ctx context.Context, systemPrompt, userPrompt string) (string, error)

	// ChatCompletionWithTools generates a chat completion with JSON response format and tools
	ChatCompletionWithTools(ctx context.Context, systemPrompt, userPrompt string, tools []openai.Tool, toolChoice *openai.ToolChoice) (string, error)

	// ChatCompletionWithHistory generates a chat completion with conversation history
	ChatCompletionWithHistory(ctx context.Context, systemPrompt string, messages []openai.ChatCompletionMessage, tools []openai.Tool, toolChoice *openai.ToolChoice) (string, error)
}

type OpenAiClient struct {
	client *openai.Client
	apiKey string
	model  string
}

func NewClient(apiKey string, model string) *OpenAiClient {
	return &OpenAiClient{
		client: openai.NewClient(apiKey),
		apiKey: apiKey,
		model:  model,
	}
}

// ChatCompletion generates a chat completion using the OpenAI API
func (c *OpenAiClient) ChatCompletion(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// Validate prompts
	if systemPrompt == "" {
		return "", errors.ErrSystemPromptEmpty
	}
	if userPrompt == "" {
		return "", errors.ErrUserPromptEmpty
	}

	// Create the chat completion request
	req := openai.ChatCompletionRequest{
		Model:            c.model,
		Temperature:      0.7,
		MaxTokens:        512,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
	}

	// Make the API call
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	// Return the first message's content
	if len(resp.Choices) == 0 {
		return "", errors.ErrNoCompletionChoicesReturned
	}

	return resp.Choices[0].Message.Content, nil
}

// ChatCompletionWithTools generates a chat completion with JSON response format and tools
func (c *OpenAiClient) ChatCompletionWithTools(ctx context.Context, systemPrompt, userPrompt string, tools []openai.Tool, toolChoice *openai.ToolChoice) (string, error) {

	// Validate prompts
	if systemPrompt == "" {
		return "", errors.ErrSystemPromptEmpty
	}
	if userPrompt == "" {
		return "", errors.ErrUserPromptEmpty
	}

	// Create messages
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}

	// Create the JSON chat completion request using the helper function
	req := CreateJSONChatCompletionRequest(c.model, messages, tools, toolChoice)

	// Make the API call
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion with tools: %w", err)
	}

	// Return the full response as JSON so MCP processor can parse tool calls
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(responseJSON), nil
}

// ChatCompletionWithHistory generates a chat completion with conversation history
func (c *OpenAiClient) ChatCompletionWithHistory(ctx context.Context, systemPrompt string, messages []openai.ChatCompletionMessage, tools []openai.Tool, toolChoice *openai.ToolChoice) (string, error) {
	// Validate system prompt
	if systemPrompt == "" {
		return "", errors.ErrSystemPromptEmpty
	}

	// Create messages with system prompt at the beginning
	allMessages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	// Add conversation history
	allMessages = append(allMessages, messages...)

	// Create the JSON chat completion request using the helper function
	req := CreateJSONChatCompletionRequest(c.model, allMessages, tools, toolChoice)

	// Make the API call
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion with tools: %w", err)
	}

	// Return the full response as JSON
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(responseJSON), nil
}

// CreateJSONChatCompletionRequest creates a chat completion request with tools
func CreateJSONChatCompletionRequest(model string, messages []openai.ChatCompletionMessage, tools []openai.Tool, toolChoice *openai.ToolChoice) openai.ChatCompletionRequest {
	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	// Only add JSON response format if no tools are provided
	// Tools and JSON response format are incompatible
	if tools == nil {
		req.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}

	if tools != nil {
		req.Tools = tools
	}

	if toolChoice != nil {
		req.ToolChoice = toolChoice
	}

	return req
}
