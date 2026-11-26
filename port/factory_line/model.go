package factory_line

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FactoryLineModel struct {
	ID           types.String `tfsdk:"id"`
	Identifier   types.String `tfsdk:"identifier"`
	Title        types.String `tfsdk:"title"`
	Description  types.String `tfsdk:"description"`
	Team         types.List   `tfsdk:"team"`
	MaxRunning   types.Int64  `tfsdk:"max_running"`
	Status       types.String `tfsdk:"status"`
	StatusReason types.String `tfsdk:"status_reason"`
	StageOrder   types.List   `tfsdk:"stage_order"`
	Stages       []StageModel `tfsdk:"stage"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

type StageModel struct {
	Identifier     types.String                         `tfsdk:"identifier"`
	Title          types.String                         `tfsdk:"title"`
	Description    types.String                         `tfsdk:"description"`
	Instructions   types.String                         `tfsdk:"instructions"`
	Team           types.List                           `tfsdk:"team"`
	Executor       types.String                         `tfsdk:"executor"`
	MaxTries       types.Int64                          `tfsdk:"max_tries"`
	AllowTakeOver  types.Bool                           `tfsdk:"allow_take_over"`
	Timeout        types.String                         `tfsdk:"timeout"`
	Checkpoints    types.String                         `tfsdk:"checkpoints"`
	FactoryAI      *FactoryAIExecutorSettingsModel      `tfsdk:"factory_ai_managed_agent"`
	Claude         *ClaudeExecutorSettingsModel         `tfsdk:"claude_managed_agent"`
	Devin          *DevinExecutorSettingsModel          `tfsdk:"devin_managed_agent"`
	CustomWorkflow *CustomWorkflowExecutorSettingsModel `tfsdk:"custom_workflow"`
	Human          *HumanExecutorSettingsModel          `tfsdk:"human"`
}

type FactoryAIExecutorSettingsModel struct {
	InstallationID  types.String `tfsdk:"installation_id"`
	EnvironmentID   types.String `tfsdk:"environment_id"`
	ComputerID      types.String `tfsdk:"computer_id"`
	Model           types.String `tfsdk:"model"`
	ReasoningEffort types.String `tfsdk:"reasoning_effort"`
}

type ClaudeExecutorSettingsModel struct {
	AgentID           types.String `tfsdk:"agent_id"`
	EnvironmentID     types.String `tfsdk:"environment_id"`
	InstallationID    types.String `tfsdk:"installation_id"`
	Session           types.String `tfsdk:"session"`
	GitAccessToken    types.String `tfsdk:"git_access_token"`
	ExecutionIdentity types.String `tfsdk:"execution_identity"`
}

type DevinExecutorSettingsModel struct {
	InstallationID    types.String `tfsdk:"installation_id"`
	PlaybookID        types.String `tfsdk:"playbook_id"`
	Repos             types.List   `tfsdk:"repos"`
	Session           types.String `tfsdk:"session"`
	ExecutionIdentity types.String `tfsdk:"execution_identity"`
}

type CustomWorkflowExecutorSettingsModel struct {
	Workflow types.String `tfsdk:"workflow"`
}

type HumanExecutorSettingsModel struct {
	Assignees types.String `tfsdk:"assignees"`
	Notify    types.String `tfsdk:"notify"`
}
