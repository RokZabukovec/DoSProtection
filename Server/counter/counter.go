package counter

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	counter map[int]Client
	mutex   sync.RWMutex
}

/**
 * Create a new MissingClientIdResponse.
 */
func NewCounter() *Counter {
	return &Counter{
		counter: make(map[int]Client),
		mutex:   sync.RWMutex{},
	}
}

/**
 * Increments client_id by one.
 */
func (c *Counter) Increment(client_id int, wg *sync.WaitGroup) {
	if !c.IsPresent(client_id) {
		c.InitializeClient(client_id)
	}

	if c.GetCount(client_id) > 4 {
		c.InitializeClient(client_id)
	}

	wg.Add(1)
	go func() {
		c.mutex.Lock()
		client := c.counter[client_id]
		client.count++
		c.counter[client_id] = client
		c.mutex.Unlock()
		wg.Done()
	}()
	wg.Wait()
	c.mutex.Lock()
	client := c.counter[client_id]
	fmt.Printf("ID: %d -> %d \n", client_id, client.count)
	c.mutex.Unlock()
}

/**
 * Returns the count of the client's requests.
 */
func (c *Counter) GetCount(client_id int) int {
	c.mutex.Lock()
	client := c.counter[client_id]
	c.mutex.Unlock()
	return client.count
}

/**
 * Returns true if client_id is present in the counter.
 */
func (c *Counter) IsPresent(client_id int) bool {
	c.mutex.Lock()
	// check if key is present in counter
	if _, ok := c.counter[client_id]; ok {
		c.mutex.Unlock()
		return true
	}
	c.mutex.Unlock()
	return false
}

/**
 * Adds the client to the counter.
 */
func (c *Counter) InitializeClient(client_id int) int {
	// check if key is present in counter
	if c.IsPresent(client_id) {
		c.ResetClientCount(client_id)
	}

	go c.StartTimer(client_id)
	return c.GetCount(client_id)
}

func (c *Counter) StartTimer(client_id int) {
	c.mutex.Lock()
	client := c.counter[client_id]
	fmt.Printf("Timer for client %d started.\n", client_id)
	timer := time.NewTimer(5 * time.Second)
	client.timerIsRunning = true
	c.counter[client_id] = client
	c.mutex.Unlock()
	<-timer.C
	timer.Stop()

	c.mutex.Lock()
	clientAfter := c.counter[client_id]
	clientAfter.timerIsRunning = false
	c.counter[client_id] = clientAfter
	c.mutex.Unlock()
	fmt.Printf("Timer for client %d expired with %d requests.\n", client_id, clientAfter.count)
}

func (c *Counter) IsTimerRunning(client_id int) bool {
	c.mutex.Lock()
	client := c.counter[client_id]
	isRunning := client.timerIsRunning
	c.mutex.Unlock()
	return isRunning
}

func (c *Counter) ResetClientCount(client_id int) {
	c.mutex.Lock()
	client := c.counter[client_id]
	client.count = 0
	c.counter[client_id] = client
	c.mutex.Unlock()
}
