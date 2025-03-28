package internal

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	tests := []struct {
		component string
		metric    string
		expected  string
	}{
		{"status", "info", "lcp_api_status_info"},
		{"cluster_discovery", "total", "lcp_api_cluster_discovery_total"},
		{"cluster_discovery", "info", "lcp_api_cluster_discovery_info"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.component, tt.metric), func(t *testing.T) {
			result := Name(tt.component)(tt.metric)
			assert.Equal(t, tt.expected, result)
		})
	}
}
