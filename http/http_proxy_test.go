package http

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	// 创建代理服务器
	proxy := &http.Server{
		Addr: ":28081", // 代理服务器的端口号
	}

	// 设置代理服务器的请求处理函数
	http.HandleFunc("/", handleRequest)

	fmt.Println("Proxy server is running on port 8080")
	proxy.ListenAndServe()
}