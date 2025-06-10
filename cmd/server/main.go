package main

import (
	"log"
	"net"
	"os"
	"strconv"

	healthpb "github.com/studyguides-com/study-guides-api/api/v1/health"
	questionpb "github.com/studyguides-com/study-guides-api/api/v1/question"
	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	tagpb "github.com/studyguides-com/study-guides-api/api/v1/tag"
	userpb "github.com/studyguides-com/study-guides-api/api/v1/user"
	"github.com/studyguides-com/study-guides-api/internal/store"

	"github.com/joho/godotenv"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/services"
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
		port = "50051"
	}

	address := ":" + port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
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

	healthpb.RegisterHealthServiceServer(grpcServer, services.NewHealthService())
	searchpb.RegisterSearchServiceServer(grpcServer, services.NewSearchService(appStore))
	userpb.RegisterUserServiceServer(grpcServer, services.NewUserService())
	tagpb.RegisterTagServiceServer(grpcServer, services.NewTagService(appStore))
	questionpb.RegisterQuestionServiceServer(grpcServer, services.NewQuestionService(appStore))

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
