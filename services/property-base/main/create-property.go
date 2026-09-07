package main

import (
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service PropertyServiceImpl) CreateProperty(ctx context.Context, req *proto.CreatePropertyRequest) (*proto.CreatePropertyResponse, error) {
	// TODO: add check to see if user owns property

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner id missing")
	}
	if req.GetPrice() == 0 {
		return nil, status.Error(codes.InvalidArgument, "property price missing")
	}
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "property name missing")
	}

	id, err := InsertPropertyInDB(ctx, service.DB, req.GetName(), req.GetDescription(), req.GetUserId(), req.GetAddress(), int(req.GetPrice()), float64(req.Location.Long), float64(req.Location.Lat))

	if err != nil {

		return nil, status.Errorf(codes.Internal, "failed insert property to database: %v", err)
	}

	return &proto.CreatePropertyResponse{Property: &proto.Property{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Price:       req.Price,
		Location:    req.Location,
		UserId:      req.UserId,
	}}, nil
}
