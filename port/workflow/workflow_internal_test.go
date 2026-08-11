package workflow

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
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
func boolPtr(b bool) *bool    { return &b }

func stringList(values ...string) types.List {
	elements := make([]types.String, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}

	list, diags := types.ListValueFrom(context.Background(), types.StringType, elements)
	if diags.HasError() {
		panic(diags.Errors())
	}
	return list
}

func TestWorkflowStateToPortBody(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier:            types.StringValue("review-pr"),
		Title:                 types.StringValue("Review PR"),
		Category:              types.StringValue("engineering"),
		AllowAnyoneToViewRuns: types.BoolValue(false),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("trigger"),
				Title:      types.StringValue("PR Trigger"),
				Verbose:    types.BoolValue(true),
				Links:      stringList("https://example.com/{{ .result.id }}"),
				EventTrigger: &EventTriggerModel{
					Type:                types.StringValue(consts.EntityUpdated),
					BlueprintIdentifier: types.StringValue("githubPullRequest"),
					Published:           types.BoolValue(true),
					Condition: &NodeConditionModel{
						Expressions: stringList(".diff.after.properties.open == true"),
						Combinator:  types.StringValue("and"),
					},
				},
			},
			{
				Identifier: types.StringValue("notify"),
				Webhook: &WebhookModel{
					Url:       types.StringValue("https://example.com/hook"),
					Method:    types.StringValue("POST"),
					Body:      types.StringValue(`{"text":"hello"}`),
					OnFailure: types.StringValue("continue"),
				},
			},
		},
		Connections: []ConnectionModel{
			{
				SourceIdentifier: types.StringValue("trigger"),
				TargetIdentifier: types.StringValue("notify"),
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	assert.Equal(t, "review-pr", w.Identifier)
	require.NotNil(t, w.Category)
	assert.Equal(t, "engineering", *w.Category)
	require.NotNil(t, w.AllowAnyoneToViewRuns)
	assert.False(t, *w.AllowAnyoneToViewRuns)

	require.Len(t, w.Nodes, 2)

	trigger := w.Nodes[0]
	assert.Equal(t, consts.EventTrigger, trigger.Config.Type)
	require.NotNil(t, trigger.Config.Event)
	assert.Equal(t, consts.EntityUpdated, trigger.Config.Event.Type)
	assert.Equal(t, "githubPullRequest", trigger.Config.Event.BlueprintIdentifier)
	require.NotNil(t, trigger.Verbose)
	assert.True(t, *trigger.Verbose)
	assert.Equal(t, []string{"https://example.com/{{ .result.id }}"}, trigger.Links)
	require.NotNil(t, trigger.Config.Condition)
	assert.Equal(t, consts.JqCondition, trigger.Config.Condition.Type)
	assert.Equal(t, []string{".diff.after.properties.open == true"}, trigger.Config.Condition.Expressions)

	webhook := w.Nodes[1]
	assert.Equal(t, consts.Webhook, webhook.Config.Type)
	require.NotNil(t, webhook.Config.Url)
	assert.Equal(t, "https://example.com/hook", *webhook.Config.Url)
	body, ok := webhook.Config.Body.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "hello", body["text"])

	require.Len(t, w.Connections, 1)
	assert.Equal(t, "trigger", w.Connections[0].SourceIdentifier)
	assert.Equal(t, "notify", w.Connections[0].TargetIdentifier)
}

func TestWorkflowStateToPortBodyIntegrationAction(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("run"),
				IntegrationAction: &IntegrationActionModel{
					InstallationId:            types.StringValue("github-install"),
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
	require.NotNil(t, config.InstallationId)
	assert.Equal(t, "github-install", *config.InstallationId)
	require.NotNil(t, config.IntegrationProvider)
	assert.Equal(t, "github", *config.IntegrationProvider)
	props, ok := config.IntegrationActionExecutionProperties.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "deploy.yml", props["workflow"])
	assert.Equal(t, "main", props["ref"])
}

