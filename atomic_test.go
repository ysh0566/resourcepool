package resourcepool

import (
	"testing"
	"time"
)

func TestNewAtomicInt32(t *testing.T) {
	ai32 := NewAtomicInt32(233)
	if ai32.Get() != 233 {
		t.Errorf("want 233, got %d", ai32.Get())
	}
	ai32.Add(2100)
	if ai32.Get() != 2333 {
		t.Errorf("want 2333, got %d", ai32.Get())
	}
	if !ai32.CompareAndSwap(2333, 1234) {
		t.Errorf("want true, got false")
	}
	if ai32.Get() != 1234 {
		t.Errorf("want 1234, got %d", ai32.Get())
	}
	if ai32.CompareAndSwap(1, 1234) {
		t.Errorf("want false, got true")
	}
}

func TestNewAtomicInt64(t *testing.T) {
	ai64 := NewAtomicInt64(233)
	if ai64.Get() != 233 {
		t.Errorf("want 233, got %d", ai64.Get())
	}
	ai64.Add(2100)
	if ai64.Get() != 2333 {
		t.Errorf("want 2333, got %d", ai64.Get())
	}
	if !ai64.CompareAndSwap(2333, 1234) {
		t.Errorf("want true, got false")
	}
	if ai64.Get() != 1234 {
		t.Errorf("want 1234, got %d", ai64.Get())
	}
	if ai64.CompareAndSwap(1, 1234) {
		t.Errorf("want false, got true")
	}
}

func TestNewAtomicDuration(t *testing.T) {
	aid := NewAtomicDuration(time.Minute)
	if aid.Get() != time.Minute {
		t.Errorf("want %d, got %d", time.Minute, aid.Get())
	}
	aid.Add(time.Minute)
	if aid.Get() != time.Minute*2 {
		t.Errorf("want %d, got %d", time.Minute*2, aid.Get())
	}
	if !aid.CompareAndSwap(time.Minute*2, time.Minute*3) {
		t.Errorf("want true, got false")
	}
	if aid.Get() != time.Minute*3 {
		t.Errorf("want %d, got %d", time.Minute*3, aid.Get())
	}
	if aid.CompareAndSwap(time.Minute, time.Minute*10) {
		t.Errorf("want false, got true")
	}
}
