package main

import (
	"context"
	"errors"
	"fmt"
)

/*
适用于按指定顺序串行化执行系列任务的场景
优势：
为使用方屏蔽了链式调用、串行执行的内部细节，使用方可以像使用单个节点一样启动责任链
组装责任链时，由于后继节点实际上是由前置节点通过 next(指针)进行调用的，
因此前置节点可以获得后继节点的执行结果，并进行一轮后处理工作(defer recovery 来执行) 这个后处理的执行切面是遍历模式所不具有的
*/

type RuleHandler func(context.Context, map[string]interface{}) error

var checkTokenRule RuleHandler = func(c context.Context, m map[string]interface{}) error {

	// 校验token是否合法
	token, ok := m["token"].(string)
	if !ok {
		return errors.New("token 不存在！")
	}

	if token != "myToken" {
		return errors.New("token 不正确！")
	}
	return nil
}

var checkAgeRule RuleHandler = func(c context.Context, m map[string]interface{}) error {
	age, _ := m["age"].(int)
	if age < 18 {
		return errors.New("age 不正确！")
	}
	return nil
}

var checkAuthorizedRule RuleHandler = func(c context.Context, m map[string]interface{}) error {
	authori, _ := m["authorized"].(bool)
	if !authori {
		return errors.New("用户未授权！")
	}
	return nil
}

func test1() {
	ruleChain := []RuleHandler{
		checkTokenRule,
		checkAgeRule,
		checkAuthorizedRule,
	}

	ctx := context.Background()
	params := map[string]interface{}{
		"token": "myToken",
		"age":   20,
	}
	for _, handler := range ruleChain {
		err := handler(ctx, params)
		if err != nil {
			// 校验未通过，终止发奖流程
			fmt.Printf("校验未通过，" + err.Error())
			return
		}
	}
}

type RuleChain interface {
	Apply(ctx context.Context, params map[string]interface{}) error
	Next() RuleChain
}

type BaseChain struct {
	next RuleChain
}

func (b *BaseChain) Apply(ctx context.Context, params map[string]interface{}) error {
	panic("no implement")
}

func (b *BaseChain) Next() RuleChain {
	return b.next
}

func (b *BaseChain) ApplyNext(ctx context.Context, params map[string]interface{}) error {
	if b.next != nil {
		return b.next.Apply(ctx, params)
	}
	return nil
}

type CheckTokenRule struct {
	BaseChain
}

func NewCheckTokenRule(next RuleChain) *CheckTokenRule {
	return &CheckTokenRule{BaseChain{next: next}}
}

func (c *CheckTokenRule) Apply(ctx context.Context, params map[string]interface{}) error {
	// 校验token是否合法
	token, ok := params["token"].(string)
	if !ok {
		return errors.New("token 不存在！")
	}

	if token != "myToken" {
		return errors.New("token 不正确！")
	}

	if err := c.ApplyNext(ctx, params); err != nil {
		// err post process
		fmt.Println("check token rule err post process...")
		return err
	}
	fmt.Println("check token rule common post process")
	return nil
}

type CheckAgeRule struct {
	BaseChain
}

func NewCheckAgeRule(next RuleChain) *CheckAgeRule {
	return &CheckAgeRule{BaseChain{next: next}}
}

func (c *CheckAgeRule) Apply(ctx context.Context, params map[string]interface{}) error {
	age, _ := params["age"].(int)
	if age < 18 {
		return errors.New("age 不正确！")
	}

	if err := c.ApplyNext(ctx, params); err != nil {
		// err post process
		fmt.Println("check age rule err post process...")
		return err
	}
	fmt.Println("check age rule common post process")
	return nil
}

type CheckAuthorizedRule struct {
	BaseChain
}

func NewCheckAuthorizedRule(next RuleChain) *CheckAuthorizedRule {
	return &CheckAuthorizedRule{BaseChain{next: next}}
}

func (c *CheckAuthorizedRule) Apply(ctx context.Context, params map[string]interface{}) error {
	authori, _ := params["authorized"].(bool)
	if !authori {
		return errors.New("用户未授权！")
	}
	if err := c.ApplyNext(ctx, params); err != nil {
		// err post process
		fmt.Println("check authorized rule err post process...")
		return err
	}
	fmt.Println("check authorized rule common post process")
	return nil
}

func test2() {
	checkAuthorizedRule := NewCheckAuthorizedRule(nil)
	checkAgeRule := NewCheckAgeRule(checkAuthorizedRule)
	checkTokenRule := NewCheckTokenRule(checkAgeRule)

	if err := checkTokenRule.Apply(context.Background(), map[string]interface{}{
		"token": "myToken",
		"age":   20,
	}); err != nil {
		// 校验未通过，终止发奖流程
		fmt.Println("校验未通过，终止发奖流程, " + err.Error())
		return
	}
}

func main() {

	//test1()
	test2()
}
