package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"net/http"
	"os"
	"user_service/config"
	"user_service/internal/db"
	"user_service/internal/service"
	"user_service/protogen/user"
)

func main() {

	err := godotenv.Load("../conf/config.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	filePath := os.Getenv("CONFIG_FILE")

	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading YAML file:", err)
		return
	}

	var cfg config.Config
	if err := yaml.Unmarshal(yamlData, &cfg); err != nil {
		log.Fatalf("Error unmarshaling YAML:", err)
		return
	}

	cassandraSvc := *db.NewCassandraDetailsSvc(&cfg)

	session := cassandraSvc.ConnectCassandra()
	defer session.Close()

	// Initialize the gRPC service
	userService := service.NewUserServiceServer(session)

	// Start the gRPC server
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userService)

	lis, err := net.Listen(cfg.GrpcDetails.Network, cfg.GrpcDetails.Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.GrpcDetails.Endpoint, opts); err != nil {
		log.Fatalf("Failed to start REST gateway: %v", err)
	}

	log.Println("Starting REST server on :8080")
	if err := http.ListenAndServe(cfg.HttpDetails.Port, mux); err != nil {
		log.Fatalf("Failed to serve REST: %v", err)
	}
}
