package main

import (
	"Schwarz--Internship--2026/services/common"
	"Schwarz--Internship--2026/services/common/database"
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"fmt"
	"time"

	"log"
	"net"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServiceImpl struct {
	proto.UnimplementedUserServiceServer
	DB *sql.DB
}

// Ping implements [proto.UserServiceServer].
func (u UserServiceImpl) Ping(context.Context, *proto.Empty) (*proto.Pong, error) {
	fmt.Println("here")
	return &proto.Pong{Message: "pong "}, nil
}

func (u UserServiceImpl) CreateUser(c context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

	// hash password
	user := req.User
	if user == nil {
		log.Printf("Empty request")
		return nil, status.Error(codes.Internal, "empty request")

	}

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

const defaultPort = 50051

func main() {

	// listen to port
	var port int = defaultPort
	p, err := common.GetPort()
	if err == nil {
		port = p
	}
	fmt.Println("Port: ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("user-base-service:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	// crete grpc server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer, UserServiceImpl{DB: db})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
