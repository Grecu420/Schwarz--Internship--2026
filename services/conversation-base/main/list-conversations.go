package main

import (
	"Schwarz--Internship--2026/services/conversation-base/main/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service ConversationServiceImpl) ListConversations(ctx context.Context, req *proto.ListConversationsRequest) (*proto.ListConversationsResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	conversations, err := SelectConversationsListInDB(ctx, service.DB, req.GetUserId())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed db query: %v", err)
	}

	return &proto.ListConversationsResponse{Conversations: conversations}, nil
}