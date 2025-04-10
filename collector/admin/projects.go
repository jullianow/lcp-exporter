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

	count *prometheus.Desc
	info  *prometheus.Desc
	age   *prometheus.Desc
}

func NewProjectsCollector(client *lcp.Client) *ProjectsCollector {
	fqName := internal.Name("projects")

	return &ProjectsCollector{
		client: client,
		count: prometheus.NewDesc(
			fqName("count"),
			"Total number of projects",
			nil, nil,
		),
		info: prometheus.NewDesc(
			fqName("info"),
			"Information about projects. 1 if the project is running, 0 otherwise",
			[]string{"commerce", "health", "id", "name", "parent_project_id", "root_project", "trial", "type"},
			nil,
		),
		age: prometheus.NewDesc(
			fqName("age"),
			"Information about projects. 1 if the project is running, 0 otherwise",
			[]string{"id"},
			nil,
		),
	}
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
		pc.count,
		prometheus.GaugeValue,
		float64(len(projects)),
	)

	for _, project := range projects {
		var status float64
		if project.Status == "running" {
			status = 1
		}

		health := project.Health == "healthy"
		rootProject := internal.IsParentProject(project)
		parentID := project.ParentProjectID
		if rootProject {
			parentID = ""
		}

		ch <- prometheus.MustNewConstMetric(
			pc.info,
			prometheus.GaugeValue,
			status,
			internal.BoolToString(project.Metadata.Commerce),
			internal.BoolToString(health),
			project.Id,
			project.ProjectID,
			parentID,
			internal.BoolToString(rootProject),
			project.Metadata.Trial,
			project.Metadata.Type,
		)

		ch <- prometheus.MustNewConstMetric(
			pc.age,
			prometheus.GaugeValue,
			float64(project.CreatedAt),
			project.Id,
		)
	}
}

func (c *ProjectsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.count
	ch <- c.info
	ch <- c.age
}
