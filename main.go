package main

import "github.com/AtaskTracker/AtaskAPI/server"

func main() {
	server := server.NewServer()
	server.Start(":5000")
}
