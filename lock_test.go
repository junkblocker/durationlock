package durationlock_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/junkblocker/durationlock"
)

func TestDurationLock_Take(t *testing.T) {
	tests := []struct {
		name     string
		lockFor  time.Duration
		waitFor  time.Duration
		wantTake bool
	}{
		{
			"Auto expire",
			time.Second * 2,
			time.Second * 3,
			true,
		},
		{
			"No expire",
			time.Second * 4,
			time.Second * 1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := durationlock.New(tt.lockFor)
			defer dl.Release()
			if got := dl.Take(); !got {
				t.Errorf("DurationLock.Take() = %v, want %v", got, true)
			}
			t.Log(dl)
			time.Sleep(tt.waitFor)
			if got := dl.Take(); got != tt.wantTake {
				t.Errorf("DurationLock.Take() = %v, want %v", got, tt.wantTake)
			}
		})
	}
}

func TestDurationLock_Release(t *testing.T) {
	tests := []struct {
		name         string
		lockFor      time.Duration
		releaseAfter time.Duration
		takeIn       time.Duration
		wantTake     bool
	}{
		{
			"Early release available",
			time.Second * 6,
			time.Second * 2,
			time.Second * 4,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := durationlock.New(tt.lockFor)
			defer dl.Release()
			go func() {
				time.Sleep(tt.releaseAfter)
				dl.Release()
			}()
			if got := dl.Take(); !got {
				t.Errorf("got = %v, want %v", got, true)
			}
			t.Log(dl)
			time.Sleep(tt.takeIn)
			if got := dl.Take(); got != tt.wantTake {
				t.Errorf("got = %v, want %v", got, tt.wantTake)
			}
		})
	}
}

func ExampleLock_Take() {
	l := durationlock.New(time.Second * 5)
	if l.Take() {
		fmt.Println("Lock acquired")
	} else {
		fmt.Println("Lock was already held")
	}
	// Output: Lock acquired
}
