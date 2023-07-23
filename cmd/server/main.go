package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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
		log.Printf("Incoming Request to [%s]", s.listenAddr)
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

func generateJsonFile(n int, starting_port int) {
	type fileStructure struct {
		Servers []string `json="servers"`
	}
	list_of_servers := fileStructure{}
	for i := 0; i < n; i++ {
		list_of_servers.Servers = append(list_of_servers.Servers, fmt.Sprintf("127.0.0.1:%d", starting_port))
		starting_port++
	}

	file, err := os.Create("servers.json")
	if err != nil {
		log.Fatal(err)
	}

	json_data, err := json.Marshal(list_of_servers)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(json_data)
	if err != nil {
		log.Fatal(err)
	}
}

//TODO graceful shutdown to make sure all servers ports are closed

func main() {
	num_of_servers := flag.Int("servers", 5, "number of servers to spawn on localhost with port starting from 5000")
	flag.Parse()
	wg := sync.WaitGroup{}
	wg.Add(*num_of_servers)
	start_port := 5000
	for i, start := 0, 5000; i < *num_of_servers; i, start = i+1, start+1 {
		go func(port int) {
			spawnServer(fmt.Sprintf("127.0.0.1:%d", port))
			wg.Done()
		}(start)
	}
	generateJsonFile(*num_of_servers, start_port)
	wg.Wait()
}
