package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
	"go-design-pattern-prac/middleware"
	"runtime"
	"sync"
	"time"
)

func handler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "pong"})
}

func main() {
	//instance := singleton.GetInstanceV1()
	//instance.Test()

	// 前缀树node的几种类型
	// static: 该node匹配的是普通路径类型
	// root: 根节点
	// param: 路径参数节点 单段动态参数 :id
	// catchAll: 匹配剩余路径的所有部分，类似*action
	// wildChild: true代表该节点的子节点有通配符，例如 :id、*action
	engine := gin.Default()
	routes := engine.Group("/aa").Use(middleware.ApiAuth())
	routes.POST("/pong", handler)
	routes.POST("/hello/aa", handler)
	routes.POST("/hello/bb", handler)
	routes.POST("/world1/*action", handler)
	routes.POST("/world3/*name", handler)
	routes.POST("/world2/:id", handler)
	engine.Run(":8888")
	//go test1()
	//for {
	//
	//}
}

func test2() {
	fmt.Println("test2方法开始执行")
	panic("发生错误！")
}

func test1() {
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println(r)
	//	}
	//}()
	fmt.Println("test1方法开始执行")
	go test2()
	fmt.Println("test1方法执行结束")
	time.Sleep(time.Second * 2)
}

func test3() {
	runtime.Gosched()
	pool := sync.Pool{}
	newPool, err := ants.NewPool()
	newPool.Submit()

}
