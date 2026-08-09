package main

import (
	frb "Schwarz--Internship--2026/services/api-rest-gateway/friend-request-base/proto"
	ub "Schwarz--Internship--2026/services/api-rest-gateway/user-base/proto"

	"Schwarz--Internship--2026/services/common"
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

	user_base_endpoint, err := common.GetRequiredEnv("USER-BASE_ENDPOINT")
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)

	}

	err = ub.RegisterUserServiceHandlerFromEndpoint(ctx, mux, user_base_endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	friend_request_base_endpoint, err := common.GetRequiredEnv("FRIEND-REQUEST-BASE_ENDPOINT")
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)

	}

	err = frb.RegisterFriendRequestServiceHandlerFromEndpoint(ctx, mux, friend_request_base_endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	log.Println("HTTP REST Gateway listening on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
