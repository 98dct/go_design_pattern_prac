package factory

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type Fruit interface {
	Eat()
}

type Orange struct {
	name string
}

func NewOrange(name string) Fruit {
	return &Orange{name: name}
}

func (o *Orange) Eat() {
	fmt.Printf("i am a orange %s, i want to be eaten \n", o.name)
}

type Strawberry struct {
	name string
}

func NewStrawberry(name string) Fruit {
	return &Strawberry{name: name}
}

func (s *Strawberry) Eat() {
	fmt.Printf("i am a Strawberry %s, i want to be eaten \n", s.name)
}

type Cherry struct {
	name string
}

func NewCherry(name string) Fruit {
	return &Cherry{name: name}
}

func (c *Cherry) Eat() {
	fmt.Printf("i am a Cherry %s, i want to be eaten \n", c.name)
}

type FruitFactory struct {
}

func NewFruitFactory() *FruitFactory {
	return &FruitFactory{}
}

func (f *FruitFactory) CreateFruitFactory(typ string) (Fruit, error) {
	src := rand.NewSource(time.Now().UnixNano())
	rander := rand.New(src)
	name := strconv.Itoa(rander.Int())

	switch typ {
	case "orange":
		return NewOrange(name), nil
	case "strawberry":
		return NewStrawberry(name), nil
	case "cherry":
		return NewCherry(name), nil
	default:
		return nil, errors.New("unsupported type: " + typ)
	}

}

type fruitCreator func(name string) Fruit

type FruitFactoryTwo struct {
	creator map[string]fruitCreator
}

func NewFruitFactoryTwo() *FruitFactoryTwo {
	return &FruitFactoryTwo{creator: map[string]fruitCreator{
		"orange":     NewOrange,
		"strawberry": NewStrawberry,
		"cherry":     NewCherry,
	}}
}

func (fTwo *FruitFactoryTwo) CreateFruitFactoryTwo(typ string) (Fruit, error) {
	v, ok := fTwo.creator[typ]
	if !ok {
		return nil, errors.New("unsupported type: " + typ)
	}
	src := rand.NewSource(time.Now().UnixNano())
	rander := rand.New(src)
	name := strconv.Itoa(rander.Int())

	return v(name), nil
}

func TestSimpleFactory(t *testing.T) {
	factoryTwo := NewFruitFactoryTwo()
	orange, _ := factoryTwo.CreateFruitFactoryTwo("orange")
	orange.Eat()

	watermelon, err := factoryTwo.CreateFruitFactoryTwo("watermelon")
	if err != nil {
		fmt.Println(err)
		return
	}
	watermelon.Eat()
}