func TestWorkflowStateToPortBodyConditionAndConnections(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("branch"),
				Condition: &ConditionModel{
					Outlets: []ConditionOutletModel{
						{
							Identifier:  types.StringValue("approved"),
							Title:       types.StringValue("Approved"),
							Expression:  types.StringValue(".outputs.review.approved"),
							StatusLabel: &StatusLabelModel{Text: types.StringValue("Approved"), Variant: types.StringValue("success")},
						},
					},
				},
			},
		},
		Connections: []ConnectionModel{
			{
				SourceIdentifier:       types.StringValue("branch"),
				TargetIdentifier:       types.StringValue("deploy"),
				SourceOutletIdentifier: types.StringValue("approved"),
			},
			{
				SourceIdentifier: types.StringValue("branch"),
				TargetIdentifier: types.StringValue("reject"),
				Fallback:         types.BoolValue(true),
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	config := w.Nodes[0].Config
	assert.Equal(t, consts.ConditionNode, config.Type)
	require.Len(t, config.Outlets, 1)
	assert.Equal(t, "approved", config.Outlets[0].Identifier)
	require.NotNil(t, config.Outlets[0].Expression)
	assert.Equal(t, ".outputs.review.approved", *config.Outlets[0].Expression)
	require.NotNil(t, config.Outlets[0].StatusLabel)
	assert.Equal(t, "Approved", config.Outlets[0].StatusLabel.Text)

	require.Len(t, w.Connections, 2)
	require.NotNil(t, w.Connections[0].SourceOutletIdentifier)
	assert.Equal(t, "approved", *w.Connections[0].SourceOutletIdentifier)
	assert.Nil(t, w.Connections[0].Fallback)
	require.NotNil(t, w.Connections[1].Fallback)
	assert.True(t, *w.Connections[1].Fallback)
}

func TestWorkflowStateToPortBodyUpsertEntity(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("upsert"),
				UpsertEntity: &UpsertEntityModel{
					BlueprintIdentifier: types.StringValue("service"),
					Mapping: &UpsertMappingModel{
						Identifier: types.StringValue("{{ .outputs.trigger.identifier }}"),
						Title:      types.StringValue("My service"),
						Teams:      stringList("platform"),
						Properties: types.StringValue(`{"owner":"platform"}`),
					},
				},
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	config := w.Nodes[0].Config
	assert.Equal(t, consts.UpsertEntity, config.Type)
	require.NotNil(t, config.Mapping)
	assert.Equal(t, []string{"platform"}, config.Mapping.Team)
	assert.Equal(t, "platform", config.Mapping.Properties["owner"])
}

func TestRefreshWorkflowStatePreservesNodeOrder(t *testing.T) {
	ctx := context.Background()
	r := &WorkflowResource{portClient: &cli.PortClient{JSONEscapeHTML: true}}

	state := &WorkflowModel{
		Identifier: types.StringValue("review-pr"),
		Nodes: []WorkflowNodeModel{
			{Identifier: types.StringValue("trigger")},
			{Identifier: types.StringValue("notify")},
		},
		Connections: []ConnectionModel{
			{SourceIdentifier: types.StringValue("trigger"), TargetIdentifier: types.StringValue("notify")},
		},
	}

	apiWorkflow := &cli.Workflow{
		Identifier:            "review-pr",
		Category:              strPtr("engineering"),
		AllowAnyoneToViewRuns: boolPtr(true),
		Nodes: []cli.WorkflowNode{
			{
				Identifier: "notify",
				Config: cli.WorkflowNodeConfig{
					Type:      consts.Webhook,
					Url:       strPtr("https://example.com/hook"),
					Method:    strPtr("POST"),
					OnFailure: strPtr("terminate"),
				},
			},
			{
				Identifier: "trigger",
				Verbose:    boolPtr(true),
				Config: cli.WorkflowNodeConfig{
					Type:      consts.EventTrigger,
					Published: boolPtr(true),
					Event: &cli.WorkflowTriggerEvent{
						Type:                consts.EntityUpdated,
						BlueprintIdentifier: "githubPullRequest",
					},
					Condition: &cli.WorkflowNodeCondition{
						Type:        consts.JqCondition,
						Expressions: []string{".diff.after.properties.open == true"},
					},
				},
			},
		},
		Connections: []cli.WorkflowConnection{
			{SourceIdentifier: "trigger", TargetIdentifier: "notify"},
		},
	}

	err := r.refreshWorkflowState(ctx, state, apiWorkflow)
	require.NoError(t, err)

	require.Len(t, state.Nodes, 2)
	assert.Equal(t, "trigger", state.Nodes[0].Identifier.ValueString())
	assert.Equal(t, "notify", state.Nodes[1].Identifier.ValueString())

	require.NotNil(t, state.Nodes[0].EventTrigger)
	assert.Equal(t, consts.EntityUpdated, state.Nodes[0].EventTrigger.Type.ValueString())
	assert.Equal(t, "githubPullRequest", state.Nodes[0].EventTrigger.BlueprintIdentifier.ValueString())
	require.NotNil(t, state.Nodes[0].EventTrigger.Condition)
	assert.True(t, state.Nodes[0].Verbose.ValueBool())

	require.NotNil(t, state.Nodes[1].Webhook)
	assert.Equal(t, "https://example.com/hook", state.Nodes[1].Webhook.Url.ValueString())

	assert.Equal(t, "engineering", state.Category.ValueString())
	assert.True(t, state.AllowAnyoneToViewRuns.ValueBool())
}

func TestRefreshWorkflowStateKeepsUnsetCollectionsNull(t *testing.T) {
	ctx := context.Background()
	r := &WorkflowResource{portClient: &cli.PortClient{JSONEscapeHTML: true}}

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes:      []WorkflowNodeModel{{Identifier: types.StringValue("trigger")}},
	}

	apiWorkflow := &cli.Workflow{
		Identifier: "wf",
		Nodes: []cli.WorkflowNode{
			{
				Identifier: "trigger",
				Links:      []string{},
				Variables:  map[string]string{},
				Config: cli.WorkflowNodeConfig{
					Type:  consts.EventTrigger,
					Event: &cli.WorkflowTriggerEvent{Type: consts.EntityCreated, BlueprintIdentifier: "service"},
				},
			},
		},
	}

	err := r.refreshWorkflowState(ctx, state, apiWorkflow)
	require.NoError(t, err)

	assert.True(t, state.Nodes[0].Links.IsNull())
	assert.True(t, state.Nodes[0].Variables.IsNull())
}

func TestRefreshWorkflowStateSelfServeTriggerUserInputs(t *testing.T) {
	ctx := context.Background()
	r := &WorkflowResource{portClient: &cli.PortClient{JSONEscapeHTML: true}}

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes:      []WorkflowNodeModel{{Identifier: types.StringValue("trigger")}},
	}

	apiWorkflow := &cli.Workflow{
		Identifier: "wf",
		Nodes: []cli.WorkflowNode{
			{
				Identifier: "trigger",
				Config: cli.WorkflowNodeConfig{
					Type:      consts.SelfServeTrigger,
					Published: boolPtr(true),
					UserInputs: &cli.WorkflowUserInputs{
						Properties: map[string]cli.ActionProperty{
							"service": {Type: "string", Title: strPtr("Service")},
						},
						Required: []any{"service"},
					},
					Permissions: &cli.WorkflowNodePermissions{Roles: []string{"Member"}},
				},
			},
		},
	}

	err := r.refreshWorkflowState(ctx, state, apiWorkflow)
	require.NoError(t, err)

	trigger := state.Nodes[0].SelfServeTrigger
	require.NotNil(t, trigger)
	require.NotNil(t, trigger.UserInputs)
	require.NotNil(t, trigger.UserInputs.UserProperties)
	require.Contains(t, trigger.UserInputs.UserProperties.StringProps, "service")
	assert.Equal(t, "Service", trigger.UserInputs.UserProperties.StringProps["service"].Title.ValueString())

	require.NotNil(t, trigger.Permissions)
	roles, err := utils.TerraformListToGoArray(ctx, trigger.Permissions.Roles, "string")
	require.NoError(t, err)
	assert.Equal(t, []any{"Member"}, roles)
}

