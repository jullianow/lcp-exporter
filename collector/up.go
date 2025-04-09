package collector

import (
	"encoding/json"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

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

	labelKeys := []string{"status"}
	return &upCollector{
		client: client,
		up: prometheus.NewDesc(
			fqName("up"),
			"1 if the API is up, 0 if it is down",
			labelKeys, nil,
		),
	}
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
	resp, err := c.client.MakeRequest("/health-check")
	if err != nil {
		logrus.Errorf("Error collecting health check data: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Warnf("Error closing response body: %v", err)
		}
	}()

	var responseData shared.HealthCheck
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		logrus.Errorf("Error decoding health check response: %v", err)
		return
	}

	// Emite o valor da métrica
	var metricValue float64
	if responseData.Status == "up" {
		metricValue = 1
	} else {
		metricValue = 0
	}

	ch <- prometheus.MustNewConstMetric(
		c.up,
		prometheus.GaugeValue,
		metricValue,
		responseData.Status,
	)
}

func (c *upCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
}
