package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/lcp"
)

var mockClusterDiscoveryResponse = map[string]internal.ClusterDiscovery{
	"cluster-1": {
		Name: "cluster-1",
		Provider: internal.Provider{
			Name:           "gcp",
			CloudProjectID: "project-123",
		},
		Location:             "us-central1",
		CustomerBackupBucket: "backup-bucket",
		Zones:                []string{"zone-a", "zone-b"},
		PlanID:               "plan-xyz",
		IsLXC:                true,
	},
}

func TestClusterDiscoveryCollector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/admin/cluster-discovery/discovered-clusters", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		response, _ := json.Marshal(mockClusterDiscoveryResponse)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "")

	collector := NewClusterDiscoveryCollector(client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	recorder := httptest.NewRecorder()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	body, _ := io.ReadAll(recorder.Body)
	metricsOutput := string(body)
	assert.Contains(t, metricsOutput, "lcp_api_cluster_discovery_info{cloud_project_id=\"project-123\",is_lxc=\"true\",location=\"us-central1\",name=\"cluster-1\",plan_id=\"plan-xyz\",provider=\"gcp\"} 1")
	assert.Contains(t, metricsOutput, "lcp_api_cluster_discovery_total 1")
}
