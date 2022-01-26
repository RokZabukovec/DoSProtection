package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Result struct {
	Status  int
	ClietId int
}

type Request struct {
	ClientId int
}

var req_number *int
var res_number *int

func sendRequest(w int, jobs <-chan Request, results chan<- Result, wg *sync.WaitGroup) {
	for req := range jobs {
		resp, err := http.Get("http://localhost:8000/?client_id=" + fmt.Sprintf("%d", w))

		if err != nil {
			fmt.Printf("Client %d: Error: %s\n", req.ClientId, err.Error())
		}
		resp.Body.Close()
		resp.Close = true
		results <- Result{Status: resp.StatusCode, ClietId: w}
		wg.Done()
	}
}

func main() {
	clients := flag.Int("clients", 1, "The number of clients to generate.")
	threads := flag.Int("threads", 1, "The number of threads per client.")
	flag.Parse()

	jobs := make(chan Request, *threads)
	results := make(chan Result, *threads)
	wg := &sync.WaitGroup{}
	SetupCloseHandler(jobs, results, wg)

	timeout := time.Duration(100 * time.Millisecond)
	_, timeout_err := net.DialTimeout("tcp", "localhost:8000", timeout)

	if timeout_err != nil {
		panic("Server not reachable. Exiting.")
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	for w := 0; w < *clients; w++ {
		go sendRequest(w, jobs, results, wg)
	}

	for client := 1; client <= *clients; client++ {
		for t := 0; t < *threads; t++ {
			go func() {
				for {
					jobs <- Request{ClientId: client}
					wg.Add(1)
					time.Sleep(time.Millisecond * 100)
				}
			}()

		}
	}

	for result := range results {
		fmt.Printf("Client: %d\tStatus: %d\n", result.ClietId, result.Status)
	}

}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler(jobs chan Request, results chan Result, wg *sync.WaitGroup) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(jobs)
		fmt.Println("Shutting down...Waiting for all requests to finish.")
		wg.Wait()
		fmt.Printf("Still processing %d requests.\n", len(jobs))
		close(results)
		fmt.Println("Bye!")
		os.Exit(0)
	}()
}
