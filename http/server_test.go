package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSer1(t *testing.T) {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI, r.Method)
		switch r.Method {
		case http.MethodGet:
			fmt.Println(r.Header.Get("header1"))
			w.Write([]byte("pong get"))
		case http.MethodPost:
			fmt.Println(r.Header.Get("header1"))
			bytes, _ := io.ReadAll(r.Body)
			var req Person
			err := json.Unmarshal(bytes, &req)
			if err != nil {
				w.Write([]byte("json.Unmarshal req failed"))
				return
			}
			fmt.Println(req)
			res, _ := json.Marshal(req)
			w.Write([]byte("pong post " + string(res)))
		}
	})

	http.ListenAndServe(":9999", nil)
}

func TestSer2(t *testing.T) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI, r.Method)
		switch r.Method {
		case http.MethodGet:
			fmt.Println(r.Header.Get("header1"))
			w.Write([]byte("pong1 get"))
		case http.MethodPost:
			fmt.Println(r.Header.Get("header1"))
			bytes, _ := io.ReadAll(r.Body)
			var req Person
			err := json.Unmarshal(bytes, &req)
			if err != nil {
				w.Write([]byte("json.Unmarshal req failed"))
				return
			}
			fmt.Println(req)
			res, _ := json.Marshal(req)
			w.Write([]byte("pong1 post " + string(res)))
		}
	})
	serveMux.Handle("/ping2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong2"))
	}))
	server := http.Server{
		Addr:              ":9998",
		Handler:           serveMux,
		ReadTimeout:       3,
		ReadHeaderTimeout: 3,
		WriteTimeout:      3,
		IdleTimeout:       3,
	}
	server.ListenAndServe()
}
