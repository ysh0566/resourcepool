[![](https://travis-ci.com/ysh0566/resourcepool.svg?branch=master)](https://travis-ci.com/ysh0566/resourcepool)
```
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

func main(){
    pool := NewResourcePool(ResourceFactory, 10, 10, time.Second)
}
```
