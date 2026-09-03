package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) CreateMessage(ctx context.Context, req *proto.CreateMessageRequest) (*proto.CreateMessageResponse, error) {
	return service.messageService.CreateMessage(ctx, req)
}

func (service GatewayServiceImpl) ListMessages(ctx context.Context, req *proto.ListMessagesRequest) (*proto.ListMessagesResponse, error) {
	return service.messageService.ListMessages(ctx, req)
}
