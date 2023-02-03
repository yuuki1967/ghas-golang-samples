package controllers

import (
	"fmt"
	"io"
	"net/http"

	"example.com/logger"
)

func hello(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	b, _ := io.ReadAll(req.Body)
	username := string(b)

	logger.Log().Errorf("user %s logged in.\n", username)
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func Controllersbody() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8090", nil)
}
