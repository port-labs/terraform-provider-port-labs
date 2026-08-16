package scorecard_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &ScorecardGroupResource{}
var _ resource.ResourceWithImportState = &ScorecardGroupResource{}

func NewScorecardGroupResource() resource.Resource {
	return &ScorecardGroupResource{}
}

type ScorecardGroupResource struct {
	portClient *cli.PortClient
}

func (r *ScorecardGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard_group"
}

func (r *ScorecardGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *ScorecardGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *ScorecardGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := scorecardGroupResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert scorecard group resource to body", err.Error())
		return
	}

	createdGroup, err := r.portClient.CreateScorecardGroup(ctx, group)
	if err != nil {
		resp.Diagnostics.AddError("failed to create scorecard group", err.Error())
		return
	}

	r.refreshScorecardGroupState(ctx, state, createdGroup)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ScorecardGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ScorecardGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, statusCode, err := r.portClient.ReadScorecardGroup(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read scorecard group", err.Error())
		return
	}

	r.refreshScorecardGroupState(ctx, state, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ScorecardGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *ScorecardGroupModel
	var previousState *ScorecardGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := scorecardGroupResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert scorecard group resource to body", err.Error())
		return
	}

	updatedGroup, err := r.portClient.UpdateScorecardGroup(ctx, previousState.Identifier.ValueString(), group)
	if err != nil {
		resp.Diagnostics.AddError("failed to update scorecard group", err.Error())
		return
	}

	r.refreshScorecardGroupState(ctx, state, updatedGroup)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ScorecardGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ScorecardGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteScorecardGroup(ctx, state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete scorecard group", err.Error())
		return
	}
}

func (r *ScorecardGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), req.ID)...)
}
