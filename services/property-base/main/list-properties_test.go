package main_test

import (
	"Schwarz--Internship--2026/services/common/pagination"
	"Schwarz--Internship--2026/services/property-base/main"
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListProperties(t *testing.T) {
	const query_p2 = `SELECT id, user_id, name, description, address, price, ST_X(location::geometry) AS lng, ST_Y(location::geometry) AS lat 
					FROM properties 
					WHERE user_id = $1 
					AND id >= $2 
					ORDER BY id ASC 
					LIMIT 3`

	const query_default = `SELECT id, user_id, name, description, address, price, ST_X(location::geometry) AS lng, ST_Y(location::geometry) AS lat 
					FROM properties 
					WHERE user_id = $1 
					AND id >= $2 
					ORDER BY id ASC 
					LIMIT 11`

	baseRequest := &proto.ListPropertiesRequest{
		OwnerId:  100,
		PageSize: 2,
	}

	owner100Hash := pagination.HashFilters(100)
	validToken, _ := pagination.BuildNextPageToken(10, owner100Hash)
	mismatchedToken, _ := pagination.BuildNextPageToken(10, "mismatched_hash")

	owner200Hash := pagination.HashFilters(200)
	differentOwnerToken, _ := pagination.BuildNextPageToken(5, owner200Hash)

	columns := []string{"id", "user_id", "name", "description", "address", "price", "lng", "lat"}

	tests := []struct {
		name         string
		request      *proto.ListPropertiesRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.ListPropertiesResponse, req *proto.ListPropertiesRequest)
	}{
		{name: "Success_SinglePage",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(1, 100, "Sunset Villa", "Cozy house", "123 Main St", 250000, 12.34, 56.78).
					AddRow(2, 100, "Ocean Apartment", "Beachfront view", "456 Beach Rd", 400000, 12.35, 56.79)

				mock.ExpectQuery(regexp.QuoteMeta(query_p2)).
					WithArgs(int64(100), int64(0)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListPropertiesResponse, req *proto.ListPropertiesRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Properties) != 2 {
					t.Errorf("expected 2 properties, got %d", len(res.Properties))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
				if res.Properties[0].Id != 1 {
					t.Errorf("unexpected first property: %v", res.Properties[0])
				}
			},
		},
		{name: "Success_HasNextPageToken",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(1, 100, "Sunset Villa", "Cozy house", "123 Main St", 250000, 12.34, 56.78).
					AddRow(2, 100, "Ocean Apartment", "Beachfront view", "456 Beach Rd", 400000, 12.35, 56.79).
					AddRow(3, 100, "Mountain Cabin", "Quiet place", "789 Forest Ln", 150000, 12.36, 56.80)

				mock.ExpectQuery(regexp.QuoteMeta(query_p2)).
					WithArgs(int64(100), int64(0)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListPropertiesResponse, req *proto.ListPropertiesRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Properties) != 2 {
					t.Errorf("expected 2 properties, got %d", len(res.Properties))
				}
				if res.NextPageToken == "" {
					t.Error("expected non-empty nextPageToken, got empty")
				}
			},
		},
		{name: "Success_WithNextPageToken",
			request: &proto.ListPropertiesRequest{
				OwnerId:       100,
				PageSize:      2,
				NextPageToken: validToken,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(10, 100, "Downtown Loft", "Modern design", "101 City Center", 500000, 12.37, 56.81)

				mock.ExpectQuery(regexp.QuoteMeta(query_p2)).
					WithArgs(int64(100), int64(10)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListPropertiesResponse, req *proto.ListPropertiesRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Properties) != 1 {
					t.Errorf("expected 1 property, got %d", len(res.Properties))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
			},
		},
		{name: "FilterMismatch_TokenHashDoesNotMatchCurrentOwner",
			request: &proto.ListPropertiesRequest{
				OwnerId:       100,
				PageSize:      2,
				NextPageToken: differentOwnerToken,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "MissingOwnerId",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "InvalidPageSize",
			request: &proto.ListPropertiesRequest{
				OwnerId:  100,
				PageSize: 0,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(1, 100, "Sunset Villa", "Cozy house", "123 Main St", 250000, 12.34, 56.78).
					AddRow(2, 100, "Ocean Apartment", "Beachfront view", "456 Beach Rd", 400000, 12.35, 56.79).
					AddRow(3, 100, "Mountain Cabin", "Quiet place", "789 Forest Ln", 150000, 12.36, 56.80)

				mock.ExpectQuery(regexp.QuoteMeta(query_default)).
					WithArgs(int64(100), int64(0)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
		},
		{name: "MalformedNextPageToken",
			request: &proto.ListPropertiesRequest{
				OwnerId:       100,
				PageSize:      2,
				NextPageToken: "invalid_base64_token",
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "NextPageToken_FilterMismatch",
			request: &proto.ListPropertiesRequest{
				OwnerId:       100,
				PageSize:      2,
				NextPageToken: mismatchedToken,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "DatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				mock.ExpectQuery(regexp.QuoteMeta(query_p2)).
					WithArgs(int64(100), int64(0)).
					WillReturnError(errors.New("db query execution failed"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.ListPropertiesResponse, req *proto.ListPropertiesRequest) {
				if res != nil {
					t.Errorf("expected nil response on error, got %v", res)
				}
			},
		},
		{name: "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
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
			res, err := svc.ListProperties(context.Background(), tt.request)

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
