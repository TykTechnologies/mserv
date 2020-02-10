// package logger provides a centralised logger resource
package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var logged = false

func logOnce(log *logrus.Logger) {
	if logged {
		return
	}
	logged = true
}

// GetLogger will provide a tagged logger by passing in a `tag` value for easier log parsing
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

	var log = logger.WithFields(logrus.Fields{
		"app": tag,
	})

	return log
}
