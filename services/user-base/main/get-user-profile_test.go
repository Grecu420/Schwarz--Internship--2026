package main_test

import (
    "Schwarz--Internship--2026/services/user-base/main"
    "Schwarz--Internship--2026/services/user-base/main/proto"
    "context"
    "database/sql"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func TestGetUserProfile(t *testing.T) {
	targetID := int64(10)

	baseRequest := &proto.GetUserProfileRequest{
		Id: targetID,
	}

	expectedSQL := `SELECT id, first_name, last_name, user_name, profile_image_url FROM users WHERE id = \$1`

	tests := []struct {
		name         string
		request      *proto.GetUserProfileRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.GetUserProfileRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.GetUserProfileResponse, req *proto.GetUserProfileRequest)
	}{
		{
			name:    "Success",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetUserProfileRequest) {
				rows := sqlmock.NewRows([]string{
					"id", "first_name", "last_name", "user_name", "profile_image_url",
				}).AddRow(10, "John", "Doe", "johndoe", "https://example.com/avatar.jpg")

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.GetUserProfileResponse, req *proto.GetUserProfileRequest) {
				if res == nil || res.User == nil {
					t.Fatalf("expected non-nil user profile response, got nil")
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
				if res.User.ProfileImageUrl != "https://example.com/avatar.jpg" {
					t.Errorf("expected ProfileImageUrl https://example.com/avatar.jpg, got %s", res.User.ProfileImageUrl)
				}
			},
		},
		{
			name:    "UserNotFound",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetUserProfileRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: codes.NotFound,
			validate: func(t *testing.T, res *proto.GetUserProfileResponse, req *proto.GetUserProfileRequest) {
				if res != nil {
					t.Errorf("expected nil user profile response on error, got %v", res)
				}
			},
		},
		{
			name: "InvalidID",
			request: &proto.GetUserProfileRequest{
				Id: 0,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetUserProfileRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetUserProfileRequest) {},
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
			res, err := svc.GetUserProfile(context.Background(), tt.request)

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