package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	var buf bytes.Buffer
	f, err := os.Create("x.tar.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	content, err := json.Marshal(map[string]string{
		"name": "uhub.service.ucloud.cn/library/nginx:1.9.7",
	})
	if err != nil {
		panic(err)
	}

	tw.WriteHeader(&tar.Header{
		Name: "IMAGE",
		Mode: 0644,
		Size: int64(len(content)),
	})

	tw.Write(content)
	tw.Close()
	gw.Close()

	hasher := sha256.New()
	hasher.Write(buf.Bytes())
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Println(hashString)
	if _, err := io.Copy(f, &buf); err != nil {
		fmt.Println(err)
	}
}

func main2() {
	// 字符串内容
	content := "This is the content of the file."

	// 创建 tar.gz 缓冲区
	var buf bytes.Buffer

	// 创建 gzip.Writer
	gzipWriter := gzip.NewWriter(&buf)

	// 创建 tar.Writer
	tarWriter := tar.NewWriter(gzipWriter)

	// 添加文件到 tar 包
	header := &tar.Header{
		Name: "example.txt",
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		fmt.Println("Error writing tar header:", err)
		return
	}
	if _, err := tarWriter.Write([]byte(content)); err != nil {
		fmt.Println("Error writing tar content:", err)
		return
	}

	// 关闭 tar.Writer 和 gzip.Writer
	if err := tarWriter.Close(); err != nil {
		fmt.Println("Error closing tar writer:", err)
		return
	}
	if err := gzipWriter.Close(); err != nil {
		fmt.Println("Error closing gzip writer:", err)
		return
	}

	// 计算 SHA-256 散列
	hasher := sha256.New()
	hasher.Write(buf.Bytes())
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	// 保存 tar.gz 文件
	file, err := os.Create("example.tar.gz")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, &buf); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// 保存 SHA-256 散列到文件
	hashFile, err := os.Create("example.sha256")
	if err != nil {
		fmt.Println("Error creating hash file:", err)
		return
	}
	defer hashFile.Close()

	if _, err := hashFile.WriteString(hashString); err != nil {
		fmt.Println("Error writing hash to file:", err)
		return
	}

	fmt.Println("File 'example.tar.gz' created successfully.")
	fmt.Println("SHA-256 hash saved in 'example.sha256'.")
}
