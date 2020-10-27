package app

import (
	"fmt"
	"net/http"
	"time"
)

func Run() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	svr := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	err := svr.ListenAndServe()
	if err != nil {
		fmt.Println("启动失败", err)
	}
}
