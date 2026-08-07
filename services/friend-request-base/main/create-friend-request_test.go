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

func TestFriendRequestEndpoint(t *testing.T) {

	tests := []struct {
		name        string
		request     *proto.CreateFriendRequestRequest
		mockDBSetup func(mock sqlmock.Sqlmock)
		wantCode    codes.Code
	}{
		{
			name: "Success - Request created successfully",
			request: &proto.CreateFriendRequestRequest{
				SenderId:   "user_1",
				ReceiverId: "user_2",
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO friend_requests`).
					WithArgs("user_1", "user_2", sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(10)))
			},
			wantCode: codes.OK,
		},
		{
			name: "Validation - Empty fields (sender_id is missing)",
			request: &proto.CreateFriendRequestRequest{
				SenderId:   "",
				ReceiverId: "user_2",
			},
			mockDBSetup: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Validation - Same user (sender == receiver)",
			request: &proto.CreateFriendRequestRequest{
				SenderId:   "user_same",
				ReceiverId: "user_same",
			},
			mockDBSetup: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "Validation - Request is nil",
			request:     nil,
			mockDBSetup: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "DB Error - The database crashed during the save operation",
			request: &proto.CreateFriendRequestRequest{
				SenderId:   "user_1",
				ReceiverId: "user_2",
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO friend_requests`).
					WithArgs("user_1", "user_2", sqlmock.AnyArg()).
					WillReturnError(errors.New("db connection timeout"))
			},
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to initialize sqlmock: %v", err)
			}
			defer db.Close()

			if tt.mockDBSetup != nil {
				tt.mockDBSetup(mock)
			}

			svc := main.FriendRequestServiceImpl{DB: db}

			res, err := svc.FriendRequestEndpoint(context.Background(), tt.request)

			if tt.wantCode == codes.OK {

				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if res == nil || res.Request == nil {
					t.Fatalf("expected valid response, got nil")
				}
				if res.Request.Id == 0 {
					t.Errorf("expected generated ID > 0, got 0")
				}
			} else {

				if err == nil {
					t.Fatalf("expected error with code %v, got nil", tt.wantCode)
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got standard error: %v", err)
				}
				if st.Code() != tt.wantCode {
					t.Errorf("expected gRPC code %v, got %v", tt.wantCode, st.Code())
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}
