package safe_map

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSafeMap(t *testing.T) {

	concurrentMap := NewConcurrentMap2()
	concurrentMap.Put(1, 11)
	concurrentMap.Put(1, 11) // 存在时直接更新，不存在时才会创建，才有信号通知机制
	val1, err1 := concurrentMap.Get(1, 5)
	fmt.Println(val1, err1)

	go func() {
		val2, err2 := concurrentMap.Get(2, 5)
		if err2 != nil {
			fmt.Println(err2)
		}
		fmt.Println("goroutine1: " + strconv.Itoa(val2))
	}()

	go func() {
		val2, err2 := concurrentMap.Get(2, 5)
		if err2 != nil {
			fmt.Println(err2)
		}
		fmt.Println("goroutine2: " + strconv.Itoa(val2))
	}()

	go func() {
		val2, err2 := concurrentMap.Get(2, 5)
		if err2 != nil {
			fmt.Println(err2)
		}
		fmt.Println("goroutine3: " + strconv.Itoa(val2))
	}()

	go func() {
		time.Sleep(time.Second * 2)
		concurrentMap.Put(2, 22)
	}()

	time.Sleep(time.Second * 10)
}

// channel中的数据被多个goroutine监听时，只能被一个goroutine读取到
func TestChanRead(t *testing.T) {

	ch := make(chan int, 10)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	ch <- 6
	ch <- 7
	ch <- 8
	ch <- 9
	ch <- 10

	go func() {
		time.Sleep(time.Second)
		for {
			select {
			case v := <-ch:
				fmt.Println("goroutine1 read: " + strconv.Itoa(v))
			}
		}
	}()
	go func() {
		time.Sleep(time.Second)
		for {
			select {
			case v := <-ch:
				fmt.Println("goroutine2 read: " + strconv.Itoa(v))
			}
		}
	}()

	time.Sleep(time.Second * 10)
}

func TestCh1(t *testing.T) {

	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		v, ok := <-ch
		fmt.Println(v, ok)
	}()

	time.Sleep(time.Second)
	close(ch)
	time.Sleep(time.Second * 2)
}

var closedCh = make(chan struct{})

func init() {
	close(closedCh)
}

// select语句不要从一个nil channel中读取数据
// 也不要向一个nil channel中写入数据，否则都会一直阻塞，造成泄露
// 即使后续修改channel指向也没有用
func TestCh2(t *testing.T) {

	var ch chan struct{}
	go func() {

		for ch == nil {
			// 阻塞，不能进行select操作，否则select会一直阻塞
			time.Sleep(time.Millisecond * 500)
		}

		select {
		case <-ch:
			fmt.Println("ch 关闭")
		}
	}()

	time.Sleep(time.Second * 2)
	ch = closedCh
	time.Sleep(time.Second * 2)
}

func Test3Map(t *testing.T) {
	var m map[string]int
	v, ok := m["aa"]
	fmt.Println(v, ok)
}

func Test4Map(t *testing.T) {
	var mp map[string]interface{}
	fmt.Println(mp, mp == nil, len(mp))
	mp = make(map[string]interface{})
	fmt.Println(mp, mp == nil, len(mp))
	mp["aa"] = "aa"
	fmt.Println(mp, len(mp))

	mp = map[string]interface{}{
		"aa": 11,
		"bb": 22,
	}
	fmt.Println(mp, len(mp))

	fmt.Println("=====================")

	var m1 map[int]int
	m1Res, ok := m1[1]
	fmt.Println(m1Res, ok)
	//m1[1] = 1

	// 增加
	m1 = make(map[int]int)
	m1[1] = 1
	fmt.Println(m1, len(m1))
	m1[2] = 2
	fmt.Println(m1, len(m1))

	// 删除
	delete(m1, 1)
	fmt.Println(m1, len(m1))

	// 修改数据
	m1[2] = 3
	fmt.Println(m1, len(m1)) // 输出：[2:3] 1

	fmt.Println("======================")
	// 遍历数据
	m1[4] = 4
	m1[5] = 5
	m1[6] = 6
	for key, val := range m1 {
		fmt.Printf("%d:%d\n", key, val)
	}

	fmt.Println("======================")

	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	p := person{
		Name: "zhangsan",
		Age:  25,
	}

	var m2 map[string]interface{}
	bytes, _ := json.Marshal(p)
	json.Unmarshal(bytes, &m2)
	fmt.Println(m2, len(m2))

	fmt.Println("=========================")

	cc := 0 & 4
	fmt.Printf("%v\n", cc)

	dd := 0 ^ 4
	fmt.Printf("%v\n", dd)
}

func TestSyncMap(t *testing.T) {

	m := sync.Map{}
	m.Store("aa", "bb")
	//m.Store([]string{}, "bb")
	//m.Store(map[string]interface{}{}, "bb")
	//m.Store(func() {}, "bb")
	//ch := make(chan int)
	//m.Store(ch, "bb")
	value, ok := m.Load("aa")
	fmt.Println(value, ok)
	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})

	m.Delete("aa")
	value, ok = m.Load("aa")
	fmt.Println(value, ok)

	closeCh := make(chan struct{})
	go func() {
		select {
		case <-closeCh:
		}
		m.Store("cc", "dd")
	}()

	go func() {
		select {
		case <-closeCh:
		}
		m.Store("ee", "ff")
	}()

	time.Sleep(500 * time.Millisecond)
	close(closeCh)
	time.Sleep(time.Second)
	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})

}
