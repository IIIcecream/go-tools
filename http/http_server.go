package http

import (
	"fmt"
	"net/http"
)

type HttpServer struct {
	port int
}

func (s *HttpServer) Start(port int) error {
	s.port = port
	http.HandleFunc("/", handler)
	http.HandleFunc("/admin", handler_admin)
	http.HandleFunc("/login", handler_login)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	if err != nil {
		return err
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func handler_admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "who are you!")
}
func handler_login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "input your username and password!")
}