func TestSelfServeTriggerUserInputsRoundTrip(t *testing.T) {
	ctx := context.Background()
	r := &WorkflowResource{portClient: &cli.PortClient{JSONEscapeHTML: true}}

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes:      []WorkflowNodeModel{{Identifier: types.StringValue("trigger")}},
	}

	apiWorkflow := &cli.Workflow{
		Identifier: "wf",
		Nodes: []cli.WorkflowNode{
			{
				Identifier: "trigger",
				Config: cli.WorkflowNodeConfig{
					Type: consts.SelfServeTrigger,
					UserInputs: &cli.WorkflowUserInputs{
						Properties: map[string]cli.ActionProperty{
							"service":  {Type: "string", Title: strPtr("Service")},
							"replicas": {Type: "number", Title: strPtr("Replicas")},
						},
						Required: []any{"service"},
					},
				},
			},
		},
	}

	require.NoError(t, r.refreshWorkflowState(ctx, state, apiWorkflow))

	body, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	userInputs := body.Nodes[0].Config.UserInputs
	require.NotNil(t, userInputs)
	require.Contains(t, userInputs.Properties, "service")
	require.Contains(t, userInputs.Properties, "replicas")
	assert.Equal(t, "string", userInputs.Properties["service"].Type)
	assert.Equal(t, "number", userInputs.Properties["replicas"].Type)
	assert.Equal(t, []string{"service"}, userInputs.Required)
}

func validateNodes(nodes []WorkflowNodeModel, connections []ConnectionModel) diag.Diagnostics {
	resp := &resource.ValidateConfigResponse{}

	nodeTypes := make(map[string]string, len(nodes))
	outlets := make(map[string]map[string]bool, len(nodes))
	for i, node := range nodes {
		nodePath := path.Root("node").AtListIndex(i)
		identifier := node.Identifier.ValueString()
		if _, duplicate := nodeTypes[identifier]; duplicate {
			resp.Diagnostics.AddAttributeError(nodePath.AtName("identifier"), "Duplicate node identifier", identifier)
		}
		nodeTypes[identifier] = validateNode(resp, nodePath, node)
		outlets[identifier] = nodeOutlets(node)
	}
	validateConnections(resp, connections, nodeTypes, outlets)

	return resp.Diagnostics
}

func errorSummaries(diags diag.Diagnostics) []string {
	summaries := make([]string, 0, len(diags.Errors()))
	for _, d := range diags.Errors() {
		summaries = append(summaries, d.Summary())
	}
	return summaries
}

func eventTriggerNode(identifier string) WorkflowNodeModel {
	return WorkflowNodeModel{
		Identifier: types.StringValue(identifier),
		EventTrigger: &EventTriggerModel{
			Type:                types.StringValue(consts.EntityCreated),
			BlueprintIdentifier: types.StringValue("service"),
		},
	}
}

func conditionNode(identifier string, outletIdentifiers ...string) WorkflowNodeModel {
	outlets := make([]ConditionOutletModel, 0, len(outletIdentifiers))
	for _, o := range outletIdentifiers {
		outlets = append(outlets, ConditionOutletModel{
			Identifier: types.StringValue(o),
			Expression: types.StringValue(".outputs.ok"),
		})
	}
	return WorkflowNodeModel{
		Identifier: types.StringValue(identifier),
		Condition:  &ConditionModel{Outlets: outlets},
	}
}

func TestValidateNodeRequiredAttributes(t *testing.T) {
	tests := map[string]struct {
		node    WorkflowNodeModel
		summary string
	}{
		"event_trigger without blueprint": {
			node: WorkflowNodeModel{
				Identifier:   types.StringValue("trigger"),
				EventTrigger: &EventTriggerModel{Type: types.StringValue(consts.EntityCreated)},
			},
			summary: "Missing required attribute",
		},
		"timer_expired without property": {
			node: WorkflowNodeModel{
				Identifier: types.StringValue("trigger"),
				EventTrigger: &EventTriggerModel{
					Type:                types.StringValue(consts.WorkflowTimerExpired),
					BlueprintIdentifier: types.StringValue("service"),
				},
			},
			summary: "Missing required attribute",
		},
		"schedule_trigger without cron": {
			node: WorkflowNodeModel{
				Identifier:      types.StringValue("trigger"),
				ScheduleTrigger: &ScheduleTriggerModel{},
			},
			summary: "Missing required attribute",
		},
		"webhook without url": {
			node:    WorkflowNodeModel{Identifier: types.StringValue("call"), Webhook: &WebhookModel{}},
			summary: "Missing required attribute",
		},
		"upsert_entity without blueprint": {
			node:    WorkflowNodeModel{Identifier: types.StringValue("upsert"), UpsertEntity: &UpsertEntityModel{}},
			summary: "Missing required attribute",
		},
		"ai without user_prompt": {
			node:    WorkflowNodeModel{Identifier: types.StringValue("summarize"), AI: &AIModel{}},
			summary: "Missing required attribute",
		},
		"ai_agent without agent_identifier": {
			node: WorkflowNodeModel{
				Identifier: types.StringValue("agent"),
				AIAgent:    &AIAgentModel{UserPrompt: types.StringValue("go")},
			},
			summary: "Missing required attribute",
		},
		"integration_action without installation_id": {
			node: WorkflowNodeModel{
				Identifier: types.StringValue("deploy"),
				IntegrationAction: &IntegrationActionModel{
					IntegrationProvider:       types.StringValue(consts.Github),
					IntegrationInvocationType: types.StringValue("GITHUB_WORKFLOW"),
				},
			},
			summary: "Missing required attribute",
		},
		"condition without outlets": {
			node:    WorkflowNodeModel{Identifier: types.StringValue("branch"), Condition: &ConditionModel{}},
			summary: "Missing required block",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			diags := validateNodes([]WorkflowNodeModel{test.node}, nil)
			assert.Contains(t, errorSummaries(diags), test.summary)
		})
	}
}

