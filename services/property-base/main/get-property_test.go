package main_test

import (
	"Schwarz--Internship--2026/services/property-base/main"
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetProperty(t *testing.T) {
	baseReq := &proto.GetPropertyRequest{Id: 10}
	expectedSQL := `SELECT id, user_id, name, description, address, price, ST_X\(location::geometry\) AS lng, ST_Y\(location::geometry\) AS lat FROM properties WHERE id = \$1`

	tests := []struct {
		name         string
		request      *proto.GetPropertyRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.GetPropertyRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.GetPropertyResponse, req *proto.GetPropertyRequest)
	}{
		{
			name:    "Success",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetPropertyRequest) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "description", "address", "price", "lng", "lat",
				}).AddRow(10, 1, "City Loft", "Downtown view", "456 Central Ave", 800, 10.5, 20.5)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.GetPropertyResponse, req *proto.GetPropertyRequest) {
				if res == nil || res.Property == nil {
					t.Fatalf("expected non-nil property response, got nil")
				}
				if res.Property.Id != 10 {
					t.Errorf("expected ID 10, got %d", res.Property.Id)
				}
				if res.Property.Name != "City Loft" {
					t.Errorf("expected Name City Loft, got %s", res.Property.Name)
				}
			},
		},
		{
			name:    "PropertyNotFound",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetPropertyRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: codes.NotFound,
		},
		{
			name:         "MissingId",
			request:      &proto.GetPropertyRequest{Id: 0},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetPropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "DatabaseError",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetPropertyRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnError(errors.New("db connection failure"))
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
			res, err := svc.GetProperty(context.Background(), tt.request)

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
