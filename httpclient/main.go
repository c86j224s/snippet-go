package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Phone struct {
	Mobile string `json:"mobile"`
	Home   string `json:"home"`
	Office string `json:"office"`
}

type Contact struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
	Phone   Phone  `json:"phone"`
}

type Root struct {
	Contacts []Contact `json:"contacts"`
}

func main() {
	var host, port, uri, function string
	flag.StringVar(&host, "host", "api.androidhive.info", "server host")
	flag.StringVar(&port, "port", "443", "server port")
	flag.StringVar(&uri, "uri", "contacts", "uri part")
	flag.StringVar(&function, "function", "defaultfunction", "function to use")
	flag.Parse()

	url := fmt.Sprint("https://", host, "/", uri)
	fmt.Println(url)

	var resp *http.Response

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("GET error : ", err)
		fmt.Println("status : ", resp.Status)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("body read error : ", err)
		return
	}

	var root Root
	if err := json.Unmarshal(body, &root); err != nil {
		fmt.Println("unmarchal error : ", err)
	}

	marshaled, err := json.Marshal(root)
	if err != nil {
		fmt.Println("marshal error : ", err)
	}

	fmt.Println(string(marshaled[:]))

}
