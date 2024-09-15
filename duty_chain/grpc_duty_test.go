package main

import (
	"context"
	"fmt"
	"testing"
)

// 执行函数
type Handler func(ctx context.Context, req []string) ([]string, error)

// 拦截器
type Interceptor func(ctx context.Context, req []string, handler Handler) ([]string, error)

// 将拦截器压缩成链
func chainInterceptors(interceptors []Interceptor) Interceptor {
	if len(interceptors) == 0 {
		return nil
	}

	// 返回一个拦截器interceptor类型的闭包执行函数
	return func(ctx context.Context, req []string, handler Handler) ([]string, error) {
		return interceptors[0](ctx, req, getChainHandler(interceptors, 0, handler))
	}
}

// 从interceptor list的index+1位置开始，结合finnal handler，形成一个增强版的handler
func getChainHandler(interceptors []Interceptor, index int, finalHandler Handler) Handler {
	if index == len(interceptors)-1 {
		return finalHandler
	}

	return func(ctx context.Context, req []string) ([]string, error) {
		return interceptors[index+1](ctx, req, getChainHandler(interceptors, index+1, finalHandler))
	}
}

// 声明好final handler
var handler Handler = func(ctx context.Context, req []string) ([]string, error) {

	fmt.Printf("final handler is running, req: %+v\n", req)
	req = append(req, "finnal_handler")
	return req, nil
}

// 声明好拦截器1
var interceptor1 Interceptor = func(ctx context.Context, req []string, handler Handler) ([]string, error) {

	fmt.Println("interceptor1 preprocess...")
	req = append(req, "interceptor1_preprocess")
	resp, err := handler(ctx, req)
	fmt.Println("interceptor1 postprocess")
	resp = append(resp, "interceptor1_postprocess")
	return resp, err

}

// 声明好拦截器2
var interceptor2 Interceptor = func(ctx context.Context, req []string, handler Handler) ([]string, error) {

	fmt.Println("interceptor2 preprocess...")
	req = append(req, "interceptor2_preprocess")
	resp, err := handler(ctx, req)
	fmt.Println("interceptor2 postprocess")
	resp = append(resp, "interceptor2_postprocess")
	return resp, err
}

func Test_interceptor_chain(t *testing.T) {
	chainedInterceptor := chainInterceptors([]Interceptor{
		interceptor1, interceptor2,
	})

	resp, err := chainedInterceptor(context.Background(), nil, handler)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("resp:%+v", resp)
}
