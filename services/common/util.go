package common

import (
	"fmt"
	"os"
	"strconv"
)

func GetPort() (int, error) {
	s, ok := os.LookupEnv("PORT")
	if !ok {
		return 0, fmt.Errorf("no port variable")
	}

	return strconv.Atoi(s)

}
