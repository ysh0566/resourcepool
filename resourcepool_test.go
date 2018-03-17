package resourcepool

import (
	"testing"
	"time"
)

var count, lastId AtomicInt64

type TestResource struct {
	id int64
	closed bool
}

func(tr TestResource) Close() {
	if !tr.closed {
		count.Add(-1)
		tr.closed = true
	}
}

func ResourceFactory() (Resource, error) {
	count.Add(1)
	return TestResource{lastId.Add(1), false}, nil
}

func TestNewResourcePool(t *testing.T) {
	pool := NewResourcePool(ResourceFactory, 10, 10, time.Second)
	var resources [10]Resource
	for i := 0 ; i < 10; i ++ {
		resource, err := pool.Get()
		if err != nil {
			t.Error(err)
		}
		resources[i] = resource
		if int(count.Get()) != i + 1 {
			t.Errorf("want %d, got %d",i+1, count.Get())
		}
		if pool.Available() != 10 - i -1{
			t.Errorf("want %d, got %d",10-i-1, pool.Available())
		}
		if pool.Active() != i + 1{
			t.Errorf("want %d, got %d", i+1, pool.Active())
		}
		if pool.InUse() != i + 1{
			t.Errorf("want %d, got %d", i+1, pool.InUse())
		}
	}
	for i:= 0; i < 10; i++ {
		pool.Put(resources[i])
		if pool.Available() != i + 1{
			t.Errorf("want %d, got %d",i+1 , pool.Available())
		}
		if pool.Active() != 10{
			t.Errorf("want %d, got %d", i+1, pool.Active())
		}
		if pool.InUse() != 10 - i - 1{
			t.Errorf("want %d, got %d", 10 - i - 1, pool.InUse())
		}
	}
	time.Sleep(time.Second*2)
	if pool.Active() != 0 {
		t.Errorf("want 0, got %d", pool.Active())
	}
}

func TestResourcePoolCapacity(t *testing.T) {
	pool := NewResourcePool(ResourceFactory, 10, 20, time.Second)
	pool.SetCapacity(15)
	if pool.Available() != 15{
		t.Errorf("want 15, got %d", pool.Available())
	}
	if len(pool.resources) != 15 {
		t.Errorf("want 15, got %d", len(pool.resources))
	}
}

