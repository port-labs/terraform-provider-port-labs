package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationToPortBody_perResourceEnableDelete(t *testing.T) {
	config := `{
		"entityDeletionThreshold": 1,
		"resources": [{
			"kind": "namespace",
			"enableDelete": false,
			"selector": { "query": "true" },
			"port": {
				"entity": {
					"mappings": [{
						"identifier": ".metadata.uid",
						"title": ".metadata.name",
						"blueprint": "'namespace'"
					}]
				}
			}
		}, {
			"kind": "application",
			"selector": { "query": "true" },
			"port": {
				"entity": {
					"mappings": [{
						"identifier": ".metadata.uid",
						"title": ".metadata.name",
						"blueprint": "'argocdApplication'"
					}]
				}
			}
		}]
	}`

	state := &IntegrationModel{
		InstallationId: types.StringValue("my-k8s-exporter"),
		Config:         types.StringValue(config),
	}

	body, err := integrationToPortBody(state)
	require.NoError(t, err)
	require.NotNil(t, body.Config)

	cfg := *body.Config
	assert.Equal(t, float64(1), cfg["entityDeletionThreshold"])

	resources, ok := cfg["resources"].([]interface{})
	require.True(t, ok)
	require.Len(t, resources, 2)

	ns, ok := resources[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "namespace", ns["kind"])
	assert.Equal(t, false, ns["enableDelete"])

	app, ok := resources[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "application", app["kind"])
	_, hasEnableDelete := app["enableDelete"]
	assert.False(t, hasEnableDelete)
}

func TestRefreshIntegrationConfig_preservesEnableDeleteInState(t *testing.T) {
	preferred := types.StringValue(`{
		"entityDeletionThreshold": 1,
		"resources": [{
			"kind": "namespace",
			"enableDelete": false,
			"selector": { "query": "true" },
			"port": {
				"entity": {
					"mappings": [{
						"identifier": ".metadata.uid",
						"title": ".metadata.name",
						"blueprint": "'namespace'"
					}]
				}
			}
		}]
	}`)

	apiConfig := map[string]interface{}{
		"entityDeletionThreshold": 1,
		"resources": []interface{}{
			map[string]interface{}{
				"kind":         "namespace",
				"enableDelete": false,
				"selector": map[string]interface{}{
					"query": "true",
				},
				"port": map[string]interface{}{
					"entity": map[string]interface{}{
						"mappings": []interface{}{
							map[string]interface{}{
								"identifier": ".metadata.uid",
								"title":      ".metadata.name",
								"blueprint":  "'namespace'",
							},
						},
					},
				},
			},
		},
	}

	got, err := utils.GoObjectToTerraformStringPreferExisting(preferred, apiConfig, false)
	require.NoError(t, err)
	assert.Equal(t, preferred, got)
}
