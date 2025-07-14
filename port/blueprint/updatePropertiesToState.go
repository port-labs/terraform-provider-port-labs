package blueprint

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/samber/lo"
)

func SetCommonProperties(v cli.BlueprintProperty, prop interface{}, jsonEscapeHTML bool) {
	properties := []string{"Description", "Icon", "Default", "Title"}
	for _, property := range properties {
		switch property {
		case "Description":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *NumberPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *BooleanPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *ArrayPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *ObjectPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			}
		case "Icon":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *NumberPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *BooleanPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *ArrayPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *ObjectPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			}
		case "Title":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *NumberPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *BooleanPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *ArrayPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *ObjectPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			}
		case "Default":
			if v.Default != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Default = types.StringValue(v.Default.(string))
				case *NumberPropModel:
					p.Default = types.Float64Value(v.Default.(float64))
				case *BooleanPropModel:
					p.Default = types.BoolValue(v.Default.(bool))
				case *ObjectPropModel:
					p.Default, _ = utils.GoObjectToTerraformString(v.Default, jsonEscapeHTML)
				}
			}
		}
	}
}

func (r *BlueprintResource) updatePropertiesToState(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) error {
	properties := &PropertiesModel{}

	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProps == nil {
				properties.StringProps = make(map[string]StringPropModel)
			}
			stringProp := AddStringPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				stringProp.Required = types.BoolValue(true)
			} else {
				stringProp.Required = types.BoolValue(false)
			}

			SetCommonProperties(v, stringProp, r.portClient.JSONEscapeHTML)

			properties.StringProps[k] = *stringProp

		case "number":
			if properties.NumberProps == nil {
				properties.NumberProps = make(map[string]NumberPropModel)
			}

			numberProp := AddNumberPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				numberProp.Required = types.BoolValue(true)
			} else {
				numberProp.Required = types.BoolValue(false)
			}

			SetCommonProperties(v, numberProp, r.portClient.JSONEscapeHTML)

			properties.NumberProps[k] = *numberProp

		case "array":
			if properties.ArrayProps == nil {
				properties.ArrayProps = make(map[string]ArrayPropModel)
			}

			arrayProp := AddArrayPropertiesToState(ctx, &v, r.portClient.JSONEscapeHTML)

			if lo.Contains(b.Schema.Required, k) {
				arrayProp.Required = types.BoolValue(true)
			} else {
				arrayProp.Required = types.BoolValue(false)
			}

			SetCommonProperties(v, arrayProp, r.portClient.JSONEscapeHTML)

			properties.ArrayProps[k] = *arrayProp

		case "boolean":
			if properties.BooleanProps == nil {
				properties.BooleanProps = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			SetCommonProperties(v, booleanProp, r.portClient.JSONEscapeHTML)

			if lo.Contains(b.Schema.Required, k) {
				booleanProp.Required = types.BoolValue(true)
			} else {
				booleanProp.Required = types.BoolValue(false)
			}

			properties.BooleanProps[k] = *booleanProp

		case "object":
			if properties.ObjectProps == nil {
				properties.ObjectProps = make(map[string]ObjectPropModel)
			}

			objectProp := AddObjectPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				objectProp.Required = types.BoolValue(true)
			} else {
				objectProp.Required = types.BoolValue(false)
			}

			SetCommonProperties(v, objectProp, r.portClient.JSONEscapeHTML)

			properties.ObjectProps[k] = *objectProp
		}
	}

	bm.Properties = properties

	return nil
}

func addRelationsToState(b *cli.Blueprint, bm *BlueprintModel) {
	for k, v := range b.Relations {
		if bm.Relations == nil {
			bm.Relations = make(map[string]RelationModel)
		}

		relationModel := &RelationModel{
			Target:      types.StringValue(*v.Target),
			Title:       flex.GoStringToFramework(v.Title),
			Description: flex.GoStringToFramework(v.Description),
			Many:        flex.GoBoolToFramework(v.Many),
			Required:    flex.GoBoolToFramework(v.Required),
		}

		bm.Relations[k] = *relationModel

	}
}

func addMirrorPropertiesToState(b *cli.Blueprint, bm *BlueprintModel) {
	if b.MirrorProperties != nil {
		for k, v := range b.MirrorProperties {

			if bm.MirrorProperties == nil {
				bm.MirrorProperties = make(map[string]MirrorPropertyModel)
			}

			mirrorPropertyModel := &MirrorPropertyModel{
				Path:  types.StringValue(v.Path),
				Title: flex.GoStringToFramework(v.Title),
			}

			bm.MirrorProperties[k] = *mirrorPropertyModel

		}
	}
}

func addCalculationPropertiesToState(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) {
	for k, v := range b.CalculationProperties {
		if bm.CalculationProperties == nil {
			bm.CalculationProperties = make(map[string]CalculationPropertyModel)
		}

		calculationPropertyModel := &CalculationPropertyModel{
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
