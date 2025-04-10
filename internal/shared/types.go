package shared

type DateRange struct {
	From string
	End  string
}

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
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type AutoscaleHistory struct {
	ProjectID                 string        `json:"projectId"`
	Availability              string        `json:"availability"`
	ServiceID                 string        `json:"serviceId"`
	NumAddedInstances         int           `json:"numAddedInstances"`
	ActiveTimePerInstanceMs   int64         `json:"activeTimePerInstanceMs"`
	ActiveTimeMs              int64         `json:"activeTimeMs"`
	BillableTimePerInstanceMs int64         `json:"billableTimePerInstanceMs"`
	BillableTimeMs            int64         `json:"billableTimeMs"`
	Cost                      AutoscaleCost `json:"cost"`
	Price                     AutoscaleCost `json:"price"`
}

type Autoscale struct {
	AutoscaleHistory        []AutoscaleHistory `json:"autoscaleHistory"`
	CurrencyCode            string             `json:"currencyCode"`
	IncludedChildProjectIds []string           `json:"includedChildProjectIds"`
	NumActiveChildProjects  int                `json:"numActiveChildProjects"`
	ProjectIds              []string           `json:"projectIds"`
	TotalActiveTimeMs       int64              `json:"totalActiveTimeMs"`
	TotalBillableTimeMs     int64              `json:"totalBillableTimeMs"`
}

type AutoscaleOverview struct {
	CurrencyCode           string  `json:"currencyCode"`
	NumActiveChildProjects int     `json:"numActiveChildProjects"`
	ParentProjectID        string  `json:"parentProjectId"`
	TotalActiveTimeMs      int64   `json:"totalActiveTimeMs"`
	TotalBillableTimeMs    int64   `json:"totalBillableTimeMs"`
	TotalCost              float64 `json:"totalCost"`
}
