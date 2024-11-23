package semaphore

import (
	"context"
	"errors"
	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"
)

// cyclicbarrier: 用于固定数量的 goroutine 等待同一个执行点”的场景中，
// 而且在放行 goroutine 之后，CyclicBarrier 可以重复利用
// sync.waitGroup：适用于一组goroutine同时完成，重复利用性比较差

// 题目：
// 有一个名叫大自然的搬运工的工厂，生产一种叫做一氧化二氢的神秘液体。
// 这种液体的分子是由一个氧原子和两个氢原子组成的，也就是水。
// 这个工厂有多条生产线，每条生产线负责生产氧原子或者是氢原子，每条生产线由一个 goroutine 负责。
// 这些生产线会通过一个栅栏，只有一个氧原子生产线和两个氢原子生产线都准备好，才能生成出一个水分子，
// 否则所有的生产线都会处于等待状态。也就是说，一个水分子必须由三个不同的生产线提供原子，而且水分子是一个一个按照顺序产生的，
// 每生产一个水分子，就会打印出 HHO、HOH、OHH 三种形式的其中一种。HHH、OOH、OHO、HOO、OOO 都是不允许的。
// 生产线中氢原子的生产线为 2N 条，氧原子的生产线为 N 条。
type H2O struct {
	semaH *semaphore.Weighted
	semaO *semaphore.Weighted
	b     cyclicbarrier.CyclicBarrier
}

func NewH2O() *H2O {
	return &H2O{
		semaH: semaphore.NewWeighted(2),
		semaO: semaphore.NewWeighted(1),
		b:     cyclicbarrier.New(3),
	}
}

func (m *H2O) productH(saveH func()) {
	m.semaH.Acquire(context.Background(), 1)
	saveH()
	m.b.Await(context.Background())
	m.semaH.Release(1)
}

func (m *H2O) productO(saveO func()) {
	m.semaO.Acquire(context.Background(), 1)
	saveO()
	m.b.Await(context.Background())
	m.semaO.Release(1)
}

// 用信号量保证生产的顺序正确性，只能一次生产两个H,一个O
// 用循环栅栏保证一次生产三个原子即一个水分子之后，才能生产下一个水分子，且可以循环生产
func Test2(t *testing.T) {
	var ch chan string

	saveH := func() {
		ch <- "H"
	}

	saveO := func() {
		ch <- "O"
	}

	var N = 100
	ch = make(chan string, N*3)
	h2O := NewH2O()

	var wg sync.WaitGroup
	wg.Add(N * 3)
	for i := 0; i < N*2; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2O.productH(saveH)
			wg.Done()
		}()
	}

	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2O.productO(saveO)
			wg.Done()
		}()
	}

	wg.Wait()
	if len(ch) != N*3 {
		panic(errors.New("len(ch) != 300"))
	}

	var s = make([]string, 3)
	for i := 0; i < N; i++ {
		s[0] = <-ch
		s[1] = <-ch
		s[2] = <-ch
		sort.Strings(s)
		water := s[0] + s[1] + s[2]
		if water != "HHO" {
			t.Fatalf("expect a water molecule but got %s", water)
		}
	}

}
