package sync

import "sync/atomic"

// Pool interface.
type Pool interface {
	// newRefCountable creates a new RefCountable.
	newRefCountable(holder Releasable) RefCountable

	// GetAndRef obtains a RefCountable and increase its reference count.
	// When you are not using it, you should call DeRef() explicitly.
	GetAndRef() RefCountable

	// Stats dump.
	Stats() string
}

// Releasable interface.
type Releasable interface {
	// Release cleanup.
	Release()
}

// RefCountable interface.
type RefCountable interface {
	// Ref increase reference count.
	Ref()

	// DeRef decrease reference count.
	DeRef()
}

// RefCountableFactory function type.
type RefCountableFactory func(holder Releasable) RefCountable

func NewRefCountable(put func(), prePut func(), postPut func()) RefCountable {
	return &__{put: put, prePut: prePut, postPut: postPut}
}

type __ struct {
	refCount uint32
	put      func()
	prePut   func()
	postPut  func()
}

func (rc *__) Ref() {
	atomic.AddUint32(&rc.refCount, 1)
}

func (rc *__) DeRef() {
	if atomic.AddUint32(&rc.refCount, ^uint32(0)) == 0 {
		if rc.prePut != nil {
			rc.prePut()
		}
		rc.put()
		if rc.postPut != nil {
			rc.postPut()
		}
	}
}
