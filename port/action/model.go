package action

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WebhookMethodModel struct {
	Url          types.String `tfsdk:"url"`
	Agent        types.Bool   `tfsdk:"agent"`
	Synchronized types.Bool   `tfsdk:"synchronized"`
	Method       types.String `tfsdk:"method"`
}

type Value struct {
	JqQuery types.String `tfsdk:"jq_query"`
}
type Rule struct {
	Blueprint types.String `tfsdk:"blueprint"`
	Property  types.String `tfsdk:"property"`
	Operator  types.String `tfsdk:"operator"`
	Value     *Value       `tfsdk:"value"`
}
type DatasetModel struct {
	Combinator types.String `tfsdk:"combinator"`
	Rules      []Rule       `tfsdk:"rules"`
}

type GithubMethodModel struct {
	Org                  types.String `tfsdk:"org"`
	Repo                 types.String `tfsdk:"repo"`
	Workflow             types.String `tfsdk:"workflow"`
	OmitPayload          types.Bool   `tfsdk:"omit_payload"`
	OmitUserInputs       types.Bool   `tfsdk:"omit_user_inputs"`
	ReportWorkflowStatus types.Bool   `tfsdk:"report_workflow_status"`
}

type AzureMethodModel struct {
	Org     types.String `tfsdk:"org"`
	Webhook types.String `tfsdk:"webhook"`
}

type GitlabMethodModel struct {
	ProjectName    types.String `tfsdk:"project_name"`
	GroupName      types.String `tfsdk:"group_name"`
	OmitPayload    types.Bool   `tfsdk:"omit_payload"`
	OmitUserInputs types.Bool   `tfsdk:"omit_user_inputs"`
	DefaultRef     types.String `tfsdk:"default_ref"`
	Agent          types.Bool   `tfsdk:"agent"`
}

type StringPropModel struct {
	Title          types.String  `tfsdk:"title"`
	Icon           types.String  `tfsdk:"icon"`
	Description    types.String  `tfsdk:"description"`
	Required       types.Bool    `tfsdk:"required"`
	DependsOn      types.List    `tfsdk:"depends_on"`
	Dataset        *DatasetModel `tfsdk:"dataset"`
	DefaultJqQuery types.String  `tfsdk:"default_jq_query"`
	Visible        types.Bool    `tfsdk:"visible"`
	VisibleJqQuery types.String  `tfsdk:"visible_jq_query"`

	Default     types.String `tfsdk:"default"`
	Blueprint   types.String `tfsdk:"blueprint"`
	Format      types.String `tfsdk:"format"`
	MaxLength   types.Int64  `tfsdk:"max_length"`
	MinLength   types.Int64  `tfsdk:"min_length"`
	Pattern     types.String `tfsdk:"pattern"`
	Enum        types.List   `tfsdk:"enum"`
	EnumJqQuery types.String `tfsdk:"enum_jq_query"`
	Encryption  types.String `tfsdk:"encryption"`
}

type NumberPropModel struct {
	Title          types.String  `tfsdk:"title"`
	Icon           types.String  `tfsdk:"icon"`
	Description    types.String  `tfsdk:"description"`
	Required       types.Bool    `tfsdk:"required"`
	DependsOn      types.List    `tfsdk:"depends_on"`
	Dataset        *DatasetModel `tfsdk:"dataset"`
	DefaultJqQuery types.String  `tfsdk:"default_jq_query"`
	Visible        types.Bool    `tfsdk:"visible"`
	VisibleJqQuery types.String  `tfsdk:"visible_jq_query"`

	Default     types.Float64 `tfsdk:"default"`
	Maximum     types.Float64 `tfsdk:"maximum"`
	Minimum     types.Float64 `tfsdk:"minimum"`
	Enum        types.List    `tfsdk:"enum"`
	EnumJqQuery types.String  `tfsdk:"enum_jq_query"`
}

type BooleanPropModel struct {
	Title          types.String  `tfsdk:"title"`
	Icon           types.String  `tfsdk:"icon"`
	Description    types.String  `tfsdk:"description"`
	Required       types.Bool    `tfsdk:"required"`
	DependsOn      types.List    `tfsdk:"depends_on"`
	Dataset        *DatasetModel `tfsdk:"dataset"`
	DefaultJqQuery types.String  `tfsdk:"default_jq_query"`
	Visible        types.Bool    `tfsdk:"visible"`
	VisibleJqQuery types.String  `tfsdk:"visible_jq_query"`

	Default types.Bool `tfsdk:"default"`
}

