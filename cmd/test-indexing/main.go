package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/mcp"
	"github.com/studyguides-com/study-guides-api/internal/mcp/indexing"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize store
	mainStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	// Initialize AI client
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}
	aiClient := ai.NewClient(openAIKey, "gpt-4o")

	// Create MCP processor
	mcpProcessor := mcp.NewMCPProcessor(aiClient)

	// Register indexing repository
	indexingRepo := indexing.NewIndexingRepositoryAdapter(mainStore.IndexingStore())
	mcpProcessor.Register(indexing.ResourceName, indexingRepo, indexing.GetResourceSchema())

	fmt.Println("🚀 Algolia Indexing Test Tool")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println()

	// Test different prompts
	testPrompts := []struct {
		name   string
		prompt string
	}{
		{
			name:   "Trigger Tag Reindex",
			prompt: "Please reindex all tags",
		},
		{
			name:   "Force Reindex Tags",
			prompt: "Force reindex tags even if they haven't changed",
		},
		{
			name:   "Check Indexing Status",
			prompt: "What is the status of current indexing jobs?",
		},
		{
			name:   "Count Running Jobs",
			prompt: "How many indexing jobs are currently running?",
		},
		{
			name:   "Sync to Algolia",
			prompt: "Sync tags to Algolia search",
		},
		{
			name:   "Index Tags",
			prompt: "Index all tags",
		},
		{
			name:   "Get Job Status",
			prompt: "Check the status of indexing operations",
		},
	}

	ctx := context.Background()

	for i, test := range testPrompts {
		fmt.Printf("\n%d. %s\n", i+1, test.name)
		fmt.Printf("   Prompt: \"%s\"\n", test.prompt)
		fmt.Println("   Processing...")

		// Process the request
		response, err := mcpProcessor.ProcessRequest(ctx, test.prompt)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
			continue
		}

		if !response.Success {
			fmt.Printf("   ❌ Failed: %s\n", response.Error)
			continue
		}

		fmt.Printf("   ✅ Success: %s\n", response.Message)

		// Pretty print the data if available
		if response.Data != nil {
			prettyPrint(response.Data)
		}

		// Add delay between tests
		if i < len(testPrompts)-1 {
			fmt.Println("\n   Waiting 2 seconds before next test...")
			time.Sleep(2 * time.Second)
		}
	}

	// Test direct API calls
	fmt.Println("\n" + string(make([]byte, 50)))
	fmt.Println("📋 Direct API Tests")
	fmt.Println(string(make([]byte, 50)))

	// Test starting an indexing job directly
	fmt.Println("\n1. Starting a direct indexing job for Tags...")
	jobID, err := mainStore.IndexingStore().StartIndexingJob(ctx, "Tag", false)
	if err != nil {
		fmt.Printf("   ❌ Failed to start job: %v\n", err)
	} else {
		fmt.Printf("   ✅ Job started with ID: %s\n", jobID)

		// Wait a moment then check status
		time.Sleep(3 * time.Second)
		
		status, err := mainStore.IndexingStore().GetJobStatus(ctx, jobID)
		if err != nil {
			fmt.Printf("   ❌ Failed to get job status: %v\n", err)
		} else {
			fmt.Printf("   📊 Job Status:\n")
			fmt.Printf("      - Status: %s\n", status.Status)
			fmt.Printf("      - Description: %s\n", status.Description)
			if status.StartedAt != nil {
				fmt.Printf("      - Started: %s\n", status.StartedAt.Format(time.RFC3339))
			}
			if status.CompletedAt != nil {
				fmt.Printf("      - Completed: %s\n", status.CompletedAt.Format(time.RFC3339))
			}
			if status.DurationSeconds != nil {
				fmt.Printf("      - Duration: %d seconds\n", *status.DurationSeconds)
			}
			if status.Metadata != nil {
				if items, ok := status.Metadata["itemsProcessed"].(float64); ok {
					fmt.Printf("      - Items Processed: %d\n", int(items))
				}
			}
		}
	}

	// Test listing running jobs
	fmt.Println("\n2. Listing running indexing jobs...")
	runningJobs, err := mainStore.IndexingStore().ListRunningJobs(ctx)
	if err != nil {
		fmt.Printf("   ❌ Failed to list running jobs: %v\n", err)
	} else {
		fmt.Printf("   📊 Found %d running jobs\n", len(runningJobs))
		for _, job := range runningJobs {
			fmt.Printf("      - Job %s: %s\n", job.ID, job.Description)
		}
	}

	// Test force reindex
	fmt.Println("\n3. Starting a force reindex job...")
	forceJobID, err := mainStore.IndexingStore().StartIndexingJob(ctx, "Tag", true)
	if err != nil {
		fmt.Printf("   ❌ Failed to start force job: %v\n", err)
	} else {
		fmt.Printf("   ✅ Force job started with ID: %s\n", forceJobID)
	}

	fmt.Println("\n✨ Test complete!")
}

// prettyPrint formats and prints data in a readable way
func prettyPrint(data interface{}) {
	// Try to marshal as JSON for pretty printing
	jsonData, err := json.MarshalIndent(data, "      ", "  ")
	if err != nil {
		fmt.Printf("      Data: %+v\n", data)
		return
	}
	
	// Parse back to check if it's an array
	var arr []interface{}
	if err := json.Unmarshal(jsonData, &arr); err == nil {
		fmt.Printf("      Results (%d items):\n", len(arr))
		for i, item := range arr {
			// Check if it's an IndexingExecution
			if execMap, ok := item.(map[string]interface{}); ok {
				fmt.Printf("      %d. Job %s:\n", i+1, execMap["id"])
				fmt.Printf("         - Object Type: %v\n", execMap["objectType"])
				fmt.Printf("         - Status: %v\n", execMap["status"])
				if msg, exists := execMap["message"]; exists && msg != "" {
					fmt.Printf("         - Message: %v\n", msg)
				}
				if force, exists := execMap["force"]; exists {
					fmt.Printf("         - Force: %v\n", force)
				}
				if items, exists := execMap["itemsProcessed"]; exists && items != 0 {
					fmt.Printf("         - Items Processed: %v\n", items)
				}
			} else {
				itemJSON, _ := json.MarshalIndent(item, "         ", "  ")
				fmt.Printf("      %d. %s\n", i+1, string(itemJSON))
			}
		}
	} else {
		// Not an array, print as is
		fmt.Printf("      Data:\n%s\n", string(jsonData))
	}
}

// Helper function to test tool call directly
func testDirectToolCall(processor *mcp.MCPProcessor) {
	fmt.Println("\n" + string(make([]byte, 50)))
	fmt.Println("🔧 Direct Tool Call Test")
	fmt.Println(string(make([]byte, 50)))
	
	// Create a tool call for indexing
	toolCall := openai.ToolCall{
		ID:   "test-1",
		Type: "function",
		Function: openai.FunctionCall{
			Name:      "indexing_find",
			Arguments: `{"triggerReindex": true, "objectType": "Tag", "force": false}`,
		},
	}
	
	fmt.Println("Executing tool call: indexing_find with triggerReindex=true")
	response, err := processor.ExecuteToolCall(context.Background(), toolCall)
	if err != nil {
		fmt.Printf("❌ Tool call failed: %v\n", err)
		return
	}
	
	if !response.Success {
		fmt.Printf("❌ Tool call unsuccessful: %s\n", response.Error)
		return
	}
	
	fmt.Printf("✅ Tool call successful: %s\n", response.Message)
	if response.Data != nil {
		prettyPrint(response.Data)
	}
}