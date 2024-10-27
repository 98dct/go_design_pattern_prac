package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewLevel(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	level := 0
	for rand.Intn(2) > 0 {
		level++
	}
	fmt.Println(level)
}

// level
// 3
// 0
// 1
// 0
// 1
// 1
// 0
// 0
// 1
// 1
// 0
// 8
// 0
// 0
// 0
// 0
