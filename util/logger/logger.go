// Package logger provides a centralised logger resource.
package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetLogger will provide a tagged logger by passing in a 'tag' value for easier log parsing.
func GetLogger(tag string) *logrus.Entry {
	lvl := os.Getenv("TYK_MSERV_LOGLEVEL")

	var level logrus.Level
	switch strings.ToLower(lvl) {
	case "debug":
		level = logrus.DebugLevel
	case "warning":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "error":
		level = logrus.ErrorLevel
	default:
		level = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(level)

	log := logger.WithField("app", tag)

	return log
}
