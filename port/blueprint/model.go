package blueprint

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WebhookChangelogDestinationModel struct {
	Url   types.String `tfsdk:"url"`
	Agent types.Bool   `tfsdk:"agent"`
}
type TeamInheritanceModel struct {
	Path types.String `tfsdk:"path"`
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
	StringProps  map[string]StringPropModel  `tfsdk:"string_props"`
	NumberProps  map[string]NumberPropModel  `tfsdk:"number_props"`
	BooleanProps map[string]BooleanPropModel `tfsdk:"boolean_props"`
	ArrayProps   map[string]ArrayPropModel   `tfsdk:"array_props"`
	ObjectProps  map[string]ObjectPropModel  `tfsdk:"object_props"`
}

type RelationModel struct {
	Target   	types.String `tfsdk:"target"`
	Title    	types.String `tfsdk:"title"`
	Description 	types.String `tfsdk:"description"`
	Required 	types.Bool   `tfsdk:"required"`
	Many     	types.Bool   `tfsdk:"many"`
}

type MirrorPropertyModel struct {
	Title types.String `tfsdk:"title"`
	Path  types.String `tfsdk:"path"`
}

type CalculationPropertyModel struct {
	Calculation types.String `tfsdk:"calculation"`
	Title       types.String `tfsdk:"title"`
	Format      types.String `tfsdk:"format"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Colorized   types.Bool   `tfsdk:"colorized"`
	Colors      types.Map    `tfsdk:"colors"`
}

type AverageEntitiesModel struct {
	AverageOf     types.String `tfsdk:"average_of"`
	MeasureTimeBy types.String `tfsdk:"measure_time_by"`
}

type AverageByProperty struct {
	MeasureTimeBy types.String `tfsdk:"measure_time_by"`
	AverageOf     types.String `tfsdk:"average_of"`
	Property      types.String `tfsdk:"property"`
}

type AggregateByPropertyModel struct {
	Func     types.String `tfsdk:"func"`
	Property types.String `tfsdk:"property"`
}

type AggregationMethodsModel struct {
	CountEntities       types.Bool                `tfsdk:"count_entities"`
	AverageEntities     *AverageEntitiesModel     `tfsdk:"average_entities"`
	AverageByProperty   *AverageByProperty        `tfsdk:"average_by_property"`
	AggregateByProperty *AggregateByPropertyModel `tfsdk:"aggregate_by_property"`
}

type AggregationPropertyModel struct {
	Title       types.String             `tfsdk:"title"`
	Icon        types.String             `tfsdk:"icon"`
	Description types.String             `tfsdk:"description"`
	Target      types.String             `tfsdk:"target"`
	Method      *AggregationMethodsModel `tfsdk:"method"`
	Query       types.String             `tfsdk:"query"`
}

type BlueprintModel struct {
	ID                          types.String                        `tfsdk:"id"`
	Identifier                  types.String                        `tfsdk:"identifier"`
	Title                       types.String                        `tfsdk:"title"`
	Icon                        types.String                        `tfsdk:"icon"`
	Description                 types.String                        `tfsdk:"description"`
	CreatedAt                   types.String                        `tfsdk:"created_at"`
	CreatedBy                   types.String                        `tfsdk:"created_by"`
	UpdatedAt                   types.String                        `tfsdk:"updated_at"`
	UpdatedBy                   types.String                        `tfsdk:"updated_by"`
	KafkaChangelogDestination   types.Object                        `tfsdk:"kafka_changelog_destination"`
	WebhookChangelogDestination *WebhookChangelogDestinationModel   `tfsdk:"webhook_changelog_destination"`
	TeamInheritance             *TeamInheritanceModel               `tfsdk:"team_inheritance"`
	Properties                  *PropertiesModel                    `tfsdk:"properties"`
	Relations                   map[string]RelationModel            `tfsdk:"relations"`
	MirrorProperties            map[string]MirrorPropertyModel      `tfsdk:"mirror_properties"`
	CalculationProperties       map[string]CalculationPropertyModel `tfsdk:"calculation_properties"`
	ForceDeleteEntities         types.Bool                          `tfsdk:"force_delete_entities"`
	CreateCatalogPage           types.Bool                          `tfsdk:"create_catalog_page"`
}
