package mutex

import (
	"golang.org/x/sync/errgroup"
	"sync"
	"testing"
	"time"
)

func TestMutex1(t *testing.T) {

	mu := sync.Mutex{}
	mu.Lock()
	mu.Unlock()

	rwMutex := sync.RWMutex{}
	rwMutex.RLock()   // 加读锁
	rwMutex.RUnlock() // 释放读锁

	rwMutex.Lock()   // 加写锁
	rwMutex.Unlock() // 释放写锁

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Minute * 3)
	}()
	wg.Wait()

	cond := sync.NewCond(&sync.Mutex{})
	cond.Wait()
	cond.Signal()

	// todo
	// 通过mutex，而不是atomic.CompareAndSwapUint32()保证了Do()函数返回后，
	// f()一定调用完成
	once := sync.Once{}
	once.Do(func() {

	})

	//todo
	eg := errgroup.Group{}
	eg.Go(func() error {
		return nil
	})

	eg.Wait()

	// todo
	//weighted := semaphore.NewWeighted(10)
	//weighted.Acquire()

	//todo
	//sg := singleflight.Group{}
	//sg.Do()

	//todo
	//cyclicbarrier

	//ctx, cancel := context.WithCancel(context.Background())

}
