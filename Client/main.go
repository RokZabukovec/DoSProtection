package main

import (
	"Client/server"
	"flag"
	"fmt"
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

func sendRequest(w int, jobs <-chan Request, results chan<- Result, wg *sync.WaitGroup, serv server.Server) {
	var url = fmt.Sprintf("%s:%d/?client_id=%d", serv.URL, serv.Port, w)
	for range jobs {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("The request could not be made.")
			fmt.Println(err)
			os.Exit(1)
		}
		resp.Body.Close()
		results <- Result{Status: resp.StatusCode, ClietId: w}
		wg.Done()
	}
}

func main() {
	// The initial data that can be modified by the user.
	clients := flag.Int("clients", 1, "The number of clients to generate.")
	threads := flag.Int("threads", 1, "The number of threads per client.")
	address := flag.String("addr", "http://localhost", "The address of the server.")
	port := flag.Int("port", 8000, "The port of the server.")
	flag.Parse()

	requestServer := server.NewServer(*address, *port)

	// Validate if the provided URL is valid.
	if !requestServer.IsValidUrl() {
		fmt.Println("The provided URL is not valid.")
		os.Exit(1)
	}

	// Validate if the server is reachable.
	if !requestServer.IsReachable() {
		fmt.Println("The server is not reachable.")
		os.Exit(1)
	}

	// Create a channels to communicate between the workers and the main thread.
	jobs := make(chan Request, *threads)
	results := make(chan Result, *threads)
	wg := &sync.WaitGroup{}

	SetupShutdownHandler(jobs, results, wg)

	for client := 1; client <= *clients; client++ {
		go sendRequest(client, jobs, results, wg, *requestServer)
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

	// Handle all the responses and print them to the console.
	for result := range results {
		fmt.Printf("Client: %d\tStatus: %d\n", result.ClietId, result.Status)
	}

}

// SetupShutdownHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupShutdownHandler(jobs chan Request, results chan Result, wg *sync.WaitGroup) {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-done
		close(jobs)
		fmt.Println("Shutting down...Waiting for all requests to finish.")
		wg.Wait()
		close(results)
		fmt.Println("Bye!")
		os.Exit(0)
	}()
}
