package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"user_service/internal/db"
	"user_service/internal/service"
	"user_service/protogen/user"
)

func main() {
	// Connect to Cassandra
	session := db.ConnectCassandra()
	defer session.Close()

	// Initialize the gRPC service
	userService := service.NewUserServiceServer(session)

	// Start the gRPC server
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start the REST gateway
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("Failed to start REST gateway: %v", err)
	}

	log.Println("Starting REST server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve REST: %v", err)
	}
}
