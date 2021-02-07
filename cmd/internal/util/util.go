package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLog() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel

	return log
}
