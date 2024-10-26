package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func handler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "pong"})
}

func main() {
	//instance := singleton.GetInstanceV1()
	//instance.Test()

	// 前缀树node的几种类型
	// static: 该node匹配的是普通路径类型    0
	// root: 根节点                       1
	// param: 路径参数节点 单段动态参数 :id    2
	// catchAll: 匹配剩余路径的所有部分，类似*action  3
	// wildChild: true代表该节点的子节点有通配符，例如 :id、*action
	//gin.SetMode(gin.ReleaseMode)
	http.HandleFunc("/bb", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("你好：bb")
		w.Write([]byte("你好：bb"))
	})
	//http.ListenAndServe("localhost:8889", nil)
	client := http.Client{}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://loclahost:8888", strings.NewReader("你好 10-07！"))
	client.Do(req)
	//engine := gin.Default()
	//routes := engine.Group("/aa").Use(middleware.ApiAuth())
	//routes.POST("/pong", handler)
	//routes.POST("/hello/aa", handler)
	//routes.POST("/hello/bb", handler)
	//routes.POST("/world1/*action", handler)
	//routes.POST("/world3/*name", handler)
	//routes.POST("/world2/:id", handler)
	//engine.POST("/aa", handler)
	//engine.Run(":8888")
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
	//runtime.Gosched()
	//pool := sync.Pool{}
	//newPool, _ := ants.NewPool()
	//newPool.Submit()

}

func test4() {
	ctxTimeout, timeCancelFunc := context.WithTimeout(context.Background(), time.Second*10000)
	timeCancelFunc()
	ctxCancel, _ := context.WithCancel(ctxTimeout)
	fmt.Println(ctxCancel)

}
