package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) getAuthServiceClient() proto.AuthServiceClient {
	return proto.NewAuthServiceClient(service.authBaseConn)
}

func (service GatewayServiceImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return service.getAuthServiceClient().Login(ctx, req)
}
