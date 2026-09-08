package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service UserServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	// Get email
	email := req.GetEmail()
	if email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing email")
	}

	// Get user from database
	resUser, err := SelectUser(ctx, service.DB, email)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Missing user (email: %s): %v", email, err)
		return nil, status.Errorf(codes.NotFound,
			"failed to get user (email: %s) from database: %v", email, err)
	} else if err != nil {
		log.Printf("Failed to get user (email: %s): %v", email, err)
		return nil, status.Errorf(codes.Internal,
			"failed to get user (email: %s) from database: %v", email, err)
	}

	return &proto.GetUserResponse{User: resUser}, nil
}
