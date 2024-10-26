package defer_recovery

import (
	"fmt"
	"testing"
)

func TestD1(t *testing.T) {
	a()
}

func a() {
	defer fmt.Println("defer a")
	b()
}

func b() {
	defer fmt.Println("defer b")
	panic("1111")
}
