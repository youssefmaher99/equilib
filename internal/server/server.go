package server

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/youssefmaher99/equilib/internal/connectionPool"
	"github.com/youssefmaher99/equilib/internal/ringBuffer"
)

type server struct {
	listenAddr string
	rngBuffer  *ringBuffer.RingBuffer
	client     *http.Client
}

type customTransport struct {
	pool              connectionPool.ConnectionPooler
	originalTransport http.RoundTripper
}

func (c *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// get connection from pool
	conn, err := c.pool.Get(r.Host)
	if err != nil {
		panic(err)
	}

	// set tcpConn with the new connection from the pool which will be read in the defaultTransport and use the connection
	r = r.WithContext(context.WithValue(r.Context(), "tcpConn", conn))
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	// release conenction back to pool
	defer c.pool.Release(conn, r.Host)

	return resp, nil
}

func New(listenAddr string, size int, servers []string) *server {
	rng := ringBuffer.New(size)
	rng.Populate(servers)

	pool := connectionPool.New()
	pool.Populate(servers)

	newCustomTransport := customTransport{pool: pool, originalTransport: http.DefaultTransport}
	client := &http.Client{Transport: &newCustomTransport}

	return &server{listenAddr: listenAddr, rngBuffer: rng, client: client}
}

func (s *server) intercept() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		server := s.rngBuffer.Next()

		request, err := http.NewRequest(r.Method, "http://"+server+"/ping", nil)
		if err != nil {
			panic(err)
		}

		resp, err := s.client.Do(request)
		if err != nil {
			panic(err)
		}

		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		for key, val := range resp.Header {
			w.Header().Add(key, val[0])
		}

		w.Write(resp_body)
	})
}

func (s *server) Start() error {
	log.Printf("equilib is running on [%s]\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.intercept())
}
