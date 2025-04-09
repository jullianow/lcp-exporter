package shared

type HealthCheck struct {
	Status string `json:"status"`
}

type Info struct {
	Version string `json:"version"`
	Domains struct {
		Infrastructure string `json:"infrastructure"`
		Service        string `json:"service"`
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
	CreatedAt       int             `json:"createdAt"`
}

type AutoscaleCost struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type AutoscaleHistory struct {
	ProjectID                 string        `json:"projectId"`
	Availability              string        `json:"availability"`
	ServiceId                 string        `json:"serviceId"`
	NumAddedInstances         int           `json:"numAddedInstances"`
	ActiveTimePerInstanceMs   int           `json:"activeTimePerInstanceMs"`
	ActiveTimeMs              int           `json:"activeTimeMs"`
	BillableTimePerInstanceMs int           `json:"billableTimePerInstanceMs"`
	BillableTimeMs            int           `json:"billableTimeMs"`
	Cost                      AutoscaleCost `json:"cost"`
	Price                     AutoscaleCost `json:"price"`
}

type Autoscale struct {
	ProjectIds              []string           `json:"projectIds"`
	IncludedChildProjectIds []string           `json:"includedChildProjectIds"`
	CurrencyCode            string             `json:"currencyCode"`
	TotalActiveTimeMs       string             `json:"totalActiveTimeMs"`
	TotalBillableTimeMs     string             `json:"totalBillableTimeMs"`
	AutoscaleHistory        []AutoscaleHistory `json:"autoscaleHistory"`
}
