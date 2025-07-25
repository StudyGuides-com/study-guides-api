package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	adminpb "github.com/studyguides-com/study-guides-api/api/v1/admin"
	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	devopspb "github.com/studyguides-com/study-guides-api/api/v1/devops"
	healthpb "github.com/studyguides-com/study-guides-api/api/v1/health"
	questionpb "github.com/studyguides-com/study-guides-api/api/v1/question"
	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	tagpb "github.com/studyguides-com/study-guides-api/api/v1/tag"
	userpb "github.com/studyguides-com/study-guides-api/api/v1/user"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/lib/webrouter"
	"github.com/studyguides-com/study-guides-api/internal/store"

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

	// Create web router for HTTP requests
	webRouter := webrouter.NewWebRouter(appStore)

	// Create unified handler for both HTTP and gRPC
	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("=== REQUEST START ===")
		log.Printf("Method: %s", r.Method)
		log.Printf("Path: %s", r.URL.Path)
		log.Printf("Protocol: %s", r.Proto)
		log.Printf("ProtoMajor: %d", r.ProtoMajor)
		log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
		log.Printf("User-Agent: %s", r.Header.Get("User-Agent"))
		log.Printf("All Headers: %v", r.Header)

		// Handle gRPC requests (both HTTP/1.1 and HTTP/2)
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			log.Printf("*** ROUTING TO GRPC: %s (Protocol: %s) ***", r.URL.Path, r.Proto)

			s.grpcServer.ServeHTTP(w, r)
			log.Printf("*** GRPC HANDLER COMPLETED ***")
			return
		}

		// Handle HTTP requests with web router
		log.Printf("*** ROUTING TO WEB ROUTER: %s ***", r.URL.Path)
		webRouter.ServeHTTP(w, r)
		log.Printf("*** WEB ROUTER COMPLETED ***")
		log.Printf("=== REQUEST END ===")
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

	// Register Chat Service with MCP system
	ai := ai.NewClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"))
	chatpb.RegisterChatServiceServer(s.grpcServer, services.NewChatService(appStore, ai))

	// Register Admin Service
	adminpb.RegisterAdminServiceServer(s.grpcServer, services.NewAdminService(appStore))

	// Register Devops Service
	devopspb.RegisterDevopsServiceServer(s.grpcServer, services.NewDevopsService(appStore))
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
