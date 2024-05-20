package team

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &TeamResource{}
var _ resource.ResourceWithImportState = &TeamResource{}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

type TeamResource struct {
	portClient *cli.PortClient
}

func (r *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *TeamModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	t, statusCode, err := r.portClient.ReadTeam(ctx, name)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read team", err.Error())
		return
	}

	err = refreshTeamState(ctx, state, t)
	if err != nil {
		resp.Diagnostics.AddError("failed writing team fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *TeamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	t, err := TeamResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert team resource to body", err.Error())
		return
	}

	tp, err := r.portClient.CreateTeam(ctx, t)
	if err != nil {
		resp.Diagnostics.AddError("failed to create team", err.Error())
		return
	}

	writeTeamComputedFieldsToState(state, tp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func writeTeamComputedFieldsToState(state *TeamModel, tp *cli.Team) {
	state.ID = types.StringValue(tp.Name)
	state.CreatedAt = types.StringValue(tp.CreatedAt.String())
	state.UpdatedAt = types.StringValue(tp.UpdatedAt.String())
	state.ProviderName = types.StringValue(tp.Provider)
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *TeamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	t, err := TeamResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert team resource to body", err.Error())
		return
	}

	tp, err := r.portClient.UpdateTeam(ctx, t.Name, t)
	if err != nil {
		resp.Diagnostics.AddError("failed to update the team", err.Error())
		return
	}

	writeTeamComputedFieldsToState(state, tp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *TeamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteTeam(ctx, state.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete team", err.Error())
		return
	}

}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("name"), req.ID,
	)...)
}
