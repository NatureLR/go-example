package socks

import (
	"flag"
	"fmt"
	"net"
)

func server() {
	l, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := l.Accept()
		go func(c net.Conn) {
			defer c.Close()
			if err != nil {
				fmt.Println(err)
			}

			buf := make([]byte, 10)
			n, err := c.Read(buf)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(buf[:n]))
			c.Write([]byte("dfdfdfdf"))
		}(c)
	}
}

func client() {
	d, err := net.Dial("tcp", "localhost:8083")
	if err != nil {
		fmt.Println(err)
	}
	d.Write([]byte("dddd"))
	buf := make([]byte, 10)
	n, err := d.Read(buf)
	if err != nil {
		fmt.Println("xxx")
	}
	fmt.Println(string(buf[:n]))
}

func main() {
	t := flag.String("t", "", "")
	flag.Parse()
	if *t == "s" {
		server()
	}
	client()
}
