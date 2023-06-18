package blueprint

import "github.com/hashicorp/terraform-plugin-framework/types"

type ChangelogDestinationModel struct {
	Type  types.String `tfsdk:"type"`
	Url   types.String `tfsdk:"url"`
	Agent types.Bool   `tfsdk:"agent"`
}

type SpecAuthenticationModel struct {
	AuthorizationUrl types.String `tfsdk:"authorization_url"`
	TokenUrl         types.String `tfsdk:"token_url"`
	ClientId         types.String `tfsdk:"client_id"`
}

type StringPropModel struct {
	Title              types.String             `tfsdk:"title"`
	Icon               types.String             `tfsdk:"icon"`
	Description        types.String             `tfsdk:"description"`
	Default            types.String             `tfsdk:"default"`
	Required           types.Bool               `tfsdk:"required"`
	Format             types.String             `tfsdk:"format"`
	MaxLength          types.Int64              `tfsdk:"max_length"`
	MinLength          types.Int64              `tfsdk:"min_length"`
	Pattern            types.String             `tfsdk:"pattern"`
	Enum               types.List               `tfsdk:"enum"`
	EnumColors         types.Map                `tfsdk:"enum_colors"`
	Spec               types.String             `tfsdk:"spec"`
	SpecAuthentication *SpecAuthenticationModel `tfsdk:"spec_authentication"`
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
	EnumColors  types.Map     `tfsdk:"enum_colors"`
}

type BooleanPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Default     types.Bool   `tfsdk:"default"`
	Required    types.Bool   `tfsdk:"required"`
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

type ArrayPropModel struct {
	Title        types.String  `tfsdk:"title"`
	Icon         types.String  `tfsdk:"icon"`
	Description  types.String  `tfsdk:"description"`
	MaxItems     types.Int64   `tfsdk:"max_items"`
	MinItems     types.Int64   `tfsdk:"min_items"`
	Required     types.Bool    `tfsdk:"required"`
	StringItems  *StringItems  `tfsdk:"string_items"`
	NumberItems  *NumberItems  `tfsdk:"number_items"`
	BooleanItems *BooleanItems `tfsdk:"boolean_items"`
	ObjectItems  *ObjectItems  `tfsdk:"object_items"`
}

type ObjectPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Required    types.Bool   `tfsdk:"required"`
	Default     types.String `tfsdk:"default"`
	Spec        types.String `tfsdk:"spec"`
}

type PropertiesModel struct {
	StringProp  map[string]StringPropModel  `tfsdk:"string_prop"`
	NumberProp  map[string]NumberPropModel  `tfsdk:"number_prop"`
	BooleanProp map[string]BooleanPropModel `tfsdk:"boolean_prop"`
	ArrayProp   map[string]ArrayPropModel   `tfsdk:"array_prop"`
	ObjectProp  map[string]ObjectPropModel  `tfsdk:"object_prop"`
}

type BlueprintModel struct {
	ID                   types.String               `tfsdk:"id"`
	Identifier           types.String               `tfsdk:"identifier"`
	Title                types.String               `tfsdk:"title"`
	Icon                 types.String               `tfsdk:"icon"`
	Description          types.String               `tfsdk:"description"`
	CreatedAt            types.String               `tfsdk:"created_at"`
	CreatedBy            types.String               `tfsdk:"created_by"`
	UpdatedAt            types.String               `tfsdk:"updated_at"`
	UpdatedBy            types.String               `tfsdk:"updated_by"`
	ChangelogDestination *ChangelogDestinationModel `tfsdk:"changelog_destination"`
	Properties           *PropertiesModel           `tfsdk:"properties"`
}
