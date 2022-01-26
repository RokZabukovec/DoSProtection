package main

import (
	"Server/counter"
	"Server/responses"
	"Server/server"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var c counter.Counter = *counter.NewCounter()

/**
 * Denial-of-service server.
 *
 * @author 					Rok Zabukovec
 * @version 					1.0
 * @since 						1.0
 */
func main() {
	port := flag.Int("port", 8000, "Port to listen on")
	flag.Parse()
	muxServer := http.NewServeMux()
	muxServer.HandleFunc("/", handle)
	err := server.Run(*port, muxServer)

	if err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
		os.Exit(0)
	}
}

/*
 * Handles the request to the "/" path.
 */
func handle(w http.ResponseWriter, r *http.Request) {

	// Only allow GET requests.
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var parameter string = "client_id"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params, ok := r.URL.Query()[parameter]

	if !ok || len(params) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		missingParameterResponse := responses.NewResponse(http.StatusBadRequest, "Missing client_id parameter.")
		body, _ := json.Marshal(missingParameterResponse)
		w.Write([]byte(body))

		return
	}

	clientIdAsInt, err := strconv.Atoi(params[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		missingParameterResponse := responses.NewResponse(http.StatusBadRequest, "Parameter client_id must be of type integer.")
		body, _ := json.Marshal(missingParameterResponse)
		w.Write([]byte(body))

		return
	}

	if c.GetCount(clientIdAsInt) > 4 && c.IsTimerRunning(clientIdAsInt) {
		w.WriteHeader(http.StatusServiceUnavailable)
		serviceUnavailableResponse := responses.NewResponse(http.StatusServiceUnavailable, "Too many requests.")
		body, _ := json.Marshal(serviceUnavailableResponse)
		w.Write([]byte(body))

		return
	} else {
		wg := sync.WaitGroup{}
		c.Increment(clientIdAsInt, &wg)
		w.WriteHeader(http.StatusOK)

		return
	}
}
