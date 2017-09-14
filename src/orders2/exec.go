package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execzip() {
	fmt.Println("打包zip并转码为gb18030")
	cmd := exec.Command("find", "data", "-type", "d", "-exec", "mkdir", "-p", "utf_data/{}", ";")
	var out bytes.Buffer
	cmd.Stdout = &out
	assert(cmd.Run())

	cmd = exec.Command("find", "data", "-type", "f", "-exec", "iconv", "-f", "utf-8", "-t", "gb18030", "{}", "-o", "utf_data/{}", ";")
	cmd.Stdout = &out
	assert(cmd.Run())

	cmd = exec.Command("zip", "-q", "-r", "data.zip", "utf_data/data")
	cmd.Stdout = &out
	assert(cmd.Run())

	cmd = exec.Command("rm", "-rf", "data")
	cmd.Stdout = &out
	assert(cmd.Run())

}
