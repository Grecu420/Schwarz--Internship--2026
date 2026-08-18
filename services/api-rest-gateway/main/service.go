package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"Schwarz--Internship--2026/services/common"

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
	userBaseConn          *grpc.ClientConn
	friendRequestBaseConn *grpc.ClientConn
	authBaseConn          *grpc.ClientConn
}

func (service GatewayServiceImpl) Close() {
	service.userBaseConn.Close()
	service.friendRequestBaseConn.Close()
	service.authBaseConn.Close()
}

func createGatewayServer() GatewayServiceImpl {

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// get endpoints
	user_base_endpoint, err := common.GetRequiredEnv("USER-BASE_ENDPOINT")
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)

	}
	friend_request_base_endpoint, err := common.GetRequiredEnv("FRIEND-REQUEST-BASE_ENDPOINT")
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)

	}
	auth_base_endpoint, err := common.GetRequiredEnv("AUTH-BASE_ENDPOINT")
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)

	}

	// create clients
	connUB, err := grpc.NewClient(user_base_endpoint, opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	connFRB, err := grpc.NewClient(friend_request_base_endpoint, opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	connAB, err := grpc.NewClient(auth_base_endpoint, opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	return GatewayServiceImpl{
		userBaseConn:          connUB,
		friendRequestBaseConn: connFRB,
		authBaseConn:          connAB,
	}
}

func main() {

	// 1. Start gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := createGatewayServer()
	defer server.Close()

	proto.RegisterGatewayServiceServer(grpcServer, server)
	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// 2. Start gRPC-Gateway HTTP Proxy
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = proto.RegisterGatewayServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	log.Println("HTTP Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
