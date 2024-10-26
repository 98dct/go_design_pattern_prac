package goroutine_channnel

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 总任务数量
const totalTasks = 100

var producerCnt = 20
var consumerCnt = 20 // 增加消费者数量明显提升程序执行效率！
var batchSize = totalTasks / producerCnt

// 生产者函数，向 channel 中发送任务数据
func producer(taskChan chan<- int, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < batchSize; i++ {
		task := rand.Intn(100) // 生成随机任务数据
		fmt.Printf("Producer %d produced: %d\n", id, task)
		taskChan <- task                   // 将任务发送到 channel 中
		time.Sleep(200 * time.Millisecond) // 模拟生产间隔
	}
}

// 消费者函数，从 channel 中接收任务并处理
func consumer(taskChan <-chan int, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan { // 从 channel 中接收任务
		fmt.Printf("Consumer %d consumed: %d\n", id, task)
		time.Sleep(200 * time.Millisecond) // 模拟处理任务的时间
	}
}

func TestPC(t *testing.T) {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子

	taskChan := make(chan int, 5) // 创建一个带缓冲的 channel
	var producerWg sync.WaitGroup // 等待生产者完成的 WaitGroup
	var consumerWg sync.WaitGroup // 等待消费者完成的 WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancel()

	st := time.Now()
	var cnt int64 = 1
	go func() {
		ticker := time.NewTicker(time.Millisecond * 50)

		for {
			select {
			case <-ticker.C:
				fmt.Printf("第%d次监听channel， channel的size是：%d\n", cnt, len(taskChan))
				atomic.AddInt64(&cnt, 1)
			case <-ctx.Done():
				fmt.Println("后台channel监控协程超时退出！")
				return
			}
		}

	}()

	// 启动 producerCnt 个生产者
	for i := 1; i <= producerCnt; i++ {
		producerWg.Add(1)
		go producer(taskChan, i, &producerWg)
	}

	// 启动 consumerCnt 个消费者
	for i := 1; i <= consumerCnt; i++ {
		consumerWg.Add(1)
		go consumer(taskChan, i, &consumerWg)
	}

	// 等待所有生产者完成任务
	producerWg.Wait()
	close(taskChan) // 所有生产者完成后关闭 channel

	// 等待所有消费者完成任务
	consumerWg.Wait()

	fmt.Println("耗时：", time.Since(st))
	fmt.Println("All tasks processed.")
}

// 测试结果展示
// 任务数    生产者数量    消费者数量   channel容量    耗时
//  100        4           4         1,5,10    5s
//  100        5           4         1,5,10    5s
//  100        4           5         1,5,10    5s
//  100        5           5         1,5,10,   4s
//  100        5           10        1,5,10   4s
//  100        10          5         1,5,10    4s
//  100        10          10        1,5,10    2s
//  100        20          20        1,5,10    1s
