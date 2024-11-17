package main

import (
	"fmt"
	"go-design-pattern-prac/reflect/t1"
	"reflect"
)

func main() {

	//rt := reflect.TypeOf(t1.T1)
	rv := reflect.ValueOf(t1.T1)
	fmt.Println(reflect.Indirect(rv).Type().Name())
}
