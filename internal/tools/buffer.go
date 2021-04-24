package tools

import (
	"bytes"
	"sync"
)

//not sure why this is a var, but I'm guessing so we can call a method on it?
var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

//This allows us to retrieve a pointer to bytes.Buffer from bufPool. Somehow
func GetBuf() *bytes.Buffer {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

//I'm guessing this puts the buffer back into the sync pool?
func PutBuf(b *bytes.Buffer) {
	bufPool.Put(b)
}
