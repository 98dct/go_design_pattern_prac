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

type myConcurrentMap struct {
	mu     sync.RWMutex
	m      map[int]int
	sigMu  sync.RWMutex
	sigM   map[int]*sync.Cond
	condMu sync.Mutex
}

func NewConcurrentMap() *myConcurrentMap {
	return &myConcurrentMap{
		m:    make(map[int]int),
		sigM: make(map[int]*sync.Cond),
	}
}

//var closedCh = make(chan struct{})
//
//func init() {
//	close(closedCh)
//}

func (m *myConcurrentMap) Put(k, v int) {

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

		m.sigMu.Lock()
		_, ok := m.sigM[k]
		if !ok {
			m.sigM[k] = sync.NewCond(&m.condMu)
		}
		m.sigM[k].Broadcast()
		m.sigMu.Unlock()
	}

}

func (m *myConcurrentMap) Get(k int, maxWaitingDuration time.Duration) (int, error) {

Done:
	m.mu.RLock()
	v, ok := m.m[k]
	m.mu.RUnlock()
	if ok {
		return v, nil
	}

	m.sigMu.Lock()
	var cond *sync.Cond
	cond, ok = m.sigM[k]
	if !ok {
		m.sigM[k] = sync.NewCond(&m.condMu)
		cond = m.sigM[k]
	}
	m.sigMu.Unlock()

	for {
		select {
		case <-time.After(maxWaitingDuration * time.Second):
			return 0, errors.New("exceed max wait time!")
		default:
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
			goto Done
		}
	}

}
