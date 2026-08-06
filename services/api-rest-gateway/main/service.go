package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/user-base/proto"
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Connect the HTTP gateway to the backend gRPC server running on localhost:50051
	err := proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "user-base-service:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	log.Println("HTTP REST Gateway listening on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
