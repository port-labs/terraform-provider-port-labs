package entity

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func refreshArrayEntityState(ctx context.Context, state *EntityModel, k string, t []interface{}) {
	if state.Properties.ArrayProp == nil {
		state.Properties.ArrayProp = &ArrayPropModel{
			StringItems:  types.MapNull(types.ListType{ElemType: types.StringType}),
			NumberItems:  types.MapNull(types.ListType{ElemType: types.NumberType}),
			BooleanItems: types.MapNull(types.ListType{ElemType: types.BoolType}),
			ObjectItems:  types.MapNull(types.ListType{ElemType: types.StringType}),
		}
	}
	switch t[0].(type) {
	case string:
		mapItems := make(map[string][]string)
		for _, item := range t {
			mapItems[k] = append(mapItems[k], item.(string))
		}
		state.Properties.ArrayProp.StringItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, mapItems)

	case float64:
		mapItems := make(map[string][]float64)
		for _, item := range t {
			mapItems[k] = append(mapItems[k], item.(float64))
		}
		state.Properties.ArrayProp.NumberItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.NumberType}, mapItems)

	case bool:
		mapItems := make(map[string][]bool)
		for _, item := range t {
			mapItems[k] = append(mapItems[k], item.(bool))
		}
		state.Properties.ArrayProp.BooleanItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.BoolType}, mapItems)

	case map[string]interface{}:
		mapItems := make(map[string][]string)
		for _, item := range t {
			js, _ := json.Marshal(&item)
			mapItems[k] = append(mapItems[k], string(js))
		}
		state.Properties.ArrayProp.ObjectItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, mapItems)

	}
}

func refreshPropertiesEntityState(ctx context.Context, state *EntityModel, e *cli.Entity) {
	state.Properties = &EntityPropertiesModel{}
	for k, v := range e.Properties {
		switch t := v.(type) {
		case float64:
			if state.Properties.NumberProp == nil {
				state.Properties.NumberProp = make(map[string]float64)
			}
			state.Properties.NumberProp[k] = float64(t)
		case string:
			if state.Properties.StringProp == nil {
				state.Properties.StringProp = make(map[string]string)
			}
			state.Properties.StringProp[k] = t

		case bool:
			if state.Properties.BooleanProp == nil {
				state.Properties.BooleanProp = make(map[string]bool)
			}
			state.Properties.BooleanProp[k] = t

		case []interface{}:
			refreshArrayEntityState(ctx, state, k, t)
		case interface{}:
			if state.Properties.ObjectProp == nil {
				state.Properties.ObjectProp = make(map[string]string)
			}

			js, _ := json.Marshal(&t)
			state.Properties.ObjectProp[k] = string(js)
		}
	}
}

func refreshRelationsEntityState(ctx context.Context, state *EntityModel, e *cli.Entity) {
	// relations := make(map[string][]string)
	relations := &RelationModel{
		SingleRelation: make(map[string]string),
		ManyRelations:  make(map[string][]string),
	}

	for identifier, r := range e.Relations {
		switch v := r.(type) {
		case []string:
			if len(v) != 0 {
				relations.ManyRelations[identifier] = v
			}

		case string:
			if len(v) != 0 {
				relations.SingleRelation[identifier] = v
			}
		}
	}

	// if len(relations) != 0 {
	// 	state.Relations = relations
	// }
}

func refreshEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, blueprint string) error {
	state.ID = types.StringValue(e.Identifier)
	state.Identifier = types.StringValue(e.Identifier)
	state.Blueprint = types.StringValue(blueprint)
	state.Title = types.StringValue(e.Title)
	state.CreatedAt = types.StringValue(e.CreatedAt.String())
	state.CreatedBy = types.StringValue(e.CreatedBy)
	state.UpdatedAt = types.StringValue(e.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(e.UpdatedBy)

	if len(e.Team) != 0 {
		state.Teams = make([]types.String, len(e.Team))
		for i, t := range e.Team {
			state.Teams[i] = types.StringValue(t)
		}
	}

	if len(e.Properties) != 0 {
		refreshPropertiesEntityState(ctx, state, e)
	}

	if len(e.Relations) != 0 {
		refreshRelationsEntityState(ctx, state, e)
	}

	return nil
}
