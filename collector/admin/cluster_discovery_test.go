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
			"customerBackupBucket": "gs://my-backup-bucket",
			"location": "us-central1",
			"planId": "plan-xyz",
			"isLXC": true
		}
	}`

	// Servidor de teste simulando a API LCP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/admin/cluster-discovery/discovered-clusters", r.URL.Path)
		_, err := fmt.Fprintln(w, mockJSON)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Cria o client apontando para o servidor fake
	client := lcp.NewClient(server.URL, "fake-token")
	collector := NewClusterDiscoveryCollector(client)

	// Registra o coletor em um registry isolado
	reg := prometheus.NewRegistry()
	require.NoError(t, reg.Register(collector))

	// Expõe as métricas e captura a saída
	serverMetrics := httptest.NewServer(promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	defer serverMetrics.Close()

	// Chamada real ao endpoint de métricas
	resp, err := http.Get(serverMetrics.URL)
	require.NoError(t, err)
	defer func() {
		closeErr := resp.Body.Close()
		require.NoError(t, closeErr)
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	output := string(body)

	require.Contains(t, output, `lcp_api_cluster_discovery_count 1`)
	require.Contains(t, output, `lcp_api_cluster_discovery_info{cloud_project_id="project-123",is_lxc="true",location="us-central1",name="cluster-1",plan_id="plan-xyz",provider="gcp"} 1`)
}
