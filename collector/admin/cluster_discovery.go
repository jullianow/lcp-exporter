package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/lcp"
)

type clusterDiscoveryCollector struct {
	client *lcp.Client
	count  *prometheus.Desc
	info   *prometheus.Desc
}

func NewClusterDiscoveryCollector(client *lcp.Client) *clusterDiscoveryCollector {
	fqName := internal.Name("cluster_discovery")

	return &clusterDiscoveryCollector{
		client: client,
		count: prometheus.NewDesc(
			fqName("total"),
			"Total number of discovered clusters",
			nil, nil,
		),
		info: prometheus.NewDesc(
			fqName("info"),
			"Information about discovered clusters",
			[]string{"name", "provider", "cloud_project_id", "location", "plan_id", "is_lxc"},
			nil,
		),
	}
}

func (c *clusterDiscoveryCollector) Collect(ch chan<- prometheus.Metric) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectMetrics(ch)
	}()
	wg.Wait()
}

func (c *clusterDiscoveryCollector) collectMetrics(ch chan<- prometheus.Metric) {
	resp, err := c.client.MakeRequest("/admin/cluster-discovery/discovered-clusters")
	if err != nil {
		logrus.Errorf("Error collecting cluster discovery data: %v", err)
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Warnf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Error accessing cluster discovery API: StatusCode %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error reading cluster discovery response: %v", err)
		return
	}

	var clusters map[string]internal.ClusterDiscovery
	if err := json.Unmarshal(body, &clusters); err != nil {
		logrus.Errorf("Error decoding cluster discovery JSON: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.count,
		prometheus.GaugeValue,
		float64(len(clusters)),
	)

	for _, cluster := range clusters {
		ch <- prometheus.MustNewConstMetric(
			c.info,
			prometheus.GaugeValue,
			1.0,
			cluster.Name,
			cluster.Provider.Name,
			cluster.Provider.CloudProjectID,
			cluster.Location,
			cluster.PlanID,
			internal.BoolToString(cluster.IsLXC),
		)
	}
}

func (c *clusterDiscoveryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.count
	ch <- c.info
}
