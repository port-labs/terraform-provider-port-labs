package system_blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/blueprint"
)

func MergeProperties(ctx context.Context, existing map[string]cli.BlueprintProperty, state *blueprint.PropertiesModel) (map[string]cli.BlueprintProperty, []string, error) {
	merged := make(map[string]cli.BlueprintProperty)
	for k, v := range existing {
		merged[k] = v
	}

	if state != nil {
		props, required, err := blueprint.PropsResourceToBody(ctx, state)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range props {
			merged[k] = v
		}
		return merged, required, nil
	}
	return merged, nil, nil
}

func MergeRelations(existing map[string]cli.Relation, state map[string]blueprint.RelationModel) map[string]cli.Relation {

	merged := make(map[string]cli.Relation)
	for k, v := range existing {
		merged[k] = v
	}

	if state != nil {
		relations := blueprint.RelationsResourceToBody(state)
		for k, v := range relations {
			merged[k] = v
		}
	}
	return merged
}

func MergeMirrorProperties(existing map[string]cli.BlueprintMirrorProperty, state map[string]blueprint.MirrorPropertyModel) map[string]cli.BlueprintMirrorProperty {
	merged := make(map[string]cli.BlueprintMirrorProperty)
	for k, v := range existing {
		merged[k] = v
	}

	if state != nil {
		mirrorProps := blueprint.MirrorPropertiesToBody(state)
		for k, v := range mirrorProps {
			merged[k] = v
		}
	}
	return merged
}

func MergeCalculationProperties(ctx context.Context, existing map[string]cli.BlueprintCalculationProperty, state map[string]blueprint.CalculationPropertyModel) map[string]cli.BlueprintCalculationProperty {
	merged := make(map[string]cli.BlueprintCalculationProperty)
	for k, v := range existing {
		merged[k] = v
	}

	if state != nil {
		calcProps := blueprint.CalculationPropertiesToBody(ctx, state)
		for k, v := range calcProps {
			merged[k] = v
		}
	}
	return merged
}

func UpdateRelationsToState(b *cli.Blueprint, state map[string]blueprint.RelationModel) error {
	if state == nil {
		state = make(map[string]blueprint.RelationModel)
	}

	bm := &blueprint.BlueprintModel{Relations: state}

	bm.Relations = make(map[string]blueprint.RelationModel)
	for k, v := range b.Relations {
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
	return nil
}

func UpdateMirrorPropertiesToState(b *cli.Blueprint, state map[string]blueprint.MirrorPropertyModel) error {
	if state == nil {
		state = make(map[string]blueprint.MirrorPropertyModel)
	}

	bm := &blueprint.BlueprintModel{MirrorProperties: state}
	for k, v := range b.MirrorProperties {
		if bm.MirrorProperties == nil {
			bm.MirrorProperties = make(map[string]blueprint.MirrorPropertyModel)
		}

		mirrorPropertyModel := &blueprint.MirrorPropertyModel{
			Path:  types.StringValue(v.Path),
			Title: flex.GoStringToFramework(v.Title),
		}

		bm.MirrorProperties[k] = *mirrorPropertyModel
	}
	return nil
}

func UpdateCalculationPropertiesToState(ctx context.Context, b *cli.Blueprint, state map[string]blueprint.CalculationPropertyModel) error {
	if state == nil {
		state = make(map[string]blueprint.CalculationPropertyModel)
	}

	bm := &blueprint.BlueprintModel{CalculationProperties: state}
	for k, v := range b.CalculationProperties {
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
	return nil
}
