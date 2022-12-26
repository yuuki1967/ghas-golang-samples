package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	logr "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	const Environment = "production"
	if Environment == "production" {
		logr.SetFormatter(&logr.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		logr.SetFormatter(&logr.TextFormatter{})
	}
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logr.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logr.SetLevel(logr.WarnLevel)
}

func hello(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	b, _ := io.ReadAll(req.Body)
	username := string(b)

	logr.WithFields(logr.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("user %s logged in.\n", username)
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	logr.WithFields(logr.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logr.WithFields(logr.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")
	/*
			logr.WithFields(logr.Fields{
				"omg":    true,
				"number": 100,
			}).Fatal("The ice breaks!")
		// A common pattern is to re-use fields between logging statements by re-using
		// the logrus.Entry returned from WithFields()
		contextLogger := logr.WithFields(logr.Fields{
			"common": "this is a common field",
			"other":  "I also should be logged always",
		})

		contextLogger.Info("I'll be logged with common and other field")
		contextLogger.Info("Me too")
	*/

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8090", nil)
}
