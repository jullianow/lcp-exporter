package admin

import (
	"fmt"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type clusterDiscoveryCollector struct {
	client             *lcp.Client
	clusterTotal       *prometheus.Desc
	labels             *prometheus.Desc
	caCreatedTimestamp *prometheus.Desc
	caExpiredTimestamp *prometheus.Desc
}

func NewClusterDiscoveryCollector(client *lcp.Client) *clusterDiscoveryCollector {
	fqName := internal.Name("cluster_discovery")

	return &clusterDiscoveryCollector{
		client: client,
		clusterTotal: prometheus.NewDesc(
			fqName("clusters_total"),
			"Total number of discovered clusters",
			nil, nil,
		),
		labels: prometheus.NewDesc(
			fqName("labels"),
			"Labels of discovered clusters",
			[]string{"lcp_cluster_name", "name", "provider", "cloud_project_id", "location", "plan_id", "is_lxc"},
			nil,
		),
		caCreatedTimestamp: prometheus.NewDesc(
			fqName("ca_created_timestamp"),
			"Timestamp of the creation of CA",
			[]string{"lcp_cluster_name"},
			nil,
		),
		caExpiredTimestamp: prometheus.NewDesc(
			fqName("ca_expired_timestamp"),
			"Timestamp of the expiration of CA",
			[]string{"lcp_cluster_name"},
			nil,
		),
	}
}

func (c *clusterDiscoveryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.caCreatedTimestamp
	ch <- c.caExpiredTimestamp
	ch <- c.clusterTotal
	ch <- c.labels
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
		c.clusterTotal,
		prometheus.GaugeValue,
		float64(len(clusters)),
	)

	for _, cluster := range clusters {
		lcpClusterName := strings.ToLower(fmt.Sprintf("%s_%s", cluster.Provider.CloudProjectID, cluster.Name))
		notBefore, notAfter, _ := internal.GetCertValidityDatesInSeconds(cluster.Kubeconfig.Cluster.CaData)

		ch <- prometheus.MustNewConstMetric(
			c.labels,
			prometheus.GaugeValue,
			1.0,
			lcpClusterName,
			cluster.Name,
			cluster.Provider.Name,
			cluster.Provider.CloudProjectID,
			cluster.Location,
			cluster.PlanID,
			internal.BoolToString(cluster.IsLXC),
		)

		ch <- prometheus.MustNewConstMetric(
			c.caCreatedTimestamp,
			prometheus.GaugeValue,
			float64(notBefore),
			lcpClusterName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.caExpiredTimestamp,
			prometheus.GaugeValue,
			float64(notAfter),
			lcpClusterName,
		)
	}
}
