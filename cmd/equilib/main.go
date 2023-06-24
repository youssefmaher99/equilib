package main

import (
	"log"
	"net/http"

	"github.com/youssefmaher99/equilib/internal/server"
)

// client libary ???
// route to different server
// response is sent from
// serve service is
func interceptionHandler(w http.ResponseWriter, r *http.Request) {

	// //CLIENT
	// fmt.Printf("forwarding to %s\n", servers[ptr])
	// res, err := http.Get("http://" + servers[ptr] + r.URL.String())
	// if err != nil {
	// 	log.Println(err)
	// }

	// ptr = (ptr + 1) % len(servers)
	// w.Write([]byte("HELLO"))
}

func main() {
	// rngBuffer := ringBuffer.New(2)
	s := server.New("127.0.0.1:8080")
	log.Fatal(s.Start())
}
