package common

import (
	"fmt"
	"os"
	"strconv"
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
