package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowSchemaIsValid(t *testing.T) {
	ctx := context.Background()
	r := &WorkflowResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, resp)

	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())

	diags := resp.Schema.ValidateImplementation(ctx)
	require.False(t, diags.HasError(), diags.Errors())
}

func strPtr(s string) *string { return &s }

func TestWorkflowStateToPortBody(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier: types.StringValue("review-pr"),
		Title:      types.StringValue("Review PR"),
		Category:   types.StringValue("engineering"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("trigger"),
				Title:      types.StringValue("PR Trigger"),
				EventTrigger: &EventTriggerModel{
					Type:                types.StringValue(consts.EntityUpdated),
					BlueprintIdentifier: types.StringValue("githubPullRequest"),
				},
			},
			{
				Identifier: types.StringValue("decision"),
				Title:      types.StringValue("Review"),
				CursorAgent: &CursorAgentModel{
					ApiKey: types.StringValue("{{ .secrets[\"CURSOR_API_KEY\"] }}"),
					Prompt: &CursorPromptModel{Text: types.StringValue("Review this PR.")},
					Source: &CursorSourceModel{PrUrl: types.StringValue("{{ .outputs.trigger.link }}")},
				},
			},
		},
		Connections: []ConnectionModel{
			{
				SourceIdentifier: types.StringValue("trigger"),
				TargetIdentifier: types.StringValue("decision"),
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	assert.Equal(t, "review-pr", w.Identifier)
	require.NotNil(t, w.Category)
	assert.Equal(t, "engineering", *w.Category)

	require.Len(t, w.Nodes, 2)
	assert.Equal(t, consts.EventTrigger, w.Nodes[0].Config.Type)
	require.NotNil(t, w.Nodes[0].Config.Event)
	assert.Equal(t, consts.EntityUpdated, w.Nodes[0].Config.Event.Type)
	assert.Equal(t, "githubPullRequest", w.Nodes[0].Config.Event.BlueprintIdentifier)

	assert.Equal(t, consts.CursorAgent, w.Nodes[1].Config.Type)
	require.NotNil(t, w.Nodes[1].Config.Prompt)
	assert.Equal(t, "Review this PR.", *w.Nodes[1].Config.Prompt.Text)
	require.NotNil(t, w.Nodes[1].Config.Source)
	assert.Equal(t, "{{ .outputs.trigger.link }}", *w.Nodes[1].Config.Source.PrUrl)

	require.Len(t, w.Connections, 1)
	assert.Equal(t, "trigger", w.Connections[0].SourceIdentifier)
	assert.Equal(t, "decision", w.Connections[0].TargetIdentifier)
}

func TestWorkflowStateToPortBodyIntegrationAction(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("run"),
				IntegrationAction: &IntegrationActionModel{
					IntegrationProvider:       types.StringValue("github"),
					IntegrationInvocationType: types.StringValue("WORKFLOW_DISPATCH"),
					ExecutionProperties:       types.StringValue(`{"workflow":"deploy.yml","ref":"main"}`),
				},
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	require.Len(t, w.Nodes, 1)
	config := w.Nodes[0].Config
	assert.Equal(t, consts.IntegrationAction, config.Type)
	require.NotNil(t, config.IntegrationProvider)
	assert.Equal(t, "github", *config.IntegrationProvider)
	require.NotNil(t, config.IntegrationInvocationType)
	assert.Equal(t, "WORKFLOW_DISPATCH", *config.IntegrationInvocationType)
	props, ok := config.IntegrationActionExecutionProperties.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "deploy.yml", props["workflow"])
	assert.Equal(t, "main", props["ref"])
}

func TestRefreshWorkflowStatePreservesOrderAndSecret(t *testing.T) {
	ctx := context.Background()
	r := &WorkflowResource{portClient: &cli.PortClient{JSONEscapeHTML: true}}

	// Prior state defines node order [trigger, decision] and holds the api_key
	// secret that the API does not return.
	state := &WorkflowModel{
		Identifier: types.StringValue("review-pr"),
		Nodes: []WorkflowNodeModel{
			{Identifier: types.StringValue("trigger")},
			{
				Identifier:  types.StringValue("decision"),
				CursorAgent: &CursorAgentModel{ApiKey: types.StringValue("super-secret")},
			},
		},
		Connections: []ConnectionModel{
			{SourceIdentifier: types.StringValue("trigger"), TargetIdentifier: types.StringValue("decision")},
		},
	}

	// API returns nodes in a different order and without the api_key.
	apiWorkflow := &cli.Workflow{
		Identifier: "review-pr",
		Category:   strPtr("engineering"),
		Nodes: []cli.WorkflowNode{
			{Identifier: "decision", Config: cli.WorkflowNodeConfig{Type: consts.CursorAgent, Prompt: &cli.CursorAgentPrompt{Text: strPtr("Review this PR.")}}},
			{Identifier: "trigger", Config: cli.WorkflowNodeConfig{Type: consts.EventTrigger, Event: &cli.WorkflowTriggerEvent{Type: consts.EntityUpdated, BlueprintIdentifier: "githubPullRequest"}}},
		},
		Connections: []cli.WorkflowConnection{
			{SourceIdentifier: "trigger", TargetIdentifier: "decision"},
		},
	}

	err := r.refreshWorkflowState(ctx, state, apiWorkflow)
	require.NoError(t, err)

	// Order preserved from prior state.
	require.Len(t, state.Nodes, 2)
	assert.Equal(t, "trigger", state.Nodes[0].Identifier.ValueString())
	assert.Equal(t, "decision", state.Nodes[1].Identifier.ValueString())

	// Event trigger classified from discriminator.
	require.NotNil(t, state.Nodes[0].EventTrigger)
	assert.Equal(t, consts.EntityUpdated, state.Nodes[0].EventTrigger.Type.ValueString())
	assert.Equal(t, "githubPullRequest", state.Nodes[0].EventTrigger.BlueprintIdentifier.ValueString())

	// Secret preserved from prior state.
	require.NotNil(t, state.Nodes[1].CursorAgent)
	assert.Equal(t, "super-secret", state.Nodes[1].CursorAgent.ApiKey.ValueString())
	require.NotNil(t, state.Nodes[1].CursorAgent.Prompt)
	assert.Equal(t, "Review this PR.", state.Nodes[1].CursorAgent.Prompt.Text.ValueString())

	// Category refreshed from the API.
	assert.Equal(t, "engineering", state.Category.ValueString())
}
