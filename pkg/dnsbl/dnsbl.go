package dnsbl

import (
	"log"
	"net"
	"sync"
)

// Checker controls the flow of package and provides a single point of
// entry for its users.
type Checker struct {
	// length is the number of providers
	// queried is the number of providers who answered
	// positive is the number of appearances reported
	length, queried, positive int
	providers                 []string

	// lookup contains the lookup function used
	lookup func(string) ([]string, error)
	// Functional control
	wg sync.WaitGroup
}

// NewChecker creates a new, default configured Checker
func NewChecker(providers []string) *Checker {
	return &Checker{lookup: net.LookupHost, providers: providers}
}

// Query contacts the providers to check if the IP is present in their
// lists
func (c *Checker) Query() *Checker {
	responses := make(chan int)
	for _, provider := range c.providers {
		c.length++
		c.wg.Add(1)
		go func(provider string) {
			responses <- c.query(provider)
		}(provider)
	}
	go func() {
		for response := range responses {
			c.positive += response
			c.queried++
			c.wg.Done()
		}
	}()
	c.wg.Wait()
	close(responses)
	return c
}

// Stats returns the number of positive results along with the amount
// of blacklists supplied and the amount that were reachable.
// length is the number of providers
// queried is the number of providers who answered
// positive is the number of appearances reported
func (c *Checker) Stats() (positive, queried, length int) {
	return c.positive, c.queried, c.length
}

// query queries a DNSBL and returns true if the argument gets a match
// in the BL.
func (c *Checker) query(address string) int {
	// We ignore errors because the providers where we are not
	// flagged can't be resolved. We can not distinguish if we are
	// not on their list or their service is broken.
	result, _ := c.lookup(address)
	if len(result) > 0 {
		log.Printf("%v returned %v\n", address, result)
	}
	return len(result)
}
