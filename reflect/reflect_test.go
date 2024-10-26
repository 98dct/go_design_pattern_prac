package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

// 根据结构体动态拼接sql字符串
func createQuery(v interface{}) string {
	if reflect.ValueOf(v).Kind() == reflect.Struct {
		// 返回结构体的名字
		// 所在包中定义的类型的名字
		name := reflect.TypeOf(v).Name()
		query := fmt.Sprintf("insert into %s values(", name)
		val := reflect.ValueOf(v)
		for i := 0; i < val.NumField(); i++ {
			switch val.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query = fmt.Sprintf("%s%d", query, val.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s, %d", query, val.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query = fmt.Sprintf("%s\"%s\"", query, val.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, \"%s\"", query, val.Field(i).String())
				}
			}
		}
		query = fmt.Sprintf("%s)", query)
		//fmt.Println(query)
		return query
	}
	return ""
}

type Student struct {
	Name string
	Age  int
}

type Trade struct {
	tradeId int
	Price   int
}

func Test2(t *testing.T) {
	p1 := Student{
		Name: "dct",
		Age:  25,
	}
	fmt.Println(createQuery(p1))

	trade1 := Trade{
		tradeId: 123,
		Price:   456,
	}
	fmt.Println(createQuery(trade1))
}

type reader interface {
	read() string
}

type readStruct struct {
}

func (r readStruct) read() string {
	return ""
}

func Test1(t *testing.T) {
	s1 := Student{}
	s1Name := reflect.TypeOf(s1).Name()
	fmt.Println(s1Name)
	fmt.Println(reflect.TypeOf(s1))        // reflect.Type返回的是实际类型, 例如：reflect.Student
	fmt.Println(reflect.TypeOf(s1).Kind()) // reflect.Kind返回的是基本类型， 例如：reflect.Struct

	fmt.Println(reflect.TypeOf("hello world").Name())

	fmt.Println("===========================")

	// reflect.value的elem()方法的第一次一定要是一个指针调用才行，用于获取value(kind是ptr或者interface)指向的元素的值
	str1 := "aaa"
	fmt.Println(reflect.ValueOf(&str1).Elem().String()) // aaa
	var v1 = 66
	var v2 interface{} = v1
	fmt.Println(reflect.ValueOf(&v2).Kind())               // ptr
	fmt.Println(reflect.ValueOf(&v2).Elem().Kind())        // interface
	fmt.Println(reflect.ValueOf(&v2).Elem().Interface())   // 66
	fmt.Println(reflect.ValueOf(&v2).Elem().Elem().Kind()) // int
	fmt.Println(reflect.ValueOf(&v2).Elem().Elem().Int())  // 66
	bb := "hello world"
	fmt.Println(reflect.ValueOf(bb).Elem().String())

}

// reflect.value的elem方法测试  只能是指针调用elem(), 用于获取指针或者接口指向的值
func Test3(t *testing.T) {
	var z = 123
	var y = &z
	var x interface{} = y
	v := reflect.ValueOf(&x)
	fmt.Println(v.Kind()) // ptr
	vx := v.Elem()
	fmt.Println(vx.Kind()) // interface
	vy := vx.Elem()
	fmt.Println(vy.Kind()) // ptr
	vz := vy.Elem()
	fmt.Println(vz.Kind()) // int
}

// reflect.type的elem()方法测试，只能是指针、数组、切片、map、channel调用，只能用来获取类型
func Test4(t *testing.T) {
	type A = [16]int16
	var c <-chan map[A][]byte
	tc := reflect.TypeOf(c)
	fmt.Println(tc.Kind())    // chan
	fmt.Println(tc.ChanDir()) // <-chan
	tm := tc.Elem()
	ta, tb := tm.Key(), tm.Elem()
	fmt.Println(tm.Kind(), ta.Kind(), tb.Kind()) // map  array  slice
	tx, ty := ta.Elem(), tb.Elem()
	fmt.Println(tx.Kind(), ty.Kind()) // int16  uint8
}

type User struct {
	Id   int
	Name string
	Age  int
}

// 反射获取结构体的字段
func Test5(t *testing.T) {

	u := User{1, "dct", 25}
	getType := reflect.TypeOf(u)
	getValue := reflect.ValueOf(u)
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v \n", field.Name, field.Type, value)
	}
}

func (u *User) RefCallArgs(age int, name string) error {
	return nil
}

// 反射调用结构体的方法
func Test6(t *testing.T) {
	u := &User{
		Id:   1,
		Name: "dct",
		Age:  25,
	}
	ref := reflect.ValueOf(u)
	rv := ref.MethodByName("RefCallArgs")
	fmt.Printf("numIn: %d, numOut: %d, numMethod: %d\n", rv.Type().NumIn(),
		rv.Type().NumOut(), rv.Type().NumMethod())
	args := []reflect.Value{reflect.ValueOf(22), reflect.ValueOf("dct")}
	res := rv.Call(args)
	fmt.Println(res[0].Interface())
}
