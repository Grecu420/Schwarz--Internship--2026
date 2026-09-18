package main_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"Schwarz--Internship--2026/services/reservation-base/main"
	"Schwarz--Internship--2026/services/reservation-base/main/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestUpdateReservation(t *testing.T) {
	customTime := time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)
	checkInTime := time.Date(2026, time.August, 15, 0, 0, 0, 0, time.UTC)
	checkOutTime := time.Date(2026, time.August, 20, 0, 0, 0, 0, time.UTC)

	baseReq := &proto.UpdateReservationRequest{
		Reservation: &proto.Reservation{
			Id:     10,
			Status: proto.ReservationStatus_RESERVATION_STATUS_CONFIRMED,
		},
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
	}

	expectedSQL := `UPDATE reservations SET status = \$1 WHERE id = \$2 RETURNING id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at`

	columns := []string{
		"id", "property_id", "user_id", "owner_id", "check_in_date", "check_out_date", "status", "created_at", "updated_at",
	}

	tests := []struct {
		name         string
		request      *proto.UpdateReservationRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.UpdateReservationResponse, req *proto.UpdateReservationRequest)
	}{
		{
			name:    "Success_UpdateStatus",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {
				rows := sqlmock.NewRows(columns).AddRow(
					10, 100, 200, 300, checkInTime, checkOutTime, "RESERVATION_STATUS_CONFIRMED", customTime, customTime,
				)

				mock.ExpectQuery(expectedSQL).
					WithArgs(req.Reservation.GetStatus().String(), req.Reservation.GetId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.UpdateReservationResponse, req *proto.UpdateReservationRequest) {
				if res == nil || res.Reservation == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if res.Reservation.Status != proto.ReservationStatus_RESERVATION_STATUS_CONFIRMED {
					t.Errorf("expected Status CONFIRMED, got %s", res.Reservation.Status)
				}
				if !res.Reservation.CreatedAt.AsTime().Equal(customTime) {
					t.Errorf("expected CreatedAt %v, got %v", customTime, res.Reservation.CreatedAt.AsTime())
				}
			},
		},
		{
			name: "Success_UpdateCheckInDate",
			request: &proto.UpdateReservationRequest{
				Reservation: &proto.Reservation{
					Id:          10,
					CheckInDate: "2026-08-16", 
				},
				FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"check_in_date"}},
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {
				expectedDateSQL := `UPDATE reservations SET check_in_date = \$1 WHERE id = \$2 RETURNING id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at`

				newCheckInTime := time.Date(2026, time.August, 16, 0, 0, 0, 0, time.UTC)
				rows := sqlmock.NewRows(columns).AddRow(
					10, 100, 200, 300, newCheckInTime, checkOutTime, "RESERVATION_STATUS_PENDING", customTime, customTime,
				)

				mock.ExpectQuery(expectedDateSQL).
					WithArgs(req.Reservation.GetCheckInDate(), req.Reservation.GetId()).
					WillReturnRows(rows)
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.UpdateReservationResponse, req *proto.UpdateReservationRequest) {
				if res == nil || res.Reservation == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
				if res.Reservation.CheckInDate != "2026-08-16" {
					t.Errorf("expected CheckInDate to be updated, got %s", res.Reservation.CheckInDate)
				}
			},
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingFieldMask",
			request: &proto.UpdateReservationRequest{
				Reservation: &proto.Reservation{Id: 10, Status: proto.ReservationStatus_RESERVATION_STATUS_CONFIRMED},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "MissingReservationId",
			request: &proto.UpdateReservationRequest{
				Reservation: &proto.Reservation{Id: 0, Status: proto.ReservationStatus_RESERVATION_STATUS_CONFIRMED},
				FieldMask:   &fieldmaskpb.FieldMask{Paths: []string{"status"}},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "InvalidFieldMaskPath",
			request: &proto.UpdateReservationRequest{
				Reservation: &proto.Reservation{Id: 10, Status: proto.ReservationStatus_RESERVATION_STATUS_CONFIRMED},
				FieldMask:   &fieldmaskpb.FieldMask{Paths: []string{"unknown_field"}},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "ProhibitedFieldMaskPathPropertyId",
			request: &proto.UpdateReservationRequest{
				Reservation: &proto.Reservation{Id: 10, PropertyId: 999},
				FieldMask:   &fieldmaskpb.FieldMask{Paths: []string{"property_id"}},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "EmptyStatusUpdateNotAllowed",
			request: &proto.UpdateReservationRequest{
				Reservation: &proto.Reservation{Id: 10, Status: proto.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED},
				FieldMask:   &fieldmaskpb.FieldMask{Paths: []string{"status"}},
			},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "ReservationNotFound",
			request: baseReq,
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {
				mock.ExpectQuery(expectedSQL).
					WithArgs(req.Reservation.GetStatus().String(), req.Reservation.GetId()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: codes.NotFound,
		},
		{
			name: "InvalidDatesConstraintError",
			request: &proto.UpdateReservationRequest{
				Reservation: &proto.Reservation{Id: 10, CheckOutDate: "2026-08-10"}, 
				FieldMask:   &fieldmaskpb.FieldMask{Paths: []string{"check_out_date"}},
			},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.UpdateReservationRequest) {
				expectedConstraintSQL := `UPDATE reservations SET check_out_date = \$1 WHERE id = \$2 RETURNING id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at`
				mock.ExpectQuery(expectedConstraintSQL).
					WithArgs(req.Reservation.GetCheckOutDate(), req.Reservation.GetId()).
					WillReturnError(errors.New("db error: chk_valid_dates violated constraint"))
			},
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

			svc := main.ReservationServiceImpl{DB: db}
			res, err := svc.UpdateReservation(context.Background(), tt.request)

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