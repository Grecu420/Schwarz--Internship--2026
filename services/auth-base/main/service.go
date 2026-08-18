package main

import (
	"Schwarz--Internship--2026/services/auth-base/main/proto"
	"Schwarz--Internship--2026/services/common"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultPort = 50053

type AuthServiceImpl struct {
	proto.UnimplementedAuthServiceServer
	UserService proto.UserServiceClient
}

func main() {

	// Listen to port
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

	// Connect to user-base endpoint
	user_base_endpoint, err := common.GetRequiredEnv("USER-BASE_ENDPOINT")
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)

	}
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	connUB, err := grpc.NewClient(user_base_endpoint, dialOpts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer connUB.Close()

	// Create grpc server
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	proto.RegisterAuthServiceServer(grpcServer, AuthServiceImpl{UserService: proto.NewUserServiceClient(connUB)})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
