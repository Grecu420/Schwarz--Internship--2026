package main

import (
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenPayload struct {
	ID         string `json:"id"`
	FilterHash string `json:"hash"`
}

func (f *FriendRequestServiceImpl) ListFriendRequests(ctx context.Context, req *proto.ListFriendRequestsRequest) (*proto.ListFriendRequestsResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	pageSize := req.GetPageSize()

	if pageSize <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page size: %v", pageSize)
	}

	return nil, status.Error(codes.Unimplemented, "method ListFriendRequests not implemented")
}
