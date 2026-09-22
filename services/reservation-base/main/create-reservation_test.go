package main_test

import (
	"context"
	"errors"
	"testing"

	"Schwarz--Internship--2026/services/reservation-base/main"
	"Schwarz--Internship--2026/services/reservation-base/main/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateReservation(t *testing.T) {
	expectedID := int64(42)

	baseReservation := func() *proto.Reservation {
		return &proto.Reservation{
			PropertyId:   100,
			UserId:       200,
			OwnerId:      300,
			CheckInDate:  "2026-08-15",
			CheckOutDate: "2026-08-20",
		}
	}

	tests := []struct {
		name         string
		getRes       func() *proto.Reservation
		setupMock    func(mock sqlmock.Sqlmock, r *proto.Reservation)
		expectedCode codes.Code
		validate     func(t *testing.T, res *proto.CreateReservationResponse, reqRes *proto.Reservation)
	}{
		{
			name:   "Success",
			getRes: baseReservation,
			setupMock: func(mock sqlmock.Sqlmock, r *proto.Reservation) {
				mock.ExpectQuery(`INSERT INTO reservations`).
					WithArgs(
						r.PropertyId,
						r.UserId,
						r.OwnerId,
						r.CheckInDate,
						r.CheckOutDate,
						proto.ReservationStatus_RESERVATION_STATUS_PENDING.String(),
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, res *proto.CreateReservationResponse, reqRes *proto.Reservation) {
				if res.Reservation.Id != expectedID {
					t.Errorf("expected ID %d, got %d", expectedID, res.Reservation.Id)
				}
				if res.Reservation.Status != proto.ReservationStatus_RESERVATION_STATUS_PENDING {
					t.Errorf("expected status PENDING, got %v", res.Reservation.Status)
				}
				if res.Reservation.CreatedAt == nil {
					t.Error("expected CreatedAt to be populated, got nil")
				}
				if res.Reservation.UpdatedAt == nil {
					t.Error("expected UpdatedAt to be populated, got nil")
				}
			},
		},
		{
			name: "Missing PropertyId",
			getRes: func() *proto.Reservation {
				r := baseReservation()
				r.PropertyId = 0
				return r
			},
			setupMock:    func(mock sqlmock.Sqlmock, r *proto.Reservation) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Missing UserId",
			getRes: func() *proto.Reservation {
				r := baseReservation()
				r.UserId = 0
				return r
			},
			setupMock:    func(mock sqlmock.Sqlmock, r *proto.Reservation) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Missing OwnerId",
			getRes: func() *proto.Reservation {
				r := baseReservation()
				r.OwnerId = 0
				return r
			},
			setupMock:    func(mock sqlmock.Sqlmock, r *proto.Reservation) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Missing CheckInDate",
			getRes: func() *proto.Reservation {
				r := baseReservation()
				r.CheckInDate = ""
				return r
			},
			setupMock:    func(mock sqlmock.Sqlmock, r *proto.Reservation) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "Missing CheckOutDate",
			getRes: func() *proto.Reservation {
				r := baseReservation()
				r.CheckOutDate = ""
				return r
			},
			setupMock:    func(mock sqlmock.Sqlmock, r *proto.Reservation) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "Nil Reservation",
			getRes:       func() *proto.Reservation { return nil },
			setupMock:    func(mock sqlmock.Sqlmock, r *proto.Reservation) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:   "Database Error",
			getRes: baseReservation,
			setupMock: func(mock sqlmock.Sqlmock, r *proto.Reservation) {
				mock.ExpectQuery(`INSERT INTO reservations`).
					WithArgs(
						r.PropertyId,
						r.UserId,
						r.OwnerId,
						r.CheckInDate,
						r.CheckOutDate,
						proto.ReservationStatus_RESERVATION_STATUS_PENDING.String(),
					).
					WillReturnError(errors.New("connection failed"))
			},
			expectedCode: codes.Internal,
			validate: func(t *testing.T, res *proto.CreateReservationResponse, reqRes *proto.Reservation) {
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

			reqRes := tt.getRes()
			if tt.setupMock != nil && reqRes != nil {
				tt.setupMock(mock, reqRes)
			}

			svc := main.ReservationServiceImpl{DB: db}
			
			var req *proto.CreateReservationRequest
			if reqRes != nil {
				req = &proto.CreateReservationRequest{Reservation: reqRes}
			} else {
				req = &proto.CreateReservationRequest{Reservation: nil}
			}

			res, err := svc.CreateReservation(context.Background(), req)

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
				tt.validate(t, res, reqRes)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}