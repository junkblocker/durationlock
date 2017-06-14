package durationlock

import (
	"sync"
	"time"
)

type DurationLock struct {
	sync.Mutex
	locked   bool
	duration time.Duration
	timeout  *time.Timer
}

// Take the lock. Returns bool to indicate success or failure.
func (tl *DurationLock) Take() bool {
	tl.Lock()
	defer tl.Unlock()
	if tl.locked {
		return false
	}
	tl.locked = true
	tl.timeout = time.AfterFunc(tl.duration, tl.Release)
	return true
}

// Release the lock unconditionally.
func (tl *DurationLock) Release() {
	tl.Lock()
	defer tl.Unlock()
	if tl.timeout != nil {
		tl.timeout.Stop()
		tl.timeout = nil
	}
	tl.locked = false
}
