package logmod

import (
	logr "github.com/sirupsen/logrus"
)

func Logprint(msg string) {
	logr.WithFields(logr.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("user %s logged in.\n", msg)
}
