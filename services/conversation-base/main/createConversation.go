package main

import (
	"Schwarz--Internship--2026/services/conversation-base/main/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service ConversationServiceImpl) CreateConversation(ctx context.Context, req *proto.CreateConversationRequest) (*proto.CreateConversationResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetUser1Id() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user1 missing")
	}
	if req.GetUser2Id() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user2 missing")
	}

	id, err := InsertConversation(ctx, service.DB, req.GetUser1Id(), req.GetUser2Id())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed insert conversation to database: %v", err)
	}
	return &proto.CreateConversationResponse{Id: id}, nil
}
