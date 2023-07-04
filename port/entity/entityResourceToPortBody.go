package entity

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func writeArrayResourceToBody(ctx context.Context, state *EntityModel, properties map[string]interface{}) error {
	if !state.Properties.ArrayProp.StringItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProp.StringItems.Elements() {
			var stringItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "string")
			if err != nil {
				return err
			}
			properties[identifier] = stringItems
		}
	}

	if !state.Properties.ArrayProp.NumberItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProp.NumberItems.Elements() {
			var numberItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "float64")
			if err != nil {
				return err
			}
			properties[identifier] = numberItems
		}
	}

	if !state.Properties.ArrayProp.BooleanItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProp.BooleanItems.Elements() {
			var booleanItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "bool")
			if err != nil {
				return err
			}
			properties[identifier] = booleanItems
		}
	}

	if !state.Properties.ArrayProp.ObjectItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProp.ObjectItems.Elements() {
			var objectItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "object")
			if err != nil {
				return err
			}
			properties[identifier] = objectItems
		}
	}
	return nil
}

func writeRelationsToBody(ctx context.Context, relations *RelationModel) (map[string]interface{}, error) {
	relationsBody := make(map[string]interface{})
	if relations != nil {
		if relations.SingleRelation != nil {
			for identifier, relation := range relations.SingleRelation {
				relationsBody[identifier] = relation
			}
		}

		if relations.ManyRelations != nil {
			for identifier, relations := range relations.ManyRelations {
				relationsBody[identifier] = relations
			}
		}
	}
	return relationsBody, nil
}

func entityResourceToBody(ctx context.Context, state *EntityModel, bp *cli.Blueprint) (*cli.Entity, error) {
	e := &cli.Entity{
		Title:     state.Title.ValueString(),
		Blueprint: bp.Identifier,
	}

	if !state.Identifier.IsUnknown() {
		e.Identifier = state.Identifier.ValueString()
	}

	if state.Teams != nil {
		e.Team = make([]string, len(state.Teams))
		for i, t := range state.Teams {
			e.Team[i] = t.ValueString()
		}
	}

	properties := make(map[string]interface{})
	if state.Properties != nil {
		if state.Properties.StringProp != nil {
			for propIdentifier, prop := range state.Properties.StringProp {
				properties[propIdentifier] = prop
			}
		}

		if state.Properties.NumberProp != nil {
			for propIdentifier, prop := range state.Properties.NumberProp {
				properties[propIdentifier] = prop
			}
		}

		if state.Properties.BooleanProp != nil {
			for propIdentifier, prop := range state.Properties.BooleanProp {
				properties[propIdentifier] = prop
			}
		}

		if state.Properties.ArrayProp != nil {
			err := writeArrayResourceToBody(ctx, state, properties)
			if err != nil {
				return nil, err
			}
		}

		if state.Properties.ObjectProp != nil {
			for identifier, prop := range state.Properties.ObjectProp {
				obj := make(map[string]interface{})
				err := json.Unmarshal([]byte(prop), &obj)
				if err != nil {
					return nil, err
				}
				properties[identifier] = obj
			}
		}
	}

	e.Properties = properties

	relations, err := writeRelationsToBody(ctx, state.Relations)
	if err != nil {
		return nil, err
	}

	e.Relations = relations
	return e, nil
}
