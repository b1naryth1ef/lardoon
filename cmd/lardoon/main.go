package main

import (
	"github.com/b1naryth1ef/lardoon"
)

func main() {
	var server lardoon.HTTPServer
	server.Run("localhost:3883")
}
