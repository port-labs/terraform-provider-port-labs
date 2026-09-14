package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationToPortBody_githubOceanPullRequestGraphQLSelector(t *testing.T) {
	config := `{
		"resources": [{
			"kind": "pull-request",
			"selector": {
				"query": "true",
				"states": ["open"],
				"api": "graphql",
				"enrichWithFirstCommit": true,
				"excludeGraphqlFields": ["additions", "deletions", "changedFiles"]
			},
			"port": {
				"entity": {
					"mappings": [{
						"identifier": ".head.repo.name + (.fullDatabaseId|tostring)",
						"blueprint": "'githubPullRequest'",
						"properties": {
							"isDraft": ".isDraft",
							"reviewDecision": ".reviewDecision",
							"firstCommitAt": ".firstCommit.committedDate"
						},
						"relations": {}
					}]
				}
			}
		}]
	}`

	state := &IntegrationModel{
		InstallationId: types.StringValue("my-github-ocean"),
		Config:         types.StringValue(config),
	}

	body, err := integrationToPortBody(state)
	require.NoError(t, err)
	require.NotNil(t, body.Config)

	resources, ok := (*body.Config)["resources"].([]interface{})
	require.True(t, ok)
	require.Len(t, resources, 1)

	resource0, ok := resources[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "pull-request", resource0["kind"])

	selector, ok := resource0["selector"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "graphql", selector["api"])
	assert.Equal(t, true, selector["enrichWithFirstCommit"])

	excluded, ok := selector["excludeGraphqlFields"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, []interface{}{"additions", "deletions", "changedFiles"}, excluded)
}
