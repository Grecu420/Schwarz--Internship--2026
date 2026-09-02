package main_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"Schwarz--Internship--2026/services/message-base/main"
	"Schwarz--Internship--2026/services/message-base/main/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetMessageList(t *testing.T) {
	baseRequest := &proto.GetMessageListRequest{
		ConversationId: 1,
	}

	expectedSQL := `SELECT id, conversation_id, sender_id, content, created_at FROM messages WHERE conversation_id = \$1`

	now := time.Now()

	tests := []struct {
		name         string
		request      *proto.GetMessageListRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.GetMessageListRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.GetMessageListResponse, req *proto.GetMessageListRequest)
	}{
		{
			name:    "Success",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetMessageListRequest) {
				rows := sqlmock.NewRows([]string{"id", "conversation_id", "sender_id", "content", "created_at"}).
					AddRow(100, 1, 2, "First message", now).
					AddRow(101, 1, 3, "Second message", now)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetConversationId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.GetMessageListResponse, req *proto.GetMessageListRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Messages) != 2 {
					t.Fatalf("expected 2 messages, got %d", len(res.Messages))
				}
				if res.Messages[0].Id != 100 || res.Messages[0].Content != "First message" {
					t.Errorf("unexpected message at index 0: %v", res.Messages[0])
				}
				if res.Messages[1].Id != 101 || res.Messages[1].Content != "Second message" {
					t.Errorf("unexpected message at index 1: %v", res.Messages[1])
				}
			},
		},
		{
			name:    "SuccessEmptyList",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetMessageListRequest) {
				rows := sqlmock.NewRows([]string{"id", "conversation_id", "sender_id", "content", "created_at"})

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetConversationId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.GetMessageListResponse, req *proto.GetMessageListRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Messages) != 0 {
					t.Errorf("expected 0 messages, got %d", len(res.Messages))
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetMessageListRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "ZeroConversationId",
			request: &proto.GetMessageListRequest{
				ConversationId: 0,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetMessageListRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "InternalDatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetMessageListRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetConversationId()).
					WillReturnError(errors.New("database query error"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.GetMessageListResponse, req *proto.GetMessageListRequest) {
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

			svc := main.MessageServiceImpl{DB: db}
			res, err := svc.GetMessageList(context.Background(), tt.request)

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
