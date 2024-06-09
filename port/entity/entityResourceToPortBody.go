package entity

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func writeArrayResourceToBody(ctx context.Context, state *EntityModel, properties map[string]interface{}) error {
	if !state.Properties.ArrayProps.StringItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProps.StringItems.Elements() {
			if !itemArray.IsNull() {
				var stringItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "string")
				if err != nil {
					return err
				}
				properties[identifier] = stringItems
			}
		}
	}

	if !state.Properties.ArrayProps.NumberItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProps.NumberItems.Elements() {
			if !itemArray.IsNull() {
				var numberItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "float64")
				if err != nil {
					return err
				}
				properties[identifier] = numberItems
			}
		}
	}

	if !state.Properties.ArrayProps.BooleanItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProps.BooleanItems.Elements() {
			if !itemArray.IsNull() {
				var booleanItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "bool")
				if err != nil {
					return err
				}
				properties[identifier] = booleanItems
			}
		}
	}

	if !state.Properties.ArrayProps.ObjectItems.IsNull() {
		for identifier, itemArray := range state.Properties.ArrayProps.ObjectItems.Elements() {
			if !itemArray.IsNull() {
				var objectItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "object")
				if err != nil {
					return err
				}
				properties[identifier] = objectItems
			}
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
		Blueprint: bp.Identifier,
	}

	if !state.Title.IsNull() {
		e.Title = state.Title.ValueString()
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
		if state.Properties.StringProps != nil {
			for propIdentifier, prop := range state.Properties.StringProps {
				if !prop.IsNull() {
					properties[propIdentifier] = prop.ValueString()
				}
			}
		}

		if state.Properties.NumberProps != nil {
			for propIdentifier, prop := range state.Properties.NumberProps {
				if !prop.IsNull() {
					properties[propIdentifier] = prop.ValueFloat64()
				}
			}
		}

		if state.Properties.BooleanProps != nil {
			for propIdentifier, prop := range state.Properties.BooleanProps {
				if !prop.IsNull() {
					properties[propIdentifier] = prop.ValueBool()
				}
			}
		}

		if state.Properties.ArrayProps != nil {
			err := writeArrayResourceToBody(ctx, state, properties)
			if err != nil {
				return nil, err
			}
		}

		if state.Properties.ObjectProps != nil {
			for identifier, prop := range state.Properties.ObjectProps {
				if !prop.IsNull() {
					obj := make(map[string]interface{})
					err := json.Unmarshal([]byte(prop.ValueString()), &obj)
					if err != nil {
						return nil, err
					}
					properties[identifier] = obj
				}
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
