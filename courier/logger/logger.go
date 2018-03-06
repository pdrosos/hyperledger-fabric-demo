package logger

import (
	"github.com/sirupsen/logrus"
)

var Revision string
var Log *logrus.Logger = NewLogger()

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	return logger
}
