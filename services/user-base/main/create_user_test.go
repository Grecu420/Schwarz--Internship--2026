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

// TestCreateUser_Success verifies user creation, password hashing, and ID assignment
func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize sqlmock: %v", err)
	}
	defer db.Close()

	svc := main.UserServiceImpl{DB: db}

	rawPassword := "supersecret123"

	reqUser := &proto.User{
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe",
		Email:     "john@example.com",
		Password:  rawPassword,
	}

	req := &proto.CreateUserRequest{User: reqUser}

	expectedID := int64(10)

	// Expect the INSERT query.
	// We use sqlmock.AnyArg() for password (bcrypt salt varies) and timestamp.
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			reqUser.FirstName,
			reqUser.LastName,
			reqUser.UserName,
			reqUser.Email,
			sqlmock.AnyArg(), // hashed password
			sqlmock.AnyArg(), // created_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	res, err := svc.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	resUser := res.User

	// 1. Verify returned ID
	if resUser.Id != expectedID {
		t.Errorf("expected ID %d, got %d", expectedID, resUser.Id)
	}

	// 2. Verify password was hashed correctly
	err = bcrypt.CompareHashAndPassword([]byte(resUser.Password), []byte(rawPassword))
	if err != nil {
		t.Errorf("returned password is not a valid hash of original password: %v", err)
	}

	// 3. Verify CreatedAt timestamp was automatically generated
	if resUser.CreatedAt == nil {
		t.Error("expected CreatedAt to be populated, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled sqlmock expectations: %s", err)
	}
}

// TestCreateUser_WithExistingCreatedAt verifies behavior when CreatedAt is already supplied
func TestCreateUser_WithExistingCreatedAt(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize sqlmock: %v", err)
	}
	defer db.Close()

	svc := main.UserServiceImpl{DB: db}

	customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	reqUser := &proto.User{
		FirstName: "Jane",
		LastName:  "Smith",
		UserName:  "janesmith",
		Email:     "jane@example.com",
		Password:  "password123",
		CreatedAt: timestamppb.New(customTime),
	}
	req := &proto.CreateUserRequest{User: reqUser}
	expectedID := int64(10)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			reqUser.FirstName,
			reqUser.LastName,
			reqUser.UserName,
			reqUser.Email,
			sqlmock.AnyArg(),
			customTime,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	res, err := svc.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	resUser := res.User

	if !resUser.CreatedAt.AsTime().Equal(customTime) {
		t.Errorf("expected CreatedAt %v, got %v", customTime, resUser.CreatedAt.AsTime())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled sqlmock expectations: %s", err)
	}
}

// TestCreateUser_DBError verifies gRPC error handling on database failures
func TestCreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize sqlmock: %v", err)
	}
	defer db.Close()

	svc := main.UserServiceImpl{DB: db}

	reqUser := &proto.User{
		FirstName: "Fail",
		LastName:  "User",
		UserName:  "failuser",
		Email:     "fail@example.com",
		Password:  "somepassword",
	}
	req := &proto.CreateUserRequest{User: reqUser}

	// Simulate database query error (e.g., unique constraint violation)
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			reqUser.FirstName,
			reqUser.LastName,
			reqUser.UserName,
			reqUser.Email,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	res, err := svc.CreateUser(context.Background(), req)

	// Assert error status
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got non-status error: %v", err)
	}

	if st.Code() != codes.Internal {
		t.Errorf("expected gRPC code Internal, got %v", st.Code())
	}

	if res != nil {
		t.Errorf("expected nil user response on error, got %v", res)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled sqlmock expectations: %s", err)
	}
}

func TestCreateUser_NilUser(t *testing.T) {
	service := main.UserServiceImpl{}
	_, err := service.CreateUser(context.Background(), &proto.CreateUserRequest{User: nil})

	if err == nil {
		t.Fatal("expected error for nil user request, got nil")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Internal {
		t.Errorf("expected gRPC status Internal, got %v", st.Code())
	}
}

func TestCreateUser_LongPasswordTruncation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	service := main.UserServiceImpl{DB: db}

	// 80-byte long password (> 72 bytes limit for bcrypt)
	longPassword := strings.Repeat("abcd", 20)
	reqUser := &proto.User{
		FirstName: "Jane",
		LastName:  "Smith",
		UserName:  "janesmith",
		Email:     "jane@example.com",
		Password:  longPassword,
	}
	req := &proto.CreateUserRequest{User: reqUser}
	expectedID := int64(10)
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			reqUser.FirstName,
			reqUser.LastName,
			reqUser.UserName,
			reqUser.Email,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	res, err := service.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify bcrypt hash matches the first 72 bytes of the original password
	err = bcrypt.CompareHashAndPassword([]byte(res.User.Password), []byte(longPassword[:72]))
	if err != nil {
		t.Errorf("password truncation test failed: %v", err)
	}
}
