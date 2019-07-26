package resourcepool

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrClosed  = errors.New("resource pool is closed")
	ErrTimeout = errors.New("resource pool timed out")
)

type Resource interface {
	Close()
}

type Factory func() (Resource, error)

type ResourceWapper struct {
	resource     Resource
	lastTimeUsed time.Time
}

type ResourcePool struct {
	resources   chan ResourceWapper
	factory     Factory
	capacity    AtomicInt64
	idleTimeout AtomicDuration //空闲关闭超时
	idleTimer   *Timer
	available   AtomicInt64    //可用资源数量
	active      AtomicInt64    //活跃资源数量
	inUse       AtomicInt64    //正在使用资源数量
	waitCount   AtomicInt64    //被记录的等待资源次数
	waitTime    AtomicDuration //记录的等待时长
	idleClosed  AtomicInt64    //空闲超时关闭数量
}

func NewResourcePool(factory Factory, capacity, maxCap int, idleTimeout time.Duration) *ResourcePool {
	if capacity <= 0 || maxCap <= 0 || maxCap < capacity {
		panic(errors.New(fmt.Sprintf("Invalid caption: %d or maxCap: %d", capacity, maxCap)))
	}
	rp := ResourcePool{
		resources:   make(chan ResourceWapper, maxCap),
		factory:     factory,
		capacity:    NewAtomicInt64(int64(capacity)),
		available:   NewAtomicInt64(int64(capacity)),
		idleTimeout: NewAtomicDuration(idleTimeout),
	}
	for i := 0; i < capacity; i++ {
		rp.resources <- ResourceWapper{}
	}
	if idleTimeout > 0 {
		timer := NewTimer(idleTimeout / 10)
		timer.Start(rp.closeIdleResource)
		rp.idleTimer = timer
	}
	return &rp
}

func (rp *ResourcePool) closeIdleResource() {
	timeout := rp.IdleTimeout()
	capacity := rp.Capacity()
	for i := 0; i < capacity; i++ {
		var wapper ResourceWapper
		select {
		case wapper = <-rp.resources:
		default:
			return
		}

		if wapper.resource != nil && timeout > 0 && time.Now().Sub(wapper.lastTimeUsed) > timeout {
			wapper.resource.Close()
			wapper.resource = nil
			rp.active.Add(-1)
			rp.idleClosed.Add(1)
		}
		rp.resources <- wapper
	}
}

func (rp *ResourcePool) Get() (resource Resource, err error) {
	var wapper ResourceWapper
	var ok bool
	select {
	case wapper, ok = <-rp.resources:
	default:
		startTime := time.Now()
		wapper, ok = <-rp.resources
		rp.waitCount.Add(1)
		rp.waitTime.Add(time.Now().Sub(startTime))
	}
	if !ok {
		return nil, ErrClosed
	}
	if wapper.resource == nil {
		wapper.resource, err = rp.factory()
		if err != nil {
			rp.resources <- ResourceWapper{}
			return nil, err
		}
		rp.active.Add(1)
	}
	rp.inUse.Add(1)
	rp.available.Add(-1)
	return wapper.resource, nil
}

func (rp *ResourcePool) Put(resource Resource) {
	var wapper ResourceWapper
	if resource != nil {
		wapper = ResourceWapper{resource, time.Now()}
	} else {
		rp.active.Add(-1)
	}
	select {
	case rp.resources <- wapper:
	default:
		panic(errors.New("Put A Resource into A Full ResoucePoll"))
	}
	rp.inUse.Add(-1)
	rp.available.Add(1)
}

func (rp *ResourcePool) Close() {
	if rp.Capacity() == 0 {
		return
	}
	if &rp.idleTimer != nil && rp.idleTimer.running {
		rp.idleTimer.Stop()
	}
	_ = rp.SetCapacity(0)
}

func (rp *ResourcePool) SetCapacity(newCap int) error {
	if newCap < 0 || newCap > cap(rp.resources) {
		return fmt.Errorf("capacity %d is out of range", newCap)
	}
	var oldCap int
	for {
		oldCap = rp.Capacity()
		if oldCap == 0 {
			return ErrClosed
		}
		if oldCap == newCap {
			return nil
		}
		if rp.capacity.CompareAndSwap(int64(oldCap), int64(newCap)) {
			break
		}
	}
	if newCap > oldCap {
		for i := 0; i < newCap-oldCap; i++ {
			rp.resources <- ResourceWapper{}
			rp.available.Add(1)
		}
	} else {
		for i := 0; i < oldCap-newCap; i++ {
			wapper := <-rp.resources
			if wapper.resource != nil {
				wapper.resource.Close()
				wapper.resource = nil
				rp.active.Add(-1)
			}
			rp.available.Add(-1)
		}
	}
	if newCap == 0 {
		close(rp.resources)
	}
	return nil
}

func (rp *ResourcePool) Capacity() int {
	return int(rp.capacity.Get())
}

func (rp *ResourcePool) IdleTimeout() time.Duration {
	return rp.idleTimeout.Get()
}

func (rp *ResourcePool) Available() int {
	return int(rp.available.Get())
}

func (rp *ResourcePool) Active() int {
	return int(rp.active.Get())
}

func (rp *ResourcePool) InUse() int {
	return int(rp.inUse.Get())
}

func (rp *ResourcePool) WaitCount() int {
	return int(rp.waitCount.Get())
}

func (rp *ResourcePool) WaitTime() time.Duration {
	return rp.waitTime.Get()
}

func (rp *ResourcePool) IdleClosed() int {
	return int(rp.idleClosed.Get())
}
