package main

import (
	"fmt"
	"net/http"
	"time"

	redis "github.com/go-redis/redis"
)

func handleHello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("hello"))
}

func handleMemo(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		handleGetMemo(resp, req)
	case "POST":
		handlePostMemo(resp, req)
	case "DELETE":
		handleDeleteMemo(resp, req)
	default:
		http.Error(resp, "not supported", 404)
	}
}

func handleGetMemo(resp http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")

	res, err := redis.NewClient(
		&redis.Options{Addr: "127.0.0.1:6379"},
	).Get(title).Result()

	resp.Write([]byte(fmt.Sprintf("GET [%s] => [%s] : %#v", title, res, err)))
}

func handlePostMemo(resp http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	content := req.FormValue("content")

	_, err := redis.NewClient(
		&redis.Options{Addr: "127.0.0.1:6379"},
	).Set(
		title,
		content,
		time.Duration(3600)*time.Second,
	).Result()
	if err != nil {
		http.Error(resp, fmt.Sprintf("POST [%s] failed : %s", title, err.Error()), 400)
		return
	}

	resp.Write([]byte(fmt.Sprintf("POST [%s] succeed", title)))
}

func handleDeleteMemo(resp http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")

	_, err := redis.NewClient(
		&redis.Options{Addr: "127.0.0.1:6379"},
	).Del(title).Result()
	if err != nil {
		http.Error(resp, fmt.Sprintf("DELETE [%s] failed : %s", title, err.Error()), 400)
	}

	resp.Write([]byte(fmt.Sprintf("DELETE [%s] succeed", title)))
}

func main() {
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/memo", handleMemo)

	static := http.FileServer(http.Dir("./static"))
	http.Handle("/", http.StripPrefix("/", static))

	http.ListenAndServe("127.0.0.1:5959", nil)
}
