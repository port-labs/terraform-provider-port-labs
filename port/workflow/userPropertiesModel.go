package workflow

import "github.com/hashicorp/terraform-plugin-framework/types"

// The models below intentionally do not reuse the ones from the action package.
// Workflow forms and action forms look alike but are validated by different
// services: workflow inputs additionally support `read_only`, give `disabled` a
// different meaning, and reject attributes such as encryption or boolean array
// items that actions accept.

type UserPropertiesModel struct {
	StringProps  map[string]StringPropModel  `tfsdk:"string_props"`
	NumberProps  map[string]NumberPropModel  `tfsdk:"number_props"`
	BooleanProps map[string]BooleanPropModel `tfsdk:"boolean_props"`
	ArrayProps   map[string]ArrayPropModel   `tfsdk:"array_props"`
	ObjectProps  map[string]ObjectPropModel  `tfsdk:"object_props"`
}

type StringPropModel struct {
	Title           types.String `tfsdk:"title"`
	Icon            types.String `tfsdk:"icon"`
	Description     types.String `tfsdk:"description"`
	Required        types.Bool   `tfsdk:"required"`
	DependsOn       types.List   `tfsdk:"depends_on"`
	Visible         types.Bool   `tfsdk:"visible"`
	VisibleJqQuery  types.String `tfsdk:"visible_jq_query"`
	ReadOnly        types.Bool   `tfsdk:"read_only"`
	ReadOnlyJqQuery types.String `tfsdk:"read_only_jq_query"`
	Disabled        types.Bool   `tfsdk:"disabled"`
	DisabledJqQuery types.String `tfsdk:"disabled_jq_query"`

	Default        types.String       `tfsdk:"default"`
	DefaultJqQuery types.String       `tfsdk:"default_jq_query"`
	Format         types.String       `tfsdk:"format"`
	Blueprint      types.String       `tfsdk:"blueprint"`
	Dataset        *DatasetModel      `tfsdk:"dataset"`
	Sort           *EntitiesSortModel `tfsdk:"sort"`
	MinLength      types.Int64        `tfsdk:"min_length"`
	MaxLength      types.Int64        `tfsdk:"max_length"`
	Pattern        types.String       `tfsdk:"pattern"`
	PatternJqQuery types.String       `tfsdk:"pattern_jq_query"`
	Enum           types.List         `tfsdk:"enum"`
	EnumColors     types.Map          `tfsdk:"enum_colors"`
	EnumJqQuery    types.String       `tfsdk:"enum_jq_query"`
}

type NumberPropModel struct {
	Title           types.String `tfsdk:"title"`
	Icon            types.String `tfsdk:"icon"`
	Description     types.String `tfsdk:"description"`
	Required        types.Bool   `tfsdk:"required"`
	DependsOn       types.List   `tfsdk:"depends_on"`
	Visible         types.Bool   `tfsdk:"visible"`
	VisibleJqQuery  types.String `tfsdk:"visible_jq_query"`
	ReadOnly        types.Bool   `tfsdk:"read_only"`
	ReadOnlyJqQuery types.String `tfsdk:"read_only_jq_query"`
	Disabled        types.Bool   `tfsdk:"disabled"`
	DisabledJqQuery types.String `tfsdk:"disabled_jq_query"`

	Default          types.Float64 `tfsdk:"default"`
	DefaultJqQuery   types.String  `tfsdk:"default_jq_query"`
	Minimum          types.Float64 `tfsdk:"minimum"`
	Maximum          types.Float64 `tfsdk:"maximum"`
	ExclusiveMinimum types.Float64 `tfsdk:"exclusive_minimum"`
	ExclusiveMaximum types.Float64 `tfsdk:"exclusive_maximum"`
	Enum             types.List    `tfsdk:"enum"`
	EnumJqQuery      types.String  `tfsdk:"enum_jq_query"`
}

