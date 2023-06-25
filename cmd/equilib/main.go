package main

import (
	"log"

	"github.com/youssefmaher99/equilib/internal/server"
)

func main() {
	servers_list := []string{"127.0.0.1:5001", "127.0.0.1:5002"}
	s := server.New("127.0.0.1:8080", 2, servers_list)
	log.Fatal(s.Start())
}
