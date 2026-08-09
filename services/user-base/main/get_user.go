package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u UserServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {

	id := req.GetId()
	if id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing id")
	}
	var resUser proto.User
	resUser.Id = id

	query := `
		SELECT first_name, last_name, user_name, email, hashed_password, created_at
		FROM users
		WHERE id = $1
	`
	var created_at time.Time

	err := u.DB.QueryRow(query, id).Scan(
		&resUser.FirstName,
		&resUser.LastName,
		&resUser.UserName,
		&resUser.Email,
		&resUser.Password,
		&created_at)

	if err == sql.ErrNoRows {
		log.Printf("Failed to get user %d: %v", id, err)
		return nil, status.Errorf(codes.NotFound, "failed to get user %d from database: %v", id, err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user %d from database: %v", id, err)
	}

	resUser.CreatedAt = timestamppb.New(created_at)

	return &proto.GetUserResponse{User: &resUser}, nil
}
