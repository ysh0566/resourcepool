package resourcepool

import (
	"sync/atomic"
	"time"
)

type AtomicInt64 struct {
	int64
}

func(n *AtomicInt64) Set(i int64) {
	atomic.StoreInt64(&n.int64, i)
}

func(n *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&n.int64)
}

func(n *AtomicInt64) Add(i int64) int64{
	return atomic.AddInt64(&n.int64, i)
}

func(n *AtomicInt64) CompareAndSwap(old, new int64) bool{
	return atomic.CompareAndSwapInt64(&n.int64, old, new)
}

func NewAtomicInt64(i int64) AtomicInt64{
	return AtomicInt64{i}
}

type AtomicInt32 struct {
	int32
}

func(n *AtomicInt32) Set(i int32) {
	atomic.StoreInt32(&n.int32, i)
}

func(n *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&n.int32)
}

func(n *AtomicInt32) Add(i int32) int32{
	return atomic.AddInt32(&n.int32, i)
}

func(n *AtomicInt32) CompareAndSwap(old, new int32) bool{
	return atomic.CompareAndSwapInt32(&n.int32, old, new)
}

func NewAtomicInt32(i int32) AtomicInt32{
	return AtomicInt32{i}
}

type AtomicDuration struct {
	int64
}

func(d *AtomicDuration) Set(t time.Duration) {
	atomic.StoreInt64(&d.int64, int64(t))
}

func(d *AtomicDuration) Get() time.Duration{
	return time.Duration(atomic.LoadInt64(&d.int64))
}

func(d *AtomicDuration) Add(t time.Duration) time.Duration{
	return time.Duration(atomic.AddInt64(&d.int64, int64(t)))
}

func(d *AtomicDuration) CompareAndSwap(old, new time.Duration) bool {
	return atomic.CompareAndSwapInt64(&d.int64, int64(old), int64(new))
}

func NewAtomicDuration(t time.Duration) AtomicDuration{
	return AtomicDuration{int64(t)}
}