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

var mockUpResponse = `{"status": "up"}`

func TestUpCollector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/health-check", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockUpResponse))
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "")

	collector := NewUpCollector(client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	recorder := httptest.NewRecorder()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics:  true,
		DisableCompression: true,
	})
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	body := recorder.Body.String()

	assert.Contains(t, body, "lcp_api_status_up{status=\"up\"} 1")

	assert.NotContains(t, body, "go_")
	assert.NotContains(t, body, "promhttp_")
}
