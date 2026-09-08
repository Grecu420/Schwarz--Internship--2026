package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return service.userService.CreateUser(ctx, req)
}

func (service GatewayServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	return service.userService.GetUser(ctx, req)
}

func (service GatewayServiceImpl) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	return service.userService.ListUsers(ctx, req)
}

func (service GatewayServiceImpl) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	return service.userService.UpdateUser(ctx, req)
}

func (service GatewayServiceImpl) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return service.userService.DeleteUser(ctx, req)
}
