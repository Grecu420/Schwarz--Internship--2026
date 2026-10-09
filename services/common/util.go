package common

import (
	"fmt"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetPort() (int, error) {
	s, err := GetRequiredEnv("Port")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)

}

func GetRequiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return val, nil
}

func DialGRPCService(envVar string) (*grpc.ClientConn, error) {
	endpoint, err := GetRequiredEnv(envVar)
	if err != nil {
		return nil, fmt.Errorf("missing environment variable %s: %w", envVar, err)
	}
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s at %s: %w", envVar, endpoint, err)
	}
	return conn, nil
}
