package collector

import (
	"encoding/json"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/lcp"
)

type infoCollector struct {
	client *lcp.Client
	info   *prometheus.Desc
}

func NewInfoCollector(client *lcp.Client) *infoCollector {
	fqName := internal.Name("status")

	labelKeys := []string{"version", "infrastructure_domain", "service_domain"}
	return &infoCollector{
		client: client,
		info: prometheus.NewDesc(
			fqName("info"),
			"1 if the API is up, 0 if it is down",
			labelKeys, nil,
		),
	}
}

func (c *infoCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectMetrics(ch)
	}()
	wg.Wait()
}

func (c *infoCollector) collectMetrics(ch chan<- prometheus.Metric) {
	resp, err := c.client.MakeRequest("/")
	if err != nil {
		logrus.Errorf("Error collecting health check data: %v", err)
		return
	}
	defer resp.Body.Close()

	var responseData internal.Info
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		logrus.Errorf("Error decoding health check response: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.info,
		prometheus.GaugeValue,
		1,
		responseData.Version,
		responseData.Domains.Infa,
		responseData.Domains.Service,
	)
}

func (c *infoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.info
}
