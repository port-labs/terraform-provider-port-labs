package workflow

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowResourceToPortBodyDelayNode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	nodes := []map[string]any{
		{
			"identifier": "wait",
			"title":      "Wait Before Next Step",
			"config": map[string]any{
				"type":    "DELAY",
				"seconds": 300,
			},
		},
	}
	nodesJSON, err := json.Marshal(nodes)
	require.NoError(t, err)

	state := &WorkflowModel{
		Identifier:  types.StringValue("delayed-action"),
		Title:       types.StringValue("Workflow with Delay"),
		Nodes:       types.StringValue(string(nodesJSON)),
		Connections: types.StringValue("[]"),
	}

	body, err := workflowResourceToPortBody(ctx, state)
	require.NoError(t, err)
	require.Len(t, body.Nodes, 1)

	config, ok := body.Nodes[0]["config"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "DELAY", config["type"])
	assert.Equal(t, float64(300), config["seconds"])
}

func TestWorkflowResourceToPortBodyDelayNodeDynamicExpression(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	nodes := []map[string]any{
		{
			"identifier": "wait",
			"title":      "Wait Before Next Step",
			"config": map[string]any{
				"type":    "DELAY",
				"seconds": "{{ .outputs.trigger.wait_seconds }}",
			},
		},
	}
	nodesJSON, err := json.Marshal(nodes)
	require.NoError(t, err)

	state := &WorkflowModel{
		Identifier:  types.StringValue("delayed-action"),
		Title:       types.StringValue("Workflow with Delay"),
		Nodes:       types.StringValue(string(nodesJSON)),
		Connections: types.StringValue("[]"),
	}

	body, err := workflowResourceToPortBody(ctx, state)
	require.NoError(t, err)
	require.Len(t, body.Nodes, 1)

	config, ok := body.Nodes[0]["config"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "DELAY", config["type"])
	assert.Equal(t, "{{ .outputs.trigger.wait_seconds }}", config["seconds"])
}

func TestRefreshWorkflowStateDelayNode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	state := &WorkflowModel{
		Identifier: types.StringValue("delayed-action"),
		Title:      types.StringValue("Workflow with Delay"),
	}

	portWorkflow := &cli.Workflow{
		Identifier: "delayed-action",
		Title:      "Workflow with Delay",
		Nodes: []map[string]any{
			{
				"identifier": "wait",
				"title":      "Wait Before Next Step",
				"config": map[string]any{
					"type":    "DELAY",
					"seconds": float64(60),
				},
			},
		},
		Connections: []map[string]any{},
	}

	err := refreshWorkflowState(ctx, state, portWorkflow)
	require.NoError(t, err)

	var nodes []map[string]any
	err = json.Unmarshal([]byte(state.Nodes.ValueString()), &nodes)
	require.NoError(t, err)
	require.Len(t, nodes, 1)

	config, ok := nodes[0]["config"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "DELAY", config["type"])
	assert.Equal(t, float64(60), config["seconds"])
}
