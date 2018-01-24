package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

/*

HTTP/1.1 Augmented BNF
https://tools.ietf.org/html/rfc2616

// non-terminals

HTTP-message = Request | Response

generic-message = start-line *( message-header CRLF ) CRLF [ message-body ]
start-line = Request-Line | Status-Line

Response = Status-Line *(( general-header | response-header | entity-header ) CRLF ) CRLF [ message-body ]

Status-Line = HTTP-Version SP Status-Code SP Reason-Phrase CRLF
HTTP-Version = "HTTP" "/" 1*DIGIT "." 1*DIGIT
Status-Code = "100" | "101" | ... | "504" | "505" | extension-code
extention-code = 3DIGIT
Reason-Phrase = *<TEXT, excluding CR, LF>

ALPHA = UPALPHA | LOALPHA
CRLF = CR LF
LWS = [CRLF] 1*( SP | HT )
TEXT = <any OCTET excepts CTLs, but including LWS>
token = 1*<any CHAR except CTLs or separators>
separators = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | "/" | "[" | "]" | "?" | "=" | "{" | "}" | SP | HT | <">
comment = "(" *( ctext | quoted-pair | comment ) ")"
ctext = <any TEXT excluding "(" and ")">
quoted-string = ( <"> *(qdtext | quoted-pair ) <"> )
qdtext = <any TEXT except <">>
quoted-pair = "\" CHAR

// terminals

OCTET = <any 8-bit sequence of data>
CHAR = <any US-ASCII character (octets 0 - 127)>
UPALPHA = <any US-ASCII uppercase letter "A".."Z">
LOALPHA  = <any US-ASCII lowercase letter "A".."Z">
DIGIT = <any US-ASCII digit "0".."9">
CTL = <any US-ASCII control character (octets 0 - 31) and DEL (127)>
CR = <US-ASCII CR, carriage return (13)>
LF = <US-ASCII LF, linefeed (10)>
SP = <US-ASCII SP, space (32)>
HT = <US-ASCII HT, horizontal-tab (9)>
<"> = <US-ASCII double quote mark (34)>

*/

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
	if i <= 0 {
		return s, "", "", "", errors.New("line separator not found")
	}

	s, s2 := s[:i], s[i+2:]

	splitted := strings.Split(s, " ")

	return s2, splitted[0], splitted[1], strings.Join(splitted[2:], " "), nil
}

func parseHeader(s string) (string, string, string, error) {
	i := strings.Index(s, "\r\n")
	if i <= 0 {
		return s, "", "", errors.New("line separator not found")
	}

	s, s2 := s[:i], s[i+2:]

	splitted := strings.Split(s, ":")

	return s2, strings.TrimSpace(splitted[0]), strings.TrimSpace(strings.Join(splitted[1:], " ")), nil
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
	for {
		response2, headerName, headerValue, err := parseHeader(response)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("headerName: " + headerName)
		fmt.Println("headerValue: " + headerValue)

		response = response2
	}

	fmt.Println(response)

	fmt.Scanln()

}
