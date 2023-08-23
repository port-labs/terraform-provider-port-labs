package scorecard

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ScorecardResource{}
var _ resource.ResourceWithImportState = &ScorecardResource{}

func NewScorecardResource() resource.Resource {
	return &ScorecardResource{}
}

type ScorecardResource struct {
	portClient *cli.PortClient
}

func (r *ScorecardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard"
}

func (r *ScorecardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *ScorecardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ScorecardModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	identifier := state.Identifier.ValueString()
	blueprintIdentifier := state.Blueprint.ValueString()
	s, statusCode, err := r.portClient.ReadScorecard(ctx, blueprintIdentifier, identifier)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read scorecard", err.Error())
		return
	}

	refreshScorecardState(ctx, state, s, blueprintIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ScorecardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *ScorecardModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := scorecardResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert scorecard resource to body", err.Error())
		return
	}

	sp, err := r.portClient.CreateScorecard(ctx, state.Blueprint.ValueString(), s)
	if err != nil {
		resp.Diagnostics.AddError("failed to create scorecard", err.Error())
		return
	}

	writeScorecardComputedFieldsToState(state, sp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func writeScorecardComputedFieldsToState(state *ScorecardModel, wp *cli.Scorecard) {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", wp.Blueprint, wp.Identifier))
	state.Identifier = types.StringValue(wp.Identifier)
	state.CreatedAt = types.StringValue(wp.CreatedAt.String())
	state.CreatedBy = types.StringValue(wp.CreatedBy)
	state.UpdatedAt = types.StringValue(wp.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(wp.UpdatedBy)
}

func (r *ScorecardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *ScorecardModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := scorecardResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert scorecard resource to body", err.Error())
		return
	}

	sp, err := r.portClient.UpdateScorecard(ctx, state.Blueprint.ValueString(), state.Identifier.ValueString(), s)
	if err != nil {
		resp.Diagnostics.AddError("failed to update the scorecard", err.Error())
		return
	}

	writeScorecardComputedFieldsToState(state, sp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ScorecardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ScorecardModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteScorecard(ctx, state.Blueprint.ValueString(), state.Identifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete scorecard", err.Error())
		return
	}

}

func (r *ScorecardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <blueprint_id>:<scorecard_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), idParts[1])...)
}
