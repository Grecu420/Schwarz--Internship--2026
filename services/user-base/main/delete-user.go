package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service UserServiceImpl) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	err := DeleteUserInDB(ctx, service.DB, req.GetId())

	if errors.Is(err, NoRowsDeleted) {
		return nil, status.Error(codes.NotFound, "user not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}
	
	return &proto.DeleteUserResponse{}, nil
}