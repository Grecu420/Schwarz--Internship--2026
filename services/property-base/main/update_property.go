package main

import (
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service PropertyServiceImpl) UpdateProperty(ctx context.Context, req *proto.UpdatePropertyRequest) (*proto.UpdatePropertyResponse, error) {

	if req == nil || req.Property == nil {
		return nil, status.Error(codes.InvalidArgument, "request and property are mandatory")
	}

	if req.FieldMask == nil || len(req.FieldMask.GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "field_mask is mandatory for updates")
	}

	requestID := req.Property.GetId()
	if requestID == 0 {
		return nil, status.Error(codes.InvalidArgument, "property id is mandatory")
	}

	build := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("properties").
		Where(sq.Eq{"id": requestID})

	hasValidFields := false
	for _, path := range req.FieldMask.GetPaths() {
		switch path {
		case "name":
			build = build.Set("name", req.Property.GetName())
			hasValidFields = true
		case "description":
			build = build.Set("description", req.Property.GetDescription())
			hasValidFields = true
		case "address":
			build = build.Set("address", req.Property.GetAddress())
			hasValidFields = true
		case "price":
			build = build.Set("price", req.Property.GetPrice())
			hasValidFields = true
		case "location":
			location := req.GetProperty().GetLocation()
			build = build.Set("location", sq.Expr("ST_MakePoint(?, ?)::geography", location.GetLong(), location.GetLat()))
			hasValidFields = true

		default:
			return nil, status.Errorf(codes.InvalidArgument, "invalid field path in field_mask: %s", path)
		}
	}

	if !hasValidFields {
		return nil, status.Error(codes.InvalidArgument, "no valid fields to update")
	}

	build = build.Suffix("RETURNING id, user_id, name, description, address, price, ST_X(location::geometry) AS lng, ST_Y(location::geometry) AS lat")

	var returnProperty = proto.Property{Location: &proto.Location{}}
	err := build.RunWith(service.DB).QueryRowContext(ctx).Scan(
		&returnProperty.Id,
		&returnProperty.UserId,
		&returnProperty.Name,
		&returnProperty.Description,
		&returnProperty.Address,
		&returnProperty.Price,
		&returnProperty.Location.Long,
		&returnProperty.Location.Lat,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "property to update not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}

	return &proto.UpdatePropertyResponse{Property: &returnProperty}, nil
}
