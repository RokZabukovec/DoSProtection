package server

import (
	"fmt"
	"log"
	"net/http"
)

/**
 * Run the server.
 */
func Run(port int, server *http.ServeMux) {
	fmt.Printf("Server is running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
