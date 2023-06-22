package entity

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
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
	var data *EntityModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := data.Blueprint.ValueString()
	e, statusCode, err := r.portClient.ReadEntity(ctx, data.Identifier.ValueString(), data.Blueprint.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading entity", err.Error())
		return
	}

	err = writeEntityFieldsToResource(ctx, data, e, blueprintIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed writing entity fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeEntityFieldsToResource(ctx context.Context, em *EntityModel, e *cli.Entity, blueprint string) error {
	em.ID = types.StringValue(e.Identifier)
	em.Identifier = types.StringValue(e.Identifier)
	em.Blueprint = types.StringValue(blueprint)
	em.Title = types.StringValue(e.Title)
	em.CreatedAt = types.StringValue(e.CreatedAt.String())
	em.CreatedBy = types.StringValue(e.CreatedBy)
	em.UpdatedAt = types.StringValue(e.UpdatedAt.String())
	em.UpdatedBy = types.StringValue(e.UpdatedBy)

	if len(e.Team) != 0 {
		em.Teams = make([]types.String, len(e.Team))
		for i, t := range e.Team {
			em.Teams[i] = types.StringValue(t)
		}
	}

	if len(e.Properties) != 0 {
		em.Properties = &EntityPropertiesModel{}
		for k, v := range e.Properties {
			switch t := v.(type) {
			// case map[string]interface{}:
			// 	js, _ := json.Marshal(&t)
			// 	propValue = string(js)
			// case []interface{}:
			// 	propValue = t
			// case float64:
			// 	propValue = strconv.FormatFloat(t, 'f', -1, 64)
			case int:
				if em.Properties.NumberProp == nil {
					em.Properties.NumberProp = make(map[string]float64)
				}
				em.Properties.NumberProp[k] = float64(t)
			case string:
				if em.Properties.StringProp == nil {
					em.Properties.StringProp = make(map[string]string)
				}
				em.Properties.StringProp[k] = t
			}
		}
	}

	return nil
}
func (r *EntityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *EntityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	e, err := entityResourceToBody(data, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	runID := ""
	if !data.RunID.IsNull() {
		runID = data.RunID.ValueString()
	}

	en, err := r.portClient.CreateEntity(ctx, e, runID)
	if err != nil {
		resp.Diagnostics.AddError("failed to create entity", err.Error())
		return
	}

	data.ID = types.StringValue(e.Identifier)

	writeEntityComputedFieldsToResource(data, en)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func entityResourceToBody(em *EntityModel, bp *cli.Blueprint) (*cli.Entity, error) {
	e := &cli.Entity{
		Identifier: em.Identifier.ValueString(),
		Blueprint:  bp.Identifier,
	}

	if !em.Title.IsNull() {
		e.Title = em.Title.ValueString()
	}

	if em.Teams != nil {
		e.Team = make([]string, len(em.Teams))
		for i, t := range em.Teams {
			e.Team[i] = t.ValueString()
		}
	}

	properties := make(map[string]interface{})
	if em.Properties != nil {
		if em.Properties.StringProp != nil {
			for propIdentifier, prop := range em.Properties.StringProp {
				properties[propIdentifier] = prop
			}
		}
		if em.Properties.NumberProp != nil {
			for propIdentifier, prop := range em.Properties.NumberProp {
				properties[propIdentifier] = prop
			}
		}
		e.Properties = properties
	}

	relations := make(map[string]interface{})
	e.Relations = relations
	return e, nil
}

func writeEntityComputedFieldsToResource(data *EntityModel, e *cli.Entity) {
	data.Identifier = types.StringValue(e.Identifier)
	data.CreatedAt = types.StringValue(e.CreatedAt.String())
	data.CreatedBy = types.StringValue(e.CreatedBy)
	data.UpdatedAt = types.StringValue(e.UpdatedAt.String())
	data.UpdatedBy = types.StringValue(e.UpdatedBy)
}

func (r *EntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *EntityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	e, err := entityResourceToBody(data, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	runID := ""
	if !data.RunID.IsNull() {
		runID = data.RunID.ValueString()
	}

	en, err := r.portClient.CreateEntity(ctx, e, runID)
	if err != nil {
		resp.Diagnostics.AddError("failed to create entity", err.Error())
		return
	}

	data.ID = types.StringValue(e.Identifier)

	writeEntityComputedFieldsToResource(data, en)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *EntityModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteEntity(ctx, data.ID.ValueString(), data.Blueprint.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete entity", err.Error())
		return
	}

}

func (r *EntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <entity_id>:<blueprint_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), idParts[1])...)
}