type BooleanPropModel struct {
	Title           types.String `tfsdk:"title"`
	Icon            types.String `tfsdk:"icon"`
	Description     types.String `tfsdk:"description"`
	Required        types.Bool   `tfsdk:"required"`
	DependsOn       types.List   `tfsdk:"depends_on"`
	Visible         types.Bool   `tfsdk:"visible"`
	VisibleJqQuery  types.String `tfsdk:"visible_jq_query"`
	ReadOnly        types.Bool   `tfsdk:"read_only"`
	ReadOnlyJqQuery types.String `tfsdk:"read_only_jq_query"`
	Disabled        types.Bool   `tfsdk:"disabled"`
	DisabledJqQuery types.String `tfsdk:"disabled_jq_query"`

	Default        types.Bool   `tfsdk:"default"`
	DefaultJqQuery types.String `tfsdk:"default_jq_query"`
}

type ObjectPropModel struct {
	Title           types.String `tfsdk:"title"`
	Icon            types.String `tfsdk:"icon"`
	Description     types.String `tfsdk:"description"`
	Required        types.Bool   `tfsdk:"required"`
	DependsOn       types.List   `tfsdk:"depends_on"`
	Visible         types.Bool   `tfsdk:"visible"`
	VisibleJqQuery  types.String `tfsdk:"visible_jq_query"`
	ReadOnly        types.Bool   `tfsdk:"read_only"`
	ReadOnlyJqQuery types.String `tfsdk:"read_only_jq_query"`
	Disabled        types.Bool   `tfsdk:"disabled"`
	DisabledJqQuery types.String `tfsdk:"disabled_jq_query"`

	Default        types.String `tfsdk:"default"`
	DefaultJqQuery types.String `tfsdk:"default_jq_query"`
	Format         types.String `tfsdk:"format"`
}

type ArrayPropModel struct {
	Title           types.String `tfsdk:"title"`
	Icon            types.String `tfsdk:"icon"`
	Description     types.String `tfsdk:"description"`
	Required        types.Bool   `tfsdk:"required"`
	DependsOn       types.List   `tfsdk:"depends_on"`
	Visible         types.Bool   `tfsdk:"visible"`
	VisibleJqQuery  types.String `tfsdk:"visible_jq_query"`
	ReadOnly        types.Bool   `tfsdk:"read_only"`
	ReadOnlyJqQuery types.String `tfsdk:"read_only_jq_query"`
	Disabled        types.Bool   `tfsdk:"disabled"`
	DisabledJqQuery types.String `tfsdk:"disabled_jq_query"`

	DefaultJqQuery  types.String       `tfsdk:"default_jq_query"`
	MinItems        types.Int64        `tfsdk:"min_items"`
	MinItemsJqQuery types.String       `tfsdk:"min_items_jq_query"`
	MaxItems        types.Int64        `tfsdk:"max_items"`
	MaxItemsJqQuery types.String       `tfsdk:"max_items_jq_query"`
	UniqueItems     types.Bool         `tfsdk:"unique_items"`
	StringItems     *StringItemsModel  `tfsdk:"string_items"`
	NumberItems     *NumberItemsModel  `tfsdk:"number_items"`
	ObjectItems     *ObjectItemsModel  `tfsdk:"object_items"`
	Sort            *EntitiesSortModel `tfsdk:"sort"`
}

type StringItemsModel struct {
	Format      types.String `tfsdk:"format"`
	Blueprint   types.String `tfsdk:"blueprint"`
	Default     types.List   `tfsdk:"default"`
	Enum        types.List   `tfsdk:"enum"`
	EnumJqQuery types.String `tfsdk:"enum_jq_query"`
	EnumColors  types.Map    `tfsdk:"enum_colors"`
	Dataset     types.String `tfsdk:"dataset"`
}

type NumberItemsModel struct {
	Default     types.List   `tfsdk:"default"`
	Enum        types.List   `tfsdk:"enum"`
	EnumJqQuery types.String `tfsdk:"enum_jq_query"`
	EnumColors  types.Map    `tfsdk:"enum_colors"`
}

type ObjectItemsModel struct {
	Format  types.String `tfsdk:"format"`
	Default types.List   `tfsdk:"default"`
}

