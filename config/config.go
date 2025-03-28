package config

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	LogFormat                     string
	LogLevel                      string
	Endpoint                      string
	Port                          string
	MetricsPath                   string
	EnableGoMetrics               bool
	EnableProcessMetrics          bool
	EnablePromHttpMetrics         bool
	EnableClusterDiscoveryMetrics bool
	Token                         string
}

func ParseFlags() *Config {
	var cfg Config

	flag.StringVar(&cfg.LogFormat, "log-format", "json", "Log format (json or text)")
	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&cfg.Endpoint, "endpoint", "", "Base endpoint for the REST API (example: https://api.liferay.st)")
	flag.StringVar(&cfg.Port, "port", "9103", "Port for the HTTP server (default: 9103)")
	flag.StringVar(&cfg.MetricsPath, "metrics-path", "/metrics", "Path for the metrics endpoint (default: /metrics)")
	flag.BoolVar(&cfg.EnableGoMetrics, "enable-go-metrics", false, "Enable Go default metrics (default: false)")
	flag.BoolVar(&cfg.EnableProcessMetrics, "enable-process-metrics", false, "Enable process metrics (default: false)")
	flag.BoolVar(&cfg.EnablePromHttpMetrics, "enable-promhttp-metrics", false, "Enable promhttp metrics (default: false)")
	flag.BoolVar(&cfg.EnableClusterDiscoveryMetrics, "enable-cluster-discovery-metrics", true, "Enable cluster discovery metrics (default: false)")

	flag.Parse()

	switch cfg.LogFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	cfg.Token = os.Getenv("TOKEN")

	if cfg.Token == "" {
		logrus.Fatal("Authentication error: Provide either TOKEN")
	}

	if cfg.Endpoint == "" {
		logrus.Fatal("The API endpoint must be provided with the -endpoint flag")
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)

	return &cfg
}
