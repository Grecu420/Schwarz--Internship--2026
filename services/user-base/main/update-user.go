package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service UserServiceImpl) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	if req == nil || req.User == nil {
		return nil, status.Error(codes.InvalidArgument, "request and user are mandatory")
	}

	if req.FieldMask == nil || len(req.FieldMask.GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "field_mask is mandatory for updates")
	}

	requestID := req.User.GetId()
	if requestID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is mandatory")
	}

	build := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("users").
		Where(sq.Eq{"id": requestID})

	hasValidFields := false
	for _, path := range req.FieldMask.GetPaths() {
		switch path {
		case "first_name":
			if req.User.GetFirstName() == "" {
				return nil, status.Error(codes.InvalidArgument, "first_name is mandatory")
			}
			build = build.Set("first_name", req.User.GetFirstName())
			hasValidFields = true
		case "last_name":
			if req.User.GetLastName() == "" {
				return nil, status.Error(codes.InvalidArgument, "last_name is mandatory")
			}
			build = build.Set("last_name", req.User.GetLastName())
			hasValidFields = true
		case "user_name":
			if req.User.GetUserName() == "" {
				return nil, status.Error(codes.InvalidArgument, "username is mandatory")
			}
			build = build.Set("user_name", req.User.GetUserName())
			hasValidFields = true
		default:
			return nil, status.Errorf(codes.InvalidArgument, "invalid field path in field_mask: %s", path)
		}
	}

	if !hasValidFields {
		return nil, status.Error(codes.InvalidArgument, "no valid fields to update")
	}

	build = build.Suffix("RETURNING id, first_name, last_name, user_name, email, hashed_password, created_at")

	var returnUser proto.User
	var createdAt time.Time

	err := build.RunWith(service.DB).QueryRowContext(ctx).Scan(
		&returnUser.Id,
		&returnUser.FirstName,
		&returnUser.LastName,
		&returnUser.UserName,
		&returnUser.Email,
		&returnUser.Password,
		&createdAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "user to update not found")
	}
	
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, status.Error(codes.AlreadyExists, "username is already taken")
		}
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}

	returnUser.CreatedAt = timestamppb.New(createdAt)

	return &proto.UpdateUserResponse{User: &returnUser}, nil
}