package main

import (
	"Schwarz--Internship--2026/services/property-list-base/main/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *PropertyListServiceImpl) CountProperties(ctx context.Context, req *proto.CountPropertiesRequest) (*proto.CountPropertiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	filters, err := parsePropertyFilters(req.GetFilters())
	if err != nil {
		return nil, err
	}

	count, err := CountPropertiesInDB(ctx, service.DB, filters.owner, filters.name, filters.price, filters.location, filters.date)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to get property count from database: %v", err)
	}

	return &proto.CountPropertiesResponse{Count: count}, nil
}