func TestValidateNodeAcceptsCompleteConfigs(t *testing.T) {
	nodes := []WorkflowNodeModel{
		eventTriggerNode("trigger"),
		{
			Identifier:      types.StringValue("nightly"),
			ScheduleTrigger: &ScheduleTriggerModel{Cron: types.StringValue("0 9 * * 1-5")},
		},
		{
			Identifier: types.StringValue("notify"),
			Webhook:    &WebhookModel{Url: types.StringValue("https://example.com/hook")},
		},
		{
			Identifier:   types.StringValue("upsert"),
			UpsertEntity: &UpsertEntityModel{BlueprintIdentifier: types.StringValue("service")},
		},
		conditionNode("branch", "approved"),
	}

	diags := validateNodes(nodes, []ConnectionModel{
		{
			SourceIdentifier:       types.StringValue("branch"),
			TargetIdentifier:       types.StringValue("notify"),
			SourceOutletIdentifier: types.StringValue("approved"),
		},
		{
			SourceIdentifier: types.StringValue("trigger"),
			TargetIdentifier: types.StringValue("branch"),
		},
	})

	assert.False(t, diags.HasError(), errorSummaries(diags))
}

func TestValidateSelfServeTriggerContexts(t *testing.T) {
	tests := map[string]struct {
		context TriggerContextModel
		summary string
	}{
		"create_entity without blueprint": {
			context: TriggerContextModel{On: types.StringValue(consts.CreateEntityContext)},
			summary: "Missing required attribute",
		},
		"create_entity with user_input": {
			context: TriggerContextModel{
				On:                  types.StringValue(consts.CreateEntityContext),
				BlueprintIdentifier: types.StringValue("service"),
				UserInput:           types.StringValue("service"),
			},
			summary: "Invalid attribute combination",
		},
		"entity without user_input": {
			context: TriggerContextModel{On: types.StringValue(consts.EntityContext)},
			summary: "Missing required attribute",
		},
		"entity with blueprint": {
			context: TriggerContextModel{
				On:                  types.StringValue(consts.EntityContext),
				UserInput:           types.StringValue("service"),
				BlueprintIdentifier: types.StringValue("service"),
			},
			summary: "Invalid attribute combination",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			node := WorkflowNodeModel{
				Identifier: types.StringValue("trigger"),
				SelfServeTrigger: &SelfServeTriggerModel{
					Contexts: []TriggerContextModel{test.context},
				},
			}
			diags := validateNodes([]WorkflowNodeModel{node}, nil)
			assert.Contains(t, errorSummaries(diags), test.summary)
		})
	}
}

func TestValidateConnections(t *testing.T) {
	nodes := []WorkflowNodeModel{
		eventTriggerNode("trigger"),
		conditionNode("branch", "approved"),
		{
			Identifier: types.StringValue("notify"),
			Webhook:    &WebhookModel{Url: types.StringValue("https://example.com/hook")},
		},
	}

	tests := map[string]struct {
		connection ConnectionModel
		summary    string
	}{
		"self connection": {
			connection: ConnectionModel{
				SourceIdentifier: types.StringValue("notify"),
				TargetIdentifier: types.StringValue("notify"),
			},
			summary: "Invalid connection",
		},
		"outlet on single outlet node": {
			connection: ConnectionModel{
				SourceIdentifier:       types.StringValue("trigger"),
				TargetIdentifier:       types.StringValue("notify"),
				SourceOutletIdentifier: types.StringValue("approved"),
			},
			summary: "Invalid connection",
		},
		"fallback on single outlet node": {
			connection: ConnectionModel{
				SourceIdentifier: types.StringValue("trigger"),
				TargetIdentifier: types.StringValue("notify"),
				Fallback:         types.BoolValue(true),
			},
			summary: "Invalid connection",
		},
		"condition without outlet or fallback": {
			connection: ConnectionModel{
				SourceIdentifier: types.StringValue("branch"),
				TargetIdentifier: types.StringValue("notify"),
			},
			summary: "Missing outlet",
		},
		"unknown outlet": {
			connection: ConnectionModel{
				SourceIdentifier:       types.StringValue("branch"),
				TargetIdentifier:       types.StringValue("notify"),
				SourceOutletIdentifier: types.StringValue("rejected"),
			},
			summary: "Unknown outlet",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			diags := validateNodes(nodes, []ConnectionModel{test.connection})
			assert.Contains(t, errorSummaries(diags), test.summary)
		})
	}
}

func TestValidateDuplicateNodeIdentifiers(t *testing.T) {
	diags := validateNodes([]WorkflowNodeModel{eventTriggerNode("trigger"), eventTriggerNode("trigger")}, nil)
	assert.Contains(t, errorSummaries(diags), "Duplicate node identifier")
}

func TestValidateInputNodeOutletsMatchButtons(t *testing.T) {
	node := WorkflowNodeModel{
		Identifier: types.StringValue("approval"),
		Input: &InputModel{
			UserInputs: &InputUserInputsModel{
				Buttons: []InputButtonModel{{Identifier: types.StringValue("approve"), Label: types.StringValue("Approve")}},
			},
			Outlets: []InputOutletModel{{Identifier: types.StringValue("reject")}},
		},
	}

	diags := validateNodes([]WorkflowNodeModel{node}, nil)
	assert.Contains(t, errorSummaries(diags), "Unknown button identifier")
}

func TestValidateInputNodeFallbackRejected(t *testing.T) {
	node := WorkflowNodeModel{
		Identifier: types.StringValue("approval"),
		Input: &InputModel{
			UserInputs: &InputUserInputsModel{
				Buttons: []InputButtonModel{{Identifier: types.StringValue("approve"), Label: types.StringValue("Approve")}},
			},
			Outlets: []InputOutletModel{{Identifier: types.StringValue("approve")}},
		},
	}

	diags := validateNodes([]WorkflowNodeModel{node}, []ConnectionModel{
		{
			SourceIdentifier: types.StringValue("approval"),
			TargetIdentifier: types.StringValue("deploy"),
			Fallback:         types.BoolValue(true),
		},
	})

	assert.Contains(t, errorSummaries(diags), "Invalid fallback")
}

