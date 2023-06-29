package entity

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &EntityResource{}
var _ resource.ResourceWithImportState = &EntityResource{}

func NewEntityResource() resource.Resource {
	return &EntityResource{}
}

type EntityResource struct {
	portClient *cli.PortClient
}

func (r *EntityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entity"
}

func (r *EntityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *EntityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntityModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.Blueprint.ValueString()
	e, statusCode, err := r.portClient.ReadEntity(ctx, state.Identifier.ValueString(), state.Blueprint.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read entity", err.Error())
		return
	}

	err = refreshEntityState(ctx, state, e, blueprintIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed writing entity fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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
			case interface{}:
				if state.Properties.ObjectProp == nil {
					state.Properties.ObjectProp = make(map[string]string)
				}

				js, _ := json.Marshal(&t)
				state.Properties.ObjectProp[k] = string(js)
			}
		}
	}

	if len(e.Relations) != 0 {
		relations := make(map[string][]string)
		for identifier, r := range e.Relations {
			switch v := r.(type) {
			case []string:
				if len(v) != 0 {
					relations[identifier] = v
				}

			case string:
				if len(v) != 0 {
					relations[identifier] = []string{v}
				}
			}
		}
		if len(relations) != 0 {
			state.Relations, _ = types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, relations)
		}
	}

	return nil
}
func (r *EntityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *EntityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, state.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	e, err := entityResourceToBody(ctx, state, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	runID := ""
	if !state.RunID.IsNull() {
		runID = state.RunID.ValueString()
	}

	en, err := r.portClient.CreateEntity(ctx, e, runID)
	if err != nil {
		resp.Diagnostics.AddError("failed to create entity", err.Error())
		return
	}

	state.ID = types.StringValue(en.Identifier)
	state.Identifier = types.StringValue(en.Identifier)

	writeEntityComputedFieldsToState(state, en)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func entityResourceToBody(ctx context.Context, state *EntityModel, bp *cli.Blueprint) (*cli.Entity, error) {
	e := &cli.Entity{
		Title:     state.Title.ValueString(),
		Blueprint: bp.Identifier,
	}

	if !state.Identifier.IsNull() {
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
			if !state.Properties.ArrayProp.StringItems.IsNull() {
				for identifier, itemArray := range state.Properties.ArrayProp.StringItems.Elements() {
					var stringItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "string")
					if err != nil {
						return nil, err
					}
					properties[identifier] = stringItems
				}
			}

			if !state.Properties.ArrayProp.NumberItems.IsNull() {
				for identifier, itemArray := range state.Properties.ArrayProp.NumberItems.Elements() {
					var numberItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "float64")
					if err != nil {
						return nil, err
					}
					properties[identifier] = numberItems
				}
			}

			if !state.Properties.ArrayProp.BooleanItems.IsNull() {
				for identifier, itemArray := range state.Properties.ArrayProp.BooleanItems.Elements() {
					var booleanItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "bool")
					if err != nil {
						return nil, err
					}
					properties[identifier] = booleanItems
				}
			}

			if !state.Properties.ArrayProp.ObjectItems.IsNull() {
				for identifier, itemArray := range state.Properties.ArrayProp.ObjectItems.Elements() {
					var objectItems, err = utils.TerraformListToGoArray(ctx, itemArray.(basetypes.ListValue), "object")
					if err != nil {
						return nil, err
					}
					properties[identifier] = objectItems
				}
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

	relations := writeRelationsToBody(ctx, state.Relations)
	e.Relations = relations
	return e, nil
}

func writeRelationsToBody(ctx context.Context, relations basetypes.MapValue) map[string]interface{} {
	relationsBody := make(map[string]interface{})
	for identifier, relation := range relations.Elements() {
		var items []tftypes.Value
		v, _ := relation.ToTerraformValue(ctx)
		v.As(&items)
		var relationsValue []string
		for _, item := range items {
			var v string
			item.As(&v)
			relationsValue = append(relationsValue, v)
		}
		relationsBody[identifier] = relationsValue
	}

	return relationsBody
}

func writeEntityComputedFieldsToState(state *EntityModel, e *cli.Entity) {
	state.Identifier = types.StringValue(e.Identifier)
	state.CreatedAt = types.StringValue(e.CreatedAt.String())
	state.CreatedBy = types.StringValue(e.CreatedBy)
	state.UpdatedAt = types.StringValue(e.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(e.UpdatedBy)
}

func (r *EntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *EntityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, state.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	e, err := entityResourceToBody(ctx, state, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	runID := ""
	if !state.RunID.IsNull() {
		runID = state.RunID.ValueString()
	}

	en, err := r.portClient.CreateEntity(ctx, e, runID)
	if err != nil {
		resp.Diagnostics.AddError("failed to create entity", err.Error())
		return
	}

	state.ID = types.StringValue(e.Identifier)

	writeEntityComputedFieldsToState(state, en)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntityModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteEntity(ctx, state.ID.ValueString(), state.Blueprint.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete entity", err.Error())
		return
	}

}

func (r *EntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <blueprint_id>:<entity_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), idParts[1])...)
}
