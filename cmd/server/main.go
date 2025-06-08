package main

import (
	"log"
	"net"
	"os"

	studypb "github.com/studyguides-com/study-guides-api/api/study"
	"github.com/studyguides-com/study-guides-api/internal/service"

	"github.com/joho/godotenv"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
		grpc.UnaryInterceptor(middleware.AuthUnaryInterceptor(os.Getenv("JWT_SECRET"))),
	)
	
	studypb.RegisterStudyServiceServer(grpcServer, service.NewStudyService())

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
