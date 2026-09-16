package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/common/pagination"
	"Schwarz--Internship--2026/services/property-base/main/proto"
)

func (service *PropertyServiceImpl) ListProperties(ctx context.Context, req *proto.ListPropertiesRequest) (*proto.ListPropertiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	var ownerID = req.GetOwnerId()
	if ownerID == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing owner id")

	}

	filterHash := pagination.HashFilters(ownerID)

	properties, nextPageToken, err := pagination.Paginate(
		req.GetPageSize(),
		req.GetNextPageToken(),
		filterHash,
		func(offsetId int64, limit int64) ([]*proto.Property, error) {
			return SelectPropertyListInDB(ctx, service.DB, offsetId, limit, ownerID)
		})
	if err != nil {
		return nil, err
	}

	return &proto.ListPropertiesResponse{
		NextPageToken: nextPageToken,
		Properties:    properties,
	}, nil
}
