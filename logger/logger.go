package logger

import (
	"fmt"
	"runtime"
	"strings"

	logrus "github.com/sirupsen/logrus"
	ecslogr "go.elastic.co/ecslogrus"
)

var (
	Log    = getLogger
	logger = logrus.New()
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&ecslogr.Formatter{
		DisableHTMLEscape: false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			arr := strings.Split(f.File, "/portfolio/")
			//there can be a parent folder named portfolio so I am taking last element
			finalStr := arr[len(arr)-1]
			//finding next / after '/portfolio', adding +1 to remove /
			res := finalStr[strings.Index(finalStr, "/")+1:]
			return funcName, fmt.Sprintf("%s:%d", res, f.Line)
		},
	})

}
func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	})
}
