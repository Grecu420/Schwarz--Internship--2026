package main_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"Schwarz--Internship--2026/services/reservation-base/main"
	"Schwarz--Internship--2026/services/reservation-base/main/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetReservation(t *testing.T) {
	customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	checkInTime := time.Date(2026, time.August, 15, 0, 0, 0, 0, time.UTC)
	checkOutTime := time.Date(2026, time.August, 20, 0, 0, 0, 0, time.UTC)
	targetId := int64(10)

	baseRequest := &proto.GetReservationRequest{
		Id: targetId,
	}

	expectedSQL := `SELECT id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at FROM reservations WHERE id = \$1`

	tests := []struct {
		name         string
		request      *proto.GetReservationRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.GetReservationRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.GetReservationResponse, req *proto.GetReservationRequest)
	}{
		{
			name:    "Success",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetReservationRequest) {
				rows := sqlmock.NewRows([]string{
					"id", "property_id", "user_id", "owner_id", "check_in_date", "check_out_date", "status", "created_at", "updated_at",
				}).AddRow(10, 100, 200, 300, checkInTime, checkOutTime, "RESERVATION_STATUS_CONFIRMED", customTime, customTime)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.GetReservationResponse, req *proto.GetReservationRequest) {
				if res == nil || res.Reservation == nil {
					t.Fatalf("expected non-nil reservation response, got nil")
				}
				if res.Reservation.Id != 10 {
					t.Errorf("expected ID %d, got %d", 10, res.Reservation.Id)
				}
				if res.Reservation.PropertyId != 100 {
					t.Errorf("expected PropertyId 100, got %d", res.Reservation.PropertyId)
				}
				if res.Reservation.UserId != 200 {
					t.Errorf("expected UserId 200, got %d", res.Reservation.UserId)
				}
				if res.Reservation.OwnerId != 300 {
					t.Errorf("expected OwnerId 300, got %d", res.Reservation.OwnerId)
				}
				if res.Reservation.CheckInDate != "2026-08-15" {
					t.Errorf("expected CheckInDate 2026-08-15, got %s", res.Reservation.CheckInDate)
				}
				if res.Reservation.CheckOutDate != "2026-08-20" {
					t.Errorf("expected CheckOutDate 2026-08-20, got %s", res.Reservation.CheckOutDate)
				}
				if res.Reservation.Status != proto.ReservationStatus_RESERVATION_STATUS_CONFIRMED {
					t.Errorf("expected Status CONFIRMED, got %s", res.Reservation.Status)
				}
				if res.Reservation.CreatedAt == nil {
					t.Fatal("expected CreatedAt to be non-nil")
				}
				if !res.Reservation.CreatedAt.AsTime().Equal(customTime) {
					t.Errorf("expected CreatedAt %v, got %v", customTime, res.Reservation.CreatedAt.AsTime())
				}
			},
		},
		{
			name:    "ReservationNotFound",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetReservationRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: codes.NotFound,
			validate: func(t *testing.T, res *proto.GetReservationResponse, req *proto.GetReservationRequest) {
				if res != nil {
					t.Errorf("expected nil response on error, got %v", res)
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "MissingId",
			request:      &proto.GetReservationRequest{Id: 0},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.GetReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "DatabaseError",
			request: baseRequest,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.GetReservationRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.GetId()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.GetReservationResponse, req *proto.GetReservationRequest) {
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

			svc := main.ReservationServiceImpl{DB: db}
			res, err := svc.GetReservation(context.Background(), tt.request)

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