package main

import (
	"Schwarz--Internship--2026/services/common"
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceImpl struct {
	proto.UnimplementedUserServiceServer
}

// Ping implements [proto.UserServiceServer].
func (u UserServiceImpl) Ping(context.Context, *proto.Empty) (*proto.Pong, error) {
	fmt.Println("here")
	return &proto.Pong{Message: "pong "}, nil
}

func (u UserServiceImpl) EndpointName(context.Context, *proto.EndpointNameRequest) (*proto.EndpointNameResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method EndpointName not implemented")
}

const defaultPort = 50051

func main() {

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

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer, UserServiceImpl{})
	grpcServer.Serve(lis)

}
