package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
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

func TestCh3(t *testing.T) {
	chanOwner := func() <-chan int {
		results := make(chan int)
		go func() {
			defer close(results)
			for i := 0; i < 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}
	results := chanOwner()
	consumer(results)
}

// 父goroutine通过done channel控制子goroutine的退出，避免泄露
func TestCh4(t *testing.T) {
	doWork := func(done <-chan struct{}, strings <-chan string) <-chan struct{} {
		terminated := make(chan struct{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan struct{})
	terminated := doWork(done, nil)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()
	<-terminated
	fmt.Println("Done.")
}

// 父goroutine通过done channel控制子goroutine的退出，避免泄露
// 父收子发
func TestCh5(t *testing.T) {
	newRandStream := func(done <-chan struct{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan struct{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	time.Sleep(time.Second)
}

// or-channel
// 多个done channel合并成一个done channel, 多个done channel中有一个channel关闭时，
// 关闭整个done channel
var or func(channels ...<-chan struct{}) <-chan struct{}

func TestCh6(t *testing.T) {
	or = func(channels ...<-chan struct{}) <-chan struct{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan struct{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):

				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan struct{} {
		c := make(chan struct{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

type Result struct {
	Error    error
	Response *http.Response
}

// 子goroutine中出现的错误应该传递到父 goroutine中去处理
func TestCh7(t *testing.T) {
	checkStatus := func(done <-chan struct{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Response: resp, Error: err}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan struct{})
	defer close(done)
	urls := []string{"https://www.baidu.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error:%v\n", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

// 子goroutine中出现的错误应该传递到父 goroutine中去处理，当出现太多错误时，及时终止
func TestCh8(t *testing.T) {
	checkStatus := func(done <-chan struct{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Response: resp, Error: err}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan struct{})
	defer close(done)

	errCount := 0
	urls := []string{"a", "https://www.baidu.com", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error:%v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Printf("too many errors, breaking!\n")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

// pipeline 流水线实践
func TestCh9(t *testing.T) {

	generator := func(done <-chan struct{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan struct{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan struct{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan struct{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	for i := range pipeline {
		fmt.Println(i)
	}

}

// pipeline常用的模式
func TestCh10(t *testing.T) {
	repeat := func(done chan struct{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan struct{})
	defer close(done)
	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
}

// 扇入模式
func TestFanIn1(t *testing.T) {
	fanIn := func(done <-chan struct{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})
		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}
	done := make(chan struct{})
	ins := make([]<-chan interface{}, 0)
	fanIn(done, ins...)
}

var fanIn func(done chan struct{}, channels ...<-chan interface{}) <-chan interface{}

func TestFanIn2(t *testing.T) {

	mergeTwo := func(chan1 <-chan interface{}, chan2 <-chan interface{}) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			for chan1 != nil || chan2 != nil {
				select {
				case v, ok := <-chan1:
					if !ok {
						chan1 = nil
						continue
					}
					c <- v
				case v, ok := <-chan2:
					if !ok {
						chan2 = nil
						continue
					}
					c <- v
				}
			}
		}()
		return c
	}

	fanIn = func(done chan struct{}, channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			c := make(chan interface{})
			close(c)
			return c
		case 1:
			return channels[0]
		case 2:
			return mergeTwo(channels[0], channels[1])
		default:
			m := len(channels) / 2
			return mergeTwo(fanIn(done, channels[:m]...), fanIn(done, channels[m:]...))
		}
	}
}

func TestFanout(t *testing.T) {
	fanOut := func(ch <-chan interface{}, out []chan interface{}, async bool) {
		go func() {
			defer func() {
				for i := 0; i < len(out); i++ {
					close(out[i])
				}
			}()

			for v := range ch {
				v := v
				for i := 0; i < len(out); i++ {
					i := i
					if async {
						go func() {
							out[i] <- v
						}()
					} else {
						out[i] <- v
					}
				}
			}

		}()
	}
	in := make(chan interface{})
	outs := make([]chan interface{}, 0)
	fanOut(in, outs, false)
}
