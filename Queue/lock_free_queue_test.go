package Queue

import (
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

type LKQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

func newLKQueue() *LKQueue {
	n := unsafe.Pointer(&node{})
	return &LKQueue{
		head: n,
		tail: n,
	}
}

// 将新元素添加到尾部
func (q *LKQueue) enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) { // 尾还是尾
			if next == nil { // 没有新数据入队
				if cas(&tail.next, next, n) { // 增加到队尾
					cas(&q.tail, tail, n) // 移动队尾
					return
				}
			} else {
				cas(&q.tail, tail, next)
			}
		}
	}
}

// 从头部出队，没有元素，返回nil
func (q *LKQueue) dequeue() interface{} {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) {
			if head == tail {
				if next == nil {
					return nil
				}
				cas(&q.tail, tail, next)
			} else {
				v := next.value
				if cas(&q.head, head, next) {
					return v
				}
			}
		}
	}
}

// 将unsafe.Pointer原子转换为*node类型
func load(p *unsafe.Pointer) *node {
	return (*node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *node) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}

type queue struct {
	mu    sync.Mutex
	store []interface{}
}

func NewQueue() *queue {
	return &queue{store: []interface{}{}}
}

func (q1 *queue) enqueue(v interface{}) {
	q1.mu.Lock()
	defer q1.mu.Unlock()
	q1.store = append(q1.store, v)
}

func (q1 *queue) dequeue() interface{} {
	q1.mu.Lock()
	defer q1.mu.Unlock()
	v := q1.store[0]
	q1.store[0] = nil
	q1.store = q1.store[1:]
	return v
}

var lkq = newLKQueue()
var q = NewQueue()

// 39.23ns serial    284ns parallel
func BenchmarkLockFreeQueue(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lkq.enqueue("aa")
			lkq.dequeue()
		}()
	}
	wg.Wait()
}

// 46.04ns  serial    505ns parallel
func BenchmarkQueue(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q.enqueue("aa")
			q.dequeue()
		}()
	}
	wg.Wait()
}
