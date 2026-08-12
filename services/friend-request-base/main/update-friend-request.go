package main

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
)

func (f *FriendRequestServiceImpl) UpdateFriendRequest(ctx context.Context, req *proto.UpdateFriendRequestRequest) (*proto.UpdateFriendRequestResponse, error) {

	if req == nil || req.FriendRequest == nil {
		return nil, status.Error(codes.InvalidArgument, "request and friend_request are mandatory")
	}

	if req.FieldMask == nil || len(req.FieldMask.GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "field_mask is mandatory for updates")
	}

	requestID := req.FriendRequest.GetId()
	if requestID == 0 {
		return nil, status.Error(codes.InvalidArgument, "friend request id is mandatory")
	}

	for _, path := range req.FieldMask.GetPaths() {
		if path != "status" {
			return nil, status.Error(codes.InvalidArgument, "permission denied: only the 'status' field can be updated")
		}
	}

	newStatus := req.FriendRequest.GetStatus()

	rowsAffected, err := UpdateFriendRequestStatusInDB(ctx, f.DB, requestID, newStatus)
	if err != nil {
		log.Printf("DB error when updating request (id: %d): %v", requestID, err)
		return nil, status.Error(codes.Internal, "an internal error occurred while updating the database")
	}

	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "friend request not found")
	}

	return &proto.UpdateFriendRequestResponse{
		FriendRequest: req.FriendRequest,
	}, nil
}
