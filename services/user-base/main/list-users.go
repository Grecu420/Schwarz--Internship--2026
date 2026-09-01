package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"Schwarz--Internship--2026/services/user-base/main/proto"
)

type TokenPayload struct {
	ID         int64  `json:"id"`
	FilterHash string `json:"hash"`
}

func HashFilter(firstName, lastName string) string {
	h := fnv.New32a()
	h.Write([]byte(firstName))
	h.Write([]byte(lastName))
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

func BuildNextPageToken(id int64, filterHash string) (string, error) {
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

func (s *UserServiceImpl) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page size: %v", pageSize)
	}

	filters := req.GetFilters()
	var firstNameFilter string
	var lastNameFilter string

	for _, f := range filters {
		switch v := f.GetFilter().(type) {
		case *proto.ListUsersFiltersOneOf_FirstName:
			if firstNameFilter != "" {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate first_name filter")
			}
			firstNameFilter = v.FirstName.GetValue()
		case *proto.ListUsersFiltersOneOf_LastName:
			if lastNameFilter != "" {
				return nil, status.Errorf(codes.InvalidArgument, "duplicate last_name filter")
			}
			lastNameFilter = v.LastName.GetValue()
		}
	}

	filterHash := HashFilter(firstNameFilter, lastNameFilter)

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

	users, err := SelectUserListInDB(ctx, s.DB, offsetId, pageSize, firstNameFilter, lastNameFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed db query: %v", err)
	}

	var nextPageToken = ""
	if len(users) == int(pageSize)+1 {
		last := users[pageSize-1]
		newOffsetId := last.GetId()

		var err error
		nextPageToken, err = BuildNextPageToken(newOffsetId, filterHash)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed nextPageToken creation: %v", err)
		}

		users = users[:pageSize]
	}

	return &proto.ListUsersResponse{NextPageToken: nextPageToken, Users: users}, nil
}