func TestInputNodeRespondersUseUsersQuery(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("approval"),
				Input: &InputModel{
					UserInputs: &InputUserInputsModel{
						Buttons: []InputButtonModel{{Identifier: types.StringValue("approve"), Label: types.StringValue("Approve")}},
					},
					Outlets: []InputOutletModel{{Identifier: types.StringValue("approve"), NumOfResponders: types.Int64Value(1)}},
					Responders: &RespondersModel{
						Roles:      stringList("Admin"),
						UsersQuery: types.StringValue(`{"combinator":"and","rules":[{"property":"$title","operator":"=","value":"SRE"}]}`),
					},
				},
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	responders := w.Nodes[0].Config.Responders
	require.NotNil(t, responders)
	assert.Equal(t, []string{"Admin"}, responders.Roles)
	assert.NotNil(t, responders.UsersQuery, "responders must send usersQuery")
	assert.Nil(t, responders.Policy, "responders must not send policy")

	body, err := json.Marshal(responders)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"usersQuery"`)
	assert.NotContains(t, string(body), `"policy"`)
}

func TestSelfServeTriggerPermissionsUsePolicy(t *testing.T) {
	ctx := context.Background()

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("trigger"),
				SelfServeTrigger: &SelfServeTriggerModel{
					Permissions: &PermissionsModel{
						Roles:  stringList("Member"),
						Policy: types.StringValue(`{"combinator":"and","rules":[]}`),
					},
				},
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	permissions := w.Nodes[0].Config.Permissions
	require.NotNil(t, permissions)
	assert.NotNil(t, permissions.Policy)
	assert.Nil(t, permissions.UsersQuery, "permissions must not send usersQuery")

	body, err := json.Marshal(permissions)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"policy"`)
	assert.NotContains(t, string(body), `"usersQuery"`)
}

func TestRespondersRoundTrip(t *testing.T) {
	ctx := context.Background()

	responders := &cli.WorkflowNodePermissions{
		Teams:      []string{"platform"},
		UsersQuery: map[string]any{"combinator": "and", "rules": []any{}},
	}

	model := respondersToModel(ctx, responders, false)
	require.NotNil(t, model)
	assert.False(t, model.UsersQuery.IsNull())

	body, err := respondersToPortBody(ctx, model)
	require.NoError(t, err)
	assert.Equal(t, []string{"platform"}, body.Teams)
	assert.NotNil(t, body.UsersQuery)
	assert.Nil(t, body.Policy)
}

func webhookNode(identifier string) WorkflowNodeModel {
	return WorkflowNodeModel{
		Identifier: types.StringValue(identifier),
		Webhook:    &WebhookModel{Url: types.StringValue("https://example.com")},
	}
}

func TestValidateSingleOutgoingConnectionPerNonBranchingNode(t *testing.T) {
	nodes := []WorkflowNodeModel{
		eventTriggerNode("trigger"),
		webhookNode("a"),
		webhookNode("b"),
	}
	connections := []ConnectionModel{
		{SourceIdentifier: types.StringValue("trigger"), TargetIdentifier: types.StringValue("a")},
		{SourceIdentifier: types.StringValue("trigger"), TargetIdentifier: types.StringValue("b")},
	}

	diags := validateNodes(nodes, connections)
	assert.Contains(t, errorSummaries(diags), "Duplicate connection source")
}

func TestValidateConditionNodeMayBranchToDistinctOutlets(t *testing.T) {
	nodes := []WorkflowNodeModel{
		eventTriggerNode("trigger"),
		conditionNode("route", "a", "b"),
		webhookNode("a"),
		webhookNode("b"),
		webhookNode("fallback"),
	}
	connections := []ConnectionModel{
		{SourceIdentifier: types.StringValue("trigger"), TargetIdentifier: types.StringValue("route")},
		{SourceIdentifier: types.StringValue("route"), TargetIdentifier: types.StringValue("a"), SourceOutletIdentifier: types.StringValue("a")},
		{SourceIdentifier: types.StringValue("route"), TargetIdentifier: types.StringValue("b"), SourceOutletIdentifier: types.StringValue("b")},
		{SourceIdentifier: types.StringValue("route"), TargetIdentifier: types.StringValue("fallback"), Fallback: types.BoolValue(true)},
	}

	assert.Empty(t, errorSummaries(validateNodes(nodes, connections)))
}

func TestValidateConditionNodeRejectsRepeatedOutlet(t *testing.T) {
	nodes := []WorkflowNodeModel{
		eventTriggerNode("trigger"),
		conditionNode("route", "a"),
		webhookNode("a"),
		webhookNode("b"),
	}
	connections := []ConnectionModel{
		{SourceIdentifier: types.StringValue("trigger"), TargetIdentifier: types.StringValue("route")},
		{SourceIdentifier: types.StringValue("route"), TargetIdentifier: types.StringValue("a"), SourceOutletIdentifier: types.StringValue("a")},
		{SourceIdentifier: types.StringValue("route"), TargetIdentifier: types.StringValue("b"), SourceOutletIdentifier: types.StringValue("a")},
	}

	assert.Contains(t, errorSummaries(validateNodes(nodes, connections)), "Duplicate connection source")
}

func TestUpsertMappingTeamIsStringArray(t *testing.T) {
	ctx := context.Background()

	mapping := &UpsertMappingModel{
		Identifier: types.StringValue("svc"),
		Teams:      stringList("{{ .outputs.trigger.diff.after.team }}"),
	}

	body, err := upsertMappingToPortBody(ctx, mapping)
	require.NoError(t, err)
	assert.Equal(t, []string{"{{ .outputs.trigger.diff.after.team }}"}, body.Team)

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "jqQuery")
}

