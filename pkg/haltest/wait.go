package haltest

import (
	"strings"
	"testing"
	"time"
)

const (
	pollInterval = 100 * time.Millisecond
	waitTimeout  = 3 * time.Second
)

// WaitFor waits for the given function to return true.
func WaitFor(t *testing.T, fn func() bool, msg ...string) {
	timeout := time.After(waitTimeout)
	for {
		select {
		case <-timeout:
			t.Fatal(strings.Join(msg, " "))
		default:
			if fn() {
				return
			}

			time.Sleep(pollInterval)
		}
	}
}
