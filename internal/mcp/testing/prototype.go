package testing

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/mcp"
	"github.com/studyguides-com/study-guides-api/internal/mcp/tag"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

// Prototype demonstrates the MCP system working with real data
type Prototype struct {
	processor *mcp.MCPProcessor
	store     store.Store
}

// NewPrototype creates a prototype instance using existing infrastructure
func NewPrototype(ctx context.Context) (*Prototype, error) {
	// Initialize existing store
	store, err := store.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}
	
	// Create AI client (using existing pattern)
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable required")
	}
	
	aiClient := ai.NewClient(openAIKey, "gpt-4")
	
	// Create MCP processor
	processor := mcp.NewMCPProcessor(aiClient)
	
	// Create and register tag repository adapter
	tagRepo := tag.NewTagRepositoryAdapter(store.TagStore())
	processor.Register(tag.ResourceName, tagRepo, tag.GetResourceSchema())
	
	return &Prototype{
		processor: processor,
		store:     store,
	}, nil
}

// TestBasicOperations runs a series of test operations to verify the system works
func (p *Prototype) TestBasicOperations(ctx context.Context) error {
	fmt.Println("🚀 Starting MCP Prototype Tests")
	fmt.Println("=====================================")
	
	// Test 1: Natural language find operation
	fmt.Println("\n📝 Test 1: Natural language tag search")
	response, err := p.processor.ProcessRequest(ctx, "find all public tags, limit to 5")
	if err != nil {
		return fmt.Errorf("test 1 failed: %w", err)
	}
	
	fmt.Printf("✅ Success: %v\n", response.Success)
	fmt.Printf("📊 Message: %s\n", response.Message)
	if response.Count != nil {
		fmt.Printf("🔢 Count: %d\n", *response.Count)
	}
	
	// Test 2: Direct command execution
	fmt.Println("\n📝 Test 2: Direct command execution")
	cmd := mcp.Command{
		Resource:  "tag",
		Operation: mcp.OperationCount,
		Payload: map[string]interface{}{
			"public": true,
		},
	}
	
	response, err = p.processor.DirectCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("test 2 failed: %w", err)
	}
	
	fmt.Printf("✅ Success: %v\n", response.Success)
	fmt.Printf("📊 Message: %s\n", response.Message)
	if response.Count != nil {
		fmt.Printf("🔢 Count: %d\n", *response.Count)
	}
	
	// Test 3: Type-specific queries
	fmt.Println("\n📝 Test 3: Find tags by type")
	response, err = p.processor.ProcessRequest(ctx, "find tags with type Category, limit to 3")
	if err != nil {
		return fmt.Errorf("test 3 failed: %w", err)
	}
	
	fmt.Printf("✅ Success: %v\n", response.Success)
	fmt.Printf("📊 Message: %s\n", response.Message)
	
	// Test 4: Root tags query
	fmt.Println("\n📝 Test 4: Find root tags")
	response, err = p.processor.ProcessRequest(ctx, "find root tags, show me 3")
	if err != nil {
		return fmt.Errorf("test 4 failed: %w", err)
	}
	
	fmt.Printf("✅ Success: %v\n", response.Success)
	fmt.Printf("📊 Message: %s\n", response.Message)
	
	// Test 5: Get unique tag types
	fmt.Println("\n📝 Test 5: Get unique tag types")
	response, err = p.processor.ProcessRequest(ctx, "what tag types are available?")
	if err != nil {
		return fmt.Errorf("test 5 failed: %w", err)
	}
	
	fmt.Printf("✅ Success: %v\n", response.Success)
	fmt.Printf("📊 Message: %s\n", response.Message)
	
	fmt.Println("\n🎉 All tests completed successfully!")
	return nil
}

// TestComparison runs the same query through both old and new systems for comparison
func (p *Prototype) TestComparison(ctx context.Context, userQuery string) error {
	fmt.Printf("\n🔍 Comparing systems for query: '%s'\n", userQuery)
	fmt.Println("=" + strings.Repeat("=", len(userQuery)+35))
	
	// Test new MCP system
	fmt.Println("\n🆕 New MCP System:")
	start := time.Now()
	newResponse, err := p.processor.ProcessRequest(ctx, userQuery)
	newDuration := time.Since(start)
	
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success: %v\n", newResponse.Success)
		fmt.Printf("📊 Message: %s\n", newResponse.Message)
		if newResponse.Count != nil {
			fmt.Printf("🔢 Count: %d\n", *newResponse.Count)
		}
		fmt.Printf("⏱️  Duration: %v\n", newDuration)
	}
	
	// For comparison with old system, you could add:
	// fmt.Println("\n🗂️  Old System:")
	// (Run through existing chat service logic)
	
	return nil
}

// GetSystemHealth returns diagnostic information about the MCP system
func (p *Prototype) GetSystemHealth() map[string]interface{} {
	health := p.processor.Health()
	health["store_available"] = p.store != nil
	return health
}

// ShowAvailableTools displays all auto-generated tools for debugging
func (p *Prototype) ShowAvailableTools() {
	fmt.Println("\n🛠️  Available Auto-Generated Tools:")
	fmt.Println("=====================================")
	
	tools := p.processor.GetAvailableTools()
	for i, tool := range tools {
		fmt.Printf("\n%d. %s\n", i+1, tool.Function.Name)
		fmt.Printf("   Description: %s\n", tool.Function.Description)
		fmt.Printf("   Parameters: [configured]\n")
	}
	
	fmt.Printf("\nTotal: %d tools generated\n", len(tools))
}

// Demo runs a complete demonstration of the MCP system
func (p *Prototype) Demo(ctx context.Context) error {
	fmt.Println("🎭 MCP System Demonstration")
	fmt.Println("============================")
	
	// Show system health
	fmt.Println("\n💊 System Health:")
	health := p.GetSystemHealth()
	for key, value := range health {
		fmt.Printf("   %s: %v\n", key, value)
	}
	
	// Show available tools
	p.ShowAvailableTools()
	
	// Run test operations
	if err := p.TestBasicOperations(ctx); err != nil {
		return err
	}
	
	// Test some comparison queries
	testQueries := []string{
		"find public tags",
		"how many tags are there?",
		"show me some category tags",
		"what are the root tags?",
	}
	
	for _, query := range testQueries {
		if err := p.TestComparison(ctx, query); err != nil {
			fmt.Printf("❌ Comparison failed for '%s': %v\n", query, err)
		}
	}
	
	fmt.Println("\n🏁 Demo completed!")
	return nil
}