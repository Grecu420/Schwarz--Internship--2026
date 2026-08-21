package main

import (
	"Schwarz--Internship--2026/services/auth-base/main/proto"
	"Schwarz--Internship--2026/services/common/auth"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service AuthServiceImpl) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	// Read email and password
	email := req.GetEmail()
	password := req.GetPassword()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "missing password")
	}

	// Get user with email
	res, err := service.UserService.GetUser(ctx, &proto.GetUserRequest{Email: email})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// Compare password with hashed password
	hashed_password := res.GetUser().GetPassword()
	err = bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, status.Error(codes.PermissionDenied, "wrong password")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "login failure: %v", err)
	}

	// Generate token
	token, err := auth.GenerateToken(res.User.Id, time.Hour*24, "auth-base-service", service.Secret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "token failure: %v", err)
	}

	return &proto.LoginResponse{JWT: token}, nil
}