// Guards against a node type being added to the model without being wired into
// the mapper: every block in nodeTypeBlockNames must serialize to its config type.
func TestEveryNodeTypeSerializesItsConfigType(t *testing.T) {
	ctx := context.Background()

	nodes := map[string]struct {
		block      string
		node       WorkflowNodeModel
		configType string
	}{
		"self_serve_trigger": {"self_serve_trigger", WorkflowNodeModel{
			Identifier:       types.StringValue("n"),
			SelfServeTrigger: &SelfServeTriggerModel{ActionCardButtonText: types.StringValue("Run")},
		}, consts.SelfServeTrigger},
		"event_trigger": {"event_trigger", eventTriggerNode("n"), consts.EventTrigger},
		"schedule_trigger": {"schedule_trigger", WorkflowNodeModel{
			Identifier:      types.StringValue("n"),
			ScheduleTrigger: &ScheduleTriggerModel{Cron: types.StringValue("0 9 * * *")},
		}, consts.ScheduleTrigger},
		"kafka": {"kafka", WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			Kafka:      &KafkaModel{Payload: types.StringValue(`{"event":"tick"}`)},
		}, consts.Kafka},
		"webhook": {"webhook", webhookNode("n"), consts.Webhook},
		"integration_action": {"integration_action", WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			IntegrationAction: &IntegrationActionModel{
				InstallationId:            types.StringValue("gh"),
				IntegrationProvider:       types.StringValue(consts.Github),
				IntegrationInvocationType: types.StringValue("GITHUB_WORKFLOW"),
			},
		}, consts.IntegrationAction},
		"upsert_entity": {"upsert_entity", WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			UpsertEntity: &UpsertEntityModel{
				BlueprintIdentifier: types.StringValue("service"),
				Mapping:             &UpsertMappingModel{Identifier: types.StringValue("svc")},
			},
		}, consts.UpsertEntity},
		"ai": {"ai", WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AI:         &AIModel{UserPrompt: types.StringValue("summarize")},
		}, consts.AI},
		"ai_agent": {"ai_agent", WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AIAgent: &AIAgentModel{
				UserPrompt:      types.StringValue("go"),
				AgentIdentifier: types.StringValue("agent"),
			},
		}, consts.AIAgent},
		"condition": {"condition", conditionNode("n", "ok"), consts.ConditionNode},
		"input": {"input", WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			Input: &InputModel{
				UserInputs: &InputUserInputsModel{
					Buttons: []InputButtonModel{{Identifier: types.StringValue("ok"), Label: types.StringValue("OK")}},
				},
				Outlets: []InputOutletModel{{Identifier: types.StringValue("ok")}},
			},
		}, consts.InputNode},
	}

	covered := make([]string, 0, len(nodes))
	for name, test := range nodes {
		t.Run(name, func(t *testing.T) {
			w, err := workflowStateToPortBody(ctx, &WorkflowModel{
				Identifier: types.StringValue("wf"),
				Nodes:      []WorkflowNodeModel{test.node},
			})
			require.NoError(t, err)
			require.Len(t, w.Nodes, 1)
			assert.Equal(t, test.configType, w.Nodes[0].Config.Type)
		})
		covered = append(covered, test.block)
	}

	assert.ElementsMatch(t, nodeTypeBlockNames, covered,
		"every node type block must have a serialization test")
}

func TestEveryEventTriggerTypeSerializes(t *testing.T) {
	ctx := context.Background()

	eventTypes := []string{
		consts.EntityCreated,
		consts.EntityUpdated,
		consts.EntityDeleted,
		consts.WorkflowTimerExpired,
		consts.AnyEntityChange,
	}

	for _, eventType := range eventTypes {
		t.Run(eventType, func(t *testing.T) {
			trigger := &EventTriggerModel{
				Type:                types.StringValue(eventType),
				BlueprintIdentifier: types.StringValue("service"),
			}
			if eventType == consts.WorkflowTimerExpired {
				trigger.PropertyIdentifier = types.StringValue("expires_at")
			}

			node := WorkflowNodeModel{Identifier: types.StringValue("trigger"), EventTrigger: trigger}

			w, err := workflowStateToPortBody(ctx, &WorkflowModel{
				Identifier: types.StringValue("wf"),
				Nodes:      []WorkflowNodeModel{node},
			})
			require.NoError(t, err)

			event := w.Nodes[0].Config.Event
			require.NotNil(t, event)
			assert.Equal(t, eventType, event.Type)
			assert.Equal(t, "service", event.BlueprintIdentifier)

			assert.Empty(t, errorSummaries(validateNodes([]WorkflowNodeModel{node}, nil)))
		})
	}
}

func TestEverySelfServeTriggerContextSerializes(t *testing.T) {
	ctx := context.Background()

	contexts := map[string]TriggerContextModel{
		consts.CreateEntityContext: {
			On:                  types.StringValue(consts.CreateEntityContext),
			BlueprintIdentifier: types.StringValue("service"),
		},
		// The ENTITY context is strict: it takes a user input, not a blueprint.
		consts.EntityContext: {
			On:        types.StringValue(consts.EntityContext),
			UserInput: types.StringValue("service"),
		},
	}

	for on, triggerContext := range contexts {
		t.Run(on, func(t *testing.T) {
			node := WorkflowNodeModel{
				Identifier:       types.StringValue("trigger"),
				SelfServeTrigger: &SelfServeTriggerModel{Contexts: []TriggerContextModel{triggerContext}},
			}

			w, err := workflowStateToPortBody(ctx, &WorkflowModel{
				Identifier: types.StringValue("wf"),
				Nodes:      []WorkflowNodeModel{node},
			})
			require.NoError(t, err)

			require.Len(t, w.Nodes[0].Config.Contexts, 1)
			assert.Equal(t, on, w.Nodes[0].Config.Contexts[0].On)

			assert.Empty(t, errorSummaries(validateNodes([]WorkflowNodeModel{node}, nil)))
		})
	}
}

