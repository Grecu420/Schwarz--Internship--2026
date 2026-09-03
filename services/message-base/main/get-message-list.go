package main

import (
	"Schwarz--Internship--2026/services/message-base/main/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service MessageServiceImpl) ListMessages(ctx context.Context, req *proto.ListMessagesRequest) (*proto.ListMessagesResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	if req.GetConversationId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid conversation id")
	}

	messages, err := SelectMessages(ctx, service.DB, req.GetConversationId())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed db query: %v", err)
	}

	return &proto.ListMessagesResponse{Messages: messages}, nil
}
