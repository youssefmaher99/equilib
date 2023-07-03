package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/youssefmaher99/equilib/internal/connectionPool"
	"github.com/youssefmaher99/equilib/internal/ringBuffer"
)

type server struct {
	listenAddr string
	rngBuffer  *ringBuffer.RingBuffer
	pool       connectionPool.ConnectionPooler
}

func New(listenAddr string, size int, servers []string) *server {
	rng := ringBuffer.New(size)
	rng.Populate(servers)

	pool := connectionPool.New()
	pool.Populate(servers)

	return &server{listenAddr: listenAddr, rngBuffer: rng, pool: pool}
}

func (s *server) intercept() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Intercept")

		server := s.rngBuffer.Next()
		// get conn from pool or create new conn
		conn, err := s.pool.Get(server)
		if err != nil {
			panic(err)
		}
		// create client
		dialFunc := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return conn, nil
		}

		// Create an http.Transport that uses the custom DialContext function.
		transport := &http.Transport{
			DialContext: dialFunc,
		}

		// Create an http.Client using the custom transport.
		client := &http.Client{
			Transport: transport,
		}

		// Make an HTTP request using the custom client.
		resp, err := client.Get("http://" + server + "/ping")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)

		// TODO release conn back to pool

		// TODO : better error handling (mark server down - forward error coming from server to client)
		// response, err := http.Get("http://" + server + r.URL.String())
		// if err != nil {
		// 	log.Println(err)
		// }

		// // forward response headers coming from a server
		// for key, val := range response.Header {
		// 	w.Header().Add(key, val[0])
		// }

		// // forward response body coming from a server
		// response_body, err := io.ReadAll(response.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// w.Write(response_body)
	})
}

func (s *server) Start() error {
	log.Printf("equilib is running on [%s]\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.intercept())
}
