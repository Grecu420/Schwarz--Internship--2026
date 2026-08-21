package main

import (
	"Schwarz--Internship--2026/services/common"
	"Schwarz--Internship--2026/services/common/database"
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"database/sql"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
)

const defaultPort = 50051

type UserServiceImpl struct {
	proto.UnimplementedUserServiceServer
	DB *sql.DB
}

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
	proto.RegisterUserServiceServer(grpcServer, &UserServiceImpl{DB: db})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
