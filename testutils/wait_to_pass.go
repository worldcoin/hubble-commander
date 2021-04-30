package testutils

import (
	"testing"
	"time"
)

const TryInterval = 250 * time.Millisecond

func WaitToPass(t *testing.T, fn func() bool, timeout time.Duration) {
	immediately := make(chan struct{}, 1)
	immediately <- struct{}{}
	defer close(immediately)

	ticker := time.NewTicker(TryInterval)
	defer ticker.Stop()

	timeoutC := time.After(timeout)
	for {
		select {
		case <-immediately:
			if fn() {
				return
			}
		case <-ticker.C:
			if fn() {
				return
			}
		case <-timeoutC:
			t.Fatal("WaitToPass: timeout")
			return
		}
	}
}
