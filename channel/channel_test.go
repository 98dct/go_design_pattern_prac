package main

import (
	"fmt"
	"testing"
	"time"
)

// channel test
func TestCh1(t *testing.T) {

	var ch chan int
	fmt.Println(ch, len(ch)) // 输出：nil, 0
	//ch <- 1                     // 向nil channel发送数据会一直阻塞，然后引发死锁 deadlock
	//c1 := <-ch                    // 从nil channel接收数据会一直阻塞，然后引发死锁 deadlock
	//close(ch)                       // 关闭一个nil channel会触发panic: close of nil channel

	// 无缓冲管道，收发goroutine是同步通信
	// 双方同时准备好才能进行发送和接收
	ch = make(chan int)
	go func() {
		c2, ok := <-ch
		fmt.Println("goroutine1", c2, ok) // 1, true
	}()
	ch <- 1
	close(ch)
	c1, ok := <-ch
	fmt.Println(c1, ok) //  0, false
	time.Sleep(time.Second)

	fmt.Println("=======================")

	// 有缓冲channel 收发goroutine之间是异步通信
	// 收发过程是非阻塞的
	ch1 := make(chan string, 2)
	go func() {
		fmt.Println("goroutine2 before", len(ch1), cap(ch1)) // goroutine2 before 1, 2 or 2, 2
		select {
		case c3 := <-ch1:
			fmt.Println("goroutine2", c3)                       // goroutine2 bb
			fmt.Println("goroutine2 after", len(ch1), cap(ch1)) // goroutine2 before 0, 2
		}

	}()
	fmt.Println(len(ch1), cap(ch1)) // 0, 2
	ch1 <- "aa"
	fmt.Println(len(ch1), cap(ch1)) // 1, 2
	ch1 <- "bb"
	fmt.Println(len(ch1), cap(ch1))      // 2, 2
	fmt.Println("main goroutine", <-ch1) // main goroutine aa
	fmt.Println(len(ch1), cap(ch1))      // 1, 2
	time.Sleep(time.Second)

}

// goroutine交替打印1a2b3c4d...
func TestCh2(t *testing.T) {

	var cnt = 1
	var c = 'a'

	ch1 := make(chan struct{})
	ch2 := make(chan struct{})

	go func() {
		for {
			select {
			case <-ch1:
				fmt.Println(cnt)
				cnt++
				time.Sleep(time.Second)
				ch2 <- struct{}{}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ch2:
				fmt.Printf("%c\n", c)
				c++
				time.Sleep(time.Second)
				ch1 <- struct{}{}
			}
		}
	}()

	ch1 <- struct{}{}
	for {

	}
}
