package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
)

func (f *FriendRequestServiceImpl) FriendRequestEndpoint(ctx context.Context, req *proto.CreateFriendRequestRequest) (*proto.CreateFriendRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "The request cannot be nil")
	}

	senderID := req.GetSenderId()
	receiverID := req.GetReceiverId()

	if senderID == "" || receiverID == "" {
		return nil, status.Error(codes.InvalidArgument, "sender_id and receiver_id are mandatory")
	}

	if senderID == receiverID {
		return nil, status.Error(codes.InvalidArgument, "You cannot send a request to yourself")
	}

	requestID := uuid.New().String()
	var reqStatus int32 = 1

	err := InsertFriendRequestInDB(ctx, f.DB, requestID, senderID, receiverID, reqStatus)
	if err != nil {
		log.Printf("Database error when inserting the request (sender: %s, receiver: %s): %v", senderID, receiverID, err)
		return nil, status.Error(codes.Internal, "An internal error occurred; please try again")
	}

	return &proto.CreateFriendRequestResponse{
		Request: &proto.FriendRequest{
			Id:         requestID,
			SenderId:   senderID,
			ReceiverId: receiverID,
			Status:     proto.RequestStatus(reqStatus),
		},
	}, nil
}
