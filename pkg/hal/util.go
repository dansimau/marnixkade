package hal

import (
	"sync/atomic"
	"time"
)

type Timer struct {
	timer *time.Timer

	d       time.Duration
	f       func()
	running atomic.Bool
}

func NewTimer(f func()) *Timer {
	return &Timer{
		f: f,
	}
}

func (t *Timer) Reset(d time.Duration) {
	t.d = d

	if t.timer == nil {
		t.timer = time.AfterFunc(d, t.f)
	} else {
		t.timer.Reset(t.d)
	}

	t.running.Store(true)

	go func(c <-chan time.Time) {
		<-c
		t.running.Store(false)
	}(t.timer.C)
}

func (t *Timer) Stop() bool {
	return t.timer.Stop()
}

func (t *Timer) IsRunning() bool {
	return t.running.Load()
}

func resetTimer(d time.Duration, f func()) *time.Timer {
	return time.AfterFunc(d, f)
}
