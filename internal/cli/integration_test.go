package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationUnmarshalJSON_StripsServerManagedSpec(t *testing.T) {
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

func TestIntegrationUnmarshalJSON_NilSpec(t *testing.T) {
	raw := `{"installationId": "my-kafka"}`

	var got Integration
	require.NoError(t, json.Unmarshal([]byte(raw), &got))
	assert.Nil(t, got.Spec)
}

func TestIntegrationUnmarshalJSON_OnlyServerManagedSpec(t *testing.T) {
	raw := `{
		"installationId": "pagerduty-prod",
		"spec": {
			"systemSpec": {"size": "M"}
		}
	}`

	var got Integration
	require.NoError(t, json.Unmarshal([]byte(raw), &got))
	assert.Nil(t, got.Spec, "spec with only server-managed fields should be nil")
}
