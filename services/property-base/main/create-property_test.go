package main_test

import (
	"Schwarz--Internship--2026/services/property-base/main"
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateProperty(t *testing.T) {
	baseReq := &proto.CreatePropertyRequest{
		UserId:      1,
		Name:        "Grand Villa",
		Description: "Spacious house with pool",
		Address:     "123 Ocean Drive",
		Price:       1500,
		Location:    &proto.Location{Long: 12.34, Lat: 56.78},
		ImageUrls:   []string{"main.png", "other.png"},
	}

	expectedSQL := `INSERT INTO properties \(user_id,name,description,address,location,price,image_urls\) VALUES \(\$1,\$2,\$3,\$4,ST_MakePoint\(\$5, \$6\)::geography,\$7,\$8\) RETURNING id`

	tests := []struct {
		name         string
		request      *proto.CreatePropertyRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.CreatePropertyRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.CreatePropertyResponse, req *proto.CreatePropertyRequest)
	}{
		{
			name:    "Success",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.CreatePropertyRequest) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(100)
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUserId(), req.GetName(), req.GetDescription(), req.GetAddress(), req.GetLocation().GetLong(), req.GetLocation().GetLat(), req.GetPrice(), pq.Array(req.GetImageUrls())).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.CreatePropertyResponse, req *proto.CreatePropertyRequest) {
				if res == nil || res.Property == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if res.Property.Id != 100 {
					t.Errorf("expected ID 100, got %d", res.Property.Id)
				}
				if res.Property.Name != req.Name {
					t.Errorf("expected Name %s, got %s", req.Name, res.Property.Name)
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingUserId",
			request: &proto.CreatePropertyRequest{
				Name:     "Villa",
				Price:    500,
				Location: &proto.Location{},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingPrice",
			request: &proto.CreatePropertyRequest{
				UserId:   1,
				Name:     "Villa",
				Location: &proto.Location{},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingName",
			request: &proto.CreatePropertyRequest{
				UserId:   1,
				Price:    500,
				Location: &proto.Location{},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.CreatePropertyRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "DatabaseError",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.CreatePropertyRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetUserId(), req.GetName(), req.GetDescription(), req.GetAddress(), req.GetLocation().GetLong(), req.GetLocation().GetLat(), req.GetPrice(), pq.Array(req.GetImageUrls())).
					WillReturnError(errors.New("db insert failure"))
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
			res, err := svc.CreateProperty(context.Background(), tt.request)

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
