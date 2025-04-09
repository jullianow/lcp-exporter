//nolint:all
package admin

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/jullianow/lcp-exporter/internal"
	// "github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type AutoscaleCollector struct {
	client          *lcp.Client
	projectProvider ProjectProvider
}

func NewAutoscaleCollector(client *lcp.Client, provider ProjectProvider) *AutoscaleCollector {
	// fqName := internal.Name("autoscale")

	return &AutoscaleCollector{
		client:          client,
		projectProvider: provider,
		// count: prometheus.NewDesc(
		// 	fqName("count"),
		// 	"Total number of projects",
		// 	nil, nil,
		// ),
		// info: prometheus.NewDesc(
		// 	fqName("info"),
		// 	"Information about projects. 1 if the project is running, 0 otherwise",
		// 	[]string{"commerce", "health", "id", "name", "parent_project_id", "root_project", "trial", "type"},
		// 	nil,
		// ),
		// age: prometheus.NewDesc(
		// 	fqName("age"),
		// 	"Information about projects. 1 if the project is running, 0 otherwise",
		// 	[]string{"id"},
		// 	nil,
		// ),
	}
}

func (ac *AutoscaleCollector) Collect(ch chan<- prometheus.Metric) {
	projects := ac.projectProvider.GetProjects()
	internal.LogInfo("AutoscaleCollector", "Found %d projects", len(projects))

}

func (ac *AutoscaleCollector) Describe(ch chan<- *prometheus.Desc) {
	// Stub method for prometheus.Collector interface
}