func TestValidateAIProviderAndModelMustBePaired(t *testing.T) {
	tests := map[string]struct {
		node      WorkflowNodeModel
		wantError bool
	}{
		"ai with provider only": {WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AI:         &AIModel{UserPrompt: types.StringValue("go"), Provider: types.StringValue("openai")},
		}, true},
		"ai with model only": {WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AI:         &AIModel{UserPrompt: types.StringValue("go"), Model: types.StringValue("gpt-4")},
		}, true},
		"ai with both": {WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AI: &AIModel{
				UserPrompt: types.StringValue("go"),
				Provider:   types.StringValue("openai"),
				Model:      types.StringValue("gpt-4"),
			},
		}, false},
		"ai with neither": {WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AI:         &AIModel{UserPrompt: types.StringValue("go")},
		}, false},
		"ai_agent with provider only": {WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AIAgent: &AIAgentModel{
				UserPrompt:      types.StringValue("go"),
				AgentIdentifier: types.StringValue("a"),
				Provider:        types.StringValue("openai"),
			},
		}, true},
		"ai_agent with both": {WorkflowNodeModel{
			Identifier: types.StringValue("n"),
			AIAgent: &AIAgentModel{
				UserPrompt:      types.StringValue("go"),
				AgentIdentifier: types.StringValue("a"),
				Provider:        types.StringValue("openai"),
				Model:           types.StringValue("gpt-4"),
			},
		}, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			summaries := errorSummaries(validateNodes([]WorkflowNodeModel{test.node}, nil))
			if test.wantError {
				assert.Contains(t, summaries, "Missing required attribute")
			} else {
				assert.Empty(t, summaries)
			}
		})
	}
}

// The event trigger is the only node type the service accepts a `condition` on.
func TestEventTriggerConditionRoundTrip(t *testing.T) {
	ctx := context.Background()
	r := &WorkflowResource{portClient: &cli.PortClient{JSONEscapeHTML: true}}

	state := &WorkflowModel{
		Identifier: types.StringValue("wf"),
		Nodes: []WorkflowNodeModel{
			{
				Identifier: types.StringValue("trigger"),
				EventTrigger: &EventTriggerModel{
					Type:                types.StringValue(consts.EntityUpdated),
					BlueprintIdentifier: types.StringValue("service"),
					Condition: &NodeConditionModel{
						Expressions: stringList(".diff.after.properties.tier == \"production\"", ".diff.before != null"),
						Combinator:  types.StringValue("or"),
					},
				},
			},
		},
	}

	w, err := workflowStateToPortBody(ctx, state)
	require.NoError(t, err)

	condition := w.Nodes[0].Config.Condition
	require.NotNil(t, condition)
	assert.Equal(t, consts.JqCondition, condition.Type)
	assert.Len(t, condition.Expressions, 2)
	require.NotNil(t, condition.Combinator)
	assert.Equal(t, "or", *condition.Combinator)

	refreshed := &WorkflowModel{Identifier: types.StringValue("wf")}
	require.NoError(t, r.refreshWorkflowState(ctx, refreshed, w))

	got := refreshed.Nodes[0].EventTrigger.Condition
	require.NotNil(t, got, "the event trigger condition must survive a refresh")
	assert.Equal(t, "or", got.Combinator.ValueString())
	assert.False(t, got.Expressions.IsNull())
}

// Walks the real schema so that exposing a condition on a second node type fails.
func TestConditionIsExposedOnlyOnEventTrigger(t *testing.T) {
	blocksWithCondition := []string{}
	for _, name := range nodeTypeBlockNames {
		block, ok := WorkflowBlocks()["node"].(schema.ListNestedBlock).
			NestedObject.Blocks[name].(schema.SingleNestedBlock)
		if !ok {
			continue
		}
		if _, has := block.Blocks["condition"]; has {
			blocksWithCondition = append(blocksWithCondition, name)
		}
	}

	assert.Equal(t, []string{"event_trigger"}, blocksWithCondition)
}

func validateWorkflow(nodes []WorkflowNodeModel, connections []ConnectionModel) diag.Diagnostics {
	resp := &resource.ValidateConfigResponse{}

	nodeTypes := make(map[string]string, len(nodes))
	outlets := make(map[string]map[string]bool, len(nodes))
	for i, node := range nodes {
		nodePath := path.Root("node").AtListIndex(i)
		nodeTypes[node.Identifier.ValueString()] = validateNode(resp, nodePath, node)
		outlets[node.Identifier.ValueString()] = nodeOutlets(node)
	}

	validateTriggerPresence(resp, nodeTypes)
	validateConnections(resp, connections, nodeTypes, outlets)

	return resp.Diagnostics
}

func connection(source, target string) ConnectionModel {
	return ConnectionModel{
		SourceIdentifier: types.StringValue(source),
		TargetIdentifier: types.StringValue(target),
	}
}

func TestValidateRejectsConnectionCycles(t *testing.T) {
	nodes := []WorkflowNodeModel{
		eventTriggerNode("trigger"),
		webhookNode("a"),
		webhookNode("b"),
		webhookNode("c"),
	}

	// trigger -> a -> b -> c -> a is a cycle the API rejects, but every node
	// still has only one outgoing connection.
	diags := validateWorkflow(nodes, []ConnectionModel{
		connection("trigger", "a"),
		connection("a", "b"),
		connection("b", "c"),
		connection("c", "a"),
	})

	assert.Contains(t, errorSummaries(diags), "Cyclic connections")
}

func TestValidateAcceptsBranchingAcyclicGraph(t *testing.T) {
	nodes := []WorkflowNodeModel{
		eventTriggerNode("trigger"),
		conditionNode("route", "left", "right"),
		webhookNode("a"),
		webhookNode("b"),
		webhookNode("join"),
	}

	// A diamond re-joins on `join`, which is acyclic despite the shared target.
	diags := validateWorkflow(nodes, []ConnectionModel{
		connection("trigger", "route"),
		{
			SourceIdentifier:       types.StringValue("route"),
			TargetIdentifier:       types.StringValue("a"),
			SourceOutletIdentifier: types.StringValue("left"),
		},
		{
			SourceIdentifier:       types.StringValue("route"),
			TargetIdentifier:       types.StringValue("b"),
			SourceOutletIdentifier: types.StringValue("right"),
		},
		connection("a", "join"),
		connection("b", "join"),
	})

	assert.Empty(t, errorSummaries(diags))
}

func TestValidateConnectionEndpointsMustBeDeclaredNodes(t *testing.T) {
	nodes := []WorkflowNodeModel{eventTriggerNode("trigger"), webhookNode("notify")}

	diags := validateWorkflow(nodes, []ConnectionModel{
		connection("trigger", "typo"),
		connection("ghost", "notify"),
	})

	summaries := errorSummaries(diags)
	assert.Equal(t, 2, countSummary(summaries, "Unknown node identifier"))
}

