package internal

type HealthCheck struct {
	Status string `json:"status"`
}

type Info struct {
	Version string `json:"version"`
	Domains struct {
		Infa    string `json:"infrastructure"`
		Service string `json:"service"`
	} `json:"domains"`
}

type Provider struct {
	Name           string `json:"name"`
	CloudProjectID string `json:"cloudProjectId"`
}

type ClusterDiscovery struct {
	Name                 string   `json:"name"`
	Provider             Provider `json:"provider"`
	Location             string   `json:"location"`
	CustomerBackupBucket string   `json:"customerBackupBucket"`
	Zones                []string `json:"zones"`
	PlanID               string   `json:"planId"`
	IsLXC                bool     `json:"isLXC"`
}

type ProjectMetadata struct {
	Commerce bool   `json:"commerce"`
	Type     string `json:"type"`
	Trial    string `json:"trial"`
}

type Projects struct {
	Id              string          `json:"id"`
	Cluster         string          `json:"cluster"`
	Health          string          `json:"health"`
	ParentProjectID string          `json:"organizationId"`
	ProjectID       string          `json:"projectId"`
	Status          string          `json:"status"`
	Metadata        ProjectMetadata `json:"metadata"`
	Type            string          `json:"type"`
}
