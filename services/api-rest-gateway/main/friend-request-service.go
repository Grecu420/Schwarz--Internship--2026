package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) CreateFriendRequest(ctx context.Context, req *proto.CreateFriendRequestRequest) (*proto.CreateFriendRequestResponse, error) {
	return service.friendRequestService.CreateFriendRequest(ctx, req)
}

func (service GatewayServiceImpl) UpdateFriendRequest(ctx context.Context, req *proto.UpdateFriendRequestRequest) (*proto.UpdateFriendRequestResponse, error) {
	return service.friendRequestService.UpdateFriendRequest(ctx, req)
}

func (service GatewayServiceImpl) ListFriendRequests(ctx context.Context, req *proto.ListFriendRequestsRequest) (*proto.ListFriendRequestsResponse, error) {
	return service.friendRequestService.ListFriendRequests(ctx, req)
}
