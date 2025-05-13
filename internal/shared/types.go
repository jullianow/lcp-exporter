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
	Kubeconfig           struct {
		Cluster struct {
			CaData string `json:"caData"`
		} `json:"cluster"`
	} `json:"kubeconfig"`
}

type ProjectCloudOptions struct {
	DatabaseEdition string `json:"gcpDatabaseEdition"`
	DatabaseVersion string `json:"gcpDatabaseVersion"`
	DiskSize        string `json:"gcpDiskSize"`
	DiskType        string `json:"gcpDiskType"`
	InstanceType    string `json:"gcpInstanceType"`
}

type ProjectMetadata struct {
	Commerce     bool   `json:"commerce"`
	DocLibStore  string `json:"documentLibraryStore"`
	Trial        string `json:"trial"`
	Subscription struct {
		Availability string `json:"availability"`
		EnvType      string `json:"envType"`
	} `json:"subscription"`
}

type Projects struct {
	CloudOptions      ProjectCloudOptions `json:"cloudOptions"`
	Cluster           string              `json:"cluster"`
	Collaborators     []string            `json:"collaborators"`
	CreatedAt         int64               `json:"createdAt"`
	Health            string              `json:"health"`
	Id                string              `json:"id"`
	Metadata          ProjectMetadata     `json:"metadata"`
	OrganizationId    string              `json:"organizationId"`
	ProjectID         string              `json:"projectId"`
	Status            string              `json:"status"`
	VolumeStorageSize int64               `json:"volumeStorageSize"`
}

type AutoscaleCost struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type AutoscaleActivationHistory struct {
	DisabledAt      int64  `json:"disabledAt"`
	DisabledByEmail string `json:"disabledByEmail"`
	EnabledAt       int64  `json:"enabledAt"`
	EnabledByEmail  string `json:"enabledByEmail"`
	ProjectID       string `json:"projectId"`
	ServiceID       string `json:"serviceId"`
	MaxInstances    int    `json:"maxInstances"`
}

type AutoscaleScalingHistory struct {
	ActiveTimePerInstanceMs int64  `json:"activeTimePerInstanceMs"`
	EndedAt                 int64  `json:"endedAt"`
	NumAdditionalInstances  int    `json:"numAdditionalInstances"`
	ProjectID               string `json:"projectId"`
	ServiceID               string `json:"serviceId"`
	StartedAt               int64  `json:"startedAt"`
}

type AutoscaleProject struct {
	Availability      string        `json:"availability"`
	BillableTimeMs    int64         `json:"billableTimeMs"`
	Cost              AutoscaleCost `json:"cost"`
	Price             AutoscaleCost `json:"price"`
	TotalActiveTimeMs int64         `json:"totalActiveTimeMs"`
}

type Autoscale struct {
	ActivationHistory       []AutoscaleActivationHistory `json:"activationHistory"`
	IncludedChildProjectIds []string                     `json:"includedChildProjectIds"`
	ScalingHistory          []AutoscaleScalingHistory    `json:"scaleHistory"`
	SubtotalsByProjectId    map[string]AutoscaleProject  `json:"subtotalsByProjectId"`
}
