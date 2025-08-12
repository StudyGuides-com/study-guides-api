package main

import (
	"context"
	"fmt"
	"log"
	"time"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := chatpb.NewChatServiceClient(conn)
	ctx := context.Background()

	// Test cases for KPI operations
	testCases := []struct {
		name    string
		message string
	}{
		{
			name:    "Check running KPIs",
			message: "check running KPIs",
		},
		{
			name:    "Run monthly interactions",
			message: "calculate monthly interactions statistics",
		},
		{
			name:    "Run user statistics",
			message: "update user stats",
		},
		{
			name:    "Run tag statistics",
			message: "update tag statistics",
		},
		{
			name:    "Run all KPIs",
			message: "run all KPIs",
		},
		{
			name:    "Check status again",
			message: "how many KPIs are running?",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n📊 Test: %s\n", tc.name)
		fmt.Printf("   Message: %s\n", tc.message)
		
		req := &chatpb.ChatRequest{
			Context: &chatpb.Context{
				UserId:    "test-user",
				SessionId: "test-session",
				Metadata:  make(map[string]string),
			},
			Message: tc.message,
		}

		resp, err := client.Chat(ctx, req)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("   ✅ Response: %s\n", resp.Answer)
		
		// For longer operations, wait a bit before next test
		if tc.name == "Run all KPIs" {
			fmt.Println("   ⏳ Waiting 5 seconds for KPIs to start...")
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Println("\n✨ All tests completed!")
}