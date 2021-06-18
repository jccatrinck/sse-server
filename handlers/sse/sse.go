package sse

import (
	"fmt"
	"net/http"

	"sse-server/sse"
)

var broker = sse.NewBroker()

func Send(msg string) {
	broker.Notifier <- []byte(msg)
}

func Handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Request: %+v\n\n", req)

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Header.Get("HTTP_ORIGIN") != "" {
		//header("Access-Control-Allow-Origin: {$_SERVER['HTTP_ORIGIN']}");
		rw.Header().Set("Access-Control-Credentials", "true")
		rw.Header().Set("Access-Control-Methods", "GET, POST, OPTIONS")
	}

	if req.Method == "OPTIONS" {
		if req.Header.Get("HTTP_ACCESS_CONTROL_REQUEST_METHOD") != "" {
			rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}

		if req.Header.Get("HTTP_ACCESS_CONTROL_REQUEST_HEADERS") != "" {
			rw.Header().Set("Access-Control-Allow-Headers", req.Header.Get("HTTP_ACCESS_CONTROL_REQUEST_HEADERS"))
		}
	}

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.NewClient(messageChan)

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.CloseClient(messageChan)
	}()

	go func() {
		<-req.Context().Done()
		broker.CloseClient(messageChan)
	}()

	for {
		// Write to the ResponseWriter
		// Server Sent Events compatible
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}

}
