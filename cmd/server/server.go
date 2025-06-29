package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	healthpb "github.com/studyguides-com/study-guides-api/api/v1/health"
	questionpb "github.com/studyguides-com/study-guides-api/api/v1/question"
	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	tagpb "github.com/studyguides-com/study-guides-api/api/v1/tag"
	userpb "github.com/studyguides-com/study-guides-api/api/v1/user"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/store"

	"github.com/studyguides-com/study-guides-api/internal/lib/router"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/services"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server manages both HTTP and gRPC traffic on a single server
type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{}
}

// Start starts the server that handles both HTTP and gRPC traffic
func (s *Server) Start(appStore store.Store) {
	port := getPort()
	address := ":" + port

	// Create gRPC server with middleware
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.ErrorUnaryInterceptor(),
			middleware.AuthUnaryInterceptor(os.Getenv("JWT_SECRET")),
			middleware.RateLimitUnaryInterceptor(
				parseEnvAsRate("RATE_LIMIT_USER_PER_SECOND", 1.0),
				parseEnvAsInt("RATE_LIMIT_USER_BURST", 5),
			),
		),
	)

	// Register services
	s.registerServices(appStore)

	// Enable gRPC reflection
	reflection.Register(s.grpcServer)

	// Create unified handler for both HTTP and gRPC
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle health check endpoint
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "ok")
			return
		}

		// Handle gRPC requests
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, r)
			return
		}

		// Handle other HTTP/2 requests
		if r.ProtoMajor == 2 {
			io.WriteString(w, "Hello HTTP/2")
			return
		}

		// Handle HTTP/1.1 requests
		io.WriteString(w, "Hello HTTP/1.1")
	}

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: h2c.NewHandler(http.HandlerFunc(handler), &http2.Server{}),
	}

	log.Printf("Server listening on %s (HTTP/1.1, HTTP/2, and gRPC)", address)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("failed to serve: %v", err)
	}
}

// registerServices registers all gRPC services
func (s *Server) registerServices(appStore store.Store) {
	// Register Health Service
	healthpb.RegisterHealthServiceServer(s.grpcServer, services.NewHealthService())

	// Register Search Service
	searchpb.RegisterSearchServiceServer(s.grpcServer, services.NewSearchService(appStore))

	// Register User Service
	userpb.RegisterUserServiceServer(s.grpcServer, services.NewUserService(appStore))

	// Register Tag Service
	tagpb.RegisterTagServiceServer(s.grpcServer, services.NewTagService(appStore))

	// Register Question Service
	questionpb.RegisterQuestionServiceServer(s.grpcServer, services.NewQuestionService(appStore))

	// Register Chat Service
	router := router.NewRouter(appStore)
	ai := ai.NewClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"))
	chatpb.RegisterChatServiceServer(s.grpcServer, services.NewChatService(router, ai))
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// Shutdown HTTP server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
			return err
		}
	}

	// Gracefully stop gRPC server
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	return nil
}

// ForceStop forcefully stops the server
func (s *Server) ForceStop() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
} 