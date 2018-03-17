package resourcepool

import (
	"sync/atomic"
	"testing"
	"time"
)

const TESTINTERVAL  = time.Millisecond * 10

var num int32

func f() {
	atomic.AddInt32(&num, 1)
}

func TestTimer_Start(t *testing.T) {
	timer := NewTimer(TESTINTERVAL)
	timer.Start(f)
	defer timer.Stop()
	time.Sleep(TESTINTERVAL/2)
	if atomic.LoadInt32(&num) != 0 {
		t.Errorf("want 0, got %d", num)
	}
	time.Sleep(TESTINTERVAL)
	if atomic.LoadInt32(&num) != 1 {
		t.Errorf("want 1, got %d", num)
	}
	time.Sleep(TESTINTERVAL)
	if atomic.LoadInt32(&num) != 2 {
		t.Errorf("want 2, got %d", num)
	}
}
