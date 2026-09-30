package integration

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationToPortBody_enableDeleteOnResource(t *testing.T) {
	config := `{
		"entityDeletionThreshold": 1,
		"resources": [
			{
				"kind": "namespace",
				"enableDelete": false,
				"selector": {"query": "true"},
				"port": {
					"entity": {
						"mappings": [{
							"identifier": ".metadata.uid",
							"title": ".metadata.name",
							"blueprint": "\"namespace\""
						}]
					}
				}
			},
			{
				"kind": "application",
				"selector": {"query": "true"},
				"port": {
					"entity": {
						"mappings": [{
							"identifier": ".metadata.uid",
							"title": ".metadata.name",
							"blueprint": "\"argocdApplication\""
						}]
					}
				}
			}
		]
	}`

	state := &IntegrationModel{
		InstallationId: types.StringValue("my-k8s-exporter"),
		Config:         types.StringValue(config),
	}

	body, err := integrationToPortBody(state, false)
	require.NoError(t, err)
	require.NotNil(t, body.Config)

	assert.Equal(t, float64(1), (*body.Config)["entityDeletionThreshold"])

	resources, ok := (*body.Config)["resources"].([]interface{})
	require.True(t, ok)
	require.Len(t, resources, 2)

	namespaceResource, ok := resources[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "namespace", namespaceResource["kind"])
	assert.Equal(t, false, namespaceResource["enableDelete"])

	applicationResource, ok := resources[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "application", applicationResource["kind"])
	_, hasEnableDelete := applicationResource["enableDelete"]
	assert.False(t, hasEnableDelete)
}

func TestRefreshIntegrationConfig_preservesEnableDelete(t *testing.T) {
	preferredConfig := `{
		"entityDeletionThreshold": 1,
		"resources": [
			{
				"kind": "namespace",
				"enableDelete": false,
				"selector": {"query": "true"},
				"port": {
					"entity": {
						"mappings": [{
							"identifier": ".metadata.uid",
							"title": ".metadata.name",
							"blueprint": "\"namespace\""
						}]
					}
				}
			}
		]
	}`

	var apiConfig map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(preferredConfig), &apiConfig))

	got, err := utils.GoObjectToTerraformStringPreferExisting(
		types.StringValue(preferredConfig),
		apiConfig,
		false,
	)
	require.NoError(t, err)
	assert.Equal(t, types.StringValue(preferredConfig), got)
}
