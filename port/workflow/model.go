package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WorkflowModel struct {
	ID          types.String        `tfsdk:"id"`
	Identifier  types.String        `tfsdk:"identifier"`
	Title       types.String        `tfsdk:"title"`
	Icon        types.String        `tfsdk:"icon"`
	Description types.String        `tfsdk:"description"`
	Tags        types.List          `tfsdk:"tags"`
	Nodes       []WorkflowNodeModel `tfsdk:"node"`
	Connections []ConnectionModel   `tfsdk:"connections"`
}

// WorkflowNodeModel models a single node of the workflow graph. Exactly one of
// the typed sub-blocks is set, enforced by an ExactlyOneOf validator.
type WorkflowNodeModel struct {
	Identifier        types.String            `tfsdk:"identifier"`
	Title             types.String            `tfsdk:"title"`
	EventTrigger      *EventTriggerModel      `tfsdk:"event_trigger"`
	CursorAgent       *CursorAgentModel       `tfsdk:"cursor_agent"`
	Delay             *DelayModel             `tfsdk:"delay"`
	IntegrationAction *IntegrationActionModel `tfsdk:"integration_action"`
}

type EventTriggerModel struct {
	Type                types.String `tfsdk:"type"`
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
}

type CursorAgentModel struct {
	ApiKey types.String       `tfsdk:"api_key"`
	Prompt *CursorPromptModel `tfsdk:"prompt"`
	Source *CursorSourceModel `tfsdk:"source"`
}

type CursorPromptModel struct {
	Text types.String `tfsdk:"text"`
}

type CursorSourceModel struct {
	PrUrl types.String `tfsdk:"pr_url"`
}

type DelayModel struct {
	// Seconds accepts a positive integer (as a string) or a dynamic Port
	// expression ("{{ ... }}").
	Seconds types.String `tfsdk:"seconds"`
}

type IntegrationActionModel struct {
	InstallationId        types.String `tfsdk:"installation_id"`
	IntegrationActionType types.String `tfsdk:"integration_action_type"`
	Org                   types.String `tfsdk:"org"`
	Repo                  types.String `tfsdk:"repo"`
	Workflow              types.String `tfsdk:"workflow"`
	WorkflowInputs        types.String `tfsdk:"workflow_inputs"`
	ReportWorkflowStatus  types.String `tfsdk:"report_workflow_status"`
}

type ConnectionModel struct {
	SourceIdentifier types.String `tfsdk:"source_identifier"`
	TargetIdentifier types.String `tfsdk:"target_identifier"`
}
