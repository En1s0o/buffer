package sync

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

// StatisticalPool wrapped into statistical struct.
type StatisticalPool struct {
	__ *sync.Pool

	// memStats memory stats, for statistics only
	memStats *runtime.MemStats
	// returned sum returns, for statistics only
	returned uint32
	// allocated sum allocates, it does not decrease when objects are returned, for statistics only
	allocated uint32
}

func NewStatisticalPool(New func(factory RefCountableFactory) RefCountable) Pool {
	pool := &StatisticalPool{
		__:       &sync.Pool{},
		memStats: &runtime.MemStats{},
	}
	pool.__.New = func() interface{} {
		return New(pool.newRefCountable)
	}
	return pool
}

func (p *StatisticalPool) newRefCountable(holder Releasable) RefCountable {
	_, ok := holder.(RefCountable)
	if !ok {
		panic(fmt.Errorf("not implements RefCountable: %#v", holder))
	}

	atomic.AddUint32(&p.allocated, 1)
	return NewRefCountable(func() {
		holder.Release()
		p.__.Put(holder)
	}, nil, func() {
		atomic.AddUint32(&p.returned, 1)
	})
}

func (p *StatisticalPool) GetAndRef() RefCountable {
	rc := p.__.Get().(RefCountable)
	rc.Ref()
	return rc
}

func (p *StatisticalPool) Stats() string {
	runtime.ReadMemStats(p.memStats)
	return fmt.Sprintf(
		"Allocated: %d, Returned: %d, TotalAlloc: %d, HeapObjects: %d",
		atomic.LoadUint32(&p.allocated),
		atomic.LoadUint32(&p.returned),
		p.memStats.TotalAlloc,
		p.memStats.HeapObjects)
}
