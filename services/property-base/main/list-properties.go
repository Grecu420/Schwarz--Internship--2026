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

	var ownerFilter *proto.FilterByOwner
	var priceFilter *proto.FilterByPriceRange
	var nameFilter *proto.FilterByName
	var locationFilter *proto.FilterByLocation

	filters := req.GetFilters()
	for _, f := range filters {
		switch v := f.GetFilter().(type) {
		case *proto.ListPropertiesFiltersOneOf_Owner:
			if ownerFilter != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate owner filter")

			}
			if v.Owner.GetValue() <= 0 {
				return nil, status.Errorf(codes.InvalidArgument, "invalid owner_id filter")

			}
			ownerFilter = v.Owner
		case *proto.ListPropertiesFiltersOneOf_Name:
			if nameFilter != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate name filter")
			}
			if v.Name.GetValue() == "" {
				return nil, status.Errorf(codes.InvalidArgument, "missing name in filter")
			}
			nameFilter = v.Name
		case *proto.ListPropertiesFiltersOneOf_Location:
			if locationFilter != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate location filter")
			}
			if v.Location.GetRadius() <= 0 {
				return nil, status.Errorf(codes.InvalidArgument, "negative radius")
			}
			if v.Location.GetCenter() == nil {
				return nil, status.Errorf(codes.InvalidArgument, "missing center")
			}
			locationFilter = v.Location
		case *proto.ListPropertiesFiltersOneOf_PriceRange:
			if priceFilter != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate priceRange filter")
			}
			if v.PriceRange.GetMin() < 0 {
				return nil, status.Errorf(codes.InvalidArgument, "invalid priceRange filter: negative min")

			}
			if v.PriceRange.GetMax() < 0 {
				return nil, status.Errorf(codes.InvalidArgument, "invalid priceRange filter: negative max")
			}
			if v.PriceRange.GetMin() > v.PriceRange.GetMax() && v.PriceRange.GetMax() > 0 {
				return nil, status.Errorf(codes.InvalidArgument, "invalid priceRange filter: min > max")
			}

			priceFilter = v.PriceRange
		}
	}

	filterHash := pagination.HashFilters(ownerFilter.GetValue(), priceFilter.GetMin(), priceFilter.GetMax(), nameFilter.GetValue(), locationFilter.GetCenter().GetLat(), locationFilter.GetCenter().GetLong(), locationFilter.GetRadius())

	properties, nextPageToken, err := pagination.Paginate(
		req.GetPageSize(),
		req.GetNextPageToken(),
		filterHash,
		func(offsetId int64, limit int64) ([]*proto.Property, error) {
			return SelectPropertyListInDB(ctx, service.DB, offsetId, limit, ownerFilter, nameFilter, priceFilter, locationFilter)
		})
	if err != nil {
		return nil, err
	}

	return &proto.ListPropertiesResponse{
		NextPageToken: nextPageToken,
		Properties:    properties,
	}, nil
}
