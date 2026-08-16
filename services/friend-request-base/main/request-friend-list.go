package main

import (
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenPayload struct {
	ID         int64  `json:"id"`
	FilterHash string `json:"hash"`
}

// get RequestStatus from string
func getStatus(status string) proto.RequestStatus {
	switch strings.ToUpper(status) {

	case "STATUS_PENDING":
		return proto.RequestStatus_STATUS_PENDING
	case "STATUS_ACCEPTED":
		return proto.RequestStatus_STATUS_ACCEPTED
	case "STATUS_REJECTED":
		return proto.RequestStatus_STATUS_REJECTED
	default:
		return proto.RequestStatus_STATUS_UNKNOWN
	}
}

func hashFilter(senderId, receiverId string, status proto.RequestStatus) string {
	h := fnv.New32a()
	h.Write([]byte(senderId))
	h.Write([]byte(receiverId))
	h.Write([]byte(status.String()))
	return fmt.Sprintf("%x", h.Sum32())
}

func getPayload(nextPageToken string) (*TokenPayload, error) {
	decoded, err := base64.StdEncoding.DecodeString(nextPageToken)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 page token: %w", err)
	}

	var payload TokenPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, fmt.Errorf("malformed page token: %w", err)
	}
	return &payload, nil

}

func buildNextPageToken(id int64, filterHash string) (string, error) {

	newToken := TokenPayload{
		ID:         id,
		FilterHash: filterHash,
	}
	tokenBytes, err := json.Marshal(newToken)
	if err != nil {
		return "", err
	}
	nextPageToken := base64.StdEncoding.EncodeToString(tokenBytes)
	return nextPageToken, nil
}

func (f *FriendRequestServiceImpl) ListFriendRequests(ctx context.Context, req *proto.ListFriendRequestsRequest) (*proto.ListFriendRequestsResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	// 1. Get page size
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page size: %v", pageSize)
	}

	// 2. Get filters
	// Currently, filters later in the list override earlier ones
	filters := req.GetFilters()
	var statusFilter proto.RequestStatus
	var senderIdFilter string
	var receiverIdFilter string
	for _, f := range filters {
		switch f.GetFilter().(type) {
		case *proto.ListFriendRequestsFiltersOneOf_ReceiverId:
			receiverIdFilter = f.GetReceiverId()
		case *proto.ListFriendRequestsFiltersOneOf_SenderId:
			senderIdFilter = f.GetSenderId()
		case *proto.ListFriendRequestsFiltersOneOf_Status:
			statusFilter = f.GetStatus()
		}
	}

	// 3. Compute filter hash
	// This is to prevent filter from being changed after the first page is sent
	filterHash := hashFilter(senderIdFilter, receiverIdFilter, statusFilter)

	// 4. Extract payload from nextPageToken
	var offsetId int64
	if req.GetNextPageToken() != "" {
		payload, err := getPayload(req.GetNextPageToken())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "next page token failure: %v", err)
		}
		if payload.FilterHash != filterHash {
			return nil, status.Error(codes.InvalidArgument, "filter criteria changed during pagination")
		}
		offsetId = payload.ID
	}

	// 5. Request list of friend requests from database
	friendRequests, err := SelectFriendListInDB(ctx, f.DB,
		offsetId, pageSize,
		senderIdFilter,
		receiverIdFilter,
		statusFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed db query: %v", err)
	}

	// 6. Build nextPageToken
	// If there are items in the next page, generate a token
	// Otherwise, leave it empty
	var nextPageToken = ""
	if len(friendRequests) == int(pageSize)+1 {
		last := friendRequests[len(friendRequests)-1]
		newOffsetId := last.GetId()
		nextPageToken, err = buildNextPageToken(newOffsetId, filterHash)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed nextPageToken creation: %v", err)
		}
	}

	return &proto.ListFriendRequestsResponse{NextPageToken: nextPageToken, Requests: friendRequests[:pageSize]}, nil
}
