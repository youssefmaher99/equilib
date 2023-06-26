package server

import (
	"io"
	"log"
	"net/http"

	"github.com/youssefmaher99/equilib/internal/ringBuffer"
)

type server struct {
	listenAddr string
	rngBuffer  *ringBuffer.RingBuffer
}

func New(listenAddr string, size int, servers []string) *server {
	rng := ringBuffer.New(size)
	rng.Populate(servers)
	return &server{listenAddr: listenAddr, rngBuffer: rng}
}

func (s *server) intercept() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Intercept")

		server := s.rngBuffer.Next()
		response, err := http.Get(server + r.URL.String())
		if err != nil {
			log.Println(err)
		}

		// forward response headers coming from a server
		for key, val := range response.Header {
			w.Header().Add(key, val[0])
		}

		// forward response body coming from a server
		response_body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(response_body)
	})
}

func (s *server) Start() error {
	log.Printf("equilib is running on [%s]\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.intercept())
}
