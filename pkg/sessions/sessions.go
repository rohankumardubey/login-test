package sessions

import (
	"fmt"	"sync"
)

//I'll probably keep this, but I don't think it's going to work the way I had hoped/intended,
//which is to keep a counter alive while a user session is active. I may need to put that
//information into a counter cookie instead
var counterpool = sync.Pool{
	New: func() interface{}{
		return new(&Counter{})
	},
}

//counter keeps track of login attempts from a user
type Counter struct {
	Count int
	Max   int
}

//NewCounter creates and returns a new counter
func GetCounter(max int) *Counter {
	c := counterpool.Get().(Counter)
	c.Max = max
	c.Reset()
	return c
}

//PutCounter puts a counter back in the pool
func PutCounter(c *Counter){
	counterpool.Put(c)
}

//Reset turns a counter's count to zero
func (c *Counter) Reset() {
	c.Count = 0
}

//CountUp adds 1 to a counter's count
func (c *Counter) CountUp() error {
	if c.Count == c.Max {
		return fmt.Errorf("Error: Maximum attempts exceded.")
	}
	c.Count++
	return nil
}

//CountDown subtracts 1 from a counter
func (c *Counter) CountDown() error {
	if c.Count == 0 {
		return fmt.Errorf("Error: count is already at 0.")
	}
	c.Count--
	return nil
}
