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
	"version": "0.0.0",
	"domains": {
		"infrastructure": "liferay.cloud",
		"service": "lfr.cloud"
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
	err := registry.Register(collector)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics:  true,
		DisableCompression: true,
	})
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	body := recorder.Body.String()

	assert.Contains(t, body, `lcp_api_status_info{infrastructure_domain="liferay.cloud",service_domain="lfr.cloud",version="0.0.0"} 1`)
	assert.NotContains(t, body, "go_")
	assert.NotContains(t, body, "promhttp_")
}

func TestInfoCollector_FetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "")

	collector := NewInfoCollector(client)

	registry := prometheus.NewRegistry()
	err := registry.Register(collector)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics:  true,
		DisableCompression: true,
	})
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	body := recorder.Body.String()

	assert.NotContains(t, body, "lcp_api_status_info")
}

func TestInfoCollector_EmptyData(t *testing.T) {
	mockEmptyResponse := `{}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockEmptyResponse))
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "")

	collector := NewInfoCollector(client)

	registry := prometheus.NewRegistry()
	err := registry.Register(collector)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics:  true,
		DisableCompression: true,
	})
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	body := recorder.Body.String()

	assert.NotContains(t, body, "lcp_api_status_info")
}
