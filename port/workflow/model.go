package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/action"
)

type WorkflowModel struct {
	ID                    types.String        `tfsdk:"id"`
	Identifier            types.String        `tfsdk:"identifier"`
	Title                 types.String        `tfsdk:"title"`
	Icon                  types.String        `tfsdk:"icon"`
	Description           types.String        `tfsdk:"description"`
	Category              types.String        `tfsdk:"category"`
	AllowAnyoneToViewRuns types.Bool          `tfsdk:"allow_anyone_to_view_runs"`
	Nodes                 []WorkflowNodeModel `tfsdk:"node"`
	Connections           []ConnectionModel   `tfsdk:"connections"`
}

type WorkflowNodeModel struct {
	Identifier  types.String `tfsdk:"identifier"`
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Verbose     types.Bool   `tfsdk:"verbose"`
	Links       types.List   `tfsdk:"links"`
	Variables   types.Map    `tfsdk:"variables"`

	SelfServeTrigger *SelfServeTriggerModel `tfsdk:"self_serve_trigger"`
	EventTrigger     *EventTriggerModel     `tfsdk:"event_trigger"`
	ScheduleTrigger  *ScheduleTriggerModel  `tfsdk:"schedule_trigger"`

	Kafka             *KafkaModel             `tfsdk:"kafka"`
	Webhook           *WebhookModel           `tfsdk:"webhook"`
	IntegrationAction *IntegrationActionModel `tfsdk:"integration_action"`
	UpsertEntity      *UpsertEntityModel      `tfsdk:"upsert_entity"`
	AI                *AIModel                `tfsdk:"ai"`
	AIAgent           *AIAgentModel           `tfsdk:"ai_agent"`
	Condition         *ConditionModel         `tfsdk:"condition"`
	Input             *InputModel             `tfsdk:"input"`
}

type SelfServeTriggerModel struct {
	UserInputs              *SelfServeUserInputsModel `tfsdk:"user_inputs"`
	ActionCardButtonText    types.String              `tfsdk:"action_card_button_text"`
	ExecuteActionButtonText types.String              `tfsdk:"execute_action_button_text"`
	Variant                 types.String              `tfsdk:"variant"`
	Published               types.Bool                `tfsdk:"published"`
	Permissions             *PermissionsModel         `tfsdk:"permissions"`
	Contexts                []TriggerContextModel     `tfsdk:"contexts"`
}

type SelfServeUserInputsModel struct {
	UserProperties  *action.UserPropertiesModel   `tfsdk:"user_properties"`
	Titles          map[string]action.ActionTitle `tfsdk:"titles"`
	RequiredJqQuery types.String                  `tfsdk:"required_jq_query"`
	OrderProperties types.List                    `tfsdk:"order_properties"`
	Steps           []action.Step                 `tfsdk:"steps"`
}

type TriggerContextModel struct {
	On                  types.String `tfsdk:"on"`
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
	UserInput           types.String `tfsdk:"user_input"`
}

type EventTriggerModel struct {
	Type                types.String        `tfsdk:"type"`
	BlueprintIdentifier types.String        `tfsdk:"blueprint_identifier"`
	PropertyIdentifier  types.String        `tfsdk:"property_identifier"`
	Published           types.Bool          `tfsdk:"published"`
	Condition           *NodeConditionModel `tfsdk:"condition"`
}

type NodeConditionModel struct {
	Expressions types.List   `tfsdk:"expressions"`
	Combinator  types.String `tfsdk:"combinator"`
}

type ScheduleTriggerModel struct {
	Cron      types.String `tfsdk:"cron"`
	Published types.Bool   `tfsdk:"published"`
}

type KafkaModel struct {
	Payload   types.String `tfsdk:"payload"`
	OnFailure types.String `tfsdk:"on_failure"`
}

type WebhookModel struct {
	Url          types.String `tfsdk:"url"`
	Agent        types.Bool   `tfsdk:"agent"`
	Synchronized types.Bool   `tfsdk:"synchronized"`
	Method       types.String `tfsdk:"method"`
	Headers      types.Map    `tfsdk:"headers"`
	Body         types.String `tfsdk:"body"`
	OnTimeout    types.String `tfsdk:"on_timeout"`
	OnFailure    types.String `tfsdk:"on_failure"`
}

type IntegrationActionModel struct {
	InstallationId            types.String `tfsdk:"installation_id"`
	IntegrationProvider       types.String `tfsdk:"integration_provider"`
	IntegrationInvocationType types.String `tfsdk:"integration_invocation_type"`
	ExecutionProperties       types.String `tfsdk:"execution_properties"`
	OnFailure                 types.String `tfsdk:"on_failure"`
}

