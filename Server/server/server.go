package server

import (
	"fmt"
	"net/http"
)

/**
 * Run the server.
 */
func Run(port int, server *http.ServeMux) error {
	fmt.Printf("Server is running on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	return err
}
