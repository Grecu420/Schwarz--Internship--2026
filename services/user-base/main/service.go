package main

import (
	"fmt"
	"services/user-base/main/proto"
)

type UserServiceImpl struct {
}

func main() {

	var x proto.Empty

	fmt.Println(x.String())
}
