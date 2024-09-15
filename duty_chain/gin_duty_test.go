package main

import (
	"fmt"
	"testing"
)

type GinHandler func(*Context, map[string]interface{}) (map[string]interface{}, error)

type Context struct {
	index    int
	handlers []GinHandler
}

func NewContext(handlers ...GinHandler) *Context {
	return &Context{
		index:    -1,
		handlers: handlers,
	}
}

const MaxHandlerCnt = 100

func (c *Context) Next(req map[string]interface{}) {
	c.index++
	for c.index < len(c.handlers) && c.index <= MaxHandlerCnt {
		c.handlers[c.index](c, req)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = MaxHandlerCnt
}

var finalHandler GinHandler = func(ctx *Context, req map[string]interface{}) (map[string]interface{}, error) {
	req["final_handler"] = true
	fmt.Printf("final_handler\n")
	return req, nil
}

var middleware1 GinHandler = func(ctx *Context, req map[string]interface{}) (map[string]interface{}, error) {
	req["middleware1_preprocess"] = true
	fmt.Printf("middleware1_preprocess\n")
	ctx.Next(req)
	req["middleware1_postprocess"] = true
	fmt.Printf("middleware1_postprocess\n")
	return req, nil
}

var middleware2 GinHandler = func(ctx *Context, req map[string]interface{}) (map[string]interface{}, error) {
	req["middleware2_preprocess"] = true
	fmt.Printf("middleware2_preprocess\n")
	ctx.Next(req)
	req["middleware2_postprocess"] = true
	fmt.Printf("middleware2_postprocess\n")
	return req, nil
}

var middleware1WithoutAbort GinHandler = func(ctx *Context, req map[string]interface{}) (map[string]interface{}, error) {
	req["middleware1_preprocess"] = true
	fmt.Printf("middleware1_preprocess\n")
	ctx.Abort()
	return req, nil
}

func TestGinHandlerChain(t *testing.T) {
	ctx := NewContext(middleware1WithoutAbort, middleware2, finalHandler)
	params := map[string]interface{}{}
	ctx.Next(params)
	t.Logf("params:%+v", params)
}
