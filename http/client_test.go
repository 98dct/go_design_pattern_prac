package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestCli1(t *testing.T) {
	// get请求
	resp, err := http.Get("http://localhost:9999/ping")
	if err != nil {
		panic(err)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("get ping", resp.Status, string(res))

	// post请求
	var req Person
	req.Name = "李四"
	req.Age = 55
	input, _ := json.Marshal(req)
	reader := bytes.NewReader(input)
	resp, err = http.Post("http://localhost:9999/ping", "application/json", reader)
	if err != nil {
		panic(err)
	}
	res, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("post ping", resp.Status, string(res))

	// 通用的请求方式-get
	client := http.Client{}
	getReq, _ := http.NewRequestWithContext(context.Background(), http.MethodGet,
		"http://localhost:9999/ping", nil)
	getReq.Header.Set("header1", "get-req")
	resp, err = client.Do(getReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	res, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("get ping", resp.Status, string(res))

	// 通用的请求方式-post
	var reqPost Person
	reqPost.Name = "王五"
	reqPost.Age = 66
	inputPost, _ := json.Marshal(reqPost)
	readerPost := bytes.NewReader(inputPost)
	postReq, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"http://localhost:9999/ping", readerPost)
	postReq.Header.Set("header1", "post-req")
	resp, _ = client.Do(postReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	res, err = io.ReadAll(resp.Body)
	fmt.Println("post ping", resp.Status, string(res))

}
