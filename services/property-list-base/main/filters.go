package main

import (
	"Schwarz--Internship--2026/services/property-list-base/main/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type parsedPropertyFilters struct {
	owner    *proto.FilterByOwner
	price    *proto.FilterByPriceRange
	name     *proto.FilterByName
	location *proto.FilterByLocation
	date     *proto.FilterByDateInterval
}

func parsePropertyFilters(filters []*proto.ListPropertiesFiltersOneOf) (*parsedPropertyFilters, error) {
	pf := &parsedPropertyFilters{}

	for _, f := range filters {
		switch v := f.GetFilter().(type) {
		case *proto.ListPropertiesFiltersOneOf_Owner:
			if pf.owner != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate owner filter")
			}
			if v.Owner.GetValue() <= 0 {
				return nil, status.Errorf(codes.InvalidArgument, "invalid owner_id filter")
			}
			pf.owner = v.Owner

		case *proto.ListPropertiesFiltersOneOf_Name:
			if pf.name != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate name filter")
			}
			if v.Name.GetValue() == "" {
				return nil, status.Errorf(codes.InvalidArgument, "missing name in filter")
			}
			pf.name = v.Name

		case *proto.ListPropertiesFiltersOneOf_Location:
			if pf.location != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate location filter")
			}
			if v.Location.GetRadius() <= 0 {
				return nil, status.Errorf(codes.InvalidArgument, "negative radius")
			}
			if v.Location.GetCenter() == nil {
				return nil, status.Errorf(codes.InvalidArgument, "missing center")
			}
			pf.location = v.Location

		case *proto.ListPropertiesFiltersOneOf_PriceRange:
			if pf.price != nil {
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
			pf.price = v.PriceRange

		case *proto.ListPropertiesFiltersOneOf_DateInterval:
			if pf.date != nil {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate dateInterval filter")
			}
			pf.date = v.DateInterval
		}
	}

	return pf, nil
}
