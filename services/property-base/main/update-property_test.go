package main_test

import (
	"Schwarz--Internship--2026/services/property-base/main"
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestUpdateProperty(t *testing.T) {
	baseReq := &proto.UpdatePropertyRequest{
		Property: &proto.Property{
			Id:   10,
			Name: "Updated Beach House",
		},
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"name"}},
	}

	expectedSQL := `UPDATE properties SET name = \$1 WHERE id = \$2 RETURNING id, user_id, name, description, address, price, ST_X\(location::geometry\) AS lng, ST_Y\(location::geometry\) AS lat`

	tests := []struct {
		name         string
		request      *proto.UpdatePropertyRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.UpdatePropertyRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.UpdatePropertyResponse, req *proto.UpdatePropertyRequest)
	}{
		{
			name:    "Success",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdatePropertyRequest) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "description", "address", "price", "lng", "lat",
				}).AddRow(10, 1, "Updated Beach House", "Near coast", "789 Beach Rd", 2000, 30.1, 40.2)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.Property.GetName(), req.Property.GetId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.UpdatePropertyResponse, req *proto.UpdatePropertyRequest) {
				if res == nil || res.Property == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if res.Property.Name != "Updated Beach House" {
					t.Errorf("expected Name Updated Beach House, got %s", res.Property.Name)
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingFieldMask",
			request: &proto.UpdatePropertyRequest{
				Property: &proto.Property{Id: 10, Name: "New Name"},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingPropertyId",
			request: &proto.UpdatePropertyRequest{
				Property:  &proto.Property{Id: 0, Name: "New Name"},
				FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"name"}},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "InvalidFieldMaskPath",
			request: &proto.UpdatePropertyRequest{
				Property:  &proto.Property{Id: 10, Name: "New Name"},
				FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"unknown_field"}},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "PropertyNotFound",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdatePropertyRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.Property.GetName(), req.Property.GetId()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: codes.NotFound,
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
			res, err := svc.UpdateProperty(context.Background(), tt.request)

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
