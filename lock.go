package durationlock

import (
	"fmt"
	"sync"
	"time"
)

// Lock is the locking API struct.
// It provides APIs to acquire a lock for a specified interval and
// a try and release facility.
type Lock struct {
	sync.Mutex
	locked   bool
	duration time.Duration
	timeout  *time.Timer
}

// New lock object for a specified duration.
// The lock needs to be 'Take()'n before it is considered held.
func New(d time.Duration) *Lock {
	return &Lock{duration: d}
}

// String implements the stringer interface.
func (dl *Lock) String() string {
	dl.Lock()
	defer dl.Unlock()
	if dl.locked {
		return fmt.Sprintf("Locked for %v", dl.duration)
	}
	return "Unlocked"
}

// Take attempts to take the lock. Returns bool to indicate success or failure.
func (dl *Lock) Take() bool {
	dl.Lock()
	defer dl.Unlock()
	if dl.locked {
		return false
	}
	dl.locked = true
	dl.timeout = time.AfterFunc(dl.duration, dl.Release)
	return true
}

// Release the lock unconditionally.
func (dl *Lock) Release() {
	dl.Lock()
	defer dl.Unlock()
	if dl.timeout != nil {
		dl.timeout.Stop()
		dl.timeout = nil
	}
	dl.locked = false
}
