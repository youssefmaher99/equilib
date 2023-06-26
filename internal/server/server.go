package server

import (
	"fmt"
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
		res, err := http.Get(server + r.URL.String())
		if err != nil {
			log.Println(err)
		}
		fmt.Println(res)
		w.Write([]byte("hello"))
	})
}

func (s *server) Start() error {
	log.Printf("equilib is running on [%s]\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.intercept())
}
