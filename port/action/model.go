package action

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WebhookMethodModel struct {
	Url   types.String `tfsdk:"url"`
	Agent types.Bool   `tfsdk:"agent"`
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

type StringPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Blueprint   types.String `tfsdk:"blueprint"`
	Description types.String `tfsdk:"description"`
	Default     types.String `tfsdk:"default"`
	Required    types.Bool   `tfsdk:"required"`
	Format      types.String `tfsdk:"format"`
	MaxLength   types.Int64  `tfsdk:"max_length"`
	MinLength   types.Int64  `tfsdk:"min_length"`
	Pattern     types.String `tfsdk:"pattern"`
	Enum        types.List   `tfsdk:"enum"`
}

type NumberPropModel struct {
	Title       types.String  `tfsdk:"title"`
	Icon        types.String  `tfsdk:"icon"`
	Description types.String  `tfsdk:"description"`
	Default     types.Float64 `tfsdk:"default"`
	Required    types.Bool    `tfsdk:"required"`
	Maximum     types.Float64 `tfsdk:"maximum"`
	Minimum     types.Float64 `tfsdk:"minimum"`
	Enum        types.List    `tfsdk:"enum"`
	Blueprint   types.String  `tfsdk:"blueprint"`
	Format      types.String  `tfsdk:"format"`
}

type BooleanPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Default     types.Bool   `tfsdk:"default"`
	Required    types.Bool   `tfsdk:"required"`
	Blueprint   types.String `tfsdk:"blueprint"`
	Format      types.String `tfsdk:"format"`
}

type StringItems struct {
	Format  types.String `tfsdk:"format"`
	Default types.List   `tfsdk:"default"`
}

type NumberItems struct {
	Default types.List `tfsdk:"default"`
}

type BooleanItems struct {
	Default types.List `tfsdk:"default"`
}

type ObjectItems struct {
	Default types.List `tfsdk:"default"`
}

type UserPropertiesModel struct {
	StringProp  map[string]StringPropModel  `tfsdk:"string_prop"`
	NumberProp  map[string]NumberPropModel  `tfsdk:"number_prop"`
	BooleanProp map[string]BooleanPropModel `tfsdk:"boolean_prop"`
	ArrayProp   map[string]ArrayPropModel   `tfsdk:"array_prop"`
	ObjectProp  map[string]ObjectPropModel  `tfsdk:"object_prop"`
}

type ArrayPropModel struct {
	Title        types.String  `tfsdk:"title"`
	Icon         types.String  `tfsdk:"icon"`
	Description  types.String  `tfsdk:"description"`
	MaxItems     types.Int64   `tfsdk:"max_items"`
	MinItems     types.Int64   `tfsdk:"min_items"`
	Required     types.Bool    `tfsdk:"required"`
	Format       types.String  `tfsdk:"format"`
	StringItems  *StringItems  `tfsdk:"string_items"`
	NumberItems  *NumberItems  `tfsdk:"number_items"`
	BooleanItems *BooleanItems `tfsdk:"boolean_items"`
	ObjectItems  *ObjectItems  `tfsdk:"object_items"`
	Blueprint    types.String  `tfsdk:"blueprint"`
}

type ObjectPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Required    types.Bool   `tfsdk:"required"`
	Default     types.String `tfsdk:"default"`
	Spec        types.String `tfsdk:"spec"`
	Blueprint   types.String `tfsdk:"blueprint"`
	Format      types.String `tfsdk:"format"`
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
	UserProperties              *UserPropertiesModel              `tfsdk:"user_properties"`
	ApprovalWebhookNotification *ApprovalWebhookNotificationModel `tfsdk:"approval_webhook_notification"`
	ApprovalEmailNotification   types.Object                      `tfsdk:"approval_email_notification"`
}