func TestValidateTriggerCannotHaveIncomingConnections(t *testing.T) {
	nodes := []WorkflowNodeModel{eventTriggerNode("trigger"), webhookNode("notify")}

	diags := validateWorkflow(nodes, []ConnectionModel{
		connection("trigger", "notify"),
		connection("notify", "trigger"),
	})

	assert.Contains(t, errorSummaries(diags), "Invalid connection")
}

func TestValidateWorkflowRequiresATriggerNode(t *testing.T) {
	diags := validateWorkflow([]WorkflowNodeModel{webhookNode("notify")}, nil)
	assert.Contains(t, errorSummaries(diags), "Missing trigger node")

	diags = validateWorkflow([]WorkflowNodeModel{eventTriggerNode("trigger")}, nil)
	assert.Empty(t, errorSummaries(diags))
}

func TestValidateRejectsDuplicateOutletIdentifiers(t *testing.T) {
	diags := validateNodes([]WorkflowNodeModel{conditionNode("route", "same", "same")}, nil)
	assert.Contains(t, errorSummaries(diags), "Duplicate outlet identifier")
}

func TestValidateRejectsDuplicateButtonIdentifiers(t *testing.T) {
	node := WorkflowNodeModel{
		Identifier: types.StringValue("approval"),
		Input: &InputModel{
			UserInputs: &InputUserInputsModel{
				Buttons: []InputButtonModel{
					{Identifier: types.StringValue("approve"), Label: types.StringValue("Approve")},
					{Identifier: types.StringValue("approve"), Label: types.StringValue("Approve again")},
				},
			},
			Outlets: []InputOutletModel{{Identifier: types.StringValue("approve")}},
		},
	}

	assert.Contains(t, errorSummaries(validateNodes([]WorkflowNodeModel{node}, nil)), "Duplicate button identifier")
}

func TestValidateEmailNotificationRequiresFields(t *testing.T) {
	inputNode := func(notification NotificationModel) WorkflowNodeModel {
		return WorkflowNodeModel{
			Identifier: types.StringValue("approval"),
			Input: &InputModel{
				UserInputs: &InputUserInputsModel{
					Buttons: []InputButtonModel{{Identifier: types.StringValue("approve"), Label: types.StringValue("Approve")}},
				},
				Outlets:       []InputOutletModel{{Identifier: types.StringValue("approve")}},
				Notifications: []NotificationModel{notification},
			},
		}
	}

	missing := inputNode(NotificationModel{Target: types.StringValue("email")})
	assert.Contains(t, errorSummaries(validateNodes([]WorkflowNodeModel{missing}, nil)), "Missing required block")

	complete := inputNode(NotificationModel{
		Target: types.StringValue("email"),
		Fields: []NotificationFieldModel{{Label: types.StringValue("Service"), Value: types.StringValue("api")}},
	})
	assert.Empty(t, errorSummaries(validateNodes([]WorkflowNodeModel{complete}, nil)))
}

func TestValidateRejectsBlankRequiredStrings(t *testing.T) {
	blankExpression := WorkflowNodeModel{
		Identifier: types.StringValue("route"),
		Condition: &ConditionModel{
			Outlets: []ConditionOutletModel{{
				Identifier: types.StringValue("outlet"),
				Expression: types.StringValue("   "),
			}},
		},
	}
	assert.Contains(t, errorSummaries(validateNodes([]WorkflowNodeModel{blankExpression}, nil)), "Missing required attribute")

	blankUserInput := WorkflowNodeModel{
		Identifier: types.StringValue("trigger"),
		SelfServeTrigger: &SelfServeTriggerModel{
			Contexts: []TriggerContextModel{{
				On:        types.StringValue(consts.EntityContext),
				UserInput: types.StringValue(""),
			}},
		},
	}
	assert.Contains(t, errorSummaries(validateNodes([]WorkflowNodeModel{blankUserInput}, nil)), "Missing required attribute")
}

func countSummary(summaries []string, summary string) int {
	count := 0
	for _, s := range summaries {
		if s == summary {
			count++
		}
	}
	return count
}

func runStringValidator(v validator.String, value string) diag.Diagnostics {
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), validator.StringRequest{
		Path:        path.Root("attribute"),
		ConfigValue: types.StringValue(value),
	}, resp)
	return resp.Diagnostics
}

func TestCronValidator(t *testing.T) {
	valid := []string{"0 9 * * 1-5", "*/5 * * * *", "0 0 1 1 * 2030", "@daily", "0 2 * * MON"}
	for _, expression := range valid {
		assert.Empty(t, errorSummaries(runStringValidator(cronValidator(), expression)), expression)
	}

	invalid := []string{"", "0 9 * *", "not a cron", "@sometimes", "0 9 * * 1-5 * *"}
	for _, expression := range invalid {
		assert.NotEmpty(t, errorSummaries(runStringValidator(cronValidator(), expression)), expression)
	}
}

func TestOutputSchemaValidator(t *testing.T) {
	assert.Empty(t, errorSummaries(runStringValidator(outputSchemaValidator(),
		`{"type":"object","properties":{"summary":{"type":"string"}}}`)))

	invalid := []string{`{"type":"string"}`, `[]`, `{"properties":{}}`, `not json`}
	for _, value := range invalid {
		assert.NotEmpty(t, errorSummaries(runStringValidator(outputSchemaValidator(), value)), value)
	}
}

func TestQueryValidator(t *testing.T) {
	assert.Empty(t, errorSummaries(runStringValidator(queryValidator("Invalid policy"),
		`{"combinator":"and","rules":[{"property":{"context":"user","property":"department"},"operator":"=","value":"engineering"}]}`)))

	invalid := []string{`{"combinator":"maybe","rules":[]}`, `{"rules":[]}`, `{"combinator":"and"}`, `"and"`, `{`}
	for _, value := range invalid {
		assert.NotEmpty(t, errorSummaries(runStringValidator(queryValidator("Invalid policy"), value)), value)
	}
}
