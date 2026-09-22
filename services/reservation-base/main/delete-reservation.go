package main

import (
	"context"
	"errors"

	"Schwarz--Internship--2026/services/reservation-base/main/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *ReservationServiceImpl) DeleteReservation(ctx context.Context, req *proto.DeleteReservationRequest) (*proto.DeleteReservationResponse, error) {
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing reservation id")
	}

	err := DeleteReservationInDB(ctx, service.DB, req.GetId())

	if errors.Is(err, errNoRowsDeleted) {
		return nil, status.Error(codes.NotFound, "reservation not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}

	return &proto.DeleteReservationResponse{}, nil
}