package main_test

import (
	"context"
	"errors"
	"testing"

	"Schwarz--Internship--2026/services/conversation-base/main"
	"Schwarz--Internship--2026/services/conversation-base/main/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateConversation(t *testing.T) {
	baseRequest := &proto.CreateConversationRequest{
		User1Id: 1,
		User2Id: 2,
	}

	expectedSQL := `INSERT INTO conversations \(user1_id,user2_id\) VALUES \(\$1,\$2\) RETURNING id`

	tests := []struct {
		name         string
		request      *proto.CreateConversationRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.CreateConversationRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.CreateConversationResponse, req *proto.CreateConversationRequest)
	}{
		{
			name:    "Success",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.CreateConversationRequest) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(100)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUser1Id(), req.GetUser2Id()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.CreateConversationResponse, req *proto.CreateConversationRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if res.Id != 100 {
					t.Errorf("expected ID 100, got %d", res.Id)
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreateConversationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingUser1",
			request: &proto.CreateConversationRequest{
				User1Id: 0,
				User2Id: 2,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreateConversationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingUser2",
			request: &proto.CreateConversationRequest{
				User1Id: 1,
				User2Id: 0,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreateConversationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "ConversationAlreadyExists",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.CreateConversationRequest) {
				pgErr := &pq.Error{Code: "23505"}
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUser1Id(), req.GetUser2Id()).
					WillReturnError(pgErr)
			},
			expectedCode: codes.AlreadyExists,
			validate: func(t *testing.T, res *proto.CreateConversationResponse, req *proto.CreateConversationRequest) {
				if res != nil {
					t.Errorf("expected nil response on error, got %v", res)
				}
			},
		},
		{
			name:    "InternalDatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.CreateConversationRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUser1Id(), req.GetUser2Id()).
					WillReturnError(errors.New("database connection error"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.CreateConversationResponse, req *proto.CreateConversationRequest) {
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
			res, err := svc.CreateConversation(context.Background(), tt.request)

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
