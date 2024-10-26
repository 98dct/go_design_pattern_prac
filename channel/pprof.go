package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	a := make(chan int)
	for {
		time.Sleep(time.Second)
		go func() {
			<-a
		}()
	}
}