type ArrayPropModel struct {
	Title          types.String  `tfsdk:"title"`
	Icon           types.String  `tfsdk:"icon"`
	Description    types.String  `tfsdk:"description"`
	Required       types.Bool    `tfsdk:"required"`
	DependsOn      types.List    `tfsdk:"depends_on"`
	Dataset        *DatasetModel `tfsdk:"dataset"`
	DefaultJqQuery types.String  `tfsdk:"default_jq_query"`
	Visible        types.Bool    `tfsdk:"visible"`
	VisibleJqQuery types.String  `tfsdk:"visible_jq_query"`

	MaxItems     types.Int64   `tfsdk:"max_items"`
	MinItems     types.Int64   `tfsdk:"min_items"`
	StringItems  *StringItems  `tfsdk:"string_items"`
	NumberItems  *NumberItems  `tfsdk:"number_items"`
	BooleanItems *BooleanItems `tfsdk:"boolean_items"`
	ObjectItems  *ObjectItems  `tfsdk:"object_items"`
}

type ObjectPropModel struct {
	Title          types.String  `tfsdk:"title"`
	Icon           types.String  `tfsdk:"icon"`
	Description    types.String  `tfsdk:"description"`
	Required       types.Bool    `tfsdk:"required"`
	DependsOn      types.List    `tfsdk:"depends_on"`
	Dataset        *DatasetModel `tfsdk:"dataset"`
	DefaultJqQuery types.String  `tfsdk:"default_jq_query"`
	Visible        types.Bool    `tfsdk:"visible"`
	VisibleJqQuery types.String  `tfsdk:"visible_jq_query"`

	Default    types.String `tfsdk:"default"`
	Encryption types.String `tfsdk:"encryption"`
}

type StringItems struct {
	Blueprint   types.String `tfsdk:"blueprint"`
	Format      types.String `tfsdk:"format"`
	Default     types.List   `tfsdk:"default"`
	Enum        types.List   `tfsdk:"enum"`
	EnumJqQuery types.String `tfsdk:"enum_jq_query"`
}

type NumberItems struct {
	Default     types.List   `tfsdk:"default"`
	Enum        types.List   `tfsdk:"enum"`
	EnumJqQuery types.String `tfsdk:"enum_jq_query"`
}

type BooleanItems struct {
	Default types.List `tfsdk:"default"`
}

type ObjectItems struct {
	Default types.List `tfsdk:"default"`
}

type UserPropertiesModel struct {
	StringProps  map[string]StringPropModel  `tfsdk:"string_props"`
	NumberProps  map[string]NumberPropModel  `tfsdk:"number_props"`
	BooleanProps map[string]BooleanPropModel `tfsdk:"boolean_props"`
	ArrayProps   map[string]ArrayPropModel   `tfsdk:"array_props"`
	ObjectProps  map[string]ObjectPropModel  `tfsdk:"object_props"`
}

type ApprovalWebhookNotificationModel struct {
	Url    types.String `tfsdk:"url"`
	Format types.String `tfsdk:"format"`
}

type ActionModel struct {
	ID                          types.String                      `tfsdk:"id"`
	Identifier                  types.String                      `tfsdk:"identifier"`
	Blueprint                   types.String                      `tfsdk:"blueprint"`
	Title                       types.String                      `tfsdk:"title"`
	Icon                        types.String                      `tfsdk:"icon"`
	Description                 types.String                      `tfsdk:"description"`
	RequiredApproval            types.Bool                        `tfsdk:"required_approval"`
	Trigger                     types.String                      `tfsdk:"trigger"`
	KafkaMethod                 types.Object                      `tfsdk:"kafka_method"`
	WebhookMethod               *WebhookMethodModel               `tfsdk:"webhook_method"`
	GithubMethod                *GithubMethodModel                `tfsdk:"github_method"`
	AzureMethod                 *AzureMethodModel                 `tfsdk:"azure_method"`
	GitlabMethod                *GitlabMethodModel                `tfsdk:"gitlab_method"`
	UserProperties              *UserPropertiesModel              `tfsdk:"user_properties"`
	ApprovalWebhookNotification *ApprovalWebhookNotificationModel `tfsdk:"approval_webhook_notification"`
	ApprovalEmailNotification   types.Object                      `tfsdk:"approval_email_notification"`
	OrderProperties             types.List                        `tfsdk:"order_properties"`
	RequiredJqQuery             types.String                      `tfsdk:"required_jq_query"`
}
