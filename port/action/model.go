package action

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

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

	Default     types.String       `tfsdk:"default"`
	Blueprint   types.String       `tfsdk:"blueprint"`
	Format      types.String       `tfsdk:"format"`
	MaxLength   types.Int64        `tfsdk:"max_length"`
	MinLength   types.Int64        `tfsdk:"min_length"`
	Pattern     types.String       `tfsdk:"pattern"`
	Enum        types.List         `tfsdk:"enum"`
	EnumJqQuery types.String       `tfsdk:"enum_jq_query"`
	Encryption  types.String       `tfsdk:"encryption"`
	Sort        *EntitiesSortModel `tfsdk:"sort"`
}

// StringPropValidationModel is a model used for the validation of StringPropModel resources
type StringPropValidationModel struct {
	Title    string
	Required *bool
}

func (e *StringPropValidationModel) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}

	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["title"].As(&e.Title)
	if err != nil {
		return err
	}

	err = v["required"].As(&e.Required)
	if err != nil {
		return err
	}

	return nil
}

// NumberPropValidationModel is a model used for the validation of StringPropModel resources
type NumberPropValidationModel struct {
	Title    string
	Required *bool
}

func (e *NumberPropValidationModel) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}

	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["title"].As(&e.Title)
	if err != nil {
		return err
	}

	err = v["required"].As(&e.Required)
	if err != nil {
		return err
	}

	return nil
}

// BooleanPropValidationModel is a model used for the validation of StringPropModel resources
type BooleanPropValidationModel struct {
	Title    string
	Required *bool
}

func (e *BooleanPropValidationModel) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}

	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["title"].As(&e.Title)
	if err != nil {
		return err
	}

	err = v["required"].As(&e.Required)
	if err != nil {
		return err
	}

	return nil
}

// ObjectPropValidationModel is a model used for the validation of StringPropModel resources
type ObjectPropValidationModel struct {
	Title    string
	Required *bool
}

func (e *ObjectPropValidationModel) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}

	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["title"].As(&e.Title)
	if err != nil {
		return err
	}

	err = v["required"].As(&e.Required)
	if err != nil {
		return err
	}

	return nil
}

// ArrayPropValidationModel is a model used for the validation of StringPropModel resources
type ArrayPropValidationModel struct {
	Title    string
	Required *bool
}

func (e *ArrayPropValidationModel) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}

	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["title"].As(&e.Title)
	if err != nil {
		return err
	}

	err = v["required"].As(&e.Required)
	if err != nil {
		return err
	}

	return nil
}

type NumberPropModel struct {
	Title          types.String `tfsdk:"title"`
	Icon           types.String `tfsdk:"icon"`
	Description    types.String `tfsdk:"description"`
	Required       types.Bool   `tfsdk:"required"`
	DependsOn      types.List   `tfsdk:"depends_on"`
	DefaultJqQuery types.String `tfsdk:"default_jq_query"`
	Visible        types.Bool   `tfsdk:"visible"`
	VisibleJqQuery types.String `tfsdk:"visible_jq_query"`

	Default     types.Float64 `tfsdk:"default"`
	Maximum     types.Float64 `tfsdk:"maximum"`
	Minimum     types.Float64 `tfsdk:"minimum"`
	Enum        types.List    `tfsdk:"enum"`
	EnumJqQuery types.String  `tfsdk:"enum_jq_query"`
}

type BooleanPropModel struct {
	Title          types.String `tfsdk:"title"`
	Icon           types.String `tfsdk:"icon"`
	Description    types.String `tfsdk:"description"`
	Required       types.Bool   `tfsdk:"required"`
	DependsOn      types.List   `tfsdk:"depends_on"`
	DefaultJqQuery types.String `tfsdk:"default_jq_query"`
	Visible        types.Bool   `tfsdk:"visible"`
	VisibleJqQuery types.String `tfsdk:"visible_jq_query"`

	Default types.Bool `tfsdk:"default"`
}

type EntitiesSortModel struct {
	Property types.String `tfsdk:"property"`
	Order    types.String `tfsdk:"order"`
}

