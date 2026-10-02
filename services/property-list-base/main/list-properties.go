package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/common/pagination"
	"Schwarz--Internship--2026/services/property-list-base/main/proto"
)

func (service *PropertyListServiceImpl) ListProperties(ctx context.Context, req *proto.ListPropertiesRequest) (*proto.ListPropertiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	filters, err := parsePropertyFilters(req.GetFilters())
	if err != nil {
		return nil, err
	}

	filterHash := pagination.HashFilters(
		filters.owner.GetValue(),
		filters.price.GetMin(),
		filters.price.GetMax(),
		filters.name.GetValue(),
		filters.location.GetCenter().GetLat(),
		filters.location.GetCenter().GetLong(),
		filters.location.GetRadius(),
	)

	properties, nextPageToken, err := pagination.Paginate(
		req.GetPageSize(),
		req.GetNextPageToken(),
		filterHash,
		func(offsetId int64, limit int64) ([]*proto.Property, error) {
			return SelectPropertyListInDB(ctx, service.DB, offsetId, limit, req.GetSortType(), filters.owner, filters.name, filters.price, filters.location, filters.date)
		})
	if err != nil {
		return nil, err
	}

	return &proto.ListPropertiesResponse{
		NextPageToken: nextPageToken,
		Properties:    properties,
	}, nil
}
