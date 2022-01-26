package main

import (
	"Server/counter"
	"Server/responses"
	"Server/server"
	"encoding/json"
	"flag"
	"net/http"
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
	server.Run(*port, muxServer)
}

func handle(w http.ResponseWriter, r *http.Request) {
	var parameter string = "client_id"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params, ok := r.URL.Query()[parameter]

	if !ok || len(params) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		missingParameterResponse := responses.NewMissingParameterResponse(parameter)
		body, _ := json.Marshal(missingParameterResponse)
		w.Write([]byte(body))
		return
	}

	clientIdAsInt, err := strconv.Atoi(params[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		missingParameterResponse := responses.NewInvalidParameterTypeResponse(parameter, "Parameter client_id must be of type integer.")
		body, _ := json.Marshal(missingParameterResponse)
		w.Write([]byte(body))
		return
	}

	if c.GetCount(clientIdAsInt) > 4 && c.IsTimerRunning(clientIdAsInt) {
		w.WriteHeader(http.StatusServiceUnavailable)
		serviceUnavailableResponse := responses.NewServiceUnavailableResponse("You reached your request limit.")
		body, _ := json.Marshal(serviceUnavailableResponse)
		w.Write([]byte(body))
		return
	} else {
		wg := sync.WaitGroup{}
		c.Increment(clientIdAsInt, &wg)
		w.WriteHeader(http.StatusOK)
		success := responses.NewSuccessfulResponse("Ok.")
		body, _ := json.Marshal(success)
		w.Write([]byte(body))
		return
	}
}
