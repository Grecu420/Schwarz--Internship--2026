package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) CreateProperty(ctx context.Context, req *proto.CreatePropertyRequest) (*proto.CreatePropertyResponse, error) {
	return service.propService.CreateProperty(ctx, req)
}

func (service GatewayServiceImpl) UpdateProperty(ctx context.Context, req *proto.UpdatePropertyRequest) (*proto.UpdatePropertyResponse, error) {
	return service.propService.UpdateProperty(ctx, req)
}

func (service GatewayServiceImpl) DeleteProperty(ctx context.Context, req *proto.DeletePropertyRequest) (*proto.DeletePropertyResponse, error) {
	return service.propService.DeleteProperty(ctx, req)
}

func (service GatewayServiceImpl) GetProperty(ctx context.Context, req *proto.GetPropertyRequest) (*proto.GetPropertyResponse, error) {
	return service.propService.GetProperty(ctx, req)
}
