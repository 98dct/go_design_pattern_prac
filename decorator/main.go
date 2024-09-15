package main

import (
	"context"
	"fmt"
)

/**
继承强调的是等级制度和子类种类，这部分架构需要一开始就明确好
装饰器模式强调的是“装饰”的过程，而不强调输入与输出，能够动态地为对象增加某种特定的附属能力，
相比于继承模式显得更加灵活，且符合开闭原则
*/

type Food interface {
	Eat() string
	Cost() float32
}

type Rice struct {
}

func NewRice() *Rice {
	return &Rice{}
}

func (r *Rice) Eat() string {
	return "开动了，一碗米饭"
}

func (r *Rice) Cost() float32 {
	return 2.0
}

type Noodle struct {
}

func NewNoodle() *Noodle {
	return &Noodle{}
}

func (n *Noodle) Eat() string {
	return "一碗香喷喷的西红柿鸡蛋面"
}

func (n *Noodle) Cost() float32 {
	return 20.0
}

type Decorator Food

func NewDecorator(f Food) Decorator {
	return f
}

type LaoGanMaDecorator struct {
	decorator Decorator
}

func NewLaoGanMaDecorator(d Decorator) *LaoGanMaDecorator {
	return &LaoGanMaDecorator{decorator: d}
}

func (l *LaoGanMaDecorator) Eat() string {
	return "加一份老干妈" + l.decorator.Eat()
}

func (l *LaoGanMaDecorator) Cost() float32 {
	return 0.5 + l.decorator.Cost()
}

type HamSausageDecorator struct {
	Decorator
}

func NewHamSausageDecorator(d Decorator) *HamSausageDecorator {
	return &HamSausageDecorator{Decorator: d}
}

func (h *HamSausageDecorator) Eat() string {
	return "加一份火腿肠" + h.Decorator.Eat()
}

func (h *HamSausageDecorator) Cost() float32 {
	return 2.0 + h.Decorator.Cost()
}

func main() {

	rice := NewRice()
	//noodle := NewNoodle()
	// 米饭加老干妈
	laoGanMaDecorator := NewLaoGanMaDecorator(rice)
	fmt.Println(laoGanMaDecorator.Eat())
	fmt.Println(laoGanMaDecorator.Cost())
	// 米饭加火腿肠
	hamSausageDecorator := NewHamSausageDecorator(rice)
	fmt.Println(hamSausageDecorator.Eat())
	fmt.Println(hamSausageDecorator.Cost())
	// 米饭加老干妈加火腿肠
	rice_laoganma_ham := NewHamSausageDecorator(laoGanMaDecorator)
	fmt.Println(rice_laoganma_ham.Eat())
	fmt.Println(rice_laoganma_ham.Cost())
}

// 闭包实现函数增强
type HandlerFunc func(ctx context.Context, param map[string]interface{}) error

func DecoratorFunc(fn HandlerFunc) HandlerFunc {
	return func(ctx context.Context, param map[string]interface{}) error {
		fmt.Println("preprocess ...")
		fn(ctx, param)
		fmt.Println("postprocess...")
		return nil
	}
}
