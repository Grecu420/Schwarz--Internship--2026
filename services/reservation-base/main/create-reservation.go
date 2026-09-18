package main

import (
	"context"
	"Schwarz--Internship--2026/services/reservation-base/main/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *ReservationServiceImpl) CreateReservation(ctx context.Context, req *proto.CreateReservationRequest) (*proto.CreateReservationResponse, error) {
	if req == nil || req.Reservation == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request or missing reservation data")
	}

	res := req.GetReservation()

	if res.GetPropertyId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "property_id is missing")
	}
	if res.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is missing")
	}
	if res.GetOwnerId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner_id is missing")
	}
	if res.GetCheckInDate() == "" {
		return nil, status.Error(codes.InvalidArgument, "check_in_date is missing")
	}
	if res.GetCheckOutDate() == "" {
		return nil, status.Error(codes.InvalidArgument, "check_out_date is missing")
	}

	res.Status = proto.ReservationStatus_RESERVATION_STATUS_PENDING

	id, err := InsertReservation(ctx, service.DB, res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert reservation into database: %v", err)
	}

	res.Id = id
	now := timestamppb.Now()
	res.CreatedAt = now
	res.UpdatedAt = now

	return &proto.CreateReservationResponse{
		Reservation: res,
	}, nil
}