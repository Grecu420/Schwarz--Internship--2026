package main_test

import (
	"Schwarz--Internship--2026/services/auth-base/main"
	"Schwarz--Internship--2026/services/auth-base/main/proto"
	"Schwarz--Internship--2026/services/common/auth"
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockUserService struct {
	getUserFunc func(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error)
}

func (m *mockUserService) GetUser(ctx context.Context, req *proto.GetUserRequest, opts ...grpc.CallOption) (*proto.GetUserResponse, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, req)
	}
	return nil, errors.New("unexpected GetUser call")
}

func (m *mockUserService) CreateUser(ctx context.Context, in *proto.CreateUserRequest, opts ...grpc.CallOption) (*proto.CreateUserResponse, error) {
	return nil, errors.New("unexpected CreateUser call")
}

func (m *mockUserService) ListUsers(ctx context.Context, in *proto.ListUsersRequest, opts ...grpc.CallOption) (*proto.ListUsersResponse, error) {
	return nil, errors.New("unexpected ListUsers call")
}

func (m *mockUserService) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest, opts ...grpc.CallOption) (*proto.DeleteUserResponse, error) {
	return nil, errors.New("unexpected DeleteUsers call")
}

func (m *mockUserService) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest, opts ...grpc.CallOption) (*proto.UpdateUserResponse, error) {
	return nil, errors.New("unexpected UpdateUsers call")
}

var secret = []byte("123")

func TestGenerateToken(t *testing.T) {
	userID := int64(100)
	duration := 24 * time.Hour

	tokenStr, err := auth.GenerateToken(userID, duration, "auth-base-service", secret)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	claims, err := auth.ParseToken(tokenStr, secret)

	if err != nil {
		t.Fatalf("failed to parse generated token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Issuer != "auth-base-service" {
		t.Errorf("expected Issuer 'auth-base-service', got %s", claims.Issuer)
	}
}

func TestLogin(t *testing.T) {
	validEmail := "john@example.com"
	validPassword := "password"
	userID := int64(100)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to setup hashed password: %v", err)
	}

	baseReq := &proto.LoginRequest{
		Email:    validEmail,
		Password: validPassword,
	}

	tests := []struct {
		name         string
		req          *proto.LoginRequest
		getUserFunc  func(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.LoginResponse)
	}{
		{
			name: "Success",
			req:  baseReq,
			getUserFunc: func(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
				return &proto.GetUserResponse{
					User: &proto.User{
						Id:       userID,
						Email:    validEmail,
						Password: string(hashedPassword),
					},
				}, nil

			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.LoginResponse) {
				if res == nil || res.JWT == "" {
					t.Fatal("expected non-empty JWT token response")
				}

				claims, err := auth.ParseToken(res.JWT, secret)

				if err != nil {
					t.Fatalf("returned JWT is invalid: %v", err)
				}
				if claims.UserID != userID {
					t.Errorf("expected token UserID %d, got %d", userID, claims.UserID)
				}
			},
		},
		{
			name:         "NilRequest",
			req:          nil,
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingEmail",
			req: &proto.LoginRequest{
				Email:    "",
				Password: validPassword,
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingPassword",
			req: &proto.LoginRequest{
				Email:    validEmail,
				Password: "",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "UserNotFound",
			req:  baseReq,
			getUserFunc: func(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
				return nil, status.Error(codes.NotFound, "user not found")
			},
			expectedCode: codes.NotFound,
		},
		{
			name: "WrongPassword",
			req: &proto.LoginRequest{
				Email:    validEmail,
				Password: "wrong",
			},
			getUserFunc: func(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
				return &proto.GetUserResponse{
					User: &proto.User{
						Id:       userID,
						Email:    validEmail,
						Password: string(hashedPassword),
					},
				}, nil
			},
			expectedCode: codes.PermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{tt.getUserFunc}

			svc := main.AuthServiceImpl{
				UserService: mockSvc,
				Secret:      secret,
			}

			res, err := svc.Login(context.Background(), tt.req)

			if tt.expectedCode != codes.OK {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got non-status error: %v", err)
				}
				if st.Code() != tt.expectedCode {
					t.Errorf("expected gRPC code %v, got %v", tt.expectedCode, st.Code())
				}
			} else if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, res)
			}
		})
	}
}
