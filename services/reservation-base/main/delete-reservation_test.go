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

func TestDeleteReservation(t *testing.T) {
	tests := []struct {
		name         string
		request      *proto.DeleteReservationRequest
		setupMock    func(mock sqlmock.Sqlmock, req *proto.DeleteReservationRequest)
		expectedCode codes.Code
	}{
		{
			name:    "Success",
			request: &proto.DeleteReservationRequest{Id: 10},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.DeleteReservationRequest) {
				mock.ExpectExec(`DELETE FROM reservations`).
					WithArgs(req.GetId()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedCode: codes.OK,
		},
		{
			name:         "NilRequest",
			request:      nil,
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.DeleteReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "MissingId",
			request:      &proto.DeleteReservationRequest{Id: 0},
			setupMock:    func(mock sqlmock.Sqlmock, req *proto.DeleteReservationRequest) {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:    "ReservationNotFound",
			request: &proto.DeleteReservationRequest{Id: 10},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.DeleteReservationRequest) {
				mock.ExpectExec(`DELETE FROM reservations`).
					WithArgs(req.GetId()).
					WillReturnResult(sqlmock.NewResult(0, 0)) 
			},
			expectedCode: codes.NotFound,
		},
		{
			name:    "DatabaseError",
			request: &proto.DeleteReservationRequest{Id: 10},
			setupMock: func(mock sqlmock.Sqlmock, req *proto.DeleteReservationRequest) {
				mock.ExpectExec(`DELETE FROM reservations`).
					WithArgs(req.GetId()).
					WillReturnError(errors.New("db connection lost"))
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

			svc := main.ReservationServiceImpl{DB: db}
			_, err = svc.DeleteReservation(context.Background(), tt.request)

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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled sqlmock expectations: %s", err)
			}
		})
	}
}