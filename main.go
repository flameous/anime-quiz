package main

import (
	"github.com/flameous/anime-quiz/server"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	s := server.NewServer()
	s.Start(":8080")
}
