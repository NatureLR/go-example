package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 验证类型的准入控制
func validate(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		msg := "空内容"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	fmt.Println("k8s请求:", string(body))

	// 请求结构体
	qar := admissionv1.AdmissionReview{}
	_, _, err := serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer().Decode(body, nil, &qar)
	checkErr(err)

	var availableLabels map[string]string
	requiredLabels := "admission"
	var errMsg error

	// 处理逻辑 从请求的结构体判断是是否满足条件
	switch qar.Request.Kind.Kind {
	case "Deployment":
		var deploy appsv1.Deployment
		if err := json.Unmarshal(qar.Request.Object.Raw, &deploy); err != nil {
			log.Println("无法解析格式:", err)
			errMsg = err
		}
		//resourceName, resourceNamespace, objectMeta = deploy.Name, deploy.Namespace, &deploy.ObjectMeta
		availableLabels = deploy.Labels
	case "Service":
		var service corev1.Service
		if err := json.Unmarshal(qar.Request.Object.Raw, &service); err != nil {
			log.Println("无法解析格式:", err)
			errMsg = err
		}
		//resourceName, resourceNamespace, objectMeta = service.Name, service.Namespace, &service.ObjectMeta
		availableLabels = service.Labels
	default:
		msg := fmt.Sprintln("不能处理的类型：", qar.Request.Kind.Kind)
		log.Println(msg)
		errMsg = errors.New(msg)
	}

	var status *apimetav1.Status
	var allowed bool

	if _, ok := availableLabels[requiredLabels]; !ok || errMsg != nil {
		msg := "不符合条件"
		if err != nil {
			msg = fmt.Sprintln(errMsg)
		}
		status = &apimetav1.Status{
			Message: msg,
			Reason:  apimetav1.StatusReason(msg),
			Code:    304,
		}
		allowed = false
	} else {
		status = &apimetav1.Status{
			Message: "通过",
			Reason:  "通过",
			Code:    200,
		}
		allowed = true
	}

	// 返回给k8s的消息
	are := &admissionv1.AdmissionReview{
		TypeMeta: apimetav1.TypeMeta{
			APIVersion: qar.APIVersion,
			Kind:       qar.Kind,
		},
		Response: &admissionv1.AdmissionResponse{
			Allowed: allowed,
			Result:  status,
			UID:     qar.Request.UID,
		},
	}

	resp, err := json.Marshal(are)
	checkErr(err)
	fmt.Println("响应:", string(resp))
	w.WriteHeader(200)
	w.Write(resp)
}

// 修改类型的准入控制
func mutate(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		msg := "空内容"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	fmt.Println("k8s请求:", string(body))

	// 请求结构体
	qar := admissionv1.AdmissionReview{}
	_, _, err := serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer().Decode(body, nil, &qar)
	checkErr(err)

	type patchOperation struct {
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value,omitempty"`
	}

	p := patchOperation{
		Op:    "add",
		Path:  "/metadata/annotations",
		Value: map[string]string{"admission-example.naturelr.cc/status": "test"},
	}
	patch, err := json.Marshal([]patchOperation{p})
	checkErr(err)

	// 返回给k8s的消息
	are := &admissionv1.AdmissionReview{
		TypeMeta: apimetav1.TypeMeta{
			APIVersion: qar.APIVersion,
			Kind:       qar.Kind,
		},
		Response: &admissionv1.AdmissionResponse{
			Allowed: true,
			Patch:   patch,
			PatchType: func() *admissionv1.PatchType {
				pt := admissionv1.PatchTypeJSONPatch
				return &pt
			}(),
			UID: qar.Request.UID,
		},
	}

	resp, err := json.Marshal(are)
	checkErr(err)
	fmt.Println("响应:", string(resp))
	w.WriteHeader(200)
	w.Write(resp)
}

func main() {
	key := flag.String("key", "", "key file")
	cert := flag.String("cert", "", "cert")
	flag.Parse()

	http.HandleFunc("/validate", validate)
	http.HandleFunc("/mutate", mutate)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	svr := http.Server{
		Addr:         ":8080",
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
		fmt.Println("https服务启动成功")
		if err := svr.ListenAndServeTLS(*cert, *key); err != nil {
			log.Fatalln(err)
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