type ArrayPropModel struct {
	Title          types.String `tfsdk:"title"`
	Icon           types.String `tfsdk:"icon"`
	Description    types.String `tfsdk:"description"`
	Required       types.Bool   `tfsdk:"required"`
	DependsOn      types.List   `tfsdk:"depends_on"`
	DefaultJqQuery types.String `tfsdk:"default_jq_query"`
	Visible        types.Bool   `tfsdk:"visible"`
	VisibleJqQuery types.String `tfsdk:"visible_jq_query"`

	MaxItems     types.Int64        `tfsdk:"max_items"`
	MinItems     types.Int64        `tfsdk:"min_items"`
	StringItems  *StringItems       `tfsdk:"string_items"`
	NumberItems  *NumberItems       `tfsdk:"number_items"`
	BooleanItems *BooleanItems      `tfsdk:"boolean_items"`
	ObjectItems  *ObjectItems       `tfsdk:"object_items"`
	Sort         *EntitiesSortModel `tfsdk:"sort"`
}

type ObjectPropModel struct {
	Title          types.String `tfsdk:"title"`
	Icon           types.String `tfsdk:"icon"`
	Description    types.String `tfsdk:"description"`
	Required       types.Bool   `tfsdk:"required"`
	DependsOn      types.List   `tfsdk:"depends_on"`
	DefaultJqQuery types.String `tfsdk:"default_jq_query"`
	Visible        types.Bool   `tfsdk:"visible"`
	VisibleJqQuery types.String `tfsdk:"visible_jq_query"`

	Default    types.String `tfsdk:"default"`
	Encryption types.String `tfsdk:"encryption"`
}

