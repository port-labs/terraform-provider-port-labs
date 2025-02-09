package system_blueprint

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/blueprint"
)

type Model struct {
	ID                    types.String                                  `tfsdk:"id"`
	Identifier            types.String                                  `tfsdk:"identifier"`
	Properties            *blueprint.PropertiesModel                    `tfsdk:"properties"`
	Relations             map[string]blueprint.RelationModel            `tfsdk:"relations"`
	MirrorProperties      map[string]blueprint.MirrorPropertyModel      `tfsdk:"mirror_properties"`
	CalculationProperties map[string]blueprint.CalculationPropertyModel `tfsdk:"calculation_properties"`
} 