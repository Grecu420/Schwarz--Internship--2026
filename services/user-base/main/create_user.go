package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"log"
	"time"

	pbf "google.golang.org/protobuf/proto"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u UserServiceImpl) CreateUser(c context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

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

	// set current time
	var createdAt time.Time
	if user.CreatedAt == nil {
		user.CreatedAt = timestamppb.Now()
	}
	createdAt = user.CreatedAt.AsTime()

	// execute query + get new id
	query := `
		INSERT INTO users (id, first_name, last_name, user_name, email, hashed_password, created_at)
		VALUES (DEFAULT, $1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int64
	err = u.DB.QueryRow(query,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Email,
		string(hashedPassword),
		createdAt).Scan(&id)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to save user to database: %v", err)
	}

	// update user for return
	user.Id = id
	user.Password = string(hashedPassword)

	return &proto.CreateUserResponse{User: user}, nil
}
