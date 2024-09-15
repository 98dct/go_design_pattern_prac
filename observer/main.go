package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
*
观察者设计模式适用于多对一的订阅/发布场景
多：多名观察者
一：有一个被观察事物
订阅：观察者时刻关注事物的动态
发布：事物状态发生变化是透明公开的，能正常进入观察者的视线
设计核心：实现观察者与被观察者对象之间的解耦，并将其设计为通用的模块，便于后续的扩展和复用
核心角色：
Event：事物的变更事件，其中topic标识了事物的身份以及变更的类型，val是变更详情，事物发生变更时，将变更情况向eventBus
统一汇报
EventBus：事件总线，位于观察者与事物之间承上启下的代理层，负责维护管理观察者与被关注事物的映射关系，在事物发生变更时，能够及时将
情况同步给每个观察者
Observer：观察者，关注事物动态的角色。需要向EventBus注册完成注册操作，注册时声明自己关心的事件类型
*/

type Event struct {
	EventType string
	Val       interface{}
}

type Observer interface {
	OnChange(ctx context.Context, e *Event) error
}

// 相当于cloud-platform的 AsyncResourceInstance 管理收到事件后的后续操作
type BaseObserver struct {
	ResourceType string // 观察者的资源类型
}

func NewBaseObserver(name string) *BaseObserver {
	return &BaseObserver{ResourceType: name}
}

func (b *BaseObserver) OnChange(c context.Context, e *Event) error {
	fmt.Printf("observer: %s, event key: %s, event val: %s\n", b.ResourceType, e.EventType, e.Val)
	// 做一些资源的变更操作 创建、更新、删除等
	return nil
}

type EventBus interface {
	Subscribe(topic string, o Observer)
	UnSubscribe(topic string, o Observer)
	Publish(ctx context.Context, e *Event)
}

type BaseEventBus struct {
	mux       sync.RWMutex
	observers map[string]map[Observer]struct{} // key 事件类型  val 观察者
}

func NewBaseEventBus() BaseEventBus {
	return BaseEventBus{
		observers: make(map[string]map[Observer]struct{}),
	}
}

func (b *BaseEventBus) Subscribe(eventType string, o Observer) {
	b.mux.Lock()
	defer b.mux.Unlock()
	_, ok := b.observers[eventType]
	if !ok {
		b.observers[eventType] = make(map[Observer]struct{})
	}
	b.observers[eventType][o] = struct{}{}
}

func (b *BaseEventBus) Unsubscribe(eventType string, o Observer) {
	b.mux.Lock()
	defer b.mux.Unlock()
	delete(b.observers[eventType], o)
}

type SyncEventBus struct {
	BaseEventBus
}

func NewSyncEventBus() *SyncEventBus {
	return &SyncEventBus{BaseEventBus: NewBaseEventBus()}
}

func (s *SyncEventBus) Publish(ctx context.Context, e *Event) {

	s.mux.RLock()
	subscribers := s.observers[e.EventType]
	s.mux.RUnlock()

	errs := make(map[Observer]error)
	for sub := range subscribers {
		if err := sub.OnChange(ctx, e); err != nil {
			errs[sub] = err
		}
	}
	s.handleErr(errs)
}

func (b *SyncEventBus) handleErr(errs map[Observer]error) {
	for o, err := range errs {
		fmt.Printf("observer: %v, err: %v", o, err)
	}
}

type ObserverWithErr struct {
	o   Observer
	err error
}

type AsyncEventBus struct {
	BaseEventBus
	errC chan *ObserverWithErr
	ctx  context.Context
	stop context.CancelFunc
}

func NewAsyncEventBus() *AsyncEventBus {
	aBus := AsyncEventBus{
		BaseEventBus: NewBaseEventBus(),
		errC:         make(chan *ObserverWithErr, 1),
	}
	aBus.ctx, aBus.stop = context.WithCancel(context.Background())
	// 处理错误的守护协程
	go aBus.handlerErr()
	return &aBus
}

func (a *AsyncEventBus) handlerErr() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case err := <-a.errC:
			fmt.Printf("observer:%v,err:%v", err.o, err.err)
		}
	}
}

func (a *AsyncEventBus) Stop() {
	a.stop()
}

func (a *AsyncEventBus) Publish(ctx context.Context, e *Event) {
	a.mux.RLock()
	subs := a.observers[e.EventType]
	a.mux.RUnlock()
	for sub := range subs {
		sub := sub
		go func() {
			if err := sub.OnChange(ctx, e); err != nil {
				select {
				case <-a.ctx.Done():
				case a.errC <- &ObserverWithErr{
					o:   sub,
					err: err,
				}:
				}
			}
		}()
	}
}

func main() {
	asyncTest()
}

func asyncTest() {
	observerA := NewBaseObserver("a")
	observerB := NewBaseObserver("b")
	observerC := NewBaseObserver("c")
	observerD := NewBaseObserver("d")

	sbus := NewAsyncEventBus()
	defer sbus.Stop()
	eventType := "order_finish"
	sbus.Subscribe(eventType, observerA)
	sbus.Subscribe(eventType, observerB)
	sbus.Subscribe(eventType, observerC)
	sbus.Subscribe(eventType, observerD)

	sbus.Publish(context.Background(), &Event{
		EventType: eventType,
		Val:       "order_id: xxx",
	})
	<-time.After(time.Second)
}

func syncTest() {
	observerA := NewBaseObserver("a")
	observerB := NewBaseObserver("b")
	observerC := NewBaseObserver("c")
	observerD := NewBaseObserver("d")

	sbus := NewSyncEventBus()
	eventType := "order_finish"
	sbus.Subscribe(eventType, observerA)
	sbus.Subscribe(eventType, observerB)
	sbus.Subscribe(eventType, observerC)
	sbus.Subscribe(eventType, observerD)

	sbus.Publish(context.Background(), &Event{
		EventType: eventType,
		Val:       "order_id: xxx",
	})
}
