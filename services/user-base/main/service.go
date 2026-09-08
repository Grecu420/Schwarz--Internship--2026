package main

import (
	"Schwarz--Internship--2026/services/common"
	"Schwarz--Internship--2026/services/common/database"
	"Schwarz--Internship--2026/services/common/rabbitmq"
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
)

const defaultPort = 50051

type UserServiceImpl struct {
	proto.UnimplementedUserServiceServer
	DB        *sql.DB
	EmailProd *rabbitmq.Producer
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

	// connect to email queue
	rabbitmq_url, err := common.GetRequiredEnv("RABBITMQ_URL")
	if err != nil {
		log.Fatalf("Failed to register url: %v", err)
	}
	ctx := context.Background()

	prod, err := rabbitmq.NewProducer(ctx, rabbitmq_url, "/queues/email_queue")
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer prod.Close(ctx)

	// crete grpc server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer, &UserServiceImpl{DB: db, EmailProd: prod})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
