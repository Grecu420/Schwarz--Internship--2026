package main

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
)

func TestUpdateFriendRequest(t *testing.T) {

	tests := []struct {
		name        string
		request     *proto.UpdateFriendRequestRequest
		mockDBSetup func(mock sqlmock.Sqlmock)
		wantCode    codes.Code
	}{
		{
			name: "Success - Status updated successfully",
			request: &proto.UpdateFriendRequestRequest{
				FriendRequest: &proto.FriendRequest{
					Id:     10,
					Status: proto.RequestStatus_STATUS_ACCEPTED,
				},
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"status"},
				},
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectExec(`UPDATE friend_requests`).
					WithArgs(proto.RequestStatus_STATUS_ACCEPTED, int64(10)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantCode: codes.OK,
		},
		{
			name: "Validation - Invalid field in mask (sender_id is not allowed)",
			request: &proto.UpdateFriendRequestRequest{
				FriendRequest: &proto.FriendRequest{
					Id:       10,
					SenderId: "hacker_user",
				},
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"status", "sender_id"},
				},
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
			name: "Validation - Field mask is missing",
			request: &proto.UpdateFriendRequestRequest{
				FriendRequest: &proto.FriendRequest{
					Id: 10,
				},
				FieldMask: nil,
			},
			mockDBSetup: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Validation - ID is missing",
			request: &proto.UpdateFriendRequestRequest{
				FriendRequest: &proto.FriendRequest{
					Id: 0,
				},
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"status"},
				},
			},
			mockDBSetup: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "DB Error - Friend request not found (0 rows affected)",
			request: &proto.UpdateFriendRequestRequest{
				FriendRequest: &proto.FriendRequest{
					Id:     99,
					Status: proto.RequestStatus_STATUS_ACCEPTED,
				},
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"status"},
				},
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE friend_requests`).
					WithArgs(proto.RequestStatus_STATUS_ACCEPTED, int64(99)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantCode: codes.NotFound,
		},
		{
			name: "DB Error - The database crashed during the update operation",
			request: &proto.UpdateFriendRequestRequest{
				FriendRequest: &proto.FriendRequest{
					Id:     10,
					Status: proto.RequestStatus_STATUS_ACCEPTED,
				},
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"status"},
				},
			},
			mockDBSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE friend_requests`).
					WithArgs(proto.RequestStatus_STATUS_ACCEPTED, int64(10)).
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

			svc := FriendRequestServiceImpl{DB: db}

			res, err := svc.UpdateFriendRequestEndpoint(context.Background(), tt.request)

			if tt.wantCode == codes.OK {

				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if res == nil || res.FriendRequest == nil {
					t.Fatalf("expected valid response, got nil")
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
