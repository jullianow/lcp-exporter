//nolint:all
package admin

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type autoscaleCollector struct {
	client          *lcp.Client
	projectProvider ProjectProvider
	overview        *prometheus.Desc
	stats           *prometheus.Desc
	dataRange       shared.DateRange
}

func NewAutoscaleCollector(client *lcp.Client, provider ProjectProvider, dataRange shared.DateRange) *autoscaleCollector {
	fqName := internal.Name("autoscale")

	return &autoscaleCollector{
		client:          client,
		projectProvider: provider,
		dataRange:       dataRange,
		overview: prometheus.NewDesc(
			fqName("overview"),
			"Overview autoscale by Parent Project",
			[]string{"currency_code", "num_active_child_projects", "parent_project_id", "total_active_time_ms", "total_billable_time_ms", "total_cost"},
			nil,
		),
		stats: prometheus.NewDesc(
			fqName("stats"),
			"Stats autoscale by Project",
			nil,
			nil,
		),
	}
}

func (ac *autoscaleCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ac.collectMetrics(ch)
	}()
	wg.Wait()
}

func (ac *autoscaleCollector) collectMetrics(ch chan<- prometheus.Metric) {
	projects := ac.projectProvider.GetProjects()
	internal.LogInfo("AutoscaleCollector", "Found %d projects", len(projects))

	if len(projects) == 0 {
		internal.LogInfo("AutoscaleCollector", "No projects found")
		return
	}

	ac.overviewMetric(ch)
	ac.statsMetric(ch, projects)
}

func (ac *autoscaleCollector) overviewMetric(ch chan<- prometheus.Metric) {
	queryParams := map[string]string{
		"start":  ac.dataRange.From,
		"end":    ac.dataRange.End,
		"format": "json",
	}

	overviews, err := lcp.FetchFrom[shared.AutoscaleOverview](ac.client, "/admin/reports/autoscale/overview", queryParams)

	if err != nil {
		internal.LogError("AutoscaleCollector", "Failed to fetch autoscale overview: %v", err)
		return
	}

	if len(overviews) == 0 {
		internal.LogInfo("AutoscaleCollector", "No autoscale overview data found")
		return
	}

	for _, overview := range overviews {
		ch <- prometheus.MustNewConstMetric(
			ac.overview,
			prometheus.GaugeValue,
			1.0,
			overview.CurrencyCode,
			internal.IntToString(overview.NumActiveChildProjects),
			overview.ParentProjectID,
			internal.IntToString(overview.TotalActiveTimeMs),
			internal.IntToString(overview.TotalBillableTimeMs),
			internal.IntToString(overview.TotalCost),
		)
	}
}

func (ac *autoscaleCollector) statsMetric(ch chan<- prometheus.Metric, projects []shared.Projects) {

	var parentProjectIDs []string

	for _, project := range projects {
		if internal.IsParentProject(project) {
			parentProjectIDs = append(parentProjectIDs, project.ProjectID)
		}
	}

	queryParams := map[string]string{
		"start":      ac.dataRange.From,
		"end":        ac.dataRange.End,
		"projectIds": internal.JoinStrings(parentProjectIDs),
	}

	stats, err := lcp.FetchFrom[shared.Autoscale](ac.client, "/admin/reports/autoscale/stats", queryParams)

	if err != nil {
		internal.LogError("AutoscaleCollector", "Failed to fetch autoscale overview: %v", err)
		return
	}

	if len(stats) == 0 {
		internal.LogInfo("AutoscaleCollector", "No autoscale overview data found")
		return
	}

	ch <- prometheus.MustNewConstMetric(
		ac.stats,
		prometheus.GaugeValue,
		1.0,
	)
}

func (ac *autoscaleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ac.overview
	ch <- ac.stats
}
