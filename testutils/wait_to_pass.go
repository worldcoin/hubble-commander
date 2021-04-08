package testutils

import (
	"log"
	"time"
)

func WaitToPass(fn func() bool, timeout time.Duration) {
	startTime := time.Now()
	for time.Since(startTime) < timeout {
		if fn() {
			return
		}
	}
	log.Fatal("WaitToPass: timeout")
}
