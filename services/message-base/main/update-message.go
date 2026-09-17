package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Schwarz--Internship--2026/services/message-base/main/proto"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service MessageServiceImpl) UpdateMessage(ctx context.Context, req *proto.UpdateMessageRequest) (*proto.UpdateMessageResponse, error) {
	if req == nil || req.Message == nil {
		return nil, status.Error(codes.InvalidArgument, "request and message are mandatory")
	}

	if req.FieldMask == nil || len(req.FieldMask.GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "field_mask is mandatory for updates")
	}

	requestID := req.Message.GetId()
	if requestID == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is mandatory")
	}

	build := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("messages").
		Where(sq.Eq{"id": requestID})

	if req.Message.GetConversationId() > 0 {
		build = build.Where(sq.Eq{"conversation_id": req.Message.GetConversationId()})
	}

	if req.Message.GetSenderId() > 0 {
		build = build.Where(sq.Eq{"sender_id": req.Message.GetSenderId()})
	}

	hasValidFields := false
	for _, path := range req.FieldMask.GetPaths() {
		switch path {
		case "is_read":
			build = build.Set("is_read", req.Message.GetIsRead())
			hasValidFields = true

		default:
			return nil, status.Errorf(codes.InvalidArgument, "invalid field path in field_mask: %s", path)
		}
	}

	if !hasValidFields {
		return nil, status.Error(codes.InvalidArgument, "no valid fields to update")
	}

	build = build.Suffix("RETURNING id, conversation_id, sender_id, content, created_at, is_read")

	var returnMessage proto.Message
	var createdAt time.Time

	err := build.RunWith(service.DB).QueryRowContext(ctx).Scan(
		&returnMessage.Id,
		&returnMessage.ConversationId,
		&returnMessage.SenderId,
		&returnMessage.Content,
		&createdAt,
		&returnMessage.IsRead,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "message to update not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}

	returnMessage.CreatedAt = timestamppb.New(createdAt)

	return &proto.UpdateMessageResponse{Message: &returnMessage}, nil
}