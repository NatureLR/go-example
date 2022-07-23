//package main
//
//import (
//	"fmt"
//	"log"
//	"net"
//)
//
//func main() {
//	l, err := net.Listen("tcp", ":7070")
//	if err != nil {
//		log.Println(err)
//	}
//	for {
//		n, err := l.Accept()
//		if err != nil {
//			log.Panicln(err)
//		}
//		fmt.Println(n.RemoteAddr())
//		n.Write([]byte("dfdfd"))
//		n.Close()
//	}
//}
package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

func main() {
	service := ":7070"
	listener, err := net.Listen("tcp", service)
	checkError(err)
	go func() {
		for {
			fmt.Println("协程数:", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println(conn.RemoteAddr())
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
