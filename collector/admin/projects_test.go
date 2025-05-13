package admin

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jullianow/lcp-exporter/lcp"
)

func TestProjectsCollector(t *testing.T) {

	mockJSON := `[
	{
		"id": "proj-1",
		"cluster": "cluster-1",
		"createdAt": 1672531199000,
		"health": "healthy",
		"organizationId": "proj-1",
		"projectId": "proj-1",
		"status": "running",
		"metadata": {
			"commerce": true,
			"documentLibraryStore": "gcs",
			"trial": "false",
			"subscription": {
			  "availability": "STD",
			  "envType": "PRODUCTION"
			}
		},
		"volumeStorageSize": 100,
		"cloudOptions": {
      "gcpDatabaseEdition": "ENTERPRISE",
      "gcpDatabaseVersion": "POSTGRES_16",
      "gcpDiskType": "PD_HDD",
      "gcpDiskSize": "10",
      "gcpInstanceType": "db-f1-micro"
    }
	},
	{
		"id": "proj-2",
		"cluster": "cluster-2",
		"createdAt": 1672531199000,
		"health": "unhealthy",
		"organizationId": "org-1",
		"projectId": "proj-2",
		"status": "stopped",
		"metadata": {
			"commerce": false,
			"documentLibraryStore": "Simplestore",
			"trial": "true",
			"subscription": {
			  "availability": "NONE",
			  "envType": "NONE"
			}
		},
		"volumeStorageSize": 100,
		"cloudOptions": {
      "gcpDatabaseEdition": "ENTERPRISE",
      "gcpDatabaseVersion": "MYSQL_5_7",
      "gcpDiskType": "PD_HDD",
      "gcpDiskSize": "10",
      "gcpInstanceType": "db-f1-micro"
    }
	}
]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/admin/projects", r.URL.Path)
		_, err := fmt.Fprintln(w, mockJSON)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "fake-token")
	collector := NewProjectsCollector(client)

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

	assert.Contains(t, output, `lcp_api_projects_create_timestamp{id="proj-1"} 1.672531199e+09`)
	assert.Contains(t, output, `lcp_api_projects_create_timestamp{id="proj-1"} 1`)
	assert.Contains(t, output, `lcp_api_projects_create_timestamp{id="proj-1"} 1`)
	assert.Contains(t, output, `lcp_api_projects_create_timestamp{id="proj-2"} 1.672531199e+09`)
	assert.Contains(t, output, `lcp_api_projects_db_labels{edition="ENTERPRISE",id="proj-2",type="db-f1-micro",version="MYSQL_5_7"} 1`)
	assert.Contains(t, output, `lcp_api_projects_db_storage_capacity_bytes{disk_type="PD_HDD",id="proj-2"} 1e+10`)
	assert.Contains(t, output, `lcp_api_projects_labels{availability="STD",cluster_name="cluster-1",commerce="true",doc_lib_store="gcs",health="true",id="proj-1",name="proj-1",parent_project_name="",root_project="true",trial="false",type="PRODUCTION"} 1`)
	assert.Contains(t, output, `lcp_api_projects_labels{availability="NONE",cluster_name="cluster-2",commerce="false",doc_lib_store="Simplestore",health="false",id="proj-2",name="proj-2",parent_project_name="org-1",root_project="false",trial="true",type="NONE"} 1`)
	assert.Contains(t, output, `lcp_api_projects_status{id="proj-1"} 1`)
	assert.Contains(t, output, `lcp_api_projects_status{id="proj-2"} 0`)
	assert.Contains(t, output, `lcp_api_projects_total 2`)
	assert.Contains(t, output, `lcp_api_projects_volume_storage_capacity_bytes{id="proj-1"} 1.073741824e+11`)
	assert.Contains(t, output, `lcp_api_projects_volume_storage_capacity_bytes{id="proj-2"} 1.073741824e+11`)
}
