package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/proxy"
)

func main() {
	p, err := proxy.SOCKS5("tcp", "localhost:1086", nil, proxy.Direct)
	if err != nil {
		log.Panicln(err)
	}
	tr := &http.Transport{Dial: p.Dial}
	hc := http.Client{Transport: tr}

	resp, err := hc.Get("http://www.google.com")
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}
