package main_test

import (
    "Schwarz--Internship--2026/services/user-base/main"
    "Schwarz--Internship--2026/services/user-base/main/proto"
    "context"
    "database/sql"
    "errors"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestUpdateUser(t *testing.T) {
    customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)

    baseReq := &proto.UpdateUserRequest{
        User: &proto.User{
            Id:        10,
            FirstName: "Jane",
        },
        FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"first_name"}},
    }

    expectedSQL := `UPDATE users SET first_name = \$1 WHERE id = \$2 RETURNING id, first_name, last_name, user_name, email, hashed_password, profile_image_url, created_at`

    tests := []struct {
        name         string
        request      *proto.UpdateUserRequest
        setupMock    func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest)
        expectedCode codes.Code
        validate     func(t *testing.T, res *proto.UpdateUserResponse, req *proto.UpdateUserRequest)
    }{
        {
            name:    "Success",
            request: baseReq,
            setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {
                rows := sqlmock.NewRows([]string{
                    "id", "first_name", "last_name", "user_name", "email", "hashed_password", "profile_image_url", "created_at",
                }).AddRow(10, "Jane", "Doe", "janedoe", "jane@example.com", "hashedpass123", "", customTime)

                mock.ExpectQuery(expectedSQL).
                    WithArgs(req.User.GetFirstName(), req.User.GetId()).
                    WillReturnRows(rows)
            },
            expectedCode: codes.OK,
            validate: func(t *testing.T, res *proto.UpdateUserResponse, req *proto.UpdateUserRequest) {
                if res == nil || res.User == nil {
                    t.Fatalf("expected non-nil response, got nil")
                }
                if res.User.FirstName != "Jane" {
                    t.Errorf("expected FirstName Jane, got %s", res.User.FirstName)
                }
                if !res.User.CreatedAt.AsTime().Equal(customTime) {
                    t.Errorf("expected CreatedAt %v, got %v", customTime, res.User.CreatedAt.AsTime())
                }
            },
        },
        {
            name: "SuccessUpdateProfileImage", 
            request: &proto.UpdateUserRequest{
                User: &proto.User{
                    Id:              10,
                    ProfileImageUrl: "https://res.cloudinary.com/demo/image/upload/noua_poza.jpg",
                },
                FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"profile_image_url"}},
            },
            setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {
                expectedImageSQL := `UPDATE users SET profile_image_url = \$1 WHERE id = \$2 RETURNING id, first_name, last_name, user_name, email, hashed_password, profile_image_url, created_at`
                
                rows := sqlmock.NewRows([]string{
                    "id", "first_name", "last_name", "user_name", "email", "hashed_password", "profile_image_url", "created_at",
                }).AddRow(10, "Jane", "Doe", "janedoe", "jane@example.com", "hashedpass123", req.User.GetProfileImageUrl(), customTime)

                mock.ExpectQuery(expectedImageSQL).
                    WithArgs(req.User.GetProfileImageUrl(), req.User.GetId()).
                    WillReturnRows(rows)
            },
            expectedCode: codes.OK,
            validate: func(t *testing.T, res *proto.UpdateUserResponse, req *proto.UpdateUserRequest) {
                if res == nil || res.User == nil {
                    t.Fatalf("expected non-nil response, got nil")
                }
                if res.User.ProfileImageUrl != "https://res.cloudinary.com/demo/image/upload/noua_poza.jpg" {
                    t.Errorf("expected ProfileImageUrl to be updated, got %s", res.User.ProfileImageUrl)
                }
            },
        },
        {
            name:         "NilRequest",
            request:      nil,
            setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {},
            expectedCode: codes.InvalidArgument,
        },
        {
            name: "MissingFieldMask",
            request: &proto.UpdateUserRequest{
                User: &proto.User{Id: 10, FirstName: "Jane"},
            },
            setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {},
            expectedCode: codes.InvalidArgument,
        },
        {
            name: "MissingUserId",
            request: &proto.UpdateUserRequest{
                User:      &proto.User{Id: 0, FirstName: "Jane"},
                FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"first_name"}},
            },
            setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {},
            expectedCode: codes.InvalidArgument,
        },
        {
            name: "InvalidFieldMaskPath",
            request: &proto.UpdateUserRequest{
                User:      &proto.User{Id: 10, FirstName: "Jane"},
                FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"unknown_field"}},
            },
            setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {},
            expectedCode: codes.InvalidArgument,
        },
        {
            name: "ProhibitedFieldMaskPathEmail",
            request: &proto.UpdateUserRequest{
                User:      &proto.User{Id: 10, Email: "new@example.com"},
                FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"email"}},
            },
            setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {},
            expectedCode: codes.InvalidArgument,
        },
        {
            name: "ProhibitedFieldMaskPathPassword",
            request: &proto.UpdateUserRequest{
                User:      &proto.User{Id: 10, Password: "new_password"},
                FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"password"}},
            },
            setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {},
            expectedCode: codes.InvalidArgument,
        },
        {
            name: "UserNotFound",
            request: baseReq,
            setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {
                mock.ExpectQuery(expectedSQL).
                    WithArgs(req.User.GetFirstName(), req.User.GetId()).
                    WillReturnError(sql.ErrNoRows)
            },
            expectedCode: codes.NotFound,
        },
        {
            name: "DuplicateUsernameError",
            request: &proto.UpdateUserRequest{
                User:      &proto.User{Id: 10, UserName: "taken_username"},
                FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"user_name"}},
            },
            setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateUserRequest) {
                expectedDuplicateSQL := `UPDATE users SET user_name = \$1 WHERE id = \$2 RETURNING id, first_name, last_name, user_name, email, hashed_password, profile_image_url, created_at`
                mock.ExpectQuery(expectedDuplicateSQL).
                    WithArgs(req.User.GetUserName(), req.User.GetId()).
                    WillReturnError(errors.New("duplicate key value violates unique constraint \"users_username_key\""))
            },
            expectedCode: codes.AlreadyExists,
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
            res, err := svc.UpdateUser(context.Background(), tt.request)

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