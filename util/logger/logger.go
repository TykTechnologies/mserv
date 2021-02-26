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
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "warning":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	default:
		level = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetReportCaller(level >= logrus.DebugLevel)

	return logger.WithField("app", tag)
}
