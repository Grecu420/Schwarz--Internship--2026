package main

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
)

func (f *FriendRequestServiceImpl) FriendRequestEndpoint(ctx context.Context, req *proto.CreateFriendRequestRequest) (*proto.CreateFriendRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	senderID := req.GetSenderId()
	receiverID := req.GetReceiverId()

	if senderID == "" || receiverID == "" {
		return nil, status.Error(codes.InvalidArgument, "sender_id and receiver_id are mandatory")
	}

	if senderID == receiverID {
		return nil, status.Error(codes.InvalidArgument, "you cannot send a request to yourself")
	}

	var reqStatus int32 = 1
	var generatedID int64

	query := `
		INSERT INTO friend_requests (sender_id, receiver_id, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := f.DB.QueryRowContext(ctx, query, senderID, receiverID, reqStatus).Scan(&generatedID)
	if err != nil {
		log.Printf("DB error when inserting the request (sender: %s, receiver: %s): %v", senderID, receiverID, err)
		return nil, status.Error(codes.Internal, "an internal error occurred; please try again")
	}

	return &proto.CreateFriendRequestResponse{
		Request: &proto.FriendRequest{
			Id:         generatedID,
			SenderId:   senderID,
			ReceiverId: receiverID,
			Status:     proto.RequestStatus(reqStatus),
		},
	}, nil
}
