package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	mcptesting "github.com/studyguides-com/study-guides-api/internal/mcp/testing"
)

func main() {
	ctx := context.Background()
	
	fmt.Println("🚀 MCP Prototype Test")
	fmt.Println("=====================")
	
	// Load environment variables from .env file
	fmt.Println("\n📂 Loading .env file...")
	if err := godotenv.Load(); err != nil {
		fmt.Printf("⚠️  Warning: Could not load .env file: %v\n", err)
		fmt.Println("   Continuing with system environment variables...")
	} else {
		fmt.Println("✅ .env file loaded successfully")
	}
	
	// Check required environment variables
	requiredEnvs := []string{
		"DATABASE_URL",
		"ALGOLIA_APP_ID", 
		"ALGOLIA_ADMIN_API_KEY",
		"OPENAI_API_KEY",
	}
	
	fmt.Println("\n🔍 Checking environment variables...")
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ Missing required environment variable: %s", env)
		}
		fmt.Printf("✅ %s: configured\n", env)
	}
	
	// Create prototype
	fmt.Println("\n🏗️  Initializing MCP prototype...")
	prototype, err := mcptesting.NewPrototype(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to create prototype: %v", err)
	}
	
	fmt.Println("✅ Prototype initialized successfully!")
	
	// Run demonstration
	if err := prototype.Demo(ctx); err != nil {
		log.Fatalf("❌ Demo failed: %v", err)
	}
	
	fmt.Println("\n🎉 MCP Prototype test completed!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review the generated tools and responses")
	fmt.Println("2. Test with your actual data")
	fmt.Println("3. Integrate with your ChatService")
	fmt.Println("4. Add more domain repositories (User, Question)")
}