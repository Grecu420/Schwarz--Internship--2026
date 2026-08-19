package main

import (
	"Schwarz--Internship--2026/services/auth-base/main/proto"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var secretKey = []byte("123")

type LoginClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT signed with a secret key that expires after duration.
func GenerateToken(userID int64, duration time.Duration, secretKey []byte) (string, error) {
	// Define claims
	claims := LoginClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-base-service",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

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
	token, err := GenerateToken(res.User.Id, time.Hour*24, secretKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "token failure: %v", err)
	}

	return &proto.LoginResponse{JWT: token}, nil
}
