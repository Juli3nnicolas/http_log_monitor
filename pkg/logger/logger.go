package logger

import (
	"github.com/sirupsen/logrus"
)

// Get returns the app logger (it is protected with a mutex for concurrent writing)
func Get() *logrus.Logger {
	return log
}
