package main_test

import (
	"Schwarz--Internship--2026/services/common/pagination"
	"Schwarz--Internship--2026/services/property-base/main"
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListProperties(t *testing.T) {
	const queryOwnerLimit3 = `SELECT id, user_id, name, description, address, price, ST_X(location::geometry) AS lng, ST_Y(location::geometry) AS lat, image_urls 
	FROM properties 
	WHERE id >= $1 
	AND user_id = $2 
	ORDER BY id ASC LIMIT 3`

	const queryOwnerLimit11 = `SELECT id, user_id, name, description, address, price, ST_X(location::geometry) AS lng, ST_Y(location::geometry) AS lat, image_urls 
	FROM properties 
	WHERE id >= $1 
	AND user_id = $2 
	ORDER BY id ASC LIMIT 11`

	const queryAllFiltersLimit3 = `SELECT id, user_id, name, description, address, price, ST_X(location::geometry) AS lng, ST_Y(location::geometry) AS lat, image_urls 
	FROM properties 
	WHERE id >= $1 
	AND user_id = $2 
	AND name LIKE $3 
	AND price >= $4 AND price <= $5 
	AND ST_DWithin(location, ST_MakePoint($6, $7)::geography, $8) 
	ORDER BY id ASC LIMIT 3`

	baseOwnerFilter := []*proto.ListPropertiesFiltersOneOf{
		{
			Filter: &proto.ListPropertiesFiltersOneOf_Owner{
				Owner: &proto.FilterByOwnerId{Value: 100},
			},
		},
	}

	baseRequest := &proto.ListPropertiesRequest{
		PageSize: 2,
		Filters:  baseOwnerFilter,
	}

	// Filter Hash order: ownerID, priceMin, priceMax, name, centerLat, centerLong, radius
	owner100Hash := pagination.HashFilters(int64(100), int64(0), int64(0), "", float32(0), float32(0), float32(0))
	validToken, _ := pagination.BuildNextPageToken(10, owner100Hash)
	mismatchedToken, _ := pagination.BuildNextPageToken(10, "mismatched_hash")

	owner200Hash := pagination.HashFilters(int64(200), int64(0), int64(0), "", float32(0), float32(0), float32(0))
	differentOwnerToken, _ := pagination.BuildNextPageToken(5, owner200Hash)

	columns := []string{"id", "user_id", "name", "description", "address", "price", "lng", "lat", "image_urls"}

	values := [][]driver.Value{
		{1, 100, "Sunset Villa", "Cozy house", "123 Main St", 250000, 12.34, 56.78, pq.Array([]string{"main.png"})},
		{2, 100, "Ocean Apartment", "Beachfront view", "456 Beach Rd", 400000, 12.35, 56.79, pq.Array([]string{"main.png"})},
		{3, 100, "Mountain Cabin", "Quiet place", "789 Forest Ln", 150000, 12.36, 56.80, pq.Array([]string{"main.png"})},
	}

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
					AddRow(values[0]...).
					AddRow(values[1]...)

				mock.ExpectQuery(regexp.QuoteMeta(queryOwnerLimit3)).
					WithArgs(int64(0), int64(100)).
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
					AddRow(values[0]...).
					AddRow(values[1]...).
					AddRow(values[2]...)

				mock.ExpectQuery(regexp.QuoteMeta(queryOwnerLimit3)).
					WithArgs(int64(0), int64(100)).
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
				PageSize:      2,
				NextPageToken: validToken,
				Filters:       baseOwnerFilter,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(values[2]...)

				mock.ExpectQuery(regexp.QuoteMeta(queryOwnerLimit3)).
					WithArgs(int64(10), int64(100)).
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
		{name: "Success_AllFilters",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_Owner{Owner: &proto.FilterByOwnerId{Value: 100}}},
					{Filter: &proto.ListPropertiesFiltersOneOf_Name{Name: &proto.FilterByName{Value: "Villa"}}},
					{Filter: &proto.ListPropertiesFiltersOneOf_PriceRange{PriceRange: &proto.FilterByPriceRange{Min: 100000, Max: 500000}}},
					{Filter: &proto.ListPropertiesFiltersOneOf_Location{Location: &proto.FilterByLocation{
						Center: &proto.Location{Lat: 12.34, Long: 56.78},
						Radius: 10.0,
					}}},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(values[0]...)

				mock.ExpectQuery(regexp.QuoteMeta(queryAllFiltersLimit3)).
					WithArgs(int64(0), int64(100), "%Villa%", int64(100000), int64(500000), float32(56.78), float32(12.34), float32(10.0)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
		},
		{name: "FilterMismatch_TokenHashDoesNotMatchCurrentOwner",
			request: &proto.ListPropertiesRequest{
				PageSize:      2,
				NextPageToken: differentOwnerToken,
				Filters:       baseOwnerFilter,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "Invalid_DuplicateOwnerFilter",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_Owner{Owner: &proto.FilterByOwnerId{Value: 100}}},
					{Filter: &proto.ListPropertiesFiltersOneOf_Owner{Owner: &proto.FilterByOwnerId{Value: 200}}},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "Invalid_OwnerIdZeroOrNegative",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_Owner{Owner: &proto.FilterByOwnerId{Value: 0}}},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "Invalid_NameEmpty",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_Name{Name: &proto.FilterByName{Value: ""}}},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "Invalid_LocationMissingCenter",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_Location{Location: &proto.FilterByLocation{Radius: 10.0}}},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "Invalid_LocationNegativeRadius",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_Location{Location: &proto.FilterByLocation{Center: &proto.Location{Lat: 1.0, Long: 2.0}, Radius: -1.0}}},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "Invalid_PriceRangeNegativeMin",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_PriceRange{PriceRange: &proto.FilterByPriceRange{Min: -10, Max: 100}}},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "Invalid_PriceRangeMinGreaterThanMax",
			request: &proto.ListPropertiesRequest{
				PageSize: 2,
				Filters: []*proto.ListPropertiesFiltersOneOf{
					{Filter: &proto.ListPropertiesFiltersOneOf_PriceRange{PriceRange: &proto.FilterByPriceRange{Min: 500, Max: 100}}},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "InvalidPageSize_DefaultsToEleven",
			request: &proto.ListPropertiesRequest{
				PageSize: 0,
				Filters:  baseOwnerFilter,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(values[0]...).
					AddRow(values[1]...).
					AddRow(values[2]...)

				mock.ExpectQuery(regexp.QuoteMeta(queryOwnerLimit11)).
					WithArgs(int64(0), int64(100)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
		},
		{name: "MalformedNextPageToken",
			request: &proto.ListPropertiesRequest{
				PageSize:      2,
				NextPageToken: "invalid_base64_token",
				Filters:       baseOwnerFilter,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "NextPageToken_FilterMismatch",
			request: &proto.ListPropertiesRequest{
				PageSize:      2,
				NextPageToken: mismatchedToken,
				Filters:       baseOwnerFilter,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{name: "DatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListPropertiesRequest) {
				mock.ExpectQuery(regexp.QuoteMeta(queryOwnerLimit3)).
					WithArgs(int64(0), int64(100)).
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
