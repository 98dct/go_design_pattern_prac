package singleton

import (
	"fmt"
	"sync"
)

/**
单例的用途：有些数据在系统中只应该保存一份数据，例如：配置信息类对象、日志类对象、协程池对象等
*/

type singleton struct {
	value int
}

var instance *singleton

func (s *singleton) Test() {
	fmt.Println(s.value)
}

// v1 懒汉式 不安全
func GetInstanceV1() *singleton {
	if instance == nil {
		instance = new(singleton)
	}
	return instance
}

// v2 懒汉式 安全 效率低
// 每次调用都要加锁
var mu sync.Mutex

func GetInstanceV2() *singleton {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = new(singleton)
	}
	return instance
}

// v3 饿汉式 安全
func init() {
	if instance == nil {
		instance = new(singleton)
	}
}

func GetInstanceV3() *singleton {
	return instance
}

// v4 双重检查单例模式
func GetInstanceV4() *singleton {
	if instance == nil {
		mu.Lock()
		if instance == nil {
			instance = new(singleton)
		}
		mu.Unlock()
	}
	return instance
}

// v5 sync.once 最简单的方式
var once sync.Once

func GetInstanceV5() *singleton {
	once.Do(func() {
		instance = new(singleton)
	})
	return instance
}
