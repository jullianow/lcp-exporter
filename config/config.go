package config

import (
	"flag"
	"os"
	"time"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Duration                      time.Duration
	EnableClusterDiscoveryMetrics bool
	EnableGoMetrics               bool
	EnableProcessMetrics          bool
	EnableProjectsMetrics         bool
	EnablePromHttpMetrics         bool
	EnableAutoscaleMetrics        bool
	Endpoint                      string
	LogFormat                     string
	LogLevel                      string
	MetricsPath                   string
	Port                          string
	Token                         string
}

func ParseFlags() *Config {
	var cfg Config

	flag.BoolVar(&cfg.EnableClusterDiscoveryMetrics, "enable-cluster-discovery-metrics", true, "Enable cluster discovery metrics")
	flag.BoolVar(&cfg.EnableAutoscaleMetrics, "enable-autoscale-metrics", true, "Enable autoscale metrics")
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
		internal.LogFatal("Config", "Authentication error: Provide either LCP_API_TOKEN")
	}

	if cfg.Endpoint == "" {
		internal.LogFatal("Config", "The API endpoint must be provided with the -endpoint flag")
	}

	if cfg.Duration.Seconds() < 0 {
		internal.LogFatal("Config", "Invalid duration: must be non-negative, got %s", cfg.Duration.String())
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		internal.LogFatal("Config", "Invalid log level: %v", err)
	}
	logrus.SetLevel(level)

	return &cfg
}
