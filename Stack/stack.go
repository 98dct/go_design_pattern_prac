package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// panic会自动打印堆栈信息
			// 但是被recovery后，会丢失堆栈信息，所以要调用debug.Stack()存下当前的goroutine的堆栈信息
			fmt.Printf("err:%v, stack: %s\n", r, debug.Stack())
			return
		}
	}()
	go func() {
		time.Sleep(10 * time.Second)
	}()
	a(99)
}

func a(num int) {
	//debug.PrintStack()
	//fmt.Printf("%s\n", debug.Stack())
	panic("测试异常")
}
