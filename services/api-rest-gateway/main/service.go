package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"Schwarz--Internship--2026/services/common"
	"fmt"
	"os"

	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayServiceImpl struct {
	proto.UnimplementedGatewayServiceServer
	userService          proto.UserServiceClient
	friendRequestService proto.FriendRequestServiceClient
	authService          proto.AuthServiceClient
	messageService       proto.MessageServiceClient
	convService          proto.ConversationServiceClient
	propService          proto.PropertyServiceClient
}

func createConnection(envVar string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	endpoint, err := common.GetRequiredEnv(envVar)
	if err != nil {
		return nil, fmt.Errorf("Failed to register endpoint: %w", err)
	}
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %w", err)
	}
	return conn, nil

}

func main() {

	// Read auth secret
	secret, err := os.ReadFile("/run/secrets/auth-key")
	if err != nil {
		log.Fatal("missing secret")
	}

	// Obtain connections
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Map backend services to environment variable keys
	envMap := map[string]string{
		"user":         "USER-BASE_ENDPOINT",
		"conversation": "CONVERSATION-BASE_ENDPOINT",
		"message":      "MESSAGE-BASE_ENDPOINT",
		"friend":       "FRIEND-REQUEST-BASE_ENDPOINT",
		"auth":         "AUTH-BASE_ENDPOINT",
		"property":     "PROPERTY-BASE_ENDPOINT",
	}

	conns := make(map[string]*grpc.ClientConn, len(envMap))
	for key, env := range envMap {
		conn, err := createConnection(env, opts)
		if err != nil {
			log.Fatalf("failed to dial service (%s): %v", env, err)
		}
		defer conn.Close()
		conns[key] = conn
	}

	server := GatewayServiceImpl{
		authService:          proto.NewAuthServiceClient(conns["auth"]),
		userService:          proto.NewUserServiceClient(conns["user"]),
		friendRequestService: proto.NewFriendRequestServiceClient(conns["friend"]),
		convService:          proto.NewConversationServiceClient(conns["conversation"]),
		messageService:       proto.NewMessageServiceClient(conns["message"]),
		propService:          proto.NewPropertyServiceClient(conns["property"]),
	}

	// Start gRPC Server
	authInterceptor := authInterceptor{secret: secret}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.AuthInterceptor),
	)
	proto.RegisterGatewayServiceServer(grpcServer, server)
	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start gRPC-Gateway HTTP Proxy
	ctx := context.Background()
	mux := runtime.NewServeMux()
	err = proto.RegisterGatewayServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	log.Println("HTTP Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
