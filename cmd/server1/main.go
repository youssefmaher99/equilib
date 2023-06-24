package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Intercept")
	w.Write([]byte("pong from 5001"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", helloHandler)
	hostAndport := "127.0.0.1:5001"
	fmt.Printf("Server is running on [%s]\n", hostAndport)
	log.Fatal(http.ListenAndServe(hostAndport, mux))
}
