package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
)

func (service GatewayServiceImpl) getFriendRequestServiceClient() proto.FriendRequestServiceClient {
	return proto.NewFriendRequestServiceClient(service.friendRequestBaseConn)
}

func (service GatewayServiceImpl) CreateFriendRequestEndpoint(ctx context.Context, req *proto.CreateFriendRequestRequest) (*proto.CreateFriendRequestResponse, error) {
	return service.getFriendRequestServiceClient().FriendRequestEndpoint(ctx, req)
}
