package main

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"Schwarz--Internship--2026/services/reservation-base/main/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *ReservationServiceImpl) GetReservation(ctx context.Context, req *proto.GetReservationRequest) (*proto.GetReservationResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing reservation id")
	}

	resReservation, err := SelectReservation(ctx, service.DB, id)
	
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Missing reservation (id: %d): %v", id, err)
		return nil, status.Errorf(codes.NotFound,
			"failed to get reservation (id: %d) from database: %v", id, err)
	} else if err != nil {
		log.Printf("Failed to get reservation (id: %d): %v", id, err)
		return nil, status.Errorf(codes.Internal,
			"failed to get reservation (id: %d) from database: %v", id, err)
	}

	return &proto.GetReservationResponse{Reservation: resReservation}, nil
}