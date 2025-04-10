package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jullianow/lcp-exporter/internal/shared"
)

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func IntToString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return ""
	}
}

func LogInfo(context string, msg string, args ...any) {
	logrus.WithField("context", context).Infof(msg, args...)
}

func LogWarn(context string, msg string, args ...any) {
	logrus.WithField("context", context).Warnf(msg, args...)
}

func LogError(context string, msg string, args ...any) {
	logrus.WithField("context", context).Errorf(msg, args...)
}

func LogDebug(context string, msg string, args ...any) {
	logrus.WithField("context", context).Debugf(msg, args...)
}

func LogFatal(context, format string, args ...interface{}) {
	logrus.Fatalf("[%s] %s", context, fmt.Sprintf(format, args...))
}

func CalculateDates(duration time.Duration) shared.DateRange {
	currentDate := time.Now()
	endDate := currentDate.Format("2006-01-02")
	fromDate := currentDate.Add(-duration).Format("2006-01-02")

	return shared.DateRange{
		From: fromDate,
		End:  endDate,
	}
}

func JoinStrings(list []string) string {
	return strings.Join(list, ",")
}

func IsParentProject(project shared.Projects) bool {
	return project.ProjectID == project.ParentProjectID
}
