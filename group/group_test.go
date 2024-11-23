package group

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"sync"
	"testing"
	"time"
)

func TestWaitgroup(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 4)
		fmt.Println("hello world")
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 3)
		fmt.Println("北京欢迎你")
	}()
	wg.Wait()
	fmt.Println("finished")
}

// 在没有错误时，errgroup和sync.Waitgroup是一样的效果，只是增加了限制goroutine数量的功能
// 一旦出现了错误会把第一个goroutine出现的错误赋值给err, 并在wait()结束时返回第一个err
// 非context的errgroup
func TestErrGroup(t *testing.T) {
	var eg errgroup.Group
	// 限制同时运行的最大goroutine数量是2
	// 也就是一批次2个goroutine，可以运行多批次
	eg.SetLimit(2)
	eg.Go(func() error {
		time.Sleep(time.Second * 3)
		fmt.Println("hello world")
		return nil
	})

	eg.Go(func() error {
		time.Sleep(time.Second * 4)
		fmt.Println("北京欢迎你")
		return errors.New("出现错误1！")
	})

	eg.Go(func() error {
		time.Sleep(time.Second * 4)
		fmt.Println("我是第三个goroutine")
		return nil
	})

	eg.Go(func() error {
		time.Sleep(time.Second * 5)
		fmt.Println("我是第四个goroutine")
		return errors.New("出现错误2！")
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("run success!")
}

// context的errgroup
func TestContextErrgroup(t *testing.T) {

}
