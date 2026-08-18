package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) getUserServiceClient() proto.UserServiceClient {
	return proto.NewUserServiceClient(service.userBaseConn)
}

func (service GatewayServiceImpl) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return service.getUserServiceClient().CreateUser(ctx, req)
}

func (service GatewayServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	return service.getUserServiceClient().GetUser(ctx, req)
}
