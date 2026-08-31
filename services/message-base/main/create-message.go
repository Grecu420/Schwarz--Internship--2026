package main

import (
	"context"
	"Schwarz--Internship--2026/services/message-base/main/proto" 

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service MessageServiceImpl) CreateMessage(ctx context.Context, req *proto.CreateMessageRequest) (*proto.CreateMessageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.GetConversationId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "conversation_id missing")
	}
	if req.GetSenderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "sender_id missing")
	}
	if req.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "content missing")
	}

	id, createdAt, err := InsertMessage(ctx, service.DB, req.GetConversationId(), req.GetSenderId(), req.GetContent())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert message to database: %v", err)
	}

	msg := &proto.Message{
		Id:             id,
		ConversationId: req.GetConversationId(),
		SenderId:       req.GetSenderId(),
		Content:        req.GetContent(),
		CreatedAt:      timestamppb.New(createdAt), 
	}

	return &proto.CreateMessageResponse{
		Message: msg,
	}, nil
}