package main

import (
	"Schwarz--Internship--2026/services/common"
	"Schwarz--Internship--2026/services/common/database"
	"Schwarz--Internship--2026/services/property-list-base/main/proto"
	"database/sql"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
)

const defaultPort = 50051

type PropertyListServiceImpl struct {
	proto.UnimplementedPropertyListServiceServer
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
	lis, err := net.Listen("tcp", fmt.Sprintf("property-list-base-service:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	// TODO: add check to verify that user_id in request matches user_id in context

	// create grpc server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterPropertyListServiceServer(grpcServer, &PropertyListServiceImpl{DB: db})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
