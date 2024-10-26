package taoyi_analysize

import (
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
)

type Struc1 struct {
	s1   []string
	name string
}

type Struc2 struct {
	name string
}

func TestTaoyi(t *testing.T) {

	sm := sync.Map{}
	sm.Store("key1", "val1")
	//v1, ok1 := sm.Load("key1")
	sm.Store("key2", "val2")
	sm.Store("key2", "val3")
	v2, ok2 := sm.Load("key2")
	//fmt.Println(v1, ok1)
	fmt.Println(v2, ok2)
}

func TestPrintStack(t *testing.T) {
	a(99)
}

func a(num int) {
	debug.PrintStack()
}
