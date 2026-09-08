package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSpec_DropsServerManagedSections(t *testing.T) {
	raw := `{
		"installationId": "pagerduty-prod",
		"installationType": "Saas",
		"spec": {
			"integrationSpec": {"token": "my-secret"},
			"appSpec": {"scheduledResyncInterval": "12h"},
			"systemSpec": {"size": "L"},
			"privateSpec": {"applierBackend": "operator"}
		}
	}`

	var got Integration
	require.NoError(t, json.Unmarshal([]byte(raw), &got))

	require.NotNil(t, got.Spec)
	assert.Equal(t, "my-secret", got.Spec.IntegrationSpec["token"])
	assert.Equal(t, "12h", got.Spec.AppSpec["scheduledResyncInterval"])

	marshaled, err := json.Marshal(got)
	require.NoError(t, err)
	assert.NotContains(t, string(marshaled), "systemSpec")
	assert.NotContains(t, string(marshaled), "privateSpec")
}

func TestIntegrationSpec_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{"absent spec", `{"installationId": "my-kafka"}`, true},
		{"only server-managed sections", `{"spec": {"systemSpec": {"size": "M"}}}`, true},
		{"manageable section present", `{"spec": {"integrationSpec": {"token": "x"}}}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Integration
			require.NoError(t, json.Unmarshal([]byte(tt.raw), &got))
			assert.Equal(t, tt.want, got.Spec.IsEmpty())
		})
	}
}
