package server

import (
	"fmt"
	"net"
	"net/url"
	"time"
)

type Server struct {
	URL  string
	Port int
}

/**
 * Create a new server.
 */
func NewServer(url string, port int) *Server {
	return &Server{
		URL:  url,
		Port: port,
	}
}

/**
 * Validate the server's address.
 */
func (s *Server) IsValidUrl() bool {
	cleanUrl, err := url.ParseRequestURI(s.URL)
	if err != nil || cleanUrl.Scheme == "" || cleanUrl.Host == "" {
		return false
	}
	return true
}

/**
 * Check if the server can be accessible.
 */
func (s *Server) IsReachable() bool {
	timeout := time.Duration(time.Second)
	cleanUrl, err := url.ParseRequestURI(s.URL)
	if err != nil || cleanUrl.Scheme == "" || cleanUrl.Host == "" {
		return false
	}
	address := fmt.Sprintf("%s:%d", cleanUrl.Hostname(), s.Port)
	connection, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return false
	}
	defer connection.Close()
	return true
}
