package entity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshArrayEntityState(ctx context.Context, state *EntityModel, arrayProperties map[string][]interface{}, blueprint *cli.Blueprint) {
	mapStringItems := make(map[string][]*string)
	mapNumberItems := make(map[string][]*float64)
	mapBooleanItems := make(map[string][]*bool)
	mapObjectItems := make(map[string][]*string)

	if state.Properties.ArrayProps == nil {
		state.Properties.ArrayProps = &ArrayPropsModel{
			StringItems:  types.MapNull(types.ListType{ElemType: types.StringType}),
			NumberItems:  types.MapNull(types.ListType{ElemType: types.Float64Type}),
			BooleanItems: types.MapNull(types.ListType{ElemType: types.BoolType}),
			ObjectItems:  types.MapNull(types.ListType{ElemType: types.StringType}),
		}
	}
	for k, t := range arrayProperties {

		switch blueprint.Schema.Properties[k].Items["type"] {
		// array without items type is array of string by default
		case "string", nil:
			if t != nil {
				for _, item := range t {
					stringItem := item.(string)
					mapStringItems[k] = append(mapStringItems[k], &stringItem)
				}
				if len(t) == 0 {
					mapStringItems[k] = []*string{}
				}
			} else {
				mapStringItems[k] = nil
			}
			state.Properties.ArrayProps.StringItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, mapStringItems)
		case "number":
			if t != nil {
				for _, item := range t {
					floatItem := item.(float64)
					mapNumberItems[k] = append(mapNumberItems[k], &floatItem)
				}
				if len(t) == 0 {
					mapNumberItems[k] = []*float64{}
				}
			} else {
				mapNumberItems[k] = nil
			}
			state.Properties.ArrayProps.NumberItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.Float64Type}, mapNumberItems)

		case "boolean":
			if t != nil {
				for _, item := range t {
					boolItem := item.(bool)
					mapBooleanItems[k] = append(mapBooleanItems[k], &boolItem)
				}
				if len(t) == 0 {
					mapBooleanItems[k] = []*bool{}
				}
			} else {
				mapBooleanItems[k] = nil
			}
			state.Properties.ArrayProps.BooleanItems, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.BoolType}, mapBooleanItems)

		case "object":
			if t != nil {
				for _, item := range t {
					js, _ := json.Marshal(&item)
					stringJs := string(js)
					mapObjectItems[k] = append(mapObjectItems[k], &stringJs)
				}
				if len(t) == 0 {
					mapObjectItems[k] = []*string{}
				}
			} else {
				mapObjectItems[k] = nil
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
			state.Properties.NumberProps[k] = types.Float64Value(t)
		case string:
			if state.Properties.StringProps == nil {
				state.Properties.StringProps = make(map[string]types.String)
			}
			state.Properties.StringProps[k] = types.StringValue(t)
		case bool:
			if state.Properties.BooleanProps == nil {
				state.Properties.BooleanProps = make(map[string]types.Bool)
			}
			state.Properties.BooleanProps[k] = types.BoolValue(t)
		case []interface{}:
			arrayProperties[k] = t
		case interface{}:
			if state.Properties.ObjectProps == nil {
				state.Properties.ObjectProps = make(map[string]types.String)
			}
			js, _ := json.Marshal(&t)
			state.Properties.ObjectProps[k] = types.StringValue(string(js))
		case nil:
			switch blueprint.Schema.Properties[k].Type {
			case "string":
				if state.Properties.StringProps == nil {
					state.Properties.StringProps = make(map[string]types.String)
				}
				state.Properties.StringProps[k] = types.StringNull()
			case "number":
				if state.Properties.NumberProps == nil {
					state.Properties.NumberProps = make(map[string]types.Float64)
				}
				state.Properties.NumberProps[k] = types.Float64Null()
			case "boolean":
				if state.Properties.BooleanProps == nil {
					state.Properties.BooleanProps = make(map[string]types.Bool)
				}
				state.Properties.BooleanProps[k] = types.BoolNull()
			case "object":
				if state.Properties.ObjectProps == nil {
					state.Properties.ObjectProps = make(map[string]types.String)
				}
				state.Properties.ObjectProps[k] = types.StringNull()
			case "array":
				arrayProperties[k] = []interface{}(nil)
			}
		}
	}
	if len(arrayProperties) != 0 {
		refreshArrayEntityState(ctx, state, arrayProperties, blueprint)
	}
}

func refreshRelationsEntityState(ctx context.Context, state *EntityModel, e *cli.Entity) {
	state.Relations = &RelationModel{}

	for identifier, r := range e.Relations {
		switch v := r.(type) {
		case []any:
			if state.Relations.ManyRelations == nil {
				state.Relations.ManyRelations = make(map[string][]string)
			}
			state.Relations.ManyRelations[identifier] = make([]string, 0, len(v))
			for _, rawValue := range v {
				if strVal, ok := rawValue.(string); ok {
					state.Relations.ManyRelations[identifier] = append(state.Relations.ManyRelations[identifier], strVal)
				}
			}
		case string:
			if state.Relations.SingleRelation == nil {
				state.Relations.SingleRelation = make(map[string]*string)
			}
			if len(v) != 0 {
				state.Relations.SingleRelation[identifier] = &v
			}
		}
	}
}

func refreshEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, blueprint *cli.Blueprint) error {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", blueprint.Identifier, e.Identifier))
	state.Identifier = types.StringValue(e.Identifier)
	state.Blueprint = types.StringValue(blueprint.Identifier)
	state.Title = types.StringValue(e.Title)

	if e.Icon != "" {
		state.Icon = types.StringValue(e.Icon)
	} else {
		state.Icon = types.StringNull()
	}

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
