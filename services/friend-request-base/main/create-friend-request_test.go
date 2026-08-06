package main_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/friend-request-base/main"
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
)

func TestCreateFriendRequest_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize sqlmock: %v", err)
	}
	defer db.Close()

	svc := main.FriendRequestServiceImpl{DB: db}

	req := &proto.CreateFriendRequestRequest{
		SenderId:   "user_sender_123",
		ReceiverId: "user_receiver_456",
	}

	mock.ExpectExec(`INSERT INTO friend_requests`).
		WithArgs(
			sqlmock.AnyArg(),
			req.SenderId,
			req.ReceiverId,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := svc.FriendRequestEndpoint(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if res == nil || res.Request == nil {
		t.Fatalf("expected non-nil response and request object")
	}

	resReq := res.Request

	if resReq.Id == "" {
		t.Error("expected generated UUID ID to be populated, got empty string")
	}

	if resReq.SenderId != req.SenderId {
		t.Errorf("expected SenderId %s, got %s", req.SenderId, resReq.SenderId)
	}

	if resReq.ReceiverId != req.ReceiverId {
		t.Errorf("expected ReceiverId %s, got %s", req.ReceiverId, resReq.ReceiverId)
	}

	if resReq.Status != proto.RequestStatus_STATUS_PENDING {
		t.Errorf("expected status STATUS_PENDING, got %v", resReq.Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled sqlmock expectations: %s", err)
	}
}

func TestCreateFriendRequest_Validation_EmptyFields(t *testing.T) {
	svc := main.FriendRequestServiceImpl{}

	req := &proto.CreateFriendRequestRequest{
		SenderId:   "",
		ReceiverId: "user_receiver_456",
	}

	_, err := svc.FriendRequestEndpoint(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty sender_id, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got non-status error: %v", err)
	}

	if st.Code() != codes.InvalidArgument && st.Code() != codes.Internal {
		t.Errorf("expected gRPC code InvalidArgument or Internal, got %v", st.Code())
	}
}

func TestCreateFriendRequest_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize sqlmock: %v", err)
	}
	defer db.Close()

	svc := main.FriendRequestServiceImpl{DB: db}

	req := &proto.CreateFriendRequestRequest{
		SenderId:   "user_sender_123",
		ReceiverId: "user_receiver_456",
	}

	mock.ExpectExec(`INSERT INTO friend_requests`).
		WithArgs(
			sqlmock.AnyArg(),
			req.SenderId,
			req.ReceiverId,
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("db connection timeout"))

	res, err := svc.FriendRequestEndpoint(context.Background(), req)

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
		t.Errorf("expected nil response on DB error, got %v", res)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled sqlmock expectations: %s", err)
	}
}

func TestCreateFriendRequest_NilRequest(t *testing.T) {
	svc := main.FriendRequestServiceImpl{}
	_, err := svc.FriendRequestEndpoint(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got non-status error: %v", err)
	}

	if st.Code() != codes.InvalidArgument && st.Code() != codes.Internal {
		t.Errorf("expected gRPC status InvalidArgument or Internal, got %v", st.Code())
	}
}
