package main

import (
	"fmt"
	"testing"
	"time"
)

type Mutex struct {
	ch chan struct{}
}

func NewLock() *Mutex {
	mu := &Mutex{ch: make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

func (m *Mutex) Lock() {
	<-m.ch
}

func (m *Mutex) UnLock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-m.ch:
		return true
	case <-timer.C:
	}
	return false
}

func (m *Mutex) IsLock() bool {
	return len(m.ch) == 0
}

func TestChannelLock(t *testing.T) {
	m := NewLock()
	ok := m.TryLock()
	fmt.Println(ok)
	ok = m.TryLock()
	fmt.Println(ok)
}
