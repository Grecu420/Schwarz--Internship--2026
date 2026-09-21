package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"Schwarz--Internship--2026/services/reservation-base/main/proto"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *ReservationServiceImpl) UpdateReservation(ctx context.Context, req *proto.UpdateReservationRequest) (*proto.UpdateReservationResponse, error) {
	if req == nil || req.Reservation == nil {
		return nil, status.Error(codes.InvalidArgument, "request and reservation are mandatory")
	}

	if req.FieldMask == nil || len(req.FieldMask.GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "field_mask is mandatory for updates")
	}

	requestID := req.Reservation.GetId()
	if requestID == 0 {
		return nil, status.Error(codes.InvalidArgument, "reservation id is mandatory")
	}

	build := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("reservations").
		Where(sq.Eq{"id": requestID})

	hasValidFields := false

	for _, path := range req.FieldMask.GetPaths() {
		switch path {
		case "status":
			if req.Reservation.GetStatus() == proto.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED {
				return nil, status.Error(codes.InvalidArgument, "valid status is mandatory")
			}
			build = build.Set("status", req.Reservation.GetStatus().String())
			hasValidFields = true
		case "check_in_date":
			if req.Reservation.GetCheckInDate() == "" {
				return nil, status.Error(codes.InvalidArgument, "check_in_date is mandatory")
			}
			build = build.Set("check_in_date", req.Reservation.GetCheckInDate())
			hasValidFields = true
		case "check_out_date":
			if req.Reservation.GetCheckOutDate() == "" {
				return nil, status.Error(codes.InvalidArgument, "check_out_date is mandatory")
			}
			build = build.Set("check_out_date", req.Reservation.GetCheckOutDate())
			hasValidFields = true
		default:
			return nil, status.Errorf(codes.InvalidArgument, "invalid field path in field_mask: %s", path)
		}
	}

	if !hasValidFields {
		return nil, status.Error(codes.InvalidArgument, "no valid fields to update")
	}

	build = build.Suffix("RETURNING id, property_id, user_id, owner_id, check_in_date, check_out_date, status, created_at, updated_at")

	var returnRes proto.Reservation
	var checkIn, checkOut, createdAt, updatedAt time.Time
	var statusStr string

	err := build.RunWith(service.DB).QueryRowContext(ctx).Scan(
		&returnRes.Id,
		&returnRes.PropertyId,
		&returnRes.UserId,
		&returnRes.OwnerId,
		&checkIn,
		&checkOut,
		&statusStr,
		&createdAt,
		&updatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "reservation to update not found")
	}

	if err != nil {
		if strings.Contains(err.Error(), "chk_valid_dates") {
			return nil, status.Error(codes.InvalidArgument, "check_out_date must be after check_in_date")
		}
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}

	returnRes.CheckInDate = checkIn.Format("2006-01-02")
	returnRes.CheckOutDate = checkOut.Format("2006-01-02")
	returnRes.Status = proto.ReservationStatus(proto.ReservationStatus_value[statusStr])
	returnRes.CreatedAt = timestamppb.New(createdAt)
	returnRes.UpdatedAt = timestamppb.New(updatedAt)

	return &proto.UpdateReservationResponse{Reservation: &returnRes}, nil
}