package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) CreateMessage(ctx context.Context, req *proto.CreateMessageRequest) (*proto.CreateMessageResponse, error) {
	return service.messageService.CreateMessage(ctx, req)
}

func (service GatewayServiceImpl) GetMessageList(ctx context.Context, req *proto.GetMessageListRequest) (*proto.GetMessageListResponse, error) {
	return service.messageService.GetMessageList(ctx, req)
}
