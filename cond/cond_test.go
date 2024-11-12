package cond

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

/**
cond相比于channel的优势：
1.cond与一个locker关联，可以对相关临界区提供保护
2.cond同时支持signal和broadcast，但是channel只支持通知一个消费者
3.cond的broadcast可以重复调用，但是channel关闭后不能open，也就是只能broadcast一次

需要注意cond的wait()方法调用前后需要加锁和解锁，signal和broadcast则不需要！

*/

// 固定有两个容量的队列，没生产一个产品，通知消费者消费
func Test1(t *testing.T) {

	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("remove from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		c.L.Unlock()
		go removeFromQueue(1 * time.Second)
	}
}

// 运动员百米跑步，所有人准备好了，裁判发令，开始跑步
func Test2(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	var ready int
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
			cond.L.Lock()
			ready++
			cond.L.Unlock()

			log.Printf("运动员%d已就绪\n", i)
			// 因为只有一个裁判，所以可以使用signal
			// 如果有多个等待者，需要使用broadcast
			cond.Signal()

		}(i)
	}

	cond.L.Lock()
	for ready != 10 {
		cond.Wait()
	}
	cond.L.Unlock()
	log.Printf("所有运动员都已就绪，准备起跑！")

}
