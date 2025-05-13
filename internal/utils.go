package internal

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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
	now := time.Now().UTC()
	start := now.Add(-duration).Truncate(24 * time.Hour)
	end := now.Truncate(24 * time.Hour).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return shared.DateRange{
		From: start.Format(time.RFC3339),
		End:  end.Format(time.RFC3339),
	}
}

func JoinStrings(list []string, separator string) string {
	return strings.Join(list, separator)
}

func RootProjectName(project shared.Projects) string {
	s := project.OrganizationId
	if project.ProjectID == project.OrganizationId {
		return ""
	}
	return s
}

func GetRootProjectIDs(projects []shared.Projects) []string {
	var rootProjectIDs []string
	for _, project := range projects {
		rootProjectName := RootProjectName(project)
		if rootProjectName != "" {
			continue
		} else {
			rootProjectIDs = append(rootProjectIDs, project.ProjectID)
		}
	}
	return rootProjectIDs
}

func MillisToSeconds(ms int64) float64 {
	return float64(ms) / 1000.0
}

func GetCertValidityDatesInSeconds(base64Cert string) (int64, int64, error) {
	certBytes, err := base64.StdEncoding.DecodeString(base64Cert)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode base64: %v", err)
	}

	block, _ := pem.Decode(certBytes)
	if block != nil {
		certBytes = block.Bytes
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert.NotBefore.Unix(), cert.NotAfter.Unix(), nil
}

func GiBToBytes(gib int64) int64 {
	return gib * 1024 * 1024 * 1024
}

func GBToBytes(gb int64) int64 {
	return gb * 1000 * 1000 * 1000
}

func StringToInt64(s string) int64 {
	result, _ := strconv.ParseInt(s, 10, 64)
	return result
}
