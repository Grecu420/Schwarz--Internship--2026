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

func TestCreateMessage(t *testing.T) {
	baseRequest := &proto.CreateMessageRequest{
		ConversationId: 1,
		SenderId:       2,
		Content:        "Hello, this is a test message!",
	}

	expectedSQL := `INSERT INTO messages \(conversation_id,sender_id,content\) VALUES \(\$1,\$2,\$3\) RETURNING id, created_at`

	now := time.Now()

	tests := []struct {
		name         string
		request      *proto.CreateMessageRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.CreateMessageRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.CreateMessageResponse, req *proto.CreateMessageRequest)
	}{
		{
			name:    "Success",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.CreateMessageRequest) {
				rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(100, now)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetConversationId(), req.GetSenderId(), req.GetContent()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.CreateMessageResponse, req *proto.CreateMessageRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if res.Message == nil {
					t.Fatalf("expected non-nil Message object inside response, got nil")
				}
				if res.Message.Id != 100 {
					t.Errorf("expected ID 100, got %d", res.Message.Id)
				}
				if res.Message.Content != req.Content {
					t.Errorf("expected Content %q, got %q", req.Content, res.Message.Content)
				}
				if res.Message.CreatedAt == nil {
					t.Errorf("expected CreatedAt to be set, got nil")
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingConversationId",
			request: &proto.CreateMessageRequest{
				ConversationId: 0,
				SenderId:       2,
				Content:        "Hello",
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingSenderId",
			request: &proto.CreateMessageRequest{
				ConversationId: 1,
				SenderId:       0,
				Content:        "Hello",
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingContent",
			request: &proto.CreateMessageRequest{
				ConversationId: 1,
				SenderId:       2,
				Content:        "",
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "InternalDatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.CreateMessageRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetConversationId(), req.GetSenderId(), req.GetContent()).
					WillReturnError(errors.New("database connection error"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.CreateMessageResponse, req *proto.CreateMessageRequest) {
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
			res, err := svc.CreateMessage(context.Background(), tt.request)

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