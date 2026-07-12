package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowResourceToPortBody(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tags, diags := types.ListValueFrom(ctx, types.StringType, []string{"lorekeeper", "cursor-agent"})
	require.False(t, diags.HasError())

	state := &WorkflowModel{
		Identifier:            types.StringValue("lorekeeper-decision"),
		Title:                 types.StringValue("LoreKeeper Decision"),
		Description:           types.StringValue("Review a PR using CURSOR_AGENT"),
		AllowAnyoneToViewRuns: types.BoolValue(false),
		Tags:                  tags,
		Nodes: types.StringValue(`[
			{
				"identifier": "decision",
				"title": "LoreKeeper Decision",
				"config": {
					"type": "CURSOR_AGENT",
					"apiKey": "{{ .secrets[\"CURSOR_API_KEY\"] }}",
					"prompt": { "text": "Review the PR" },
					"source": { "prUrl": "{{ .outputs.trigger.diff.after.properties.link }}" }
				}
			}
		]`),
		Connections: types.StringValue(`[
			{
				"sourceIdentifier": "trigger",
				"targetIdentifier": "decision"
			}
		]`),
	}

	body, err := workflowResourceToPortBody(ctx, state)
	require.NoError(t, err)

	assert.Equal(t, "lorekeeper-decision", body.Identifier)
	assert.Equal(t, "LoreKeeper Decision", body.Title)
	require.NotNil(t, body.Description)
	assert.Equal(t, "Review a PR using CURSOR_AGENT", *body.Description)
	require.NotNil(t, body.AllowAnyoneToViewRuns)
	assert.False(t, *body.AllowAnyoneToViewRuns)
	assert.Equal(t, []string{"lorekeeper", "cursor-agent"}, body.Tags)
	require.Len(t, body.Nodes, 1)
	assert.Equal(t, "CURSOR_AGENT", body.Nodes[0]["config"].(map[string]any)["type"])
	require.Len(t, body.Connections, 1)
	assert.Equal(t, "trigger", body.Connections[0]["sourceIdentifier"])
}

func TestRefreshWorkflowState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	state := &WorkflowModel{}
	allowAnyoneToViewRuns := true
	description := "LoreKeeper workflow"
	icon := "Robot"

	err := refreshWorkflowState(ctx, state, &cli.Workflow{
		Identifier:            "lorekeeper-decision",
		Title:                 "LoreKeeper Decision",
		Icon:                  &icon,
		Description:           &description,
		Tags:                  []string{"lorekeeper"},
		AllowAnyoneToViewRuns: &allowAnyoneToViewRuns,
		Nodes: []map[string]any{
			{
				"identifier": "decision",
				"config": map[string]any{
					"type": "CURSOR_AGENT",
				},
			},
		},
		Connections: []map[string]any{
			{
				"sourceIdentifier": "trigger",
				"targetIdentifier": "decision",
			},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, "lorekeeper-decision", state.ID.ValueString())
	assert.Equal(t, "LoreKeeper Decision", state.Title.ValueString())
	assert.Equal(t, "Robot", state.Icon.ValueString())
	assert.Equal(t, "LoreKeeper workflow", state.Description.ValueString())
	assert.True(t, state.AllowAnyoneToViewRuns.ValueBool())
	assert.Contains(t, state.Nodes.ValueString(), "CURSOR_AGENT")
	assert.Contains(t, state.Connections.ValueString(), "targetIdentifier")
}
