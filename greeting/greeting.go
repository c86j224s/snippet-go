package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

/*
type method int

const (
	get method = iota
	post
	option
	head
	delete
	put
)

type resource string

type version string

const (
	http10 version = "HTTP/1.0"
	http11 version = "HTTP/1.1"
	http20 version = "HTTP/2.0"
)

type headerName string

const (
	host headerName = "Host"
)
*/

func parseStatusLine(s string) (string, string, string, string, error) {
	i := strings.Index(s, "\r\n")
	if i == -1 {
		return s, "", "", "", errors.New("line separator not found")
	}

	s, s2 := s[:i], s[i+2:]

	splitted := strings.Split(s, " ")

	return s2, splitted[0], splitted[1], strings.Join(splitted[2:], " "), nil
}

func sendGetIndex(c net.Conn) {
	s := "GET / HTTP/1.1\r\n\r\n"

	n, err := c.Write([]byte(s))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(n)
}

func recvGetIndex(c net.Conn) (response string, err error) {
	data := make([]byte, 8192)

	response = ""

	for {
		n, err := c.Read(data)
		if err != nil {
			break
		}

		if n == 0 {
			break
		}

		response += string(data[:n])

		c.SetReadDeadline(time.Now().Add(3 * time.Second))
	}

	return
}

func main() {
	client, err := net.Dial("tcp", "ncsoft.com:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	sendGetIndex(client)
	time.Sleep(1 * time.Second)
	response, err := recvGetIndex(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	response, version, statuscode, text, err := parseStatusLine(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("version: " + version)
	fmt.Println("statuscode: " + statuscode)
	fmt.Println("text: " + text)
	fmt.Println(response)

	fmt.Scanln()

}
