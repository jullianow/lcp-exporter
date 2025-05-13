package admin

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"

	"github.com/jullianow/lcp-exporter/lcp"
)

func TestClusterDiscoveryCollector(t *testing.T) {
	mockJSON := `{
		"123": {
			"name": "cluster-1",
			"provider": {
				"name": "gcp",
				"cloudProjectId": "project-123"
			},
			"kubeconfig": {
				"cluster": {
					"caData": "base64-ca-data"
				}
			},
			"customerBackupBucket": "gs://my-backup-bucket",
			"location": "us-central1",
			"planId": "plan-xyz",
			"isLXC": true
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/admin/cluster-discovery/discovered-clusters", r.URL.Path)
		_, err := fmt.Fprintln(w, mockJSON)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "fake-token")
	collector := NewClusterDiscoveryCollector(client)

	reg := prometheus.NewRegistry()
	require.NoError(t, reg.Register(collector))

	serverMetrics := httptest.NewServer(promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	defer serverMetrics.Close()

	resp, err := http.Get(serverMetrics.URL)
	require.NoError(t, err)
	defer func() {
		closeErr := resp.Body.Close()
		require.NoError(t, closeErr)
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	output := string(body)

	require.Contains(t, output, `lcp_api_cluster_discovery_ca_created_timestamp{lcp_cluster_name="project-123_cluster-1"} 0`)
	require.Contains(t, output, `lcp_api_cluster_discovery_clusters_total 1`)
	require.Contains(t, output, `lcp_api_cluster_discovery_labels{cloud_project_id="project-123",is_lxc="true",lcp_cluster_name="project-123_cluster-1",location="us-central1",name="cluster-1",plan_id="plan-xyz",provider="gcp"} 1`)
}
