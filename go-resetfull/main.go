package main

import (
	"crypto/tls"
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"k8s.io/apiserver/pkg/server"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

func main() {
	// 创建 API Server 的 Handler 处理器
	handler := server.NewAPIServerHandler(
		"test-server",
		legacyscheme.Codecs,
		func(apiHandler http.Handler) http.Handler {
			return apiHandler
		},
		nil)

	// 注册路由
	testApisV1 := new(restful.WebService).Path("/apis/test/v1")
	{
		testApisV1.Route(testApisV1.GET("hello").To(
			func(req *restful.Request, resp *restful.Response) {
				resp.WriteAsJson(map[string]interface{}{"k": "v"})
			},
		)).Doc("hello endpoint")
	}

	// 路由添加到 GoRestfulContainer
	handler.GoRestfulContainer.Add(testApisV1)

	// 启动监听服务

	tlsConfig := &tls.Config{}
	secureServer := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConfig,
	}

	secureServer.ListenAndServe()

}

// $ curl 127.0.0.1:8080/apis/test/v1/hello
// {
// "k": "v"
// }
