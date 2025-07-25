package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
)

// MCPProcessor handles Model Context Protocol request processing
type MCPProcessor struct {
	registry     *RepositoryRegistry
	handler      *CommandHandler
	toolGen      *SimpleToolGenerator
	aiClient     ai.AiClient
	systemPrompt string
}

// NewMCPProcessor creates a new MCP processor
func NewMCPProcessor(aiClient ai.AiClient) *MCPProcessor {
	registry := NewRepositoryRegistry()
	handler := NewCommandHandler(registry)
	toolGen := NewSimpleToolGenerator(registry)

	return &MCPProcessor{
		registry:     registry,
		handler:      handler,
		toolGen:      toolGen,
		aiClient:     aiClient,
		systemPrompt: defaultSystemPrompt,
	}
}

// Register adds a repository to the processor
func (s *MCPProcessor) Register(resourceName string, repository interface{}, schema ResourceSchema) {
	s.registry.Register(resourceName, repository, schema)
}

// SetSystemPrompt allows customizing the AI system prompt
func (s *MCPProcessor) SetSystemPrompt(prompt string) {
	s.systemPrompt = prompt
}

// ProcessRequest processes a natural language request through the AI pipeline
func (s *MCPProcessor) ProcessRequest(ctx context.Context, userPrompt string) (*Response, error) {
	// Generate tools from registered repositories
	tools := s.toolGen.GenerateTools()
	
	if len(tools) == 0 {
		return &Response{
			Success: false,
			Error:   "No tools available. Please register repositories first.",
		}, nil
	}

	fmt.Printf("🤖 DEBUG: Processing request: '%s'\n", userPrompt)
	fmt.Printf("🛠️  DEBUG: Available tools: %d\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("   %d. %s - %s\n", i+1, tool.Function.Name, tool.Function.Description)
	}

	// Use AI to select and call appropriate tool
	fmt.Printf("🧠 DEBUG: Calling AI with system prompt length: %d chars\n", len(s.systemPrompt))
	rawResponse, err := s.aiClient.ChatCompletionWithTools(ctx, s.systemPrompt, userPrompt, tools, nil)
	if err != nil {
		fmt.Printf("❌ DEBUG: AI completion error: %v\n", err)
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("AI completion failed: %v", err),
		}, nil
	}

	fmt.Printf("📥 DEBUG: Raw AI response length: %d chars\n", len(rawResponse))
	fmt.Printf("📥 DEBUG: Raw AI response: %s\n", rawResponse)

	// Parse AI response
	var chatResp openai.ChatCompletionResponse
	if err := json.Unmarshal([]byte(rawResponse), &chatResp); err != nil {
		fmt.Printf("❌ DEBUG: Failed to parse AI response as JSON: %v\n", err)
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse AI response: %v", err),
		}, nil
	}

	fmt.Printf("🔍 DEBUG: AI response has %d choices\n", len(chatResp.Choices))
	if len(chatResp.Choices) == 0 {
		fmt.Printf("❌ DEBUG: AI returned no choices\n")
		return &Response{
			Success: false,
			Error:   "AI returned no choices",
		}, nil
	}

	fmt.Printf("🔍 DEBUG: First choice has %d tool calls\n", len(chatResp.Choices[0].Message.ToolCalls))
	fmt.Printf("🔍 DEBUG: First choice message content: '%s'\n", chatResp.Choices[0].Message.Content)
	
	if len(chatResp.Choices[0].Message.ToolCalls) == 0 {
		fmt.Printf("❌ DEBUG: AI did not select any tools\n")
		fmt.Printf("🔍 DEBUG: Message role: %s\n", chatResp.Choices[0].Message.Role)
		return &Response{
			Success: false,
			Error:   "AI did not select any tools",
		}, nil
	}

	// Execute the selected tool
	toolCall := chatResp.Choices[0].Message.ToolCalls[0]
	fmt.Printf("🔧 DEBUG: Executing tool call: %s with args: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)
	return s.ExecuteToolCall(ctx, toolCall)
}

