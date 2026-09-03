package main_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"Schwarz--Internship--2026/services/conversation-base/main"
	"Schwarz--Internship--2026/services/conversation-base/main/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListConversations(t *testing.T) {
	baseRequest := &proto.ListConversationsRequest{
		UserId: 1,
	}

	expectedSQL := `SELECT id, user1_id, user2_id, created_at, updated_at FROM conversations WHERE \(user1_id = \$1 OR user2_id = \$2\) ORDER BY updated_at DESC`

	now := time.Now()

	tests := []struct {
		name         string
		request      *proto.ListConversationsRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.ListConversationsRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.ListConversationsResponse, req *proto.ListConversationsRequest)
	}{
		{
			name:    "Success",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListConversationsRequest) {
				rows := sqlmock.NewRows([]string{"id", "user1_id", "user2_id", "created_at", "updated_at"}).
					AddRow(100, 1, 2, now, now).
					AddRow(101, 1, 3, now, now)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUserId(), req.GetUserId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListConversationsResponse, req *proto.ListConversationsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				
				if len(res.Conversations) != 2 {
					t.Fatalf("expected 2 conversations, got %d", len(res.Conversations))
				}
				if res.Conversations[0].Id != 100 || res.Conversations[0].User2Id != 2 {
					t.Errorf("unexpected conversation at index 0: %v", res.Conversations[0])
				}
				if res.Conversations[1].Id != 101 || res.Conversations[1].User2Id != 3 {
					t.Errorf("unexpected conversation at index 1: %v", res.Conversations[1])
				}
			},
		},
		{
			name:    "SuccessEmptyList",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListConversationsRequest) {
				rows := sqlmock.NewRows([]string{"id", "user1_id", "user2_id", "created_at", "updated_at"})

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUserId(), req.GetUserId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListConversationsResponse, req *proto.ListConversationsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Conversations) != 0 {
					t.Errorf("expected 0 conversations, got %d", len(res.Conversations))
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListConversationsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "ZeroUserId",
			request: &proto.ListConversationsRequest{
				UserId: 0,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListConversationsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "InternalDatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListConversationsRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUserId(), req.GetUserId()).
					WillReturnError(errors.New("database query error"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.ListConversationsResponse, req *proto.ListConversationsRequest) {
				if res != nil {
					t.Errorf("expected nil response on error, got %v", res)
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

			if tt.setupMock != nil {
				tt.setupMock(mock, tt.request)
			}

			svc := main.ConversationServiceImpl{DB: db}
			res, err := svc.ListConversations(context.Background(), tt.request)

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