package system_blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/blueprint"
	"github.com/samber/lo"
)

func updatePropertiesToState(ctx context.Context, b *cli.Blueprint, systemBp *cli.Blueprint, bm *SystemBlueprintModel) error {
	properties := &blueprint.PropertiesModel{}

	for k, v := range b.Schema.Properties {
		// Skip if the property exists in systemBp
		if systemBp != nil {
			if _, exists := systemBp.Schema.Properties[k]; exists {
				continue
			}
		}

		switch v.Type {
		case "string":
			if properties.StringProps == nil {
				properties.StringProps = make(map[string]blueprint.StringPropModel)
			}
			stringProp := blueprint.AddStringPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				stringProp.Required = types.BoolValue(true)
			} else {
				stringProp.Required = types.BoolValue(false)
			}

			blueprint.SetCommonProperties(v, stringProp)

			properties.StringProps[k] = *stringProp

		case "number":
			if properties.NumberProps == nil {
				properties.NumberProps = make(map[string]blueprint.NumberPropModel)
			}

			numberProp := blueprint.AddNumberPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				numberProp.Required = types.BoolValue(true)
			} else {
				numberProp.Required = types.BoolValue(false)
			}

			blueprint.SetCommonProperties(v, numberProp)

			properties.NumberProps[k] = *numberProp

		case "array":
			if properties.ArrayProps == nil {
				properties.ArrayProps = make(map[string]blueprint.ArrayPropModel)
			}

			arrayProp := blueprint.AddArrayPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				arrayProp.Required = types.BoolValue(true)
			} else {
				arrayProp.Required = types.BoolValue(false)
			}

			blueprint.SetCommonProperties(v, arrayProp)

			properties.ArrayProps[k] = *arrayProp

		case "boolean":
			if properties.BooleanProps == nil {
				properties.BooleanProps = make(map[string]blueprint.BooleanPropModel)
			}

			booleanProp := &blueprint.BooleanPropModel{}

			blueprint.SetCommonProperties(v, booleanProp)

			if lo.Contains(b.Schema.Required, k) {
				booleanProp.Required = types.BoolValue(true)
			} else {
				booleanProp.Required = types.BoolValue(false)
			}

			properties.BooleanProps[k] = *booleanProp

		case "object":
			if properties.ObjectProps == nil {
				properties.ObjectProps = make(map[string]blueprint.ObjectPropModel)
			}

			objectProp := blueprint.AddObjectPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				objectProp.Required = types.BoolValue(true)
			} else {
				objectProp.Required = types.BoolValue(false)
			}

			blueprint.SetCommonProperties(v, objectProp)

			properties.ObjectProps[k] = *objectProp
		}
	}

	bm.Properties = properties

	return nil
}


func addRelationsToState(b *cli.Blueprint, systemBp *cli.Blueprint, bm *SystemBlueprintModel) {
	for k, v := range b.Relations {
		// Skip if the relation exists in systemBp
		if systemBp != nil {
			if _, exists := systemBp.Relations[k]; exists {
				continue
			}
		}

		if bm.Relations == nil {
			bm.Relations = make(map[string]blueprint.RelationModel)
		}

		relationModel := &blueprint.RelationModel{
			Target:      types.StringValue(*v.Target),
			Title:       flex.GoStringToFramework(v.Title),
			Description: flex.GoStringToFramework(v.Description),
			Many:        flex.GoBoolToFramework(v.Many),
			Required:    flex.GoBoolToFramework(v.Required),
		}

		bm.Relations[k] = *relationModel
	}
}

func addMirrorPropertiesToState(b *cli.Blueprint, systemBp *cli.Blueprint, bm *SystemBlueprintModel) {
	if b.MirrorProperties != nil {
		for k, v := range b.MirrorProperties {
			// Skip if the mirror property exists in systemBp
			if systemBp != nil {
				if _, exists := systemBp.MirrorProperties[k]; exists {
					continue
				}
			}

			if bm.MirrorProperties == nil {
				bm.MirrorProperties = make(map[string]blueprint.MirrorPropertyModel)
			}

			mirrorPropertyModel := &blueprint.MirrorPropertyModel{
				Path:  types.StringValue(v.Path),
				Title: flex.GoStringToFramework(v.Title),
			}

			bm.MirrorProperties[k] = *mirrorPropertyModel
		}
	}
}

func addCalculationPropertiesToState(ctx context.Context, b *cli.Blueprint, systemBp *cli.Blueprint, bm *SystemBlueprintModel) {
	for k, v := range b.CalculationProperties {
		// Skip if the calculation property exists in systemBp
		if systemBp != nil {
			if _, exists := systemBp.CalculationProperties[k]; exists {
				continue
			}
		}
		if bm.CalculationProperties == nil {
			bm.CalculationProperties = make(map[string]blueprint.CalculationPropertyModel)
		}

		calculationPropertyModel := &blueprint.CalculationPropertyModel{
			Calculation: types.StringValue(v.Calculation),
			Type:        types.StringValue(v.Type),
			Title:       flex.GoStringToFramework(v.Title),
			Icon:        flex.GoStringToFramework(v.Icon),
			Description: flex.GoStringToFramework(v.Description),
			Format:      flex.GoStringToFramework(v.Format),
			Colorized:   flex.GoBoolToFramework(v.Colorized),
		}

		if v.Colors != nil {
			calculationPropertyModel.Colors, _ = types.MapValueFrom(ctx, types.StringType, v.Colors)
		} else {
			calculationPropertyModel.Colors = types.MapNull(types.StringType)
		}

		bm.CalculationProperties[k] = *calculationPropertyModel
	}
}
