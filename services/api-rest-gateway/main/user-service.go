package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"Schwarz--Internship--2026/services/common"
	"context"
	"log"

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

func (service GatewayServiceImpl) getUserServiceClient() proto.UserServiceClient {
	return proto.NewUserServiceClient(service.userBaseConn)
}

func (service GatewayServiceImpl) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return service.getUserServiceClient().CreateUser(ctx, req)
}

func (service GatewayServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	return service.getUserServiceClient().GetUser(ctx, req)
}
