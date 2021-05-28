package sync

import (
	"fmt"
	"sync"
)

// SimplePool wrapped struct.
type SimplePool struct {
	__ *sync.Pool
}

func NewSimplePool(newFunc func(RefCountableFactory) RefCountable) Pool {
	pool := &SimplePool{__: &sync.Pool{}}
	pool.__.New = func() interface{} {
		return newFunc(pool.newRefCountable)
	}
	return pool
}

func (p *SimplePool) newRefCountable(holder Releasable) RefCountable {
	_, ok := holder.(RefCountable)
	if !ok {
		panic(fmt.Errorf("not implements RefCountable: %#v", holder))
	}

	return NewRefCountable(func() {
		holder.Release()
		p.__.Put(holder)
	}, nil, nil)
}

func (p *SimplePool) GetAndRef() RefCountable {
	rc := p.__.Get().(RefCountable)
	rc.Ref()
	return rc
}

func (p *SimplePool) Stats() string {
	return ""
}
