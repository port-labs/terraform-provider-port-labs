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
	Category    types.String        `tfsdk:"category"`
	Nodes       []WorkflowNodeModel `tfsdk:"node"`
	Connections []ConnectionModel   `tfsdk:"connections"`
}

type WorkflowNodeModel struct {
	Identifier        types.String            `tfsdk:"identifier"`
	Title             types.String            `tfsdk:"title"`
	EventTrigger      *EventTriggerModel      `tfsdk:"event_trigger"`
	CursorAgent       *CursorAgentModel       `tfsdk:"cursor_agent"`
	IntegrationAction *IntegrationActionModel `tfsdk:"integration_action"`
}

type EventTriggerModel struct {
	Type                types.String `tfsdk:"type"`
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
	PropertyIdentifier  types.String `tfsdk:"property_identifier"`
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
	Repository types.String `tfsdk:"repository"`
	Ref        types.String `tfsdk:"ref"`
	PrUrl      types.String `tfsdk:"pr_url"`
}

type IntegrationActionModel struct {
	InstallationId               types.String `tfsdk:"installation_id"`
	IntegrationProvider          types.String `tfsdk:"integration_provider"`
	IntegrationInvocationType    types.String `tfsdk:"integration_invocation_type"`
	ExecutionProperties          types.String `tfsdk:"execution_properties"`
	DeferIntegrationInstallation types.Bool   `tfsdk:"defer_integration_installation"`
	OnFailure                    types.String `tfsdk:"on_failure"`
}

type ConnectionModel struct {
	SourceIdentifier types.String `tfsdk:"source_identifier"`
	TargetIdentifier types.String `tfsdk:"target_identifier"`
}
