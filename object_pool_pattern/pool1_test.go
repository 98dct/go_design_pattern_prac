package object_pool_pattern

import (
	"fmt"
	"sync"
	"testing"
)

type Pool struct {
	mu        sync.Mutex
	Inuse     []interface{}
	Available []interface{}
	new       func() interface{}
}

func NewPool(new func() interface{}) *Pool {
	return &Pool{new: new}
}

func (p *Pool) Acquire() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	var obj interface{}
	if len(p.Available) != 0 {
		obj = p.Available[0]
		p.Available = append(p.Available[:0], p.Available[1:]...)
		p.Inuse = append(p.Inuse, obj)
	} else {
		obj = p.new()
		p.Inuse = append(p.Inuse, obj)
	}
	return obj
}

func (p *Pool) Release(obj interface{}) {
	p.mu.Lock()
	p.mu.Unlock()
	p.Available = append(p.Available, obj)
	for i, v := range p.Inuse {
		if v == obj {
			p.Inuse = append(p.Inuse[:i], p.Inuse[i+1:]...)
			break
		}
	}
}

func TestA(t *testing.T) {
	num := func() interface{} {
		return 10.0
	}

	pool := NewPool(num)
	obj := pool.Acquire()
	fmt.Println(pool.Inuse)
	fmt.Println(pool.Available)

	pool.Release(obj)
	fmt.Println(pool.Inuse)
	fmt.Println(pool.Available)
}
