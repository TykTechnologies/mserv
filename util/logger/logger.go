// package logger provides a centralised logger resource
package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
)

var logged = false

func logOnce(log *logrus.Logger) {
	if logged {
		return
	}
	//log.Info("tracer hook initialising")
	logged = true
}

// Certain apps shou.d not have mongo tracing enabled
var mgoExclusions = make([]string, 0)

func isExcluded(app string) bool {
	for _, exc := range mgoExclusions {
		//fmt.Printf("checking %v against %v\n", app, exc)
		if app == exc {
			return true
		}
	}

	return false
}

func GetAndExcludeLoggerFromTrace(tag string) *logrus.Entry {
	for _, v := range mgoExclusions {
		if v == tag {
			return GetLogger(tag)
		}
	}

	mgoExclusions = append(mgoExclusions, tag)
	return GetLogger(tag)
}

// GetLogger will provide a tagged logger by passing in a `tag` value for easier log parsing
func GetLogger(tag string) *logrus.Entry {
	lvl := os.Getenv("TYK_CONTROLLER_LOGLEVEL")

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

	logrus.SetLevel(level)

	logger := logrus.New()

	var log = logger.WithFields(logrus.Fields{
		"app": tag,
	})

	return log
}
