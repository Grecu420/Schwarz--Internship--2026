package main

import (
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *PropertyServiceImpl) CountProperties(ctx context.Context, req *proto.CountPropertiesRequest) (*proto.CountPropertiesResponse, error) {

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

	count, err := CountPropertiesInDB(ctx, service.DB, ownerFilter, nameFilter, priceFilter, locationFilter)

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to get property count from database: %v", err)

	}

	return &proto.CountPropertiesResponse{Count: count}, nil
}
