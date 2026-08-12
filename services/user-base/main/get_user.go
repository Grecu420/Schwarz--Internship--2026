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

func (u UserServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {

	id := req.GetId()
	if id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing id")
	}

	resUser, err := SelectUser(u.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Missing user %d: %v", id, err)
		return nil, status.Errorf(codes.NotFound,
			"failed to get user %d from database: %v", id, err)
	} else if err != nil {
		log.Printf("Failed to get user %d: %v", id, err)
		return nil, status.Errorf(codes.Internal,
			"failed to get user %d from database: %v", id, err)
	}

	return &proto.GetUserResponse{User: resUser}, nil
}
