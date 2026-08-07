package main

import (
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"

	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type FriendRequestServiceImpl struct {
	proto.UnimplementedFriendRequestServiceServer
	DB *sql.DB
}

func Connect() (*sql.DB, error) {
	passwordBytes, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, err
	}

	password := strings.TrimSpace(string(passwordBytes))
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	service := os.Getenv("POSTGRES_SERVICE")

	connStr := fmt.Sprintf("postgres://postgres:%s@%s:%s/%s?sslmode=disable", password, service, port, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("The database is not responding to pings: %w", err)
	}

	return db, nil
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

	dbConn, err := Connect()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer dbConn.Close()

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
	myService := &FriendRequestServiceImpl{
		DB: dbConn,
	}

	proto.RegisterFriendRequestServiceServer(grpcServer, myService)

	log.Printf("gRPC server successfully started on port %d...", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
