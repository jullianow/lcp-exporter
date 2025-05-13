package admin

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

func TestAutoscaleCollector(t *testing.T) {
	projectProvider := &ProjectsCollector{
		projects: []shared.Projects{
			{
				ProjectID: "proj-1",
				Cluster:   "project-123_cluster-1",
				Metadata: shared.ProjectMetadata{
					Subscription: struct {
						Availability string `json:"availability"`
						EnvType      string `json:"envType"`
					}{
						Availability: "HA",
						EnvType:      "prod",
					},
					Commerce:    true,
					DocLibStore: "doclib",
					Trial:       "false",
				},
				CreatedAt: 1740926966862,
				Status:    "running",
				Health:    "healthy",
			},
			{
				ProjectID: "proj-2",
				Cluster:   "project-123_cluster-1",
				Metadata: shared.ProjectMetadata{
					Subscription: struct {
						Availability string `json:"availability"`
						EnvType      string `json:"envType"`
					}{
						Availability: "HA",
						EnvType:      "prod",
					},
					Commerce:    true,
					DocLibStore: "doclib",
					Trial:       "false",
				},
				CreatedAt: 1740926966862,
				Status:    "running",
				Health:    "healthy",
			},
		},
	}

	mockJSON := `{
		"activationHistory": [
			{
				"projectId": "proj-1",
				"serviceId": "liferay",
				"availability": "HA",
				"enabledAt": 1740926966862,
				"enabledByEmail": "user-proj-1@liferay.com",
				"minInstances": 0,
				"maxInstances": 0,
				"disabledAt": 1740926969477,
				"disabledByEmail": "user-proj-1@liferay.com",
				"scaleLimitsChanges": []
			},
			{
				"projectId": "proj-2",
				"serviceId": "liferay",
				"availability": "HA",
				"enabledAt": 1740926966862,
				"enabledByEmail": "user-proj-2@liferay.com",
				"minInstances": 0,
				"maxInstances": 0,
				"disabledAt": 1740926969477,
				"disabledByEmail": "user-proj-2@liferay.com",
				"scaleLimitsChanges": []
			}
		],
		"includedChildProjectIds": [
			"proj-1",
			"proj-2"
		],
		"scaleHistory": [
			{
				"projectId": "proj-1",
				"availability": "HA",
				"serviceId": "liferay",
				"numAdditionalInstances": 1,
				"startedAt": 1740927029636,
				"endedAt": 1740927226281,
				"activeTimePerInstanceMs": 196645,
				"activeTimeMs": 196645
			},
			{
				"projectId": "proj-2",
				"availability": "HA",
				"serviceId": "liferay",
				"numAdditionalInstances": 2,
				"startedAt": 1740927736767,
				"endedAt": 1740927803626,
				"activeTimePerInstanceMs": 66859,
				"activeTimeMs": 66859
			}
		],
		"subtotalsByProjectId": {
			"proj-1": {
				"availability": "HA",
				"totalActiveTimeMs": 263504,
				"billableTimeMs": 3600000,
				"cost": {
					"amount": 12,
					"currency": "USD"
				},
				"price": {
					"amount": 12,
					"currency": "USD"
				}
			},
			"proj-2": {
				"availability": "STD",
				"totalActiveTimeMs": 163501,
				"billableTimeMs": 2000000,
				"cost": {
					"amount": 36,
					"currency": "USD"
				},
				"price": {
					"amount": 12,
					"currency": "USD"
				}
			}
		}
	}`

	dataRange := internal.CalculateDates(1 * time.Hour)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/admin/reports/autoscale/stats", r.URL.Path)
		_, err := fmt.Fprintln(w, mockJSON)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := lcp.NewClient(server.URL, "fake-token")
	collector := NewAutoscaleCollector(client, projectProvider, dataRange)

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

	require.Contains(t, output, `lcp_api_autoscale_activation_history_count{disabled_at="1740926969477",disabled_by="user-proj-1@liferay.com",enabled_at="1740926966862",enabled_by="user-proj-1@liferay.com",project_name="proj-1",service_id="liferay"} 0`)
	require.Contains(t, output, `lcp_api_autoscale_activation_history_count{disabled_at="1740926969477",disabled_by="user-proj-2@liferay.com",enabled_at="1740926966862",enabled_by="user-proj-2@liferay.com",project_name="proj-2",service_id="liferay"} 0`)
	require.Contains(t, output, `lcp_api_autoscale_billable_duration_ms{project_name="proj-1"} 3.6e+06`)
	require.Contains(t, output, `lcp_api_autoscale_billable_duration_ms{project_name="proj-2"} 2e+06`)
	require.Contains(t, output, `lcp_api_autoscale_cost_amount{currency_code="USD",project_name="proj-1"} 12`)
	require.Contains(t, output, `lcp_api_autoscale_cost_amount{currency_code="USD",project_name="proj-2"} 36`)
	require.Contains(t, output, `lcp_api_autoscale_cost_duration_ms{project_name="proj-1"} 263504`)
	require.Contains(t, output, `lcp_api_autoscale_cost_duration_ms{project_name="proj-2"} 163501`)
	require.Contains(t, output, `lcp_api_autoscale_price_amount{currency_code="USD",project_name="proj-1"} 12`)
	require.Contains(t, output, `lcp_api_autoscale_price_amount{currency_code="USD",project_name="proj-2"} 12`)
	require.Contains(t, output, `lcp_api_autoscale_scaling_history_duration_ms{ended_at="1740927226281",instances="1",project_name="proj-1",service_id="liferay",started_at="1740927029636"} 196645`)
	require.Contains(t, output, `lcp_api_autoscale_scaling_history_duration_ms{ended_at="1740927803626",instances="2",project_name="proj-2",service_id="liferay",started_at="1740927736767"} 66859`)
}