type StringItems struct {
	Blueprint   types.String `tfsdk:"blueprint"`
	Format      types.String `tfsdk:"format"`
	Default     types.List   `tfsdk:"default"`
	Enum        types.List   `tfsdk:"enum"`
	EnumJqQuery types.String `tfsdk:"enum_jq_query"`
	Dataset     types.String `tfsdk:"dataset"`
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

type Step struct {
	Title types.String   `tfsdk:"title"`
	Order []types.String `tfsdk:"order"`
}

type SelfServiceTriggerModel struct {
	BlueprintIdentifier types.String         `tfsdk:"blueprint_identifier"`
	Operation           types.String         `tfsdk:"operation"`
	UserProperties      *UserPropertiesModel `tfsdk:"user_properties"`
	RequiredJqQuery     types.String         `tfsdk:"required_jq_query"`
	OrderProperties     types.List           `tfsdk:"order_properties"`
	Steps               []Step               `tfsdk:"steps"`
	Condition           types.String         `tfsdk:"condition"`
}

type EntityCreatedEventModel struct {
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
}

type EntityUpdatedEventModel struct {
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
}

type EntityDeletedEventModel struct {
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
}

type AnyEntityChangeEventModel struct {
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
}

type TimerPropertyExpiredEventModel struct {
	BlueprintIdentifier types.String `tfsdk:"blueprint_identifier"`
	PropertyIdentifier  types.String `tfsdk:"property_identifier"`
}

type RunCreatedEvent struct {
	ActionIdentifier types.String `tfsdk:"action_identifier"`
}

type RunUpdatedEvent struct {
	ActionIdentifier types.String `tfsdk:"action_identifier"`
}

type AnyRunChangeEvent struct {
	ActionIdentifier types.String `tfsdk:"action_identifier"`
}

type JqConditionModel struct {
	Expressions []types.String `tfsdk:"expressions"`
	Combinator  types.String   `tfsdk:"combinator"`
}

type AutomationTriggerModel struct {
	EntityCreatedEvent        *EntityCreatedEventModel        `tfsdk:"entity_created_event"`
	EntityUpdatedEvent        *EntityUpdatedEventModel        `tfsdk:"entity_updated_event"`
	EntityDeletedEvent        *EntityDeletedEventModel        `tfsdk:"entity_deleted_event"`
	AnyEntityChangeEvent      *AnyEntityChangeEventModel      `tfsdk:"any_entity_change_event"`
	TimerPropertyExpiredEvent *TimerPropertyExpiredEventModel `tfsdk:"timer_property_expired_event"`
	RunCreatedEvent           *RunCreatedEvent                `tfsdk:"run_created_event"`
	RunUpdatedEvent           *RunUpdatedEvent                `tfsdk:"run_updated_event"`
	AnyRunChangeEvent         *AnyRunChangeEvent              `tfsdk:"any_run_change_event"`
	JqCondition               *JqConditionModel               `tfsdk:"jq_condition"`
}

type KafkaMethodModel struct {
	Payload types.String `tfsdk:"payload"`
}

type WebhookMethodModel struct {
	Url          types.String `tfsdk:"url"`
	Agent        types.String `tfsdk:"agent"`
	Synchronized types.String `tfsdk:"synchronized"`
	Method       types.String `tfsdk:"method"`
	Headers      types.Map    `tfsdk:"headers"`
	Body         types.String `tfsdk:"body"`
}

type GithubMethodModel struct {
	Org                  types.String `tfsdk:"org"`
	Repo                 types.String `tfsdk:"repo"`
	Workflow             types.String `tfsdk:"workflow"`
	WorkflowInputs       types.String `tfsdk:"workflow_inputs"`
	ReportWorkflowStatus types.String `tfsdk:"report_workflow_status"`
}

type GitlabMethodModel struct {
	ProjectName       types.String `tfsdk:"project_name"`
	GroupName         types.String `tfsdk:"group_name"`
	DefaultRef        types.String `tfsdk:"default_ref"`
	PipelineVariables types.String `tfsdk:"pipeline_variables"`
}

type AzureMethodModel struct {
	Org     types.String `tfsdk:"org"`
	Webhook types.String `tfsdk:"webhook"`
	Payload types.String `tfsdk:"payload"`
}

type MappingModel struct {
	Properties types.String           `tfsdk:"properties"`
	Relations  types.String           `tfsdk:"relations"`
	Identifier types.String           `tfsdk:"identifier"`
	Teams      basetypes.DynamicValue `tfsdk:"teams"`
	Icon       types.String           `tfsdk:"icon"`
}

type UpsertEntityMethodModel struct {
	Title               types.String  `tfsdk:"title"`
	BlueprintIdentifier types.String  `tfsdk:"blueprint_identifier"`
	Mapping             *MappingModel `tfsdk:"mapping"`
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
	SelfServiceTrigger          *SelfServiceTriggerModel          `tfsdk:"self_service_trigger"`
	AutomationTrigger           *AutomationTriggerModel           `tfsdk:"automation_trigger"`
	KafkaMethod                 *KafkaMethodModel                 `tfsdk:"kafka_method"`
	WebhookMethod               *WebhookMethodModel               `tfsdk:"webhook_method"`
	GithubMethod                *GithubMethodModel                `tfsdk:"github_method"`
	GitlabMethod                *GitlabMethodModel                `tfsdk:"gitlab_method"`
	AzureMethod                 *AzureMethodModel                 `tfsdk:"azure_method"`
	UpsertEntityMethod          *UpsertEntityMethodModel          `tfsdk:"upsert_entity_method"`
	RequiredApproval            types.String                      `tfsdk:"required_approval"`
	ApprovalWebhookNotification *ApprovalWebhookNotificationModel `tfsdk:"approval_webhook_notification"`
	ApprovalEmailNotification   types.Object                      `tfsdk:"approval_email_notification"`
	Publish                     types.Bool                        `tfsdk:"publish"`
}

// ActionValidationModel is a model used for the validation of ActionModel resources
type ActionValidationModel struct {
	ID                          types.String `tfsdk:"id"`
	Identifier                  types.String `tfsdk:"identifier"`
	Blueprint                   types.String `tfsdk:"blueprint"`
	Title                       types.String `tfsdk:"title"`
	Icon                        types.String `tfsdk:"icon"`
	Description                 types.String `tfsdk:"description"`
	SelfServiceTrigger          types.Object `tfsdk:"self_service_trigger"`
	AutomationTrigger           types.Object `tfsdk:"automation_trigger"`
	KafkaMethod                 types.Object `tfsdk:"kafka_method"`
	WebhookMethod               types.Object `tfsdk:"webhook_method"`
	GithubMethod                types.Object `tfsdk:"github_method"`
	GitlabMethod                types.Object `tfsdk:"gitlab_method"`
	AzureMethod                 types.Object `tfsdk:"azure_method"`
	UpsertEntityMethod          types.Object `tfsdk:"upsert_entity_method"`
	RequiredApproval            types.String `tfsdk:"required_approval"`
	ApprovalWebhookNotification types.Object `tfsdk:"approval_webhook_notification"`
	ApprovalEmailNotification   types.Object `tfsdk:"approval_email_notification"`
	Publish                     types.Bool   `tfsdk:"publish"`
}
