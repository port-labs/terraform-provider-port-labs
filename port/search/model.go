package search

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type ArrayPropsModel struct {
	StringItems  types.Map `tfsdk:"string_items"`
	NumberItems  types.Map `tfsdk:"number_items"`
	BooleanItems types.Map `tfsdk:"boolean_items"`
	ObjectItems  types.Map `tfsdk:"object_items"`
}

type EntityPropertiesModel struct {
	StringProps  map[string]types.String  `tfsdk:"string_props"`
	NumberProps  map[string]types.Float64 `tfsdk:"number_props"`
	BooleanProps map[string]types.Bool    `tfsdk:"boolean_props"`
	ObjectProps  map[string]types.String  `tfsdk:"object_props"`
	ArrayProps   *ArrayPropsModel         `tfsdk:"array_props"`
}

type RelationModel struct {
	SingleRelation map[string]*string  `tfsdk:"single_relations"`
	ManyRelations  map[string][]string `tfsdk:"many_relations"`
}

type ScorecardRulesModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Status     types.String `tfsdk:"status"`
	Level      types.String `tfsdk:"level"`
}

type ScorecardModel struct {
	Rules []ScorecardRulesModel `tfsdk:"rules"`
	Level types.String          `tfsdk:"level"`
}

type EntityModel struct {
	Identifier types.String               `tfsdk:"identifier"`
	Blueprint  types.String               `tfsdk:"blueprint"`
	Title      types.String               `tfsdk:"title"`
	Icon       types.String               `tfsdk:"icon"`
	RunID      types.String               `tfsdk:"run_id"`
	CreatedAt  types.String               `tfsdk:"created_at"`
	CreatedBy  types.String               `tfsdk:"created_by"`
	UpdatedAt  types.String               `tfsdk:"updated_at"`
	UpdatedBy  types.String               `tfsdk:"updated_by"`
	Properties *EntityPropertiesModel     `tfsdk:"properties"`
	Teams      []types.String             `tfsdk:"teams"`
	Scorecards *map[string]ScorecardModel `tfsdk:"scorecards"`
	Relations  *RelationModel             `tfsdk:"relations"`
}

type SearchDataModel struct {
	ID                          types.String   `tfsdk:"id"`
	Query                       types.String   `tfsdk:"query"`
	ExcludeCalculatedProperties types.Bool     `tfsdk:"exclude_calculated_properties"`
	Include                     []types.String `tfsdk:"include"`
	Exclude                     []types.String `tfsdk:"exclude"`
	AttachTitleToRelation       types.Bool     `tfsdk:"attach_title_to_relation"`
	MatchingBlueprints          []types.String `tfsdk:"matching_blueprints"`
	Entities                    []EntityModel  `tfsdk:"entities"`
}

func (m *SearchDataModel) GenerateID() string {
	// Concatenate the model fields into a single string
	var sb strings.Builder
	sb.WriteString(m.Query.ValueString())
	sb.WriteString(fmt.Sprintf("%t", m.ExcludeCalculatedProperties.ValueBool()))
	for _, include := range m.Include {
		sb.WriteString(include.ValueString())
	}
	for _, exclude := range m.Exclude {
		sb.WriteString(exclude.ValueString())
	}
	sb.WriteString(fmt.Sprintf("%t", m.AttachTitleToRelation.ValueBool()))

	// Compute the SHA-256 hash of the concatenated string
	hash := sha256.Sum256([]byte(sb.String()))

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash[:])

	return hashString
}
