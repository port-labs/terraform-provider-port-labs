package search

import (
	"context"
	"fmt"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func (d *SearchDataSource) refreshArrayEntityState(ctx context.Context, state *EntityModel, arrayProperties map[string][]interface{}, blueprint *cli.Blueprint) {
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
		case "string":
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
					stringJs, _ := utils.GoObjectToTerraformString(&item, d.portClient.JSONEscapeHTML)
					mapObjectItems[k] = append(mapObjectItems[k], stringJs.ValueStringPointer())
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

func (d *SearchDataSource) refreshPropertiesEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, blueprint *cli.Blueprint) {
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
			state.Properties.ObjectProps[k], _ = utils.GoObjectToTerraformString(&t, d.portClient.JSONEscapeHTML)
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
		d.refreshArrayEntityState(ctx, state, arrayProperties, blueprint)
	}
}

func refreshRelationsEntityState(state *EntityModel, e *cli.Entity) {
	relations := &RelationModel{
		SingleRelation: make(map[string]*string),
		ManyRelations:  make(map[string][]string),
	}

	for identifier, r := range e.Relations {
		switch v := r.(type) {
		case []interface{}:
			if len(v) != 0 {
				switch v[0].(type) {
				case string:
					relations.ManyRelations[identifier] = make([]string, len(v))
					for i, s := range v {
						relations.ManyRelations[identifier][i] = s.(string)
					}
				}
			}

		case interface{}:
			if v != nil {
				value := fmt.Sprintf("%v", v)
				relations.SingleRelation[identifier] = &value
			}
		}
	}

	state.Relations = relations
}

func refreshScorecardsEntityState(state *EntityModel, e *cli.Entity) {
	if len(e.Scorecards) != 0 {
		state.Scorecards = &map[string]ScorecardModel{}
		*state.Scorecards = make(map[string]ScorecardModel)

		for k, v := range e.Scorecards {
			rules := make([]ScorecardRulesModel, len(v.Rules))
			for i, r := range v.Rules {
				rules[i] = ScorecardRulesModel{
					Identifier: types.StringValue(r.Identifier),
					Status:     types.StringValue(r.Status),
					Level:      types.StringValue(r.Level),
				}
			}
			(*state.Scorecards)[k] = ScorecardModel{
				Rules: rules,
				Level: types.StringValue(v.Level),
			}
		}
	}
}

func (d *SearchDataSource) refreshEntityState(ctx context.Context, e *cli.Entity, b *cli.Blueprint) *EntityModel {
	state := &EntityModel{}
	state.Identifier = types.StringValue(e.Identifier)
	state.Blueprint = types.StringValue(e.Blueprint)
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
		d.refreshPropertiesEntityState(ctx, state, e, b)
	}

	if len(e.Relations) != 0 {
		refreshRelationsEntityState(state, e)
	}

	if len(e.Scorecards) != 0 {
		refreshScorecardsEntityState(state, e)
	}

	return state
}
