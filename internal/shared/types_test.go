package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	hc := HealthCheck{
		Status: "ok",
	}

	assert.Equal(t, "ok", hc.Status)
}

func TestInfo(t *testing.T) {
	info := Info{
		Version: "v1.0.0",
		Domains: struct {
			Infrastructure string `json:"infrastructure"`
			Service        string `json:"service"`
		}{
			Infrastructure: "infrastructure.com",
			Service:        "service.com",
		},
	}

	assert.Equal(t, "v1.0.0", info.Version)
	assert.Equal(t, "infrastructure.com", info.Domains.Infrastructure)
	assert.Equal(t, "service.com", info.Domains.Service)
}

func TestClusterDiscovery(t *testing.T) {
	cluster := ClusterDiscovery{
		Name:                 "cluster-1",
		Provider:             Provider{Name: "gcp", CloudProjectID: "project-123"},
		Location:             "us-central1",
		CustomerBackupBucket: "backup-bucket",
		PlanID:               "plan-xyz",
		IsLXC:                true,
	}

	assert.Equal(t, "cluster-1", cluster.Name)
	assert.Equal(t, "gcp", cluster.Provider.Name)
	assert.Equal(t, "project-123", cluster.Provider.CloudProjectID)
	assert.Equal(t, "us-central1", cluster.Location)
	assert.Equal(t, "backup-bucket", cluster.CustomerBackupBucket)
	assert.Equal(t, "plan-xyz", cluster.PlanID)
	assert.True(t, cluster.IsLXC)
}

func TestProvider(t *testing.T) {
	provider := Provider{
		Name:           "gcp",
		CloudProjectID: "project-123",
	}

	assert.Equal(t, "gcp", provider.Name)
	assert.Equal(t, "project-123", provider.CloudProjectID)
}
