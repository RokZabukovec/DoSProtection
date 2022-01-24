package counter

type Client struct {
	count          int
	timerIsRunning bool
}

/**
 * Create a new client.
 */
func NewClient() *Client {
	return &Client{
		count:          0,
		timerIsRunning: true,
	}
}
