package main

import (
	"context"
	"fmt"
	"friend-request-base/main/proto"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

type FriendRequestServiceImpl struct {
	proto.UnimplementedFriendRequestServiceServer
}

// Ping implements [proto.UserServiceServer].
func (f FriendRequestServiceImpl) Ping(context.Context, *proto.Empty) (*proto.Pong, error) {
	fmt.Println("here")
	return &proto.Pong{Message: "pong "}, nil
}

const defaultPort = 50052

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
	lis, err := net.Listen("tcp", fmt.Sprintf("friend-request-base-service:%d", port))
	//lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterFriendRequestServiceServer(grpcServer, FriendRequestServiceImpl{})
	grpcServer.Serve(lis)

}
