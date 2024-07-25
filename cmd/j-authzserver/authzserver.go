package main

import (
	"math/rand"
	"time"

	"j-iam/internal/authzserver"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	authzserver.NewApp("j-authz-server").Run()
}
