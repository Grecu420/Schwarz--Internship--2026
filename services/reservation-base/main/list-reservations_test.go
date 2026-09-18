package main_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/reservation-base/main"
	"Schwarz--Internship--2026/services/reservation-base/main/proto"
)

func TestListReservations(t *testing.T) {
	baseRequest := &proto.ListReservationsRequest{
		PageSize: 2,
	}

	customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	checkInTime := time.Date(2026, time.August, 15, 0, 0, 0, 0, time.UTC)
	checkOutTime := time.Date(2026, time.August, 20, 0, 0, 0, 0, time.UTC)

	emptyHash := main.HashReservationFilter(nil, 0, 0, proto.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED)
	validToken, _ := main.BuildNextPageToken(10, emptyHash)
	mismatchedToken, _ := main.BuildNextPageToken(10, "mismatched_hash")

	userIdHash := main.HashReservationFilter(nil, 200, 0, proto.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED)
	filteredValidToken, _ := main.BuildNextPageToken(5, userIdHash)

	columns := []string{"id", "property_id", "user_id", "owner_id", "check_in_date", "check_out_date", "status", "created_at", "updated_at"}

	tests := []struct {
		name         string
		request      *proto.ListReservationsRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.ListReservationsResponse, req *proto.ListReservationsRequest)
	}{
		{
			name:    "Success_SinglePage",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(1, 100, 200, 300, checkInTime, checkOutTime, "RESERVATION_STATUS_PENDING", customTime, customTime).
					AddRow(2, 101, 201, 301, checkInTime, checkOutTime, "RESERVATION_STATUS_CONFIRMED", customTime, customTime)

				query := `SELECT id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at FROM reservations ORDER BY id ASC LIMIT $1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListReservationsResponse, req *proto.ListReservationsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Reservations) != 2 {
					t.Errorf("expected 2 reservations, got %d", len(res.Reservations))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
				if res.Reservations[0].Id != 1 || res.Reservations[0].PropertyId != 100 {
					t.Errorf("unexpected first reservation content: %v", res.Reservations[0])
				}
			},
		},
		{
			name:    "Success_HasNextPageToken",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(1, 100, 200, 300, checkInTime, checkOutTime, "RESERVATION_STATUS_PENDING", customTime, customTime).
					AddRow(2, 101, 201, 301, checkInTime, checkOutTime, "RESERVATION_STATUS_CONFIRMED", customTime, customTime).
					AddRow(3, 102, 202, 302, checkInTime, checkOutTime, "RESERVATION_STATUS_CANCELLED", customTime, customTime)

				query := `SELECT id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at FROM reservations ORDER BY id ASC LIMIT $1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListReservationsResponse, req *proto.ListReservationsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Reservations) != 2 {
					t.Errorf("expected 2 sliced reservations, got %d", len(res.Reservations))
				}
				if res.NextPageToken == "" {
					t.Error("expected non-empty nextPageToken, got empty")
				}
			},
		},
		{
			name: "Success_WithNextPageToken",
			request: &proto.ListReservationsRequest{
				PageSize:      2,
				NextPageToken: validToken,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(15, 100, 200, 300, checkInTime, checkOutTime, "RESERVATION_STATUS_PENDING", customTime, customTime)

				query := `SELECT id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at FROM reservations WHERE id > $1 ORDER BY id ASC LIMIT $2`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(10), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListReservationsResponse, req *proto.ListReservationsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Reservations) != 1 {
					t.Errorf("expected 1 reservation, got %d", len(res.Reservations))
				}
				if res.NextPageToken != "" {
					t.Errorf("expected empty nextPageToken, got %s", res.NextPageToken)
				}
			},
		},
		{
			name: "Success_FilterByUserId",
			request: &proto.ListReservationsRequest{
				PageSize: 2,
				Filters: []*proto.ListReservationsFiltersOneOf{
					{
						Filter: &proto.ListReservationsFiltersOneOf_UserId{
							UserId: &proto.FilterByUserId{Value: 200},
						},
					},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(1, 100, 200, 300, checkInTime, checkOutTime, "RESERVATION_STATUS_PENDING", customTime, customTime).
					AddRow(4, 105, 200, 305, checkInTime, checkOutTime, "RESERVATION_STATUS_CONFIRMED", customTime, customTime)

				query := `SELECT id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at FROM reservations WHERE user_id = $1 ORDER BY id ASC LIMIT $2`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(200), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListReservationsResponse, req *proto.ListReservationsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Reservations) != 2 {
					t.Fatalf("expected 2 reservations, got %d", len(res.Reservations))
				}
				for _, r := range res.Reservations {
					if r.UserId != 200 {
						t.Errorf("expected UserId 200, got %d", r.UserId)
					}
				}
			},
		},
		{
			name: "Success_CombinedFiltersAndPaginationToken",
			request: &proto.ListReservationsRequest{
				PageSize: 2,
				Filters: []*proto.ListReservationsFiltersOneOf{
					{
						Filter: &proto.ListReservationsFiltersOneOf_UserId{
							UserId: &proto.FilterByUserId{Value: 200},
						},
					},
				},
				NextPageToken: filteredValidToken,
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {
				rows := sqlmock.NewRows(columns).
					AddRow(6, 100, 200, 300, checkInTime, checkOutTime, "RESERVATION_STATUS_PENDING", customTime, customTime)

				query := `SELECT id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at FROM reservations WHERE user_id = $1 AND id > $2 ORDER BY id ASC LIMIT $3`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(200), int64(5), int64(3)).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.ListReservationsResponse, req *proto.ListReservationsRequest) {
				if res == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if len(res.Reservations) != 1 {
					t.Fatalf("expected 1 reservation, got %d", len(res.Reservations))
				}
				if res.Reservations[0].Id != 6 {
					t.Errorf("expected reservation ID 6, got %d", res.Reservations[0].Id)
				}
			},
		},
		{
			name: "FilterMismatch_TokenHashDoesNotMatchCurrentFilters",
			request: &proto.ListReservationsRequest{
				PageSize: 2,
				Filters: []*proto.ListReservationsFiltersOneOf{
					{
						Filter: &proto.ListReservationsFiltersOneOf_UserId{
							UserId: &proto.FilterByUserId{Value: 200},
						},
					},
				},
				NextPageToken: validToken, 
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "FilterMismatch_DuplicateUserIdFilter",
			request: &proto.ListReservationsRequest{
				PageSize: 2,
				Filters: []*proto.ListReservationsFiltersOneOf{
					{
						Filter: &proto.ListReservationsFiltersOneOf_UserId{
							UserId: &proto.FilterByUserId{Value: 200},
						},
					},
					{
						Filter: &proto.ListReservationsFiltersOneOf_UserId{
							UserId: &proto.FilterByUserId{Value: 201},
						},
					},
				},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "DatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {
				query := `SELECT id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at FROM reservations ORDER BY id ASC LIMIT $1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(3)).
					WillReturnError(errors.New("db query execution failed"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.ListReservationsResponse, req *proto.ListReservationsRequest) {
				if res != nil {
					t.Errorf("expected nil response on error, got %v", res)
				}
			},
		},
		{
			name: "InvalidPageSize_Zero",
			request: &proto.ListReservationsRequest{
				PageSize: 0,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MalformedNextPageToken",
			request: &proto.ListReservationsRequest{
				PageSize:      2,
				NextPageToken: "invalid_base64_token",
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "NextPageToken_FilterMismatch",
			request: &proto.ListReservationsRequest{
				PageSize:      2,
				NextPageToken: mismatchedToken,
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.ListReservationsRequest) {},
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
			svc := &main.ReservationServiceImpl{DB: db}
			res, err := svc.ListReservations(context.Background(), tt.request)

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