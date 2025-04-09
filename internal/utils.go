package internal

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// LogInfo logs informational messages.
func LogInfo(context string, msg string, args ...any) {
	logrus.WithField("context", context).Infof(msg, args...)
}

// LogWarn logs warning messages.
func LogWarn(context string, msg string, args ...any) {
	logrus.WithField("context", context).Warnf(msg, args...)
}

// LogError logs error messages.
func LogError(context string, msg string, args ...any) {
	logrus.WithField("context", context).Errorf(msg, args...)
}

// LogDebug logs debug messages.
func LogDebug(context string, msg string, args ...any) {
	logrus.WithField("context", context).Debugf(msg, args...)
}

func LogFatal(context, format string, args ...interface{}) {
	logrus.Fatalf("[%s] %s", context, fmt.Sprintf(format, args...))
}
