package marshal

import (
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"
)

// json序列化后是文本格式的二进制字节切片，包含 { } ，： 占用空间更大
// protobuf序列化后是二进制格式的字节切片
func TestMarshal(t *testing.T) {
	bytes, err := json.Marshal(struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "dct",
		Age:  25,
	})
	fmt.Println(bytes, err)
}

func TestInterface(t *testing.T) {
	var i interface{}
	i = 55
	fmt.Println(unsafe.Sizeof(i))
	i = []string{"aa", "bb"}
	fmt.Println(unsafe.Sizeof(i))
}

type Stringr interface {
	String() string
}

type Binary struct {
	uint64
}

func (i Binary) String() string {
	return "hello world"
}

func TestInterfaceImplement(t *testing.T) {
	a := &Binary{54}
	// 指针可以调用值方法，值不可以调用指针方法，因为值赋值给接口会在堆内存中有值拷贝！如果在方法中修改了对象，
	// 会对原先的对象不生效
	b := Stringr(a)
	b.String()
}
