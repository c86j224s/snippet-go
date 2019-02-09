package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
)

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test")
	})
	http.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		res := make(chan string)
		go func(res chan<- string) {
			cli := redis.NewClient(&redis.Options{
				Addr:     "127.0.0.1:6379",
				Password: "",
				DB:       0,
			})
			result, err := cli.Incr("counter").Result()
			res <- fmt.Sprintf("redis incr %v %v", result, err)
		}(res)
		fmt.Fprintln(w, <-res)
	})
	http.ListenAndServe("0.0.0.0:5959", nil)
}
