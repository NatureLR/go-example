package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	var t string
	switch t {
	case "1", "2":
		fmt.Println()
	}

	key := flag.String("key", "", "key file")
	cert := flag.String("cert", "", "cert")
	flag.Parse()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	svr := http.Server{
		Addr:         ":8888",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	go func() {
		if *key == "" || *cert == "" {
			fmt.Println("http服务启动成功")
			if err := svr.ListenAndServe(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	// 优雅的关闭
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	<-ctx.Done()

	stop()

	timeoutCTX, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := svr.Shutdown(timeoutCTX); err != nil {
		fmt.Println(err)
	}
}
