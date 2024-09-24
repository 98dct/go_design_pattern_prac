package marshal

import (
	"encoding/json"
	"fmt"
	"testing"
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
