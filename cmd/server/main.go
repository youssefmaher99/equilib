package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type server struct {
	listenAddr string
}

func New(listenAddr string) *server {
	return &server{listenAddr: listenAddr}
}

func (s *server) pongHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := "pong from " + s.listenAddr + "\n"
		w.Write([]byte(msg))
	})
}

func (s *server) Start() error {
	log.Printf("Server is running on [%s]\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.pongHandler())
}

func spawnServer(address string) {
	s := New(address)
	s.Start()
}

//TODO graceful shutdown

func main() {
	num_of_servers := flag.Int("servers", 10, "number of servers to spawn on localhost with port starting from 5001")
	flag.Parse()
	wg := sync.WaitGroup{}
	wg.Add(*num_of_servers)
	start_port := 5000
	for i := 0; i < *num_of_servers; i++ {
		go func(port int) {
			spawnServer(fmt.Sprintf("127.0.0.1:%d", port))
			wg.Done()
		}(start_port)
		start_port++
	}
	wg.Wait()
}
