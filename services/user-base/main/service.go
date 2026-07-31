package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

type UserServiceImpl struct {
	proto.UnimplementedUserServiceServer
}

// Ping implements [proto.UserServiceServer].
func (u UserServiceImpl) Ping(context.Context, *proto.Empty) (*proto.Pong, error) {
	fmt.Println("here")
	return &proto.Pong{Message: "pong "}, nil
}

const defaultPort = 50051

func getPort() (int, error) {
	s, ok := os.LookupEnv("PORT")
	if !ok {
		return 0, fmt.Errorf("no port variable")
	}

	return strconv.Atoi(s)

}

func main() {

	var port int = defaultPort
	p, err := getPort()
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
