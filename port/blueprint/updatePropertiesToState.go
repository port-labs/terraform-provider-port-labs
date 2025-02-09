package blueprint

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/samber/lo"
)

func setCommonProperties(v cli.BlueprintProperty, prop interface{}) {
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
					js, _ := json.Marshal(v.Default)
					value := string(js)
					p.Default = types.StringValue(value)
				}
			}
		}
	}
}

func updatePropertiesToState(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) error {
	properties := &PropertiesModel{}

	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProps == nil {
				properties.StringProps = make(map[string]StringPropModel)
			}
			stringProp := addStringPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				stringProp.Required = types.BoolValue(true)
			} else {
				stringProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, stringProp)

			properties.StringProps[k] = *stringProp

		case "number":
			if properties.NumberProps == nil {
				properties.NumberProps = make(map[string]NumberPropModel)
			}

			numberProp := addNumberPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				numberProp.Required = types.BoolValue(true)
			} else {
				numberProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, numberProp)

			properties.NumberProps[k] = *numberProp

		case "array":
			if properties.ArrayProps == nil {
				properties.ArrayProps = make(map[string]ArrayPropModel)
			}

			arrayProp := addArrayPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				arrayProp.Required = types.BoolValue(true)
			} else {
				arrayProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, arrayProp)

			properties.ArrayProps[k] = *arrayProp

		case "boolean":
			if properties.BooleanProps == nil {
				properties.BooleanProps = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			setCommonProperties(v, booleanProp)

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

			objectProp := addObjectPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				objectProp.Required = types.BoolValue(true)
			} else {
				objectProp.Required = types.BoolValue(false)
			}

			setCommonProperties(v, objectProp)

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

func UpdatePropertiesToState(ctx context.Context, b *cli.Blueprint, state *PropertiesModel) error {
	if state == nil {
		state = &PropertiesModel{}
	}
	bm := &BlueprintModel{Properties: state}
	return updatePropertiesToState(ctx, b, bm)
}

func UpdateRelationsToState(b *cli.Blueprint, state map[string]RelationModel) error {
	if state == nil {
		state = make(map[string]RelationModel)
	}

	// Instead of merging, just set the state to what's in the plan
	bm := &BlueprintModel{Relations: state}
	
	// Clear existing relations and add only the ones from the plan
	bm.Relations = make(map[string]RelationModel)
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
	return nil
}

func UpdateMirrorPropertiesToState(b *cli.Blueprint, state map[string]MirrorPropertyModel) error {
	if state == nil {
		state = make(map[string]MirrorPropertyModel)
	}

	bm := &BlueprintModel{MirrorProperties: state}
	addMirrorPropertiesToState(b, bm)
	return nil
}

func UpdateCalculationPropertiesToState(ctx context.Context, b *cli.Blueprint, state map[string]CalculationPropertyModel) error {
	if state == nil {
		state = make(map[string]CalculationPropertyModel)
	}

	bm := &BlueprintModel{CalculationProperties: state}
	addCalculationPropertiesToState(ctx, b, bm)
	return nil
}
