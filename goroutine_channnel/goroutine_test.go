package goroutine_channnel

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// 交替打印1azb3c
func Test1(t *testing.T) {

	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)

	ch1 <- struct{}{}
	var i = 1
	var j = 'a'

	go func() {
		for {
			select {
			case <-ch1:
				fmt.Println(i)
				i++
				time.Sleep(time.Second)
				ch2 <- struct{}{}
			}
		}

	}()

	go func() {
		for {
			select {
			case <-ch2:
				fmt.Printf("%c\n", j)
				j = j + 1
				time.Sleep(time.Second)
				ch1 <- struct{}{}
			}
		}
	}()

	for {

	}
}

func cal() {
	for i := 0; i < 1000000; i++ {
		runtime.Gosched()
	}
}

// goroutine切换
func Test2(t *testing.T) {

	runtime.GOMAXPROCS(1)
	currentTime := time.Now()
	fmt.Println(currentTime)
	go cal()
	for i := 0; i < 1000000; i++ {
		runtime.Gosched()
	}
	fmt.Println(time.Now().Sub(currentTime) / 2000000)

}

func TestSwitchTime(t *testing.T) {
	runtime.GOMAXPROCS(1) // 限制为单线程调度

	ch := make(chan struct{})
	start := time.Now()

	// 启动第一个 goroutine
	go func() {
		for i := 0; i < 1000000; i++ {
			ch <- struct{}{} // 发送消息
		}
	}()

	// 在主 goroutine 中接收消息
	for i := 0; i < 1000000; i++ {
		<-ch // 接收消息
	}

	// 计算切换时间
	elapsed := time.Since(start)
	fmt.Printf("Average switch time: %v\n", elapsed/2000000)
}

// channel操作出现异常的几种情况：
// 1.向已经关闭的channel写数据   panic: send on closed channel
// 2.重复关闭channel           panic: close of closed channel
// 3.关闭一个nil channel       panic: close a nil channel
// 4.nil channel 无论是读取还是写入都会阻塞
func TestChannel1(t *testing.T) {

	var ch = make(chan int)
	go func() {
		for {
			// channel关闭后，可以从channel中多次读取，并且每次读取到的数据都是零值
			data, ok := <-ch
			fmt.Println("goroutine 1", data, ok)
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			// channel关闭后，可以从channel中多次读取，并且每次读取到的数据都是零值
			data, ok := <-ch
			fmt.Println("goroutine 2", data, ok)
			time.Sleep(time.Second * 1)
		}
	}()
	close(ch)
	time.Sleep(time.Second * 5)

}
