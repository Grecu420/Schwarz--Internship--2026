package main

import (
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service PropertyServiceImpl) GetProperty(ctx context.Context, req *proto.GetPropertyRequest) (*proto.GetPropertyResponse, error) {

	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing id")
	}
	id := req.GetId()
	property, err := SelectPropertyInDB(ctx, service.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Missing property (id: %d): %v", id, err)
		return nil, status.Errorf(codes.NotFound,
			"failed to find property (id: %d) ", id)
	} else if err != nil {
		log.Printf("Failed to get property (id: %d): %v", id, err)
		return nil, status.Errorf(codes.Internal,
			"failed to get property (id: %d) from database: %v", id, err)
	}

	return &proto.GetPropertyResponse{Property: property}, nil
}
