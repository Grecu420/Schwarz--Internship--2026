package main

import (
    "Schwarz--Internship--2026/services/user-base/main/proto"
    "context"
    "log"
    "strings"

    pbf "google.golang.org/protobuf/proto"

    "golang.org/x/crypto/bcrypt"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func (u UserServiceImpl) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

    if req.GetUser() == nil {
        log.Printf("Empty request")
        return nil, status.Error(codes.InvalidArgument, "empty request")
    }
    user := pbf.Clone(req.User).(*proto.User)

    // hash password
    password := user.Password
    bytePassword := []byte(password)

    if len(bytePassword) > 72 {
        // truncate password to 72 bytes
        bytePassword = bytePassword[:72]
    }

    hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return nil, status.Error(codes.Internal, "failed to process password")
    }
    user.Password = string(hashedPassword)

    // set current time
    if user.CreatedAt == nil {
        user.CreatedAt = timestamppb.Now()
    }

	id, err := InsertUser(ctx, u.DB, user)
    if err != nil {
        log.Printf("Failed to insert user: %v", err)
        
        errMsg := err.Error()
        if strings.Contains(errMsg, "duplicate key value") || strings.Contains(errMsg, "unique constraint") {
            if strings.Contains(errMsg, "email") {
                return nil, status.Error(codes.AlreadyExists, "This email address is already in use")
            }
            if strings.Contains(errMsg, "username") || strings.Contains(errMsg, "user_name") {
                return nil, status.Error(codes.AlreadyExists, "This username is already taken")
            }
        }
        
        return nil, status.Errorf(codes.Internal, "failed to save user to database: %v", err)
    }

    // update user for return
    user.Id = id

    return &proto.CreateUserResponse{User: user}, nil
}