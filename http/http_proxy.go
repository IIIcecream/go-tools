package http

import (
	"io"
	"log"
	"net"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 创建与目标服务器的连接
	destConn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		log.Printf("Error connecting to destination server: %v", err)
		http.Error(w, "Error connecting to destination server", http.StatusInternalServerError)
		return
	}
	defer destConn.Close()

	// 将客户端请求发送给目标服务器
	err = r.Write(destConn)
	if err != nil {
		log.Printf("Error sending request to destination server: %v", err)
		http.Error(w, "Error sending request to destination server", http.StatusInternalServerError)
		return
	}

	// 将目标服务器的响应发送给客户端
	_, err = io.Copy(w, destConn)
	if err != nil {
		log.Printf("Error sending response to client: %v", err)
		http.Error(w, "Error sending response to client", http.StatusInternalServerError)
		return
	}
}
