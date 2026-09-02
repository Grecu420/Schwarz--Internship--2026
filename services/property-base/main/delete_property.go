package main

import (
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service PropertyServiceImpl) DeleteProperty(ctx context.Context, req *proto.DeletePropertyRequest) (*proto.DeletePropertyResponse, error) {
	// TODO: add check to see if user owns property

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	err := DeletePropertyInDB(ctx, service.DB, req.GetId())

	if errors.Is(err, NoRowsDeleted) {
		return nil, status.Error(codes.NotFound, "property not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}
	return &proto.DeletePropertyResponse{}, nil
}
