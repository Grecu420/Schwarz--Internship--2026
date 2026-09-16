package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service UserServiceImpl) GetUserProfile(ctx context.Context, req *proto.GetUserProfileRequest) (*proto.GetUserProfileResponse, error) {
    // Get id
    id := req.GetId()
    if id <= 0 {
        return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
    }

    // Get user profile from database
    resProfile, err := SelectPublicUserProfileInDB(ctx, service.DB, id)
    if errors.Is(err, sql.ErrNoRows) {
        log.Printf("Missing user profile (id: %d): %v", id, err)
        return nil, status.Errorf(codes.NotFound,
            "failed to get user profile (id: %d) from database: %v", id, err)
    } else if err != nil {
        log.Printf("Failed to get user profile (id: %d): %v", id, err)
        return nil, status.Errorf(codes.Internal,
            "failed to get user profile (id: %d) from database: %v", id, err)
    }

    return &proto.GetUserProfileResponse{User: resProfile}, nil
}