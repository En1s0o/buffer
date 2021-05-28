package v2

import (
	"bytes"

	"github.com/En1s0o/buffer/pkg/sync"
)

const defaultBufferSize = 8192

var pool_ sync.Pool

func init() {
	newFunc := func(factory sync.RefCountableFactory) sync.RefCountable {
		buffer := &Buffer{inner: &bytes.Buffer{}}
		buffer.inner.Grow(defaultBufferSize)
		buffer.RefCountable = factory(buffer)
		return buffer
	}
	// pool_ = sync.NewSimplePool(newFunc)
	pool_ = sync.NewStatisticalPool(newFunc)
}

func NewBuffer() *Buffer {
	return pool_.GetAndRef().(*Buffer)
}

func Stats() string {
	return pool_.Stats()
}
