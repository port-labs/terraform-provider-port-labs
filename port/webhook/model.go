package webhook

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SecurityModel struct {
	Secret                types.String `tfsdk:"secret"`
	SignatureHeaderName   types.String `tfsdk:"signature_header_name"`
	SignatureAlgorithm    types.String `tfsdk:"signature_algorithm"`
	SignaturePrefix       types.String `tfsdk:"signature_prefix"`
	RequestIdentifierPath types.String `tfsdk:"request_identifier_path"`
}

type EntityModel struct {
	Identifier types.String      `tfsdk:"identifier"`
	Title      types.String      `tfsdk:"title"`
	Icon       types.String      `tfsdk:"icon"`
	Team       types.String      `tfsdk:"team"`
	Properties map[string]string `tfsdk:"properties"`
	Relations  map[string]string `tfsdk:"relations"`
}

type PropertiesModel struct {
	Identifier   types.String `tfsdk:"identifier"`
	Filter       types.String `tfsdk:"filter"`
	ItemsToParse types.String `tfsdk:"items_to_parse"`
	Entity       *EntityModel `tfsdk:"entity"`
}

type MappingsModel struct {
	Properties *PropertiesModel `tfsdk:"properties"`
}

type WebhookModel struct {
	ID          types.String    `tfsdk:"id"`
	Icon        types.String    `tfsdk:"icon"`
	Identifier  types.String    `tfsdk:"identifier"`
	Title       types.String    `tfsdk:"title"`
	Description types.String    `tfsdk:"description"`
	Enabled     types.Bool      `tfsdk:"enabled"`
	Security    *SecurityModel  `tfsdk:"security"`
	Mappings    []MappingsModel `tfsdk:"mappings"`
	CreatedAt   types.String    `tfsdk:"created_at"`
	CreatedBy   types.String    `tfsdk:"created_by"`
	UpdatedAt   types.String    `tfsdk:"updated_at"`
	UpdatedBy   types.String    `tfsdk:"updated_by"`
}
