package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-redis/redis"
)

func connectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func handler(conn net.Conn) {
	data := make([]byte, 4096)

	for {
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = conn.Write(data[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func ping(conn net.Conn) {
	pingmsg := make(chan string)

	go func() {
		for {
			pingmsg <- "ping"
			time.Sleep(3000 * time.Millisecond)
		}
	}()

	for {
		msg := <-pingmsg
		conn.Write([]byte(msg))
	}
}

func betterping(conn net.Conn) {
	ticker := time.NewTicker(3000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn.Write([]byte("ping"))
		}
	}

	// more better example : https://jacking75.github.io/go_time_scheduler/
}

func main() {
	fmt.Println("hello go!")

	//connectRedis()

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue // ??
		}
		defer conn.Close()

		go handler(conn)
		//go ping(conn)
		go betterping(conn)
	}
}
