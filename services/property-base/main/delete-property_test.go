package main_test

import (
	"Schwarz--Internship--2026/services/property-base/main"
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteProperty(t *testing.T) {
	baseReq := &proto.DeletePropertyRequest{Id: 10}
	expectedSQL := `DELETE FROM properties WHERE id = \$1`

	tests := []struct {
		name         string
		request      *proto.DeletePropertyRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.DeletePropertyRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.DeletePropertyResponse, req *proto.DeletePropertyRequest)
	}{
		{
			name:    "Success",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.DeletePropertyRequest) {
				mock.ExpectExec(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.DeletePropertyResponse, req *proto.DeletePropertyRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
			},
		},
		{
			name:    "PropertyNotFound",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.DeletePropertyRequest) {
				mock.ExpectExec(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedCode: codes.NotFound,
		},
		{
			name:         "MissingId",
			request:      &proto.DeletePropertyRequest{},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.DeletePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "DatabaseError",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.DeletePropertyRequest) {
				mock.ExpectExec(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnError(errors.New("db error"))
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

			svc := main.PropertyServiceImpl{DB: db}
			res, err := svc.DeleteProperty(context.Background(), tt.request)

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