type UpsertEntityModel struct {
	BlueprintIdentifier types.String        `tfsdk:"blueprint_identifier"`
	Mapping             *UpsertMappingModel `tfsdk:"mapping"`
	OnFailure           types.String        `tfsdk:"on_failure"`
}

type UpsertMappingModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
	Icon       types.String `tfsdk:"icon"`
	Teams      types.List   `tfsdk:"teams"`
	Properties types.String `tfsdk:"properties"`
	Relations  types.String `tfsdk:"relations"`
}

type AIModel struct {
	UserPrompt   types.String     `tfsdk:"user_prompt"`
	SystemPrompt types.String     `tfsdk:"system_prompt"`
	Provider     types.String     `tfsdk:"provider"`
	Model        types.String     `tfsdk:"model"`
	Tools        types.List       `tfsdk:"tools"`
	McpServers   []McpServerModel `tfsdk:"mcp_servers"`
	OutputSchema types.String     `tfsdk:"output_schema"`
}

type McpServerModel struct {
	Identifier types.String `tfsdk:"identifier"`
}

type AIAgentModel struct {
	UserPrompt      types.String `tfsdk:"user_prompt"`
	AgentIdentifier types.String `tfsdk:"agent_identifier"`
	Provider        types.String `tfsdk:"provider"`
	Model           types.String `tfsdk:"model"`
	OutputSchema    types.String `tfsdk:"output_schema"`
}

type ConditionModel struct {
	Outlets []ConditionOutletModel `tfsdk:"outlets"`
}

type ConditionOutletModel struct {
	Identifier          types.String      `tfsdk:"identifier"`
	Title               types.String      `tfsdk:"title"`
	Expression          types.String      `tfsdk:"expression"`
	StatusLabel         *StatusLabelModel `tfsdk:"status_label"`
	WorkflowStatusLabel *StatusLabelModel `tfsdk:"workflow_status_label"`
}

type StatusLabelModel struct {
	Text    types.String `tfsdk:"text"`
	Variant types.String `tfsdk:"variant"`
}

type InputModel struct {
	Description   types.String          `tfsdk:"description"`
	UserInputs    *InputUserInputsModel `tfsdk:"user_inputs"`
	Outlets       []InputOutletModel    `tfsdk:"outlets"`
	Notifications []NotificationModel   `tfsdk:"notifications"`
	Responders    *RespondersModel      `tfsdk:"responders"`
}

type InputUserInputsModel struct {
	UserProperties  *action.UserPropertiesModel   `tfsdk:"user_properties"`
	Titles          map[string]action.ActionTitle `tfsdk:"titles"`
	RequiredJqQuery types.String                  `tfsdk:"required_jq_query"`
	OrderProperties types.List                    `tfsdk:"order_properties"`
	Steps           []action.Step                 `tfsdk:"steps"`
	Buttons         []InputButtonModel            `tfsdk:"buttons"`
}

type InputButtonModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Label      types.String `tfsdk:"label"`
	Variant    types.String `tfsdk:"variant"`
	Icon       types.String `tfsdk:"icon"`
}

type InputOutletModel struct {
	Identifier          types.String      `tfsdk:"identifier"`
	Title               types.String      `tfsdk:"title"`
	NumOfResponders     types.Int64       `tfsdk:"num_of_responders"`
	StatusLabel         *StatusLabelModel `tfsdk:"status_label"`
	WorkflowStatusLabel *StatusLabelModel `tfsdk:"workflow_status_label"`
}

type NotificationModel struct {
	Target  types.String             `tfsdk:"target"`
	Fields  []NotificationFieldModel `tfsdk:"fields"`
	Url     types.String             `tfsdk:"url"`
	Method  types.String             `tfsdk:"method"`
	Headers types.Map                `tfsdk:"headers"`
	Body    types.String             `tfsdk:"body"`
	Agent   types.Bool               `tfsdk:"agent"`
}

type NotificationFieldModel struct {
	Label types.String `tfsdk:"label"`
	Value types.String `tfsdk:"value"`
}

type PermissionsModel struct {
	Users  types.List   `tfsdk:"users"`
	Roles  types.List   `tfsdk:"roles"`
	Teams  types.List   `tfsdk:"teams"`
	Policy types.String `tfsdk:"policy"`
}

type RespondersModel struct {
	Users      types.List   `tfsdk:"users"`
	Roles      types.List   `tfsdk:"roles"`
	Teams      types.List   `tfsdk:"teams"`
	UsersQuery types.String `tfsdk:"users_query"`
}

type ConnectionModel struct {
	SourceIdentifier       types.String `tfsdk:"source_identifier"`
	TargetIdentifier       types.String `tfsdk:"target_identifier"`
	Description            types.String `tfsdk:"description"`
	SourceOutletIdentifier types.String `tfsdk:"source_outlet_identifier"`
	Fallback               types.Bool   `tfsdk:"fallback"`
}
