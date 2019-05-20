package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//http路由
func startServers() {
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/test/", testHandler)

	svr := http.Server{
		Addr:         ":" + "8080",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	assert(svr.ListenAndServe())
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	apiErr(w, "")
	apiOk(w, "")
}

//统一返回结构
func apiOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	assert(enc.Encode(map[string]interface{}{
		"stat": 0,
		"msg":  "请求成功",
		"data": data,
	}))
}

func apiErr(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	assert(enc.Encode(map[string]interface{}{
		"stat": 1,
		"msg":  msg,
	}))
}

//结构体方式
type APIReply struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Reply(w http.ResponseWriter, r *http.Request) (trcid string, reply func()) {
	trcid = UUID(8)
	Dbg(trcid, "%s => %s", r.RemoteAddr, r.URL.Path)
	reply = func() {
		var ar *APIReply
		e := recover()
		switch e.(type) {
		case *APIReply:
			ar = e.(*APIReply)
		case nil:
			ar = &APIReply{Code: 0}
		default:
			ar = &APIReply{Code: 100, Msg: fmt.Sprintf("%v", e)} //generic error
		}
		if ar.Data == nil {
			ar.Data = map[string]string{}
		}
		Dbg(trcid, ar.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintln(w, ar.Error())
	}
	return
}

func (ar *APIReply) Error() string {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(ar)
	return buf.String()
}

func Ret(ar *APIReply) {
	if ar == nil {
		panic(&APIReply{Code: 0})
	}
	panic(ar)
}
