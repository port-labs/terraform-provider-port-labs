package entity

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func refreshArrayEntityState(ctx context.Context, state *EntityModel, arrayProperties map[string][]interface{}, blueprint *cli.Blueprint) {
	mapStringItems := make(map[string][]string)
	mapNumberItems := make(map[string][]float64)
	mapBooleanItems := make(map[string][]bool)
	mapObjectItems := make(map[string][]string)

	if state.Properties.ArrayProps == nil {
		state.Properties.ArrayProps = &ArrayPropsModel{
			StringItems:  types.MapNull(types.ListType{ElemType: types.StringType}),
			NumberItems:  types.MapNull(types.ListType{ElemType: types.NumberType}),
			BooleanItems: types.MapNull(types.ListType{ElemType: types.BoolType}),
			ObjectItems:  types.MapNull(types.ListType{ElemType: types.StringType}),
		}
	}
	for k, t := range arrayProperties {

		switch blueprint.Schema.Properties[k].Items["type"] {
		case "string":
			for _, item := range t {
				mapStringItems[k] = append(mapStringItems[k], item.(string))
			}
			if len(t) == 0 {
				mapStringItems[k] = []string{}
			}
			state.Properties.ArrayProps.StringItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, mapStringItems)

		case "number":
			for _, item := range t {
				mapNumberItems[k] = append(mapNumberItems[k], item.(float64))
			}
			if len(t) == 0 {
				mapNumberItems[k] = []float64{}
			}
			state.Properties.ArrayProps.NumberItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.NumberType}, mapNumberItems)

		case "boolean":
			for _, item := range t {
				mapBooleanItems[k] = append(mapBooleanItems[k], item.(bool))
			}
			if len(t) == 0 {
				mapBooleanItems[k] = []bool{}
			}
			state.Properties.ArrayProps.BooleanItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.BoolType}, mapBooleanItems)

		case "object":
			for _, item := range t {
				js, _ := json.Marshal(&item)
				mapObjectItems[k] = append(mapObjectItems[k], string(js))
			}
			if len(t) == 0 {
				mapObjectItems[k] = []string{}
			}
			state.Properties.ArrayProps.ObjectItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, mapObjectItems)

		}
	}
}

func refreshPropertiesEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, blueprint *cli.Blueprint) {
	state.Properties = &EntityPropertiesModel{}
	arrayProperties := make(map[string][]interface{})
	for k, v := range e.Properties {
		switch t := v.(type) {
		case float64:
			if state.Properties.NumberProps == nil {
				state.Properties.NumberProps = make(map[string]types.Float64)
			}
			state.Properties.NumberProps[k] = basetypes.NewFloat64Value(t)
		case string:
			if state.Properties.StringProps == nil {
				state.Properties.StringProps = make(map[string]string)
			}
			state.Properties.StringProps[k] = t

		case bool:
			if state.Properties.BooleanProps == nil {
				state.Properties.BooleanProps = make(map[string]bool)
			}
			state.Properties.BooleanProps[k] = t

		case []interface{}:
			arrayProperties[k] = t
		case interface{}:
			if state.Properties.ObjectProps == nil {
				state.Properties.ObjectProps = make(map[string]string)
			}
			js, _ := json.Marshal(&t)
			state.Properties.ObjectProps[k] = string(js)
		}
	}
	if len(arrayProperties) != 0 {
		refreshArrayEntityState(ctx, state, arrayProperties, blueprint)
	}
}

func refreshRelationsEntityState(ctx context.Context, state *EntityModel, e *cli.Entity) {
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
}

func refreshEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, blueprint *cli.Blueprint) error {
	state.ID = types.StringValue(e.Identifier)
	state.Identifier = types.StringValue(e.Identifier)
	state.Blueprint = types.StringValue(blueprint.Identifier)
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
		refreshPropertiesEntityState(ctx, state, e, blueprint)
	}

	if len(e.Relations) != 0 {
		refreshRelationsEntityState(ctx, state, e)
	}

	return nil
}
