package collector

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	"github.com/jullianow/lcp-exporter/lcp"
)

var mockInfoResponse = `{
	"version": "v1.0.0",
	"domains": {
		"infa": "infrastructure.com",
		"service": "service.com"
	}
}`

func TestInfoCollector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockInfoResponse))
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "")

	collector := NewInfoCollector(client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	recorder := httptest.NewRecorder()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics:  true,
		DisableCompression: true,
	})
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	body := recorder.Body.String()

	assert.Contains(t, body, "lcp_api_status_info{infrastructure_domain=\"\",service_domain=\"service.com\",version=\"v1.0.0\"} 1")

	assert.NotContains(t, body, "go_")
	assert.NotContains(t, body, "promhttp_")
}
