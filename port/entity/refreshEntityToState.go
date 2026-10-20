package entity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// configuredArrayKeys returns a set of property identifiers that are configured in oldProperties.ArrayProps
func configuredArrayKeys(oldProperties *EntityPropertiesModel) map[string]struct{} {
	keys := make(map[string]struct{})
	if oldProperties == nil || oldProperties.ArrayProps == nil {
		return keys
	}
	if !oldProperties.ArrayProps.StringItems.IsNull() {
		for k := range oldProperties.ArrayProps.StringItems.Elements() {
			keys[k] = struct{}{}
		}
	}
	if !oldProperties.ArrayProps.NumberItems.IsNull() {
		for k := range oldProperties.ArrayProps.NumberItems.Elements() {
			keys[k] = struct{}{}
		}
	}
	if !oldProperties.ArrayProps.BooleanItems.IsNull() {
		for k := range oldProperties.ArrayProps.BooleanItems.Elements() {
			keys[k] = struct{}{}
		}
	}
	if !oldProperties.ArrayProps.ObjectItems.IsNull() {
		for k := range oldProperties.ArrayProps.ObjectItems.Elements() {
			keys[k] = struct{}{}
		}
	}
	return keys
}

func (r *EntityResource) refreshArrayEntityState(ctx context.Context, state *EntityModel, arrayProperties map[string][]interface{}, blueprint *cli.Blueprint, oldProperties *EntityPropertiesModel) {
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
				var oldList types.List
				if oldProperties != nil && oldProperties.ArrayProps != nil && !oldProperties.ArrayProps.ObjectItems.IsNull() {
					if oldListVal, ok := oldProperties.ArrayProps.ObjectItems.Elements()[k]; ok {
						if objectList, ok := oldListVal.(types.List); ok {
							oldList = objectList
						}
					}
				}
				for i, item := range t {
					stringJs, _ := utils.GoObjectToTerraformStringPreferExisting(
						utils.TerraformStringAtList(oldList, i), &item, r.portClient.JSONEscapeHTML)
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

// isPropertyConfigured checks if a property key was configured in oldProperties
func isPropertyConfigured(k string, oldProperties *EntityPropertiesModel) bool {
	if oldProperties == nil {
		return false
	}
	if oldProperties.StringProps != nil {
		if _, ok := oldProperties.StringProps[k]; ok {
			return true
		}
	}
	if oldProperties.NumberProps != nil {
		if _, ok := oldProperties.NumberProps[k]; ok {
			return true
		}
	}
	if oldProperties.BooleanProps != nil {
		if _, ok := oldProperties.BooleanProps[k]; ok {
			return true
		}
	}
	if oldProperties.ObjectProps != nil {
		if _, ok := oldProperties.ObjectProps[k]; ok {
			return true
		}
	}
	if oldProperties.ArrayProps != nil {
		configuredArrays := configuredArrayKeys(oldProperties)
		if _, ok := configuredArrays[k]; ok {
			return true
		}
	}
	return false
}

func (r *EntityResource) refreshPropertiesEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, blueprint *cli.Blueprint) {
	oldProperties := state.Properties
	state.Properties = &EntityPropertiesModel{}
	arrayProperties := make(map[string][]interface{})
	for k, v := range e.Properties {
		// The API returns {} for an unset object property; treat it like null unless state holds an empty object.
		if isEmptyJSONObject(v) && !hasEmptyObjectInState(k, oldProperties) {
			v = nil
		}
		if v == nil {
			// The API returns null for every unset blueprint property. Only keep a null that is in state (the user
			// wrote it explicitly) so unset properties don't cause a perpetual diff. Non-null values are always
			// kept so out-of-band changes show up as drift and import gets every value.
			if !isPropertyConfigured(k, oldProperties) {
				continue
			}
			// The user wrote {} and the API normalized it to null: keep the user's value.
			if hasEmptyObjectInState(k, oldProperties) {
				if state.Properties.ObjectProps == nil {
					state.Properties.ObjectProps = make(map[string]types.String)
				}
				state.Properties.ObjectProps[k] = oldProperties.ObjectProps[k]
				continue
			}
		}

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
			var oldValue types.String
			if oldProperties != nil && oldProperties.ObjectProps != nil {
				oldValue = oldProperties.ObjectProps[k]
			}
			state.Properties.ObjectProps[k], _ = utils.GoObjectToTerraformStringPreferExisting(oldValue, &t, r.portClient.JSONEscapeHTML)
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
		r.refreshArrayEntityState(ctx, state, arrayProperties, blueprint, oldProperties)
	}
	if oldProperties == nil && isEmptyPropertiesModel(state.Properties) {
		state.Properties = nil
	}
}

func isEmptyJSONObject(v interface{}) bool {
	m, ok := v.(map[string]interface{})
	return ok && len(m) == 0
}

// hasEmptyObjectInState reports whether state holds k as an object property whose value is the empty JSON object.
func hasEmptyObjectInState(k string, oldProperties *EntityPropertiesModel) bool {
	if oldProperties == nil || oldProperties.ObjectProps == nil {
		return false
	}
	old, ok := oldProperties.ObjectProps[k]
	if !ok || old.IsNull() || old.IsUnknown() {
		return false
	}
	var m map[string]interface{}
	return json.Unmarshal([]byte(old.ValueString()), &m) == nil && m != nil && len(m) == 0
}

func isEmptyPropertiesModel(p *EntityPropertiesModel) bool {
	return p == nil || (len(p.StringProps) == 0 && len(p.NumberProps) == 0 && len(p.BooleanProps) == 0 &&
		len(p.ObjectProps) == 0 && p.ArrayProps == nil)
}

// isRelationConfigured checks if a relation identifier was configured in oldRelations
func isRelationConfigured(identifier string, oldRelations *RelationModel) bool {
	if oldRelations == nil {
		return false
	}
	if oldRelations.SingleRelation != nil {
		if _, ok := oldRelations.SingleRelation[identifier]; ok {
			return true
		}
	}
	if oldRelations.ManyRelations != nil {
		if _, ok := oldRelations.ManyRelations[identifier]; ok {
			return true
		}
	}
	return false
}

func refreshRelationsEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, oldRelations *RelationModel) {
	state.Relations = &RelationModel{
		SingleRelation: make(map[string]*string),
		ManyRelations:  make(map[string][]string),
	}

	for identifier, r := range e.Relations {
		// The API returns null for every unset relation; only keep a null that is in state. Non-null relations are
		// always kept so drift and import work.
		if r == nil && !isRelationConfigured(identifier, oldRelations) {
			continue
		}

		switch v := r.(type) {
		case []any:
			values := make([]string, 0, len(v))
			for _, rawValue := range v {
				if strVal, ok := rawValue.(string); ok {
					values = append(values, strVal)
				}
			}
			state.Relations.ManyRelations[identifier] = values
		case string:
			if len(v) != 0 {
				state.Relations.SingleRelation[identifier] = &v
			}
		default:
			if r == nil {
				state.Relations.SingleRelation[identifier] = nil
			}
		}
	}
	if len(state.Relations.SingleRelation) == 0 {
		state.Relations.SingleRelation = nil
	}
	if len(state.Relations.ManyRelations) == 0 {
		state.Relations.ManyRelations = nil
	}
	if oldRelations == nil && state.Relations.SingleRelation == nil && state.Relations.ManyRelations == nil {
		state.Relations = nil
	}
}

func (r *EntityResource) refreshEntityState(ctx context.Context, state *EntityModel, e *cli.Entity, blueprint *cli.Blueprint) error {
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

	if len(e.Team) == 0 {
		state.Teams = nil
	} else {
		state.Teams = make([]types.String, len(e.Team))
		for i, t := range e.Team {
			state.Teams[i] = types.StringValue(t)
		}
	}

	if len(e.Properties) == 0 {
		state.Properties = nil
	} else {
		r.refreshPropertiesEntityState(ctx, state, e, blueprint)
	}

	oldRelations := state.Relations
	if len(e.Relations) == 0 {
		state.Relations = nil
	} else {
		refreshRelationsEntityState(ctx, state, e, oldRelations)
	}

	return nil
}
