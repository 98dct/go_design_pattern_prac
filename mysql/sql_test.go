package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
	// 注册mysql驱动
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserId int64
}

func Test_sql(t *testing.T) {
	db, err := sql.Open("mysql", "username:password@(localhost:3306)/database")
	if err != nil {
		t.Error(err)
		return
	}

	db.SetConnMaxLifetime(time.Second)
	// 执行sql
	ctx := context.Background()
	row := db.QueryRowContext(ctx, "select user_id from user where order by created_at desc limit 1")
	if row.Err() != nil {
		t.Error(row.Err())
		return
	}

	// 解析结果
	var u User
	if err = row.Scan(&u.UserId); err != nil {
		t.Error(err)
		return
	}
}

var ErrNotFound = errors.New("not found")

func Test_error(t *testing.T) {

	// Wrapping ErrNotFound
	err := fmt.Errorf("something went wrong: %w", ErrNotFound)

	// Checking if err contains ErrNotFound
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Error is ErrNotFound")
	} else {
		fmt.Println("Error is not ErrNotFound")
	}
}

func handlerReq(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("handler request done.")
			return
		default:
			fmt.Println("handler request running ,paramenter:", ctx.Value("parameter"))
			time.Sleep(2 * time.Second)
		}
	}

}

func TestValueCtx(t *testing.T) {

	cancelCtx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	cancelCtx2, cancel2 := context.WithCancel(cancelCtx1)
	defer cancel2()
	valueCtx := context.WithValue(cancelCtx2, "parameter", "1")
	go handlerReq(valueCtx)

	time.Sleep(3 * time.Second)
}

func TestValueCtx1(t *testing.T) {

	timeoutCtx1, cancel1 := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancel1()
	timeoutCtx2, cancel2 := context.WithCancel(timeoutCtx1)
	defer cancel2()
	valueCtx := context.WithValue(timeoutCtx2, "parameter", "1")
	go handlerReq(valueCtx)

	time.Sleep(3 * time.Second)
}

func TestValueCtx2(t *testing.T) {

	valueCtx1 := context.WithValue(context.Background(), "aa", "bb")
	cancelCtx, cancel := context.WithCancel(valueCtx1)
	defer cancel()
	valueCtx := context.WithValue(cancelCtx, "parameter", "1")
	go handlerReq(valueCtx)

	time.Sleep(3 * time.Second)
}

func TestCtx3(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		select {
		case <-cancelCtx.Done():
			fmt.Println("结束了！")
			return
		}
	}()
	cancel()

	fmt.Println(cancelCtx.Value("aa"))
	time.Sleep(10 * time.Second)
}

func Test4(t *testing.T) {
	valueCtx1 := context.WithValue(context.Background(), "aa", "bb")
	timeoutCtx, cancel := context.WithTimeout(valueCtx1, 10*time.Second*100000)
	defer cancel()
	go handlerReq(timeoutCtx)

	time.Sleep(3 * time.Second)
}

func Test5(t *testing.T) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	timeoutCtx, _ := context.WithTimeout(cancelCtx, 10*time.Second*100000)
	defer cancelFunc()
	go handlerReq(timeoutCtx)

	time.Sleep(3 * time.Second)
}

func Test6(t *testing.T) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second*100000)
	defer cancel()
	go handlerReq(timeoutCtx)

	time.Sleep(3 * time.Second)
}

func Test7(t *testing.T) {
	timeoutCtx1, _ := context.WithTimeout(context.Background(), 200000*time.Second)
	timeoutCtx2, _ := context.WithTimeout(timeoutCtx1, time.Second*100000)
	go handlerReq(timeoutCtx2)

	time.Sleep(3 * time.Second)
}

func TestEqual(t *testing.T) {

	// reflect.DeepEqual的比较策略
	// 1.基础类型直接比较值
	// 2.切片和数组类型，比较长度和每个值，不同返回false
	// 3.map比较key和value，即使容量不同
	// 4.reflect.DeepEqual不支持比较channel和func()类型,直接返回false, 其中channel可以通过==比较，但是func()完全不能比较
	// 5.结构体逐个字段比较，参数名称和类型必须相等

	c1 := make(chan int, 10)
	c2 := make(chan int)
	c3 := c1
	fmt.Println(c1 == c2) // 比较的是两个channel的地址是否相等
	fmt.Println(reflect.DeepEqual(c1, c2))
	fmt.Println(c1 == c3)

	s1 := make([]struct{}, 10)
	s2 := make([]struct{}, 20)
	var s3 []struct{}
	//s3 := s1
	//fmt.Println(s1 == s2)
	fmt.Println(reflect.DeepEqual(s1, s2))
	fmt.Println(reflect.ValueOf(s1).IsZero())
	fmt.Println(reflect.ValueOf(s3).IsZero())
	//fmt.Println(s1 == s3)

	m1 := make(map[string]struct{}, 10)
	m2 := make(map[string]struct{}, 20)
	var m3 map[string]struct{}
	//m3 := s1
	//fmt.Println(m1 == m2)
	fmt.Println(reflect.DeepEqual(m1, m2))
	fmt.Println(reflect.ValueOf(m1).IsZero())
	fmt.Println(reflect.ValueOf(m3).IsZero())
	//fmt.Println(m1 == m3)

	f1 := func(a string, b int) {}
	f2 := func(c int, d string) {}
	var f3 func(string, int)
	f4 := f1
	f5 := func(e string, f int) {}
	//fmt.Println(f1 == f2)
	fmt.Println(reflect.DeepEqual(f1, f2))
	fmt.Println(reflect.DeepEqual(f1, f4))
	fmt.Println(reflect.DeepEqual(f1, f5))
	fmt.Println(reflect.ValueOf(f3).IsZero())
}

