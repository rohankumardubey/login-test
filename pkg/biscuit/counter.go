package biscuit

import (
	"fmt"
	"sync"
)

//I'll probably keep this, but I don't think it's going to work the way I had hoped/intended,
//which is to keep a counter alive while a user session is active. I may need to put that
//information into a counter cookie instead
var counterpool = sync.Pool{
	New: func() interface{} {
		return new(&counter{})
	},
}

//counter keeps track of login attempts from a user
type counter struct {
	Count int
	Max   int
}

//NewCounter creates and returns a new counter
func getCounter(max int) *counter {
	c := counterpool.Get().(*counter)
	c.Max = max
	c.Reset()
	return c
}

//PutCounter puts a counter back in the pool
func putCounter(c *counter) {
	counterpool.Put(c)
}

//Reset turns a counter's count to zero
func (c *counter) Reset() {
	c.Count = 0
}

//CountUp adds 1 to a counter's count
func (c *counter) CountUp() error {
	if c.Count == c.Max {
		return fmt.Errorf("Error: Maximum attempts exceded.")
	}
	c.Count++
	return nil
}

//CountDown subtracts 1 from a counter
func (c *counter) CountDown() error {
	if c.Count == 0 {
		return fmt.Errorf("Error: count is already at 0.")
	}
	c.Count--
	return nil
}
