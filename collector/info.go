package collector

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type infoCollector struct {
	client *lcp.Client
	info   *prometheus.Desc
}

func NewInfoCollector(client *lcp.Client) *infoCollector {
	fqName := internal.Name("status")

	return &infoCollector{
		client: client,
		info: prometheus.NewDesc(
			fqName("info"),
			"1 if the API is up, 0 if it is down",
			[]string{"version", "infrastructure_domain", "service_domain"},
			nil,
		),
	}
}

func (c *infoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.info
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
	info, err := lcp.FetchFrom[shared.Info](c.client, "/", nil)
	if err != nil {
		internal.LogError("InfoCollector", "Failed to fetch status info: %v", err)
		return
	}
	if len(info) == 0 {
		internal.LogWarn("InfoCollector", "No status info returned from API")
		return
	}

	data := info[0]
	ch <- prometheus.MustNewConstMetric(
		c.info,
		prometheus.GaugeValue,
		1,
		data.Version,
		data.Domains.Infrastructure,
		data.Domains.Service,
	)
}
