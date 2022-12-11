package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"log"
	"github.com/sirupsen/logrus"
)

func ExampleJSONFormatter_CallerPrettyfier() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}
	l.Info("example of custom format caller")
	// Output:
	// {"file":"example_custom_caller_test.go","func":"ExampleJSONFormatter_CallerPrettyfier","level":"info","msg":"example of custom format caller"}
}

func handler(ss string) {
    username := ss
    log.Printf("user %s logged in.\n", username)
}

func main() {
	fmt.Println("Welcome to the playground!")
	handler("test")
	fmt.Println("The time is", time.Now())
	ExampleJSONFormatter_CallerPrettyfier()
}
