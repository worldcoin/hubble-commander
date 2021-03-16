package testutils

import (
	"log"
	"time"
)

var timeout = 1 * time.Second

func WaitToPass(fn func() bool) {
	startTime := time.Now()
	for time.Since(startTime) < timeout {
		if fn() {
			return
		}
	}
	log.Fatal("WaitToPass: timeout")
}
