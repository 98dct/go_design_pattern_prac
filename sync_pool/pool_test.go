package sync_pool

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestPool1(t *testing.T) {
	pool := sync.Pool{}
	pool.New = func() any {
		return new(bytes.Buffer)
	}

	buffer := pool.Get().(*bytes.Buffer)
	buffer.WriteString("hello dct")
	fmt.Println(buffer.String()) // hello dct
	var res = make([]byte, 10)
	// 只有buffer没有数据时，从buffer中读取的err报io.EOF
	n, err := buffer.Read(res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res[:n]))               // h
	fmt.Println(buffer.Len(), buffer.String()) // 0, ""

	// 回收buffer对象
	buffer.Reset()
	pool.Put(buffer)
}

func TestBuffer(t *testing.T) {
	var b bytes.Buffer
	b.Grow(64)
	b.Write([]byte("abcde"))
	rdbuf := make([]byte, 1)
	n, err := b.Read(rdbuf)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)             // 1
	fmt.Println(b.String())    // bcde
	fmt.Println(string(rdbuf)) // a
}
