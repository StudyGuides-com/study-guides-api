package mcp

import (
	"context"
	"fmt"
	"strings"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/mcp/tag"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

// MCPIntegratedChatService demonstrates how to integrate MCP with your existing ChatService
type MCPIntegratedChatService struct {
	chatpb.UnimplementedChatServiceServer
	
	// New MCP system
	mcpProcessor *MCPProcessor
	
	// Keep existing system for fallback/comparison
	// router router.Router  // Your existing router
	// ai     ai.AiClient    // Your existing AI client
	
	// Feature flag to switch between systems
	useMCP bool
}

// NewMCPIntegratedChatService creates a chat service with MCP integration
func NewMCPIntegratedChatService(store store.Store, aiClient ai.AiClient, useMCP bool) (*MCPIntegratedChatService, error) {
	service := &MCPIntegratedChatService{
		useMCP: useMCP,
	}
	
	// Initialize MCP processor
	mcpProcessor := NewMCPProcessor(aiClient)
	
	// Register repositories
	tagRepo := tag.NewTagRepositoryAdapter(store.TagStore())
	mcpProcessor.Register(tag.ResourceName, tagRepo, tag.GetResourceSchema())
	
	service.mcpProcessor = mcpProcessor
	
	// TODO: Initialize existing router/ai for fallback
	// service.router = router.NewRouter(store)
	// service.ai = aiClient
	
	return service, nil
}

// Chat handles chat requests using either MCP or legacy system
func (s *MCPIntegratedChatService) Chat(ctx context.Context, req *chatpb.ChatRequest) (*chatpb.ChatResponse, error) {
	if s.useMCP {
		return s.handleWithMCP(ctx, req)
	}
	
	// Fallback to existing system
	return s.handleWithLegacySystem(ctx, req)
}

// handleWithMCP processes the request using the new MCP system
func (s *MCPIntegratedChatService) handleWithMCP(ctx context.Context, req *chatpb.ChatRequest) (*chatpb.ChatResponse, error) {
	// Process request through MCP
	response, err := s.mcpProcessor.ProcessRequest(ctx, req.Message)
	if err != nil {
		return &chatpb.ChatResponse{
			Answer: fmt.Sprintf("MCP Error: %v", err),
			Context: req.Context, // Preserve context
		}, nil
	}
	
	// Format response for chat
	answer := s.formatMCPResponse(response)
	
	// Create updated context (simplified - you may want to preserve conversation history)
	updatedContext := &chatpb.Context{
		UserId:    req.Context.GetUserId(),
		SessionId: req.Context.GetSessionId(),
		Metadata:  req.Context.GetMetadata(), // Keep existing metadata for now
	}
	
	return &chatpb.ChatResponse{
		Answer:     answer,
		Context:    updatedContext,
		Operation:  "mcp_processed", // Indicate this was processed by MCP
		Parameters: map[string]string{"system": "mcp"},
	}, nil
}

// handleWithLegacySystem processes the request using your existing system
func (s *MCPIntegratedChatService) handleWithLegacySystem(ctx context.Context, req *chatpb.ChatRequest) (*chatpb.ChatResponse, error) {
	// TODO: Implement fallback to your existing chat logic
	// This would be your current Chat method implementation
	return &chatpb.ChatResponse{
		Answer: "Legacy system not implemented in this prototype",
		Context: req.Context,
	}, fmt.Errorf("legacy system fallback not implemented")
}

// formatMCPResponse converts MCP response to human-readable format
func (s *MCPIntegratedChatService) formatMCPResponse(response *Response) string {
	if !response.Success {
		return fmt.Sprintf("I encountered an error: %s", response.Error)
	}
	
	var result strings.Builder
	
	// Add the main message
	if response.Message != "" {
		result.WriteString(response.Message)
	}
	
	// Add count information if available
	if response.Count != nil {
		if result.Len() > 0 {
			result.WriteString(" ")
		}
		result.WriteString(fmt.Sprintf("(%d items)", *response.Count))
	}
	
	// Add data if it's a reasonable size for chat
	if response.Data != nil {
		// For now, just indicate data is available
		// You could format specific types differently
		if result.Len() > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString("📊 Data retrieved successfully. ")
		
		// You could add specific formatting based on the type of data
		// For example, format tag lists, user info, etc.
	}
	
	return result.String()
}

// ToggleMCP allows switching between MCP and legacy systems at runtime
func (s *MCPIntegratedChatService) ToggleMCP(useMCP bool) {
	s.useMCP = useMCP
}

// GetMCPHealth returns health information about the MCP system
func (s *MCPIntegratedChatService) GetMCPHealth() map[string]interface{} {
	if s.mcpProcessor != nil {
		return s.mcpProcessor.Health()
	}
	return map[string]interface{}{
		"status": "not_initialized",
	}
}

// CompareResponses processes the same query through both systems for comparison
func (s *MCPIntegratedChatService) CompareResponses(ctx context.Context, req *chatpb.ChatRequest) (*ComparisonResult, error) {
	result := &ComparisonResult{
		Query: req.Message,
	}
	
	// Test MCP system
	originalUseMCP := s.useMCP
	s.useMCP = true
	
	mcpResp, mcpErr := s.Chat(ctx, req)
	result.MCPResponse = mcpResp
	result.MCPError = mcpErr
	
	// Test legacy system
	s.useMCP = false
	legacyResp, legacyErr := s.Chat(ctx, req)
	result.LegacyResponse = legacyResp
	result.LegacyError = legacyErr
	
	// Restore original setting
	s.useMCP = originalUseMCP
	
	return result, nil
}

// ComparisonResult holds the results of comparing both systems
type ComparisonResult struct {
	Query          string                 `json:"query"`
	MCPResponse    *chatpb.ChatResponse   `json:"mcpResponse"`
	MCPError       error                  `json:"mcpError"`
	LegacyResponse *chatpb.ChatResponse   `json:"legacyResponse"`
	LegacyError    error                  `json:"legacyError"`
}

// String formats the comparison result for display
func (cr *ComparisonResult) String() string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("Query: %s\n", cr.Query))
	result.WriteString("=" + strings.Repeat("=", len(cr.Query)+7) + "\n\n")
	
	// MCP Results
	result.WriteString("🆕 MCP System:\n")
	if cr.MCPError != nil {
		result.WriteString(fmt.Sprintf("   ❌ Error: %v\n", cr.MCPError))
	} else if cr.MCPResponse != nil {
		result.WriteString(fmt.Sprintf("   ✅ Answer: %s\n", cr.MCPResponse.Answer))
		result.WriteString(fmt.Sprintf("   🔧 Operation: %s\n", cr.MCPResponse.Operation))
	}
	
	result.WriteString("\n")
	
	// Legacy Results  
	result.WriteString("🗂️  Legacy System:\n")
	if cr.LegacyError != nil {
		result.WriteString(fmt.Sprintf("   ❌ Error: %v\n", cr.LegacyError))
	} else if cr.LegacyResponse != nil {
		result.WriteString(fmt.Sprintf("   ✅ Answer: %s\n", cr.LegacyResponse.Answer))
		result.WriteString(fmt.Sprintf("   🔧 Operation: %s\n", cr.LegacyResponse.Operation))
	}
	
	return result.String()
}

// Example usage function
func ExampleIntegration() {
	fmt.Println(`
Integration Example:

1. In your server initialization (cmd/server/main.go):

   // Create store and AI client as usual
   store, _ := store.NewStore()
   aiClient := ai.NewOpenAIClient(apiKey)
   
   // Create integrated chat service
   chatService, _ := mcp.NewMCPIntegratedChatService(store, aiClient, true)
   
   // Register with gRPC server
   chatpb.RegisterChatServiceServer(grpcServer, chatService)

2. For testing/comparison:

   // Toggle between systems
   chatService.ToggleMCP(false) // Use legacy
   chatService.ToggleMCP(true)  // Use MCP
   
   // Compare responses
   req := &chatpb.ChatRequest{Message: "find public tags"}
   comparison, _ := chatService.CompareResponses(ctx, req)
   fmt.Println(comparison.String())

3. For gradual migration:

   // Start with MCP disabled in production
   useMCP := os.Getenv("USE_MCP") == "true"
   chatService, _ := mcp.NewMCPIntegratedChatService(store, aiClient, useMCP)
   
   // Enable for specific users or operations
   if req.Context.UserId == "test-user" {
       chatService.ToggleMCP(true)
   }
`)
}