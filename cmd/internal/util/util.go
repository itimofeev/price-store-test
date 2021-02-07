package util

import (
	"io"
	"os"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

func NewLog() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel

	return log
}

func RandomID() string {
	return xid.New().String()
}

type SimpleReadCloser struct {
	io.Reader
}

func (rc *SimpleReadCloser) Close() error {
	return nil
}