// ExecuteToolCall executes a specific tool call
func (s *MCPProcessor) ExecuteToolCall(ctx context.Context, toolCall openai.ToolCall) (*Response, error) {
	// Parse tool call into command
	fmt.Printf("🔧 DEBUG: Parsing tool call: %s\n", toolCall.Function.Name)
	cmd, err := s.toolGen.ParseToolCall(toolCall)
	if err != nil {
		fmt.Printf("❌ DEBUG: Failed to parse tool call: %v\n", err)
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse tool call: %v", err),
		}, nil
	}

	fmt.Printf("🔧 DEBUG: Parsed command - Resource: %s, Operation: %s, Payload: %+v\n", cmd.Resource, cmd.Operation, cmd.Payload)

	// Execute command
	result, err := s.handler.Handle(ctx, *cmd)
	if err != nil {
		fmt.Printf("❌ DEBUG: Command execution failed: %v\n", err)
		return &Response{
			Success: false,
			Error:   fmt.Sprintf("Command execution failed: %v", err),
		}, nil
	}

	fmt.Printf("✅ DEBUG: Command executed successfully. Success: %v, Data length: %d\n", result.Success, len(fmt.Sprintf("%v", result.Data)))
	return result, nil
}

// DirectCommand allows direct execution of commands without AI processing
func (s *MCPProcessor) DirectCommand(ctx context.Context, cmd Command) (*Response, error) {
	return s.handler.Handle(ctx, cmd)
}

// GetAvailableTools returns all available tools for debugging/inspection
func (s *MCPProcessor) GetAvailableTools() []openai.Tool {
	return s.toolGen.GenerateTools()
}

// GetRegisteredResources returns a list of all registered resources
func (s *MCPProcessor) GetRegisteredResources() []string {
	return s.registry.ListResources()
}

// GetResourceSchema returns the schema for a specific resource
func (s *MCPProcessor) GetResourceSchema(resourceName string) (ResourceSchema, bool) {
	return s.registry.GetSchema(resourceName)
}

// Batch processing support
type BatchRequest struct {
	Commands []Command `json:"commands"`
}

type BatchResponse struct {
	Responses []Response `json:"responses"`
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
}

// ProcessBatch executes multiple commands in a single request
func (s *MCPProcessor) ProcessBatch(ctx context.Context, batch BatchRequest) (*BatchResponse, error) {
	responses := make([]Response, len(batch.Commands))
	successCount := 0

	for i, cmd := range batch.Commands {
		resp, err := s.handler.Handle(ctx, cmd)
		if err != nil {
			responses[i] = Response{
				Success: false,
				Error:   err.Error(),
			}
		} else {
			responses[i] = *resp
			if resp.Success {
				successCount++
			}
		}
	}

	return &BatchResponse{
		Responses: responses,
		Success:   successCount == len(batch.Commands),
		Message:   fmt.Sprintf("Executed %d commands, %d successful", len(batch.Commands), successCount),
	}, nil
}

// Health check
func (s *MCPProcessor) Health() map[string]interface{} {
	resources := s.registry.ListResources()
	tools := s.toolGen.GenerateTools()

	return map[string]interface{}{
		"status":           "healthy",
		"registeredResources": len(resources),
		"availableTools":   len(tools),
		"resources":        resources,
	}
}

// Default system prompt for the MCP server
const defaultSystemPrompt = `You are a data access assistant that MUST use tools to respond to user requests. You have access to various resources through specialized function tools.

CRITICAL: You MUST always call one of the available tools. Do not provide explanations or text responses - only use function calls.

Available operations:
- find: Search for entities with filters (use tag_find, user_find, etc.)
- findById: Get a specific entity by ID (use tag_findById, etc.)
- create: Create new entities
- update: Modify existing entities
- delete: Remove entities
- count: Count entities matching criteria (use tag_count, etc.)

Guidelines:
1. ALWAYS call a function tool - never respond with plain text
2. Use the most specific tool for the user's request
3. For searches, use appropriate filters to narrow results
4. For counts, use the count tool for that resource type
5. Map natural language to proper filter parameters

Examples:
- "find public tags" → call tag_find with {"filter": {"public": true}}
- "how many tags are there" → call tag_count with {}
- "find category tags" → call tag_find with {"filter": {"type": "Category"}}
- "find root tags" → call tag_find with {"filter": {"isRoot": true}}

Common tag types: Category, Topic, UserContent, UserTopic, Branch, Instruction_Type, Instruction_Group, Instruction, Chapter, Section

You must use function calling for every response.`