package main

import (
	"html/template"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/jullianow/lcp-exporter/collector"
	"github.com/jullianow/lcp-exporter/collector/admin"
	"github.com/jullianow/lcp-exporter/config"
	"github.com/jullianow/lcp-exporter/lcp"
)

var VERSION string

const (
	rootTemplate = `<!DOCTYPE html>
<html>
<head>
	<title>LCP exporter</title>
</head>
<body>
	<h2>LCP exporter. Version: {{.VERSION}}</h2>
	<ul>
		<li><a href="{{.MetricsPath}}">Metrics</a></li>
		<li><a href="/healthz">Health Check</a></li>
	</ul>
</body>
</html>`
)

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		msg := "error writing healthz handler"
		logrus.Errorf("[handleHealthz] %s: %v", msg, err)
	}
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.New("root").Parse(rootTemplate))
	data := struct {
		MetricsPath string
	}{
		MetricsPath: "/metrics", // Assuming MetricsPath is a global variable
	}

	if err := tmpl.Execute(w, data); err != nil {
		msg := "error rendering root template"
		logrus.Errorf("[handleRoot] %s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func main() {
	cfg := config.ParseFlags()

	client := lcp.NewClient(cfg.Endpoint, cfg.Token)

	registry := prometheus.NewRegistry()

	if !cfg.EnableGoMetrics {
		logrus.Info("Disabling Go default metrics")
		prometheus.Unregister(collectors.NewGoCollector())
	}

	if !cfg.EnableProcessMetrics {
		logrus.Info("Disabling process metrics")
		prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	if !cfg.EnablePromHttpMetrics {
		logrus.Info("Disabling promhttp metrics")
	} else {
		http.Handle(cfg.MetricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	}

	collectorConfigs := map[string]struct {
		collector prometheus.Collector
		enable    bool
	}{
		"cluster_discovery": {
			collector: admin.NewclusterDiscoveryCollector(client),
			enable:    cfg.EnableClusterDiscoveryMetrics,
		},
		"info": {
			collector: collector.NewInfoCollector(client),
			enable:    true,
		},
		"up": {
			collector: collector.NewUpCollector(client),
			enable:    true,
		},
	}

	for name, config := range collectorConfigs {
		if config.enable {
			logrus.Infof("Registering collector: %s", name)
			registry.MustRegister(config.collector)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/healthz", http.HandlerFunc(handleHealthz))
	mux.Handle(cfg.MetricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	logrus.Infof("Starting HTTP server on Port %s, serving metrics at %s. Version: %s", cfg.Port, cfg.MetricsPath, VERSION)

	logrus.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
