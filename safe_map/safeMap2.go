package safe_map

import (
	"errors"
	"sync"
	"time"
)

/**
1.面向高并发
2.只存在插入和查询操作
3.查询时，key存在，直接返回val;key不存在阻塞等待key ,val被放入后，返回val,等待指定时长未放入，返回超时错误
*/

type myConcurrentMap2 struct {
	mu    sync.Mutex
	m     map[int]int
	sigMu sync.RWMutex
	sigM  map[int]chan struct{}
}

func NewConcurrentMap2() *myConcurrentMap2 {
	return &myConcurrentMap2{
		m:    make(map[int]int),
		sigM: make(map[int]chan struct{}),
	}
}

type MyChan struct {
	ch chan struct{}
	sync.Once
}

func NewMyChan() *MyChan {
	return &MyChan{
		ch: make(chan struct{}),
	}
}

// 只会关闭一次
func (m *MyChan) Close() {
	m.Once.Do(func() {
		close(m.ch)
	})
}

func (m *myConcurrentMap2) Put(k, v int) {

	// 判断K存不存在
	m.mu.Lock()
	if _, ok := m.m[k]; ok {
		// 存在更新
		m.m[k] = v
		m.mu.Unlock()
	} else {
		// 不存在创建, 发送创建k的通知
		m.m[k] = v
		m.mu.Unlock()

		m.sigMu.RLock()
		defer m.sigMu.RUnlock()
		ch, ok := m.sigM[k]
		if !ok {
			// 直接返回，没有人关心这个key，不用发通知
			return
		}

		select {
		case <-ch:
			// 已经关闭过
			return
		default:
			// 通过关闭会通知所有下游的接受者，也不会阻塞发送者
			close(ch)
		}

	}

}

func (m *myConcurrentMap2) Get(k int, maxWaitingDuration time.Duration) (int, error) {

Done:
	m.mu.Lock()
	v, ok := m.m[k]
	m.mu.Unlock()
	if ok {
		return v, nil
	}

	m.sigMu.Lock()
	ch, ok := m.sigM[k]
	if !ok {
		ch = make(chan struct{})
		m.sigM[k] = ch
	}
	m.sigMu.Unlock()

	for {
		select {
		case <-time.After(maxWaitingDuration * time.Second):
			return 0, errors.New("exceed max wait time!")
		case <-ch:
			goto Done
		}
	}

}
