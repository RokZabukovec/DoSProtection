package server

type Server struct {
	URL  string
	Port int
}

/**
 * Create a new client.
 */
func NewServer(url string, port int) *Server {
	return &Server{
		URL:  url,
		Port: port,
	}
}
