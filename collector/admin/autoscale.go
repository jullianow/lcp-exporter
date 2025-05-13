package admin

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type autoscaleCollector struct {
	client          *lcp.Client
	projectProvider ProjectProvider
	dataRange       shared.DateRange

	activationHistory        *prometheus.Desc
	billableDurationMs       *prometheus.Desc
	costAmount               *prometheus.Desc
	scalingHistoryDurationMs *prometheus.Desc
	priceAmount              *prometheus.Desc
	totalCostDurationMs      *prometheus.Desc
}

func NewAutoscaleCollector(client *lcp.Client, provider ProjectProvider, dataRange shared.DateRange) *autoscaleCollector {
	fqName := internal.Name("autoscale")

	return &autoscaleCollector{
		client:          client,
		projectProvider: provider,
		dataRange:       dataRange,
		activationHistory: prometheus.NewDesc(
			fqName("activation_history_count"),
			"History total instances activated by project and service",
			[]string{"project_name", "service_id", "disabled_at", "disabled_by", "enabled_at", "enabled_by"},
			nil,
		),
		billableDurationMs: prometheus.NewDesc(
			fqName("billable_duration_ms"),
			"Billable time by project in milliseconds",
			[]string{"project_name"},
			nil,
		),
		costAmount: prometheus.NewDesc(
			fqName("cost_amount"),
			"Cost of the autoscale by project",
			[]string{"project_name", "currency_code"},
			nil,
		),
		scalingHistoryDurationMs: prometheus.NewDesc(
			fqName("scaling_history_duration_ms"),
			"History scaling time in milliseconds by project and service",
			[]string{"project_name", "service_id", "started_at", "ended_at", "instances"},
			nil,
		),
		priceAmount: prometheus.NewDesc(
			fqName("price_amount"),
			"Price per hour of autoscale per project",
			[]string{"project_name", "currency_code"},
			nil,
		),
		totalCostDurationMs: prometheus.NewDesc(
			fqName("cost_duration_ms"),
			"Total cost of the autoscale by project",
			[]string{"project_name"},
			nil,
		),
	}
}

func (ac *autoscaleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ac.activationHistory
	ch <- ac.billableDurationMs
	ch <- ac.costAmount
	ch <- ac.priceAmount
	ch <- ac.scalingHistoryDurationMs
	ch <- ac.totalCostDurationMs
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

	if len(projects) == 0 {
		internal.LogWarn("AutoscaleCollector", "No projects found")
		return
	}

	internal.LogDebug("AutoscaleCollector", "Found %d projects", len(projects))

	queryParams := map[string]string{
		"start":      ac.dataRange.From,
		"end":        ac.dataRange.End,
		"projectIds": internal.JoinStrings(internal.GetRootProjectIDs(projects), ","),
	}

	internal.LogDebug(
		"AutoscaleCollector",
		"Fetched autoscale data for projects: %s | start: %s | end: %s",
		queryParams["projectIds"],
		queryParams["start"],
		queryParams["end"],
	)

	stats, err := lcp.FetchFrom[shared.Autoscale](ac.client, "/admin/reports/autoscale/stats", queryParams)

	if err != nil {
		internal.LogError("AutoscaleCollector", "Failed to fetch autoscale overview: %v", err)
		return
	}

	stat := stats[0]
	childProjectIds := stat.IncludedChildProjectIds
	subtotalsByProjectIds := stat.SubtotalsByProjectId
	totalChildProjectIds := len(childProjectIds)
	totalSubtotalsByProjectIds := len(subtotalsByProjectIds)

	if (totalChildProjectIds+totalSubtotalsByProjectIds)%2 == 0 && (totalChildProjectIds+totalSubtotalsByProjectIds) > 0 {
		for _, childProjectId := range childProjectIds {
			subtotalsByProjectId := subtotalsByProjectIds[childProjectId]

			ch <- prometheus.MustNewConstMetric(
				ac.billableDurationMs,
				prometheus.GaugeValue,
				float64(subtotalsByProjectId.BillableTimeMs),
				childProjectId,
			)

			ch <- prometheus.MustNewConstMetric(
				ac.costAmount,
				prometheus.GaugeValue,
				float64(subtotalsByProjectId.Cost.Amount),
				childProjectId,
				subtotalsByProjectId.Cost.Currency,
			)

			ch <- prometheus.MustNewConstMetric(
				ac.priceAmount,
				prometheus.GaugeValue,
				float64(subtotalsByProjectId.Price.Amount),
				childProjectId,
				subtotalsByProjectId.Price.Currency,
			)

			ch <- prometheus.MustNewConstMetric(
				ac.totalCostDurationMs,
				prometheus.GaugeValue,
				float64(subtotalsByProjectId.TotalActiveTimeMs),
				childProjectId,
			)
		}

		internal.LogDebug("AutoscaleCollector", "Found %d scaling history events", len(stat.ScalingHistory))
		for _, event := range stat.ScalingHistory {
			ch <- prometheus.MustNewConstMetric(
				ac.scalingHistoryDurationMs,
				prometheus.GaugeValue,
				float64(event.ActiveTimePerInstanceMs),
				event.ProjectID,
				event.ServiceID,
				internal.IntToString(event.StartedAt),
				internal.IntToString(event.EndedAt),
				internal.IntToString(event.NumAdditionalInstances),
			)
		}

		internal.LogDebug("AutoscaleCollector", "Found %d activation history events", len(stat.ActivationHistory))
		for _, event := range stat.ActivationHistory {
			numInstances := 0.0
			if event.DisabledAt > 0 {
				numInstances = float64(event.MaxInstances)
			}
			ch <- prometheus.MustNewConstMetric(
				ac.activationHistory,
				prometheus.GaugeValue,
				numInstances,
				event.ProjectID,
				event.ServiceID,
				internal.IntToString(event.DisabledAt),
				event.DisabledByEmail,
				internal.IntToString(event.EnabledAt),
				event.EnabledByEmail,
			)
		}

	} else {
		internal.LogWarn(
			"AutoscaleCollector",
			"includedChildProjectIds=%d and subtotalsByProjectIds=%d do not match",
			len(childProjectIds),
			len(subtotalsByProjectIds),
		)
	}
}