type EntitiesSortModel struct {
	Property types.String `tfsdk:"property"`
	Order    types.String `tfsdk:"order"`
}

type DatasetModel struct {
	Combinator types.String       `tfsdk:"combinator"`
	Rules      []DatasetRuleModel `tfsdk:"rules"`
}

type DatasetRuleModel struct {
	Blueprint types.String `tfsdk:"blueprint"`
	Property  types.String `tfsdk:"property"`
	Operator  types.String `tfsdk:"operator"`
	Value     *ValueModel  `tfsdk:"value"`
	ValueJson types.String `tfsdk:"value_json"`
	// Group rules carry a combinator and nested rules instead of an operator.
	Combinator types.String       `tfsdk:"combinator"`
	Rules      []DatasetRuleModel `tfsdk:"rules"`
}

type ValueModel struct {
	JqQuery types.String `tfsdk:"jq_query"`
}

// propCommon points at the attributes every input type shares so the mappers can
// read and write them once instead of switching on the property type. The
// framework needs a flat struct per property, so the fields cannot simply be
// embedded.
type propCommon struct {
	Title           *types.String
	Icon            *types.String
	Description     *types.String
	Required        *types.Bool
	DependsOn       *types.List
	Visible         *types.Bool
	VisibleJqQuery  *types.String
	ReadOnly        *types.Bool
	ReadOnlyJqQuery *types.String
	Disabled        *types.Bool
	DisabledJqQuery *types.String
}

func (p *StringPropModel) common() propCommon {
	return propCommon{
		Title: &p.Title, Icon: &p.Icon, Description: &p.Description, Required: &p.Required, DependsOn: &p.DependsOn,
		Visible: &p.Visible, VisibleJqQuery: &p.VisibleJqQuery,
		ReadOnly: &p.ReadOnly, ReadOnlyJqQuery: &p.ReadOnlyJqQuery,
		Disabled: &p.Disabled, DisabledJqQuery: &p.DisabledJqQuery,
	}
}

func (p *NumberPropModel) common() propCommon {
	return propCommon{
		Title: &p.Title, Icon: &p.Icon, Description: &p.Description, Required: &p.Required, DependsOn: &p.DependsOn,
		Visible: &p.Visible, VisibleJqQuery: &p.VisibleJqQuery,
		ReadOnly: &p.ReadOnly, ReadOnlyJqQuery: &p.ReadOnlyJqQuery,
		Disabled: &p.Disabled, DisabledJqQuery: &p.DisabledJqQuery,
	}
}

func (p *BooleanPropModel) common() propCommon {
	return propCommon{
		Title: &p.Title, Icon: &p.Icon, Description: &p.Description, Required: &p.Required, DependsOn: &p.DependsOn,
		Visible: &p.Visible, VisibleJqQuery: &p.VisibleJqQuery,
		ReadOnly: &p.ReadOnly, ReadOnlyJqQuery: &p.ReadOnlyJqQuery,
		Disabled: &p.Disabled, DisabledJqQuery: &p.DisabledJqQuery,
	}
}

func (p *ObjectPropModel) common() propCommon {
	return propCommon{
		Title: &p.Title, Icon: &p.Icon, Description: &p.Description, Required: &p.Required, DependsOn: &p.DependsOn,
		Visible: &p.Visible, VisibleJqQuery: &p.VisibleJqQuery,
		ReadOnly: &p.ReadOnly, ReadOnlyJqQuery: &p.ReadOnlyJqQuery,
		Disabled: &p.Disabled, DisabledJqQuery: &p.DisabledJqQuery,
	}
}

func (p *ArrayPropModel) common() propCommon {
	return propCommon{
		Title: &p.Title, Icon: &p.Icon, Description: &p.Description, Required: &p.Required, DependsOn: &p.DependsOn,
		Visible: &p.Visible, VisibleJqQuery: &p.VisibleJqQuery,
		ReadOnly: &p.ReadOnly, ReadOnlyJqQuery: &p.ReadOnlyJqQuery,
		Disabled: &p.Disabled, DisabledJqQuery: &p.DisabledJqQuery,
	}
}
