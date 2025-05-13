package main

import (
	"html/template"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jullianow/lcp-exporter/collector"
	"github.com/jullianow/lcp-exporter/collector/admin"
	"github.com/jullianow/lcp-exporter/config"
	"github.com/jullianow/lcp-exporter/internal"
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
		internal.LogError("handleHealthz", "%s: %v", msg, err)
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
		internal.LogError("handleRoot", "%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func main() {
	cfg := config.ParseFlags()

	client := lcp.NewClient(cfg.Endpoint, cfg.Token)

	registry := prometheus.NewRegistry()

	if !cfg.EnableGoMetrics {
		internal.LogInfo("Main", "Disabling Go default metrics")
		prometheus.Unregister(collectors.NewGoCollector())
	}

	if !cfg.EnableProcessMetrics {
		internal.LogInfo("Main", "Disabling process metrics")
		prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	if !cfg.EnablePromHttpMetrics {
		internal.LogInfo("Main", "Disabling promhttp metrics")
	} else {
		http.Handle(cfg.MetricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	}

	dataRange := internal.CalculateDates(cfg.Duration)
	projectsCollector := admin.NewProjectsCollector(client)

	collectorConfigs := []struct {
		name      string
		collector prometheus.Collector
		enable    bool
	}{
		{
			name:      "projects",
			collector: projectsCollector,
			enable:    true,
		},
		{
			name:      "autoscale",
			collector: admin.NewAutoscaleCollector(client, projectsCollector, dataRange),
			enable:    true,
		},
		{
			name:      "cluster_discovery",
			collector: admin.NewClusterDiscoveryCollector(client),
			enable:    cfg.EnableClusterDiscoveryMetrics,
		},
		{
			name:      "info",
			collector: collector.NewInfoCollector(client),
			enable:    true,
		},
		{
			name:      "up",
			collector: collector.NewUpCollector(client),
			enable:    true,
		},
	}

	for _, config := range collectorConfigs {
		if config.enable {
			internal.LogInfo("Main", "Registering collector: %s", config.name)
			registry.MustRegister(config.collector)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/healthz", http.HandlerFunc(handleHealthz))
	mux.Handle(cfg.MetricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	internal.LogInfo("Main", "Starting HTTP server on Port %s, serving metrics at %s. Version: %s", cfg.Port, cfg.MetricsPath, VERSION)

	internal.LogFatal("Main", "%s", http.ListenAndServe(":"+cfg.Port, mux))
}
