package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) ListProperties(ctx context.Context, req *proto.ListPropertiesRequest) (*proto.ListPropertiesResponse, error) {
	return service.propListService.ListProperties(ctx, req)
}

func (service GatewayServiceImpl) CountProperties(ctx context.Context, req *proto.CountPropertiesRequest) (*proto.CountPropertiesResponse, error) {
	return service.propListService.CountProperties(ctx, req)
}
