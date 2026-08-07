package main_test

import (
	"Schwarz--Internship--2026/services/user-base/main"
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateUser(t *testing.T) {
	customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	longPassword := strings.Repeat("abcd", 20) // 80 bytes (> 72-byte limit)
	expectedID := int64(10)

	// Base user generator
	baseUser := func() *proto.User {
		return &proto.User{
			FirstName: "John",
			LastName:  "Doe",
			UserName:  "johndoe",
			Email:     "john@example.com",
			Password:  "supersecret123",
		}
	}

	tests := []struct {
		name         string
		getUser      func() *proto.User
		setupMock    func(mock sqlmock.Sqlmock, u *proto.User)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.CreateUserResponse, reqUser *proto.User)
	}{
		{name: "Success",
			getUser: baseUser,
			setupMock: func(mock sqlmock.Sqlmock, u *proto.User) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(u.FirstName, u.LastName, u.UserName, u.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.CreateUserResponse, reqUser *proto.User) {

				if res.User.FirstName != reqUser.FirstName {
					t.Errorf("expected FirstName %s, got %s", reqUser.FirstName, res.User.FirstName)
				}
				if res.User.LastName != reqUser.LastName {
					t.Errorf("expected LastName %s, got %s", reqUser.LastName, res.User.LastName)
				}
				if res.User.UserName != reqUser.UserName {
					t.Errorf("expected UserName %s, got %s", reqUser.UserName, res.User.UserName)
				}
				if res.User.Email != reqUser.Email {
					t.Errorf("expected Email %s, got %s", reqUser.Email, res.User.Email)
				}
				if res.User.Id != expectedID {
					t.Errorf("expected ID %d, got %d", expectedID, res.User.Id)
				}
				if err := bcrypt.CompareHashAndPassword([]byte(res.User.Password), []byte(reqUser.Password)); err != nil {
					t.Errorf("returned password is not a valid hash: %v", err)
				}
				if res.User.CreatedAt == nil {
					t.Error("expected CreatedAt to be populated, got nil")
				}
			},
		},
		{name: "WithExistingCreatedAt",
			getUser: func() *proto.User {
				u := baseUser()
				u.CreatedAt = timestamppb.New(customTime)
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock, u *proto.User) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(u.FirstName, u.LastName, u.UserName, u.Email, sqlmock.AnyArg(), customTime).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.CreateUserResponse, reqUser *proto.User) {
				if !res.User.CreatedAt.AsTime().Equal(customTime) {
					t.Errorf("expected CreatedAt %v, got %v", customTime, res.User.CreatedAt.AsTime())
				}
			},
		},
		{name: "DBError",
			getUser: baseUser,
			setupMock: func(mock sqlmock.Sqlmock, u *proto.User) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(u.FirstName, u.LastName, u.UserName, u.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.CreateUserResponse, reqUser *proto.User) {
				if res != nil {
					t.Errorf("expected nil user response on error, got %v", res)
				}
			},
		},
		{name: "NilUser",
			getUser:      func() *proto.User { return nil },
			setupMock:    func(mock sqlmock.Sqlmock, u *proto.User) {},
			expectedCode: codes.Internal,
		},
		{name: "LongPasswordTruncation",
			getUser: func() *proto.User {
				u := baseUser()
				u.Password = longPassword
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock, u *proto.User) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(u.FirstName, u.LastName, u.UserName, u.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.CreateUserResponse, reqUser *proto.User) {
				err := bcrypt.CompareHashAndPassword([]byte(res.User.Password), []byte(longPassword[:72]))
				if err != nil {
					t.Errorf("password truncation test failed: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to initialize sqlmock: %v", err)
			}
			defer db.Close()

			reqUser := tt.getUser()
			if tt.setupMock != nil {
				tt.setupMock(mock, reqUser)
			}

			svc := main.UserServiceImpl{DB: db}
			res, err := svc.CreateUser(context.Background(), &proto.CreateUserRequest{User: reqUser})

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
				tt.validate(t, res, reqUser)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}
