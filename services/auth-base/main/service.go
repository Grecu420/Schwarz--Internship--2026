package main

import (
	"Schwarz--Internship--2026/services/auth-base/main/proto"
	"Schwarz--Internship--2026/services/common"
	"Schwarz--Internship--2026/services/common/database"
	"database/sql"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
)

const defaultPort = 50053

type AuthServiceImpl struct {
	proto.UnimplementedAuthServiceServer
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
	lis, err := net.Listen("tcp", fmt.Sprintf("auth-base-service:%d", port))
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
	proto.RegisterAuthServiceServer(grpcServer, AuthServiceImpl{DB: db})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
