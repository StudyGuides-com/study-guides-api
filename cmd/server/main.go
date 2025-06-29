package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	healthpb "github.com/studyguides-com/study-guides-api/api/v1/health"
	questionpb "github.com/studyguides-com/study-guides-api/api/v1/question"
	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	tagpb "github.com/studyguides-com/study-guides-api/api/v1/tag"
	userpb "github.com/studyguides-com/study-guides-api/api/v1/user"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/store"

	"github.com/joho/godotenv"
	"github.com/studyguides-com/study-guides-api/internal/lib/router"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/services"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func parseEnvAsInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func parseEnvAsRate(key string, fallback rate.Limit) rate.Limit {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fallback
	}
	return rate.Limit(val)
}

func main() {

	_ = godotenv.Load() // Ignore errors; default to environment

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.ErrorUnaryInterceptor(),
			middleware.AuthUnaryInterceptor(os.Getenv("JWT_SECRET")),
			middleware.RateLimitUnaryInterceptor(
				parseEnvAsRate("RATE_LIMIT_USER_PER_SECOND", 1.0),
				parseEnvAsInt("RATE_LIMIT_USER_BURST", 5),
			),
		),
	)

	appStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	// Register Health Service
	healthpb.RegisterHealthServiceServer(grpcServer, services.NewHealthService())

	// Register Search Service
	searchpb.RegisterSearchServiceServer(grpcServer, services.NewSearchService(appStore))

	// Register User Service
	userpb.RegisterUserServiceServer(grpcServer, services.NewUserService(appStore))

	// Register Tag Service
	tagpb.RegisterTagServiceServer(grpcServer, services.NewTagService(appStore))

	// Register Question Service
	questionpb.RegisterQuestionServiceServer(grpcServer, services.NewQuestionService(appStore))

	// Register Chat Service
	router := router.NewRouter(appStore)
	ai := ai.NewClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"))
	chatpb.RegisterChatServiceServer(grpcServer, services.NewChatService(router, ai))

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	// Start a separate HTTP health check server
	go func() {
		httpPort := "8080"
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		log.Printf("HTTP health server listening on :%s", httpPort)
		if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
			log.Printf("HTTP health server error: %v", err)
		}
	}()

	// h2c handler for gRPC
	h2cHandler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		grpcServer.ServeHTTP(w, r)
	}), &http2.Server{})

	server := &http.Server{
		Addr:    address,
		Handler: h2cHandler,
	}

	// Create a channel to receive shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the h2c gRPC server in a goroutine
	go func() {
		log.Printf("gRPC (h2c) server listening on %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to serve: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Printf("\n")
	log.Println("‼️ detected shutdown signal, starting shutdown...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully stop the server
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		log.Println("✅ server gracefully stopped")
	case <-ctx.Done():
		log.Println("❌ graceful shutdown timeout, forcing stop")
		grpcServer.Stop()
	}

	log.Println("🙌 server shutdown complete")
}
