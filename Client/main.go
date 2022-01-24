package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gosuri/uilive"
)

type Result struct {
	Status  int
	ClietId int
}

type Request struct {
	ClientId int
}

func sendRequest(w int, jobs <-chan Request, results chan<- Result) {

	for req := range jobs {

		resp, err := http.Get("http://localhost:8000/?client_id=" + fmt.Sprintf("%d", w))

		if err != nil {
			fmt.Printf("Client %d: Error: %s\n", req.ClientId, err.Error())
		}
		resp.Body.Close()
		resp.Close = true
		results <- Result{Status: resp.StatusCode, ClietId: w}
	}
}

func main() {
	clients := flag.Int("clients", 1, "The number of clients to generate.")
	threads := flag.Int("threads", 1, "The number of threads per client.")
	flag.Parse()

	jobs := make(chan Request, 100)
	results := make(chan Result, 100)
	timeout := time.Duration(100 * time.Millisecond)
	_, timeout_err := net.DialTimeout("tcp", "localhost:8000", timeout)
	if timeout_err != nil {
		panic("Server not reachable. Exiting.")
	}
	for w := 0; w < *clients; w++ {
		go sendRequest(w, jobs, results)
	}

	for client := 1; client <= *clients; client++ {
		for t := 0; t < *threads; t++ {
			go func() {
				for {
					jobs <- Request{ClientId: client}
					time.Sleep(time.Millisecond * 5000)
				}
			}()

		}
	}

	writer := uilive.New()
	writer.Start()
	for result := range results {
		fmt.Fprintf(writer, "%d: %d\n", result.ClietId, result.Status)
	}
	writer.Stop()
}
