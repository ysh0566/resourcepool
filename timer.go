package resourcepool

import (
	"sync"
	"time"
)

type timerAction = int

const (
	stopMsg = iota
	resetMsg
	triggerMsg
)

type Timer struct {
	interval AtomicDuration
	mu sync.Mutex
	msg chan timerAction
	running bool
}

func NewTimer(interval time.Duration) *Timer {
	t := Timer{
		msg: make(chan timerAction),
	}
	t.interval.Set(interval)
	return &t
}

func(t *Timer) Start(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		return
	}
	t.running = true
	go t.Run(f)
}

func(t *Timer) Run(f func()) {
	for {
		var ch <-chan time.Time
		interval := t.interval.Get()
		if interval <= 0 {
			ch = nil
		} else {
			ch = time.After(interval)
		}
		select {
		case msg := <-t.msg:
			switch msg {
			case stopMsg:
				return
			case resetMsg:
				continue
			}
		case <-ch:
		}
		f()
	}
}

func(t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		t.msg <- stopMsg
		t.running = false
	}
}

func (t *Timer) Trigger() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		t.msg <- triggerMsg
	}
}

func (t *Timer) Interval() time.Duration {
	return t.interval.Get()
}