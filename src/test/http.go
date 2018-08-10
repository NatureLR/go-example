package main

import (
	"bytes"
	"encoding/json"
	"expvar"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//http相关功能测试 文件
func getCache() {
	url := fmt.Sprintf("http://localhost:3852/shop_parent")
	fmt.Println("getCache:", url)
	hc := http.Client{Timeout: 10 * time.Second}
	resp, err := hc.Get(url)
	assert(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
	}
	var bb bytes.Buffer
	io.Copy(&bb, resp.Body)

	var data map[string]map[string]string

	/*jd := json.NewDecoder(resp.Body)
	err=jd.Decode(&data)
	assert(err)*/

	buf := bb.Bytes()
	e := json.Unmarshal(buf, &data)
	if e != nil {
		fmt.Println(e)
		fmt.Println(len(buf), string(buf))
		return
	}
	return
}

type tableInfo struct {
	Domain    string
	Source    string
	Name      string
	DataCount int
	DataSize  int
	Mode      string
}

type mobiusJob struct {
	Table       tableInfo
	Storage     string
	JobID       int
	Type        string
	Submitted   time.Time
	Completed   time.Time
	Elapsed     float64
	Status      string
	TotalTasks  int
	ActiveTasks int
}

//"http://123.59.135.51:23779/jobs?status=SUCCEEDED"
func getSparkStatus() (idle bool) {
	url := fmt.Sprintf("http://123.59.135.51:23779/jobs?type=VIEW")
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR:", e.(error).Error())
			time.Sleep(time.Minute)
			idle = true
		}
	}()
	type mobiusJob struct {
		Table struct {
			Name   string
			Domain string
		}
		Submitted time.Time
		Elapsed   float64
	}
	url = fmt.Sprintf("http://%s:%s/jobs?type=VIEW")
	hc := http.Client{Timeout: 10 * time.Second}
	resp, err := hc.Get(url)
	assert(err)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(fmt.Errorf(resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	assert(err)
	var sparkStatus []mobiusJob
	assert(json.Unmarshal(body, &sparkStatus))
	for _, s := range sparkStatus {
		if s.Table.Domain == "likeitbi" {
			fmt.Printf("SPARK BUSY: table=%s; started=%v elapsed=%f\n",
				s.Table.Name, s.Submitted, s.Elapsed)
			return false
		}
	}
	return true
}

type HeartBeat time.Time

func (hb *HeartBeat) String() string {
	return `"` + time.Time(*hb).Format(time.RFC3339Nano) + `"`
}

var (
	SELF string
	HB   HeartBeat
)

func http_server() {
	http.HandleFunc("/test/", testHandler)
	svr := http.Server{
		Addr:         ":9658",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	assert(svr.ListenAndServe())
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	a := args(r)
	for k, v := range a {
		fmt.Println(k, v)
	}
}

func args(r *http.Request) []string {
	ps := strings.Split(r.URL.Path, "/")
	for i := len(ps) - 1; i > 0; i-- {
		if ps[i] != "" {
			return ps[2 : i+1]
		}
	}
	return nil
}

func sysInit(p string) {

	expvar.NewString("epoch").Set(time.Now().Format(time.RFC3339))
	expvar.NewString("version").Set(_G_REVS + "." + _G_HASH)
	expvar.NewInt("pid").Set(int64(os.Getpid()))
	svr := http.Server{
		Addr:         ":" + p,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	assert(svr.ListenAndServe())
}
