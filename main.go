package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"sse-server/handlers"
	"sse-server/middlewares"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/send", handlers.Send).Methods("GET").Name("healthcheck")
	router.HandleFunc("/sse", handlers.SSE).Methods("GET").Name("healthcheck")

	server := &http.Server{
		Addr:         ":8081",
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      middlewares.CORSMiddleware(router),
	}

	log.Fatal("HTTP server error: ", server.ListenAndServe())

}
