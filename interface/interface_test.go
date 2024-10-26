package _interface

import (
	"errors"
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	// 接口的比较
	var e1 error
	var e2 = errors.New("hello world")
	var e3 error
	fmt.Println(e1 == e2) // 只有一个接口为false，比较结果即为false
	fmt.Println(e1 == e3) // nil接口比较都为true

	fmt.Println("===========================")

	var e4 = errors.New("hello1")
	var e5 = errors.New("hello2")
	var e6 = errors.New("hello1")
	e7 := e6
	fmt.Println(e4 == e5)
	fmt.Println(e4 == e6)
	fmt.Println(e7 == e6)
	var p1, p2 wr
	p1 = person{
		name: "dct",
		age:  25,
	}
	p2 = person{
		name: "dct",
		age:  25,
	}
	fmt.Println(p1 == p2)

}

type wr interface {
	writ()
}

type person struct {
	name string
	age  int
}

func (p person) writ() {

}

func tt(p person) {

}
