package server

import (
	"log"
	"net/http"

	"github.com/youssefmaher99/equilib/internal/strategy"
)

type server struct {
	listenAddr string
	strategy   strategy.Strategy
}

func New(listenAddr string) *server {
	return &server{listenAddr: listenAddr}
}

func (s *server) intercept() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("intercepted")
		w.Write([]byte("hello"))
	})
}

func (s *server) Start() error {
	log.Printf("equilib is running on [%s]\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.intercept())
}
