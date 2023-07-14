package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

type Dialer struct {
	conn *websocket.Conn
}

func NewDialer(r *http.Request, w http.ResponseWriter) (*Dialer, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	d := &Dialer{
		conn: ws,
	}
	return d, nil
}

func (d *Dialer) Close() error {
	return d.conn.Close()
}

func (d *Dialer) Output(stop chan struct{}, w io.Writer) {
	go func() {
		for {
			msgType, msg, err := d.conn.ReadMessage()
			fmt.Println("收到消息：", "类型是:", msgType, "内容是：", string(msg), "错误是：", err)
			if err != nil {
				// 前端关闭
				if msgType == -1 {
					close(stop)
					break
				}
			}
			msg = append(msg, '\n')
			if _, err := w.Write(msg); err != nil {
				break
			}
		}
	}()
}

func (d *Dialer) Input(r io.Reader) {
	go func(r io.Reader) {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if err := d.conn.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
				break
			}
		}
		if s.Err() != nil {
			fmt.Println("scan:", s.Err())
		}
		d.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}(r)
}

func h(w http.ResponseWriter, r *http.Request) {
	ws, err := NewDialer(r, w)
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})

	ws.Input(nil)

	ws.Output(stop, nil)

}

func main() {

	http.HandleFunc("/", h)

	svr := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	svr.ListenAndServe()

}
