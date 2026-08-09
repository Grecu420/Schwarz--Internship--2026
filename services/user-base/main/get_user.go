package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u UserServiceImpl) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	id := req.GetId()
	var resUser proto.User
	resUser.Id = id

	query := `
		SELECT first_name, last_name, user_name, email, hashed_password, created_at
		FROM users
		WHERE id == $1
	`
	var created_at time.Time

	err := u.DB.QueryRow(query, id).Scan(
		&resUser.FirstName,
		&resUser.LastName,
		&resUser.UserName,
		&resUser.Email,
		&resUser.Password,
		&created_at)

	if err != nil {
		log.Printf("Failed to get user %d: %v", id, err)
		return nil, status.Errorf(codes.Internal, "failed to get user %d from database: %v", id, err)
	}

	return &proto.GetUserResponse{User: &resUser}, nil
}
