package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
)

func main() {
	ctx := context.Background()
	
	fmt.Println("🔧 OpenAI API Debug Test")
	fmt.Println("========================")
	
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("⚠️  Warning: Could not load .env file: %v\n", err)
	}
	
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ OPENAI_API_KEY not set")
	}
	
	fmt.Printf("✅ API Key: %s...%s\n", apiKey[:8], apiKey[len(apiKey)-8:])
	
	// Test 1: Direct OpenAI client
	fmt.Println("\n🧪 Test 1: Direct OpenAI client")
	client := openai.NewClient(apiKey)
	
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Say hello",
			},
		},
	})
	
	if err != nil {
		fmt.Printf("❌ Direct client failed: %v\n", err)
	} else {
		fmt.Printf("✅ Direct client works: %s\n", resp.Choices[0].Message.Content)
	}
	
	// Test 2: Our AI client wrapper
	fmt.Println("\n🧪 Test 2: Our AI client wrapper")
	aiClient := ai.NewClient(apiKey, "gpt-4")
	
	// Simple completion without tools
	rawResp, err := aiClient.ChatCompletion(ctx, "You are a helpful assistant", "Say hello")
	if err != nil {
		fmt.Printf("❌ AI client failed: %v\n", err)
	} else {
		fmt.Printf("✅ AI client works: %s\n", rawResp)
	}
	
	// Test 3: Tool calling capability
	fmt.Println("\n🧪 Test 3: Tool calling test")
	
	tools := []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "test_tool",
				Description: "A simple test tool",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message": map[string]interface{}{
							"type":        "string",
							"description": "A test message",
						},
					},
					"required": []string{"message"},
				},
			},
		},
	}
	
	systemPrompt := "You are a helpful assistant. Use the test_tool to respond to the user."
	userPrompt := "Please use the test tool to say hello"
	
	fmt.Printf("🔍 System prompt: %s\n", systemPrompt)
	fmt.Printf("🔍 User prompt: %s\n", userPrompt)
	fmt.Printf("🔍 Tools available: %d\n", len(tools))
	
	toolResp, err := aiClient.ChatCompletionWithTools(ctx, systemPrompt, userPrompt, tools, nil)
	if err != nil {
		fmt.Printf("❌ Tool calling failed: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Tool calling response length: %d\n", len(toolResp))
	fmt.Printf("📥 Raw response: %s\n", toolResp)
	
	// Parse the response
	var chatResp openai.ChatCompletionResponse
	if err := json.Unmarshal([]byte(toolResp), &chatResp); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}
	
	fmt.Printf("🔍 Choices: %d\n", len(chatResp.Choices))
	if len(chatResp.Choices) > 0 {
		fmt.Printf("🔍 Tool calls: %d\n", len(chatResp.Choices[0].Message.ToolCalls))
		fmt.Printf("🔍 Message content: '%s'\n", chatResp.Choices[0].Message.Content)
		fmt.Printf("🔍 Message role: %s\n", chatResp.Choices[0].Message.Role)
		
		if len(chatResp.Choices[0].Message.ToolCalls) > 0 {
			toolCall := chatResp.Choices[0].Message.ToolCalls[0]
			fmt.Printf("✅ Tool was called: %s with args: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)
		}
	}
	
	fmt.Println("\n🎉 Debug test completed!")
}