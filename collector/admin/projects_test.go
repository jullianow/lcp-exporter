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

var mockProjectsResponse = []internal.Projects{
	{
		Id:              "proj-1",
		Cluster:         "cluster-1",
		Health:          "healthy",
		ParentProjectID: "proj-1",
		ProjectID:       "proj-1",
		Status:          "running",
		Metadata: internal.ProjectMetadata{
			Commerce: true,
			Type:     "dev",
			Trial:    "false",
		},
	},
	{
		Id:              "proj-2",
		Cluster:         "cluster-2",
		Health:          "unhealthy",
		ParentProjectID: "org-1",
		ProjectID:       "proj-2",
		Status:          "stopped",
		Metadata: internal.ProjectMetadata{
			Commerce: false,
			Type:     "prod",
			Trial:    "true",
		},
	},
}

func TestProjectsCollector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/admin/projects", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(mockProjectsResponse)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "")
	collector := NewProjectsCollector(client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	recorder := httptest.NewRecorder()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	body, _ := io.ReadAll(recorder.Body)
	metricsOutput := string(body)

	assert.Contains(t, metricsOutput, `lcp_api_projects_info{commerce="true",health="true",id="proj-1",name="proj-1",parent_project_id="",root_project="true",trial="false",type="dev"} 1`)
	assert.Contains(t, metricsOutput, `lcp_api_projects_info{commerce="false",health="false",id="proj-2",name="proj-2",parent_project_id="org-1",root_project="false",trial="true",type="prod"} 0`)
	assert.Contains(t, metricsOutput, `lcp_api_projects_count 2`)
}
