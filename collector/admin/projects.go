package admin

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type ProjectProvider interface {
	GetProjects() []shared.Projects
}

type ProjectsCollector struct {
	client   *lcp.Client
	projects []shared.Projects
	mu       sync.RWMutex

	collaborators      *prometheus.Desc
	create             *prometheus.Desc
	dbLabels           *prometheus.Desc
	dbStorageBytes     *prometheus.Desc
	labels             *prometheus.Desc
	status             *prometheus.Desc
	total              *prometheus.Desc
	volumeStorageBytes *prometheus.Desc
}

func NewProjectsCollector(client *lcp.Client) *ProjectsCollector {
	fqName := internal.Name("projects")

	return &ProjectsCollector{
		client: client,
		collaborators: prometheus.NewDesc(
			fqName("collaborators"),
			"Number of collaborators per project",
			[]string{"id"},
			nil,
		),
		create: prometheus.NewDesc(
			fqName("create_timestamp"),
			"Timestamp of the creation per project",
			[]string{"id"},
			nil,
		),
		dbLabels: prometheus.NewDesc(
			fqName("db_labels"),
			"Labels about database per project",
			[]string{
				"id",
				"type",
				"version",
				"edition",
			}, nil,
		),
		dbStorageBytes: prometheus.NewDesc(
			fqName("db_storage_capacity_bytes"),
			"Total database storage capacity per project in bytes",
			[]string{
				"disk_type",
				"id",
			},
			nil,
		),
		labels: prometheus.NewDesc(
			fqName("labels"),
			"Labels about project",
			[]string{
				"availability",
				"cluster_name",
				"commerce",
				"doc_lib_store",
				"health",
				"id",
				"name",
				"parent_project_name",
				"root_project",
				"trial",
				"type",
			},
			nil,
		),
		status: prometheus.NewDesc(
			fqName("status"),
			"Status of project. 1 if is running, 0 otherwise",
			[]string{"id"},
			nil,
		),
		total: prometheus.NewDesc(
			fqName("total"),
			"Total number of projects",
			nil, nil,
		),
		volumeStorageBytes: prometheus.NewDesc(
			fqName("volume_storage_capacity_bytes"),
			"Total storage capacity per project in bytes",
			[]string{"id"},
			nil,
		),
	}
}

func (c *ProjectsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.collaborators
	ch <- c.create
	ch <- c.dbLabels
	ch <- c.dbStorageBytes
	ch <- c.labels
	ch <- c.status
	ch <- c.total
	ch <- c.volumeStorageBytes
}

func (pc *ProjectsCollector) GetProjects() []shared.Projects {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.projects
}

func (pc *ProjectsCollector) FetchInitial() {
	projects := pc.fetch()
	if len(projects) == 0 {
		internal.LogWarn("ProjectsCollector", "Initial fetch returned 0 projects")
	}

	pc.mu.Lock()
	pc.projects = projects
	pc.mu.Unlock()
}

func (pc *ProjectsCollector) fetch() []shared.Projects {
	projects, err := lcp.FetchFrom[shared.Projects](pc.client, "/admin/projects", nil)
	if err != nil {
		internal.LogError("ProjectsCollector", "Failed to fetch projects: %v", err)
		return nil
	}
	return projects
}

func (pc *ProjectsCollector) Collect(ch chan<- prometheus.Metric) {
	projects := pc.fetch()

	pc.mu.Lock()
	pc.projects = projects
	pc.mu.Unlock()

	ch <- prometheus.MustNewConstMetric(
		pc.total,
		prometheus.GaugeValue,
		float64(len(projects)),
	)

	for _, project := range projects {
		var status float64
		if project.Status == "running" {
			status = 1
		}

		health := project.Health == "healthy"
		isRootProject := false
		rootProjectName := internal.RootProjectName(project)
		if rootProjectName == "" {
			isRootProject = true
		}

		ch <- prometheus.MustNewConstMetric(
			pc.collaborators,
			prometheus.GaugeValue,
			float64(len(project.Collaborators)),
			project.Id,
		)

		ch <- prometheus.MustNewConstMetric(
			pc.create,
			prometheus.GaugeValue,
			internal.MillisToSeconds(project.CreatedAt),
			project.Id,
		)

		ch <- prometheus.MustNewConstMetric(
			pc.labels,
			prometheus.GaugeValue,
			1.0,
			project.Metadata.Subscription.Availability,
			project.Cluster,
			internal.BoolToString(project.Metadata.Commerce),
			project.Metadata.DocLibStore,
			internal.BoolToString(health),
			project.Id,
			project.ProjectID,
			rootProjectName,
			internal.BoolToString(isRootProject),
			project.Metadata.Trial,
			project.Metadata.Subscription.EnvType,
		)

		ch <- prometheus.MustNewConstMetric(
			pc.status,
			prometheus.GaugeValue,
			status,
			project.Id,
		)

		ch <- prometheus.MustNewConstMetric(
			pc.volumeStorageBytes,
			prometheus.GaugeValue,
			float64(internal.GiBToBytes(project.VolumeStorageSize)),
			project.Id,
		)

		if project.CloudOptions != (shared.ProjectCloudOptions{}) && !isRootProject {
			ch <- prometheus.MustNewConstMetric(
				pc.dbLabels,
				prometheus.GaugeValue,
				1.0,
				project.Id,
				project.CloudOptions.InstanceType,
				project.CloudOptions.DatabaseVersion,
				project.CloudOptions.DatabaseEdition,
			)

			if project.CloudOptions.DiskSize != "" {
				ch <- prometheus.MustNewConstMetric(
					pc.dbStorageBytes,
					prometheus.GaugeValue,
					float64(internal.GBToBytes(internal.StringToInt64(project.CloudOptions.DiskSize))),
					project.CloudOptions.DiskType,
					project.Id,
				)
			}
		}
	}
}
