package main_test

import (
	"Schwarz--Internship--2026/services/user-base/main"
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetUser(t *testing.T) {
	customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	targetEmail := "john@example.com"

	baseRequest := &proto.GetUserRequest{
		Email: targetEmail,
	}

	expectedSQL := `SELECT id, first_name, last_name, user_name, email, hashed_password, created_at FROM users WHERE email = \$1`

	tests := []struct {
		name         string
		request      *proto.GetUserRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.GetUserRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.GetUserResponse, req *proto.GetUserRequest)
	}{
		{
			name:    "Success",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetUserRequest) {
				rows := sqlmock.NewRows([]string{
					"id", "first_name", "last_name", "user_name", "email", "hashed_password", "created_at",
				}).AddRow(10, "John", "Doe", "johndoe", "john@example.com", "hashed_pass_123", customTime)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetEmail()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.GetUserResponse, req *proto.GetUserRequest) {
				if res == nil || res.User == nil {
					t.Fatalf("expected non-nil user response, got nil")
				}
				if res.User.Id != 10 {
					t.Errorf("expected ID %d, got %d", 10, res.User.Id)
				}
				if res.User.FirstName != "John" {
					t.Errorf("expected FirstName John, got %s", res.User.FirstName)
				}
				if res.User.LastName != "Doe" {
					t.Errorf("expected LastName Doe, got %s", res.User.LastName)
				}
				if res.User.UserName != "johndoe" {
					t.Errorf("expected UserName johndoe, got %s", res.User.UserName)
				}
				if res.User.Email != "john@example.com" {
					t.Errorf("expected Email john@example.com, got %s", res.User.Email)
				}
				if res.User.Password != "hashed_pass_123" {
					t.Errorf("expected Password hashed_pass_123, got %s", res.User.Password)
				}
				if res.User.CreatedAt == nil {
					t.Fatal("expected CreatedAt to be non-nil")
				}
				if !res.User.CreatedAt.AsTime().Equal(customTime) {
					t.Errorf("expected CreatedAt %v, got %v", customTime, res.User.CreatedAt.AsTime())
				}
			},
		},
		{
			name:    "UserNotFound",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetUserRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetEmail()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: codes.NotFound,
			validate: func(t *testing.T, res *proto.GetUserResponse, req *proto.GetUserRequest) {
				if res != nil {
					t.Errorf("expected nil user response on error, got %v", res)
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetUserRequest) {},
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to initialize sqlmock: %v", err)
			}
			defer db.Close()

			if tt.setupMock != nil {
				tt.setupMock(mock, tt.request)
			}

			svc := main.UserServiceImpl{DB: db}
			res, err := svc.GetUser(context.Background(), tt.request)

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
				tt.validate(t, res, tt.request)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}
