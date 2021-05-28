package v1

import (
	"bytes"
	"runtime"
	"sync"
)

const defaultBufferSize = 8192

var pool sync.Pool

func init() {
	pool.New = func() interface{} {
		buffer := &bytes.Buffer{}
		buffer.Grow(defaultBufferSize)
		return buffer
	}
}

// NewBuffer creates a new Buffer with automatic recycling mechanism.
func NewBuffer() *Buffer {
	buffer := &Buffer{inner: pool.Get().(*bytes.Buffer)}
	runtime.SetFinalizer(buffer, func(buffer *Buffer) {
		buffer.Reset()
		pool.Put(buffer.inner)
		buffer.inner = nil
	})
	return buffer
}