func TestCompare(t *testing.T) {

	// go中可以比较的类型：基本类型、channel、指针类型、接口类型, 只包含可比较字段的结构体
	// 不可比较的类型：map、slice、func()和包含不可比较字段的结构体
	// 可比较的大前提一定是相同类型，否则在编译阶段就报错，reflect.DeepEqual会在运行时报错

	key1 := []string{}
	key2 := map[string]interface{}{}
	key3 := make(chan string, 0) // channel是可比较的类型！ 比较地址
	key4 := func() {}
	var key5 interface{}
	key5 = 5
	var key51 interface{}
	key51 = 5
	var key6 *int
	var key61 *int
	k62 := 62
	k63 := 62
	key61 = &k62
	key6 = &k63
	var key7 struct {
		age  int
		name string
	}
	var key8 struct {
		age  int
		name string
		f    func()
		//c1 chan struct{}
	}
	fmt.Println(reflect.TypeOf(key1).Comparable())
	fmt.Println(reflect.TypeOf(key2).Comparable())
	fmt.Println(reflect.TypeOf(key3).Comparable())
	fmt.Println(reflect.TypeOf(key4).Comparable())
	fmt.Println(reflect.TypeOf(key5).Comparable())
	fmt.Println(reflect.TypeOf(key6).Comparable())
	fmt.Println(reflect.TypeOf(key7).Comparable())
	fmt.Println(reflect.TypeOf(key8).Comparable())

	fmt.Println("---------------------")
	fmt.Println(key51 == key5)
	fmt.Println(key6 == key61)

}

func TestSlice(t *testing.T) {
	s := []int{1, 2, 3}
	fmt.Println(s, len(s), cap(s))
	s = append(s, 4)
	fmt.Println(s, len(s), cap(s))
}
func Test1(t *testing.T) {
	var s []int
	fmt.Println(s, len(s), cap(s))
}

func Test2(t *testing.T) {
	s := make([]int, 0, 10)
	fmt.Println(s, len(s), cap(s))
}

func Test3(t *testing.T) {
	s := make([]int, 0, 4)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 4)
	fmt.Println(s, len(s), cap(s), &s[0]) // [1 2 3 4] 4 4 0xc000016340
	s = append(s, 5)                      // 触发扩容
	fmt.Println(s, len(s), cap(s), &s[0]) // [1 2 3 4 5] 5 8 0xc000012580 首地址发生了变化
}

func Test4_1(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	s1 := s[2:4]                      // s[low:high]  low和high都可以省略
	fmt.Println(s1, len(s1), cap(s1)) // s1:[3,4]   len:2 (high-low) cap:3 (cap-low)
	s2 := s[2:4:4]                    // s[low:high:max]  只有low省略
	fmt.Println(s2, len(s2), cap(s2)) // s2:[3,4]   len:2 (high-low)  cap:2 (max-low)
}

func Test5_1(t *testing.T) {
	// 切片的遍历
	s := []int{1, 2, 3, 4, 5}
	for i, v := range s {
		fmt.Println(i, v, s[i])
		fmt.Println(&i)
		s[i] = i
	}
}

func changeSlice1(s []int) {
	s[0] = 6
}

func changeSlice2(s []int) {
	s = append(s, 6)
	fmt.Println(s, len(s), cap(s), &s[0])
}

func Test6_1(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	changeSlice1(s)
	fmt.Println(s, len(s), cap(s), &s[0])

	changeSlice2(s)
	fmt.Println(s, len(s), cap(s), &s[0])
}

func Test6_2(t *testing.T) {
	s := make([]int, 3, 5)
	s[0] = 0
	s[1] = 1
	s[2] = 2
	changeSlice1(s)
	fmt.Println(s, len(s), cap(s), &s[0])

	changeSlice2(s)
	fmt.Println(s, len(s), cap(s), &s[0])
	fmt.Println(s[2:5])
}

func Test7_1(t *testing.T) {

	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// 切片的深拷贝
	s1 := make([]int, len(s))
	copy(s1, s)
	fmt.Println(s1, len(s1), cap(s1))

	// 切片的删除
	// 删除第6个元素
	s = append(s[:5], s[6:]...)
	fmt.Println(s, len(s), cap(s))

}
