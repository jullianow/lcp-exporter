package admin

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
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
			fqName("count"),
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
	clusters, err := lcp.FetchFrom[shared.ClusterDiscovery](c.client, "/admin/cluster-discovery/discovered-clusters", nil)
	if err != nil {
		internal.LogError("ClusterDiscoveryCollector", "Failed to fetch discovered clusters: %v", err)
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
