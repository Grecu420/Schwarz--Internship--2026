package main_test

import (
	"Schwarz--Internship--2026/services/message-base/main"
	"Schwarz--Internship--2026/services/message-base/main/proto"
	"context"
	"database/sql"
	"testing"
	"time"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestUpdateMessage(t *testing.T) {
	now := time.Now()

	baseMessage := &proto.Message{
		Id:             10,
		ConversationId: 1,
		SenderId:       2,
		IsRead:         true,
	}

	baseFieldMask := &fieldmaskpb.FieldMask{
		Paths: []string{"is_read"},
	}

	expectedSQL := `UPDATE messages SET is_read = \$1 WHERE id = \$2 AND conversation_id = \$3 AND sender_id = \$4 RETURNING id, conversation_id, sender_id, content, created_at, is_read`

	tests := []struct {
		name         string
		request      *proto.UpdateMessageRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.UpdateMessageResponse)
	}{
		{
			name: "Success",
			request: &proto.UpdateMessageRequest{
				Message:   baseMessage,
				FieldMask: baseFieldMask,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {
				rows := sqlmock.NewRows([]string{"id", "conversation_id", "sender_id", "content", "created_at", "is_read"}).
					AddRow(10, 1, 2, "Hello world", now, true)

				mock.ExpectQuery(expectedSQL).
					WithArgs(true, req.Message.GetId(), req.Message.GetConversationId(), req.Message.GetSenderId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.UpdateMessageResponse) {
				if res == nil || res.Message == nil {
					t.Fatalf("expected non-nil response and message, got nil")
				}
				if res.Message.Id != 10 {
					t.Errorf("expected ID 10, got %d", res.Message.Id)
				}
				if !res.Message.IsRead {
					t.Errorf("expected IsRead to be true, got false")
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingMessage",
			request: &proto.UpdateMessageRequest{
				Message:   nil,
				FieldMask: baseFieldMask,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingFieldMask",
			request: &proto.UpdateMessageRequest{
				Message:   baseMessage,
				FieldMask: nil,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "EmptyFieldMaskPaths",
			request: &proto.UpdateMessageRequest{
				Message:   baseMessage,
				FieldMask: &fieldmaskpb.FieldMask{Paths: []string{}},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingMessageId",
			request: &proto.UpdateMessageRequest{
				Message: &proto.Message{
					Id:             0,
					ConversationId: 1,
				},
				FieldMask: baseFieldMask,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "InvalidFieldMaskPath",
			request: &proto.UpdateMessageRequest{
				Message: baseMessage,
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"invalid_path"},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "NotFound",
			request: &proto.UpdateMessageRequest{
				Message:   baseMessage,
				FieldMask: baseFieldMask,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(true, req.Message.GetId(), req.Message.GetConversationId(), req.Message.GetSenderId()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: codes.NotFound,
		},
		{
			name: "InternalDatabaseError",
			request: &proto.UpdateMessageRequest{
				Message:   baseMessage,
				FieldMask: baseFieldMask,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateMessageRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(true, req.Message.GetId(), req.Message.GetConversationId(), req.Message.GetSenderId()).
					WillReturnError(errors.New("db connection timeout"))
			},
			expectedCode: codes.Internal,
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
			res, err := svc.UpdateMessage(context.Background(), tt.request)

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
				tt.validate(t, res)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}