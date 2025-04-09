package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/jullianow/lcp-exporter/internal"
	"github.com/jullianow/lcp-exporter/internal/shared"
	"github.com/jullianow/lcp-exporter/lcp"
)

type projectsCollector struct {
	client *lcp.Client
	count  *prometheus.Desc
	info   *prometheus.Desc
}

func NewProjectsCollector(client *lcp.Client) *projectsCollector {
	fqName := internal.Name("projects")

	return &projectsCollector{
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
	}
}

func (c *projectsCollector) Collect(ch chan<- prometheus.Metric) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectMetrics(ch)
	}()
	wg.Wait()
}

func (c *projectsCollector) collectMetrics(ch chan<- prometheus.Metric) {
	resp, err := c.client.MakeRequest("/admin/projects")
	if err != nil {
		logrus.Errorf("Error collecting projects data: %v", err)
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Warnf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Error accessing projects API: StatusCode %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error reading projects response: %v", err)
		return
	}

	var projects []shared.Projects
	if err := json.Unmarshal(body, &projects); err != nil {
		logrus.Errorf("Error decoding cluster discovery JSON: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.count,
		prometheus.GaugeValue,
		float64(len(projects)),
	)

	for _, project := range projects {
		var status = 0.0
		if project.Status == "running" {
			status = 1.0
		}

		var health = false
		if project.Health == "healthy" {
			health = true
		}

		var parent_project_id = project.ParentProjectID
		var root_project = false
		if project.ProjectID == project.ParentProjectID {
			parent_project_id = ""
			root_project = true
		}

		ch <- prometheus.MustNewConstMetric(
			c.info,
			prometheus.GaugeValue,
			status,
			internal.BoolToString(project.Metadata.Commerce),
			internal.BoolToString(health),
			project.Id,
			project.ProjectID,
			parent_project_id,
			internal.BoolToString(root_project),
			project.Metadata.Trial,
			project.Metadata.Type,
		)
	}
}

func (c *projectsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.count
	ch <- c.info
}
