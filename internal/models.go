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

type ClusterDiscovery struct {
	Name                 string   `json:"name"`
	Provider             Provider `json:"provider"`
	Location             string   `json:"location"`
	CustomerBackupBucket string   `json:"customerBackupBucket"`
	Zones                []string `json:"zones"`
	PlanID               string   `json:"planId"`
	IsLXC                bool     `json:"isLXC"`
}

type Provider struct {
	Name           string `json:"name"`
	CloudProjectID string `json:"cloudProjectId"`
}
