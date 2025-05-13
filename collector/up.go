package collector

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type upCollector struct {
	client *lcp.Client
	up     *prometheus.Desc
}

func NewUpCollector(client *lcp.Client) *upCollector {
	fqName := internal.Name("status")

	return &upCollector{
		client: client,
		up: prometheus.NewDesc(
			fqName("up"),
			"1 if the API is up, 0 if it is down",
			[]string{"status"},
			nil,
		),
	}
}

func (c *upCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
}

func (c *upCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectMetrics(ch)
	}()
	wg.Wait()
}

func (c *upCollector) collectMetrics(ch chan<- prometheus.Metric) {
	health, err := lcp.FetchFrom[shared.HealthCheck](c.client, "/health-check", nil)
	if err != nil {
		internal.LogError("UpCollector", "Failed to fetch health check: %v", err)
		return
	}
	if len(health) == 0 {
		internal.LogWarn("UpCollector", "No health check data returned")
		return
	}

	status := health[0].Status
	var value float64
	if status == "up" {
		value = 1
	} else {
		value = 0
	}

	ch <- prometheus.MustNewConstMetric(
		c.up,
		prometheus.GaugeValue,
		value,
		status,
	)
}
