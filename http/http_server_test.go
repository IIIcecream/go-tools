package http

import "testing"

func TestHttpServerStart(t *testing.T) {
	server := HttpServer{0}
	if err := server.Start(8080); err != nil {
		t.Fatal(err)
	}
}
