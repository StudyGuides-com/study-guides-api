package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
)

type Runner struct{}

func (r *Runner) runREPL(ctx context.Context, cmd *cobra.Command, systemPrompt string) error {
	fmt.Println("🤖 AI Chat REPL - Type 'exit', 'quit', or 'bye' to end the session")
	fmt.Println("Type your message and press Enter:")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	conn, err := grpc.Dial("localhost:1973", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	rpcClient := chatpb.NewChatServiceClient(conn)

	// Initialize conversation context
	conversationContext := &chatpb.Context{
		UserId:    "user123",
		SessionId: "session456",
		Metadata:  make(map[string]string),
	}

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Check for exit commands
		if input == "" {
			continue
		}

		lowerInput := strings.ToLower(input)
		if lowerInput == "exit" || lowerInput == "quit" || lowerInput == "bye" {
			fmt.Println("👋 Goodbye!")
			break
		}

		// Create context with JWT token in metadata
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjMiLCJuYW1lIjoiQnJ1Y2UgU3RvY2t3ZWxsIiwiZW1haWwiOiJicnVjZS5zdG9ja3dlbGxAZ21haWwuY29tIiwicm9sZXMiOlsidXNlciIsImFkbWluIl0sImlhdCI6MTc0OTU2MjMwMywiZXhwIjoxNzQ5NjQ4NzAzfQ.I6AdQ14oeQC7jVP4M3teR5_YZVI62YRsdhb4EKOjR_g",
		})
		ctxWithAuth := metadata.NewOutgoingContext(ctx, md)

		// Call gRPC service with authenticated context and conversation context
		resp, err := rpcClient.Chat(ctxWithAuth, &chatpb.ChatRequest{
			Context: conversationContext,
			Message: input,
		})
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		// Update conversation context with the response context for next request
		conversationContext = resp.Context

		fmt.Println("🤖 " + resp.Answer)
	}

	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "test-client",
		Short: "Test client for the chat service",
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := &Runner{}
			return runner.runREPL(context.Background(), cmd, "")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
