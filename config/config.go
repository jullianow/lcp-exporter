package config

import (
	"flag"
	"os"
	"time"

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
	Duration                      time.Duration
	EnableProjectsMetrics         bool
}

func ParseFlags() *Config {
	var cfg Config

	flag.BoolVar(&cfg.EnableClusterDiscoveryMetrics, "enable-cluster-discovery-metrics", true, "Enable cluster discovery metrics")
	flag.BoolVar(&cfg.EnableProjectsMetrics, "enable-projects-metrics", true, "Enable projects metrics")
	flag.BoolVar(&cfg.EnableGoMetrics, "enable-go-metrics", false, "Enable Go default metrics")
	flag.BoolVar(&cfg.EnableProcessMetrics, "enable-process-metrics", false, "Enable process metrics")
	flag.BoolVar(&cfg.EnablePromHttpMetrics, "enable-promhttp-metrics", false, "Enable promhttp metrics")
	flag.DurationVar(&cfg.Duration, "duration", 0, "Duration to shift from now (e.g. 24h, -48h)")
	flag.StringVar(&cfg.Endpoint, "endpoint", "", "Base endpoint for the REST API")
	flag.StringVar(&cfg.LogFormat, "log-format", "json", "Log format (json or text)")
	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&cfg.MetricsPath, "metrics-path", "/metrics", "Path for the metrics endpoint")
	flag.StringVar(&cfg.Port, "port", "9103", "Port for the HTTP server")

	flag.Parse()

	switch cfg.LogFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	cfg.Token = os.Getenv("LCP_API_TOKEN")

	if cfg.Token == "" {
		logrus.Fatal("Authentication error: Provide either LCP_API_TOKEN")
	}

	if cfg.Endpoint == "" {
		logrus.Fatal("The API endpoint must be provided with the -endpoint flag")
	}

	if cfg.Duration.Seconds() < 0 {
		logrus.Fatalf("Invalid duration: must be non-negative, got %s", cfg.Duration.String())
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)

	return &cfg
}
