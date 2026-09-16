package scorecard_group

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &ScorecardGroupScorecardResource{}
var _ resource.ResourceWithImportState = &ScorecardGroupScorecardResource{}

func NewScorecardGroupScorecardResource() resource.Resource {
	return &ScorecardGroupScorecardResource{}
}

type ScorecardGroupScorecardResource struct {
	portClient *cli.PortClient
}

type ScorecardGroupScorecardModel struct {
	ID                  types.String `tfsdk:"id"`
	GroupIdentifier     types.String `tfsdk:"group_identifier"`
	ScorecardIdentifier types.String `tfsdk:"scorecard_identifier"`
	OverrideGroup       types.Bool   `tfsdk:"override_group"`
}

func (r *ScorecardGroupScorecardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard_group_scorecard"
}

func (r *ScorecardGroupScorecardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *ScorecardGroupScorecardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Associates an existing scorecard with a scorecard group using the Port `/v1/scorecard-groups/{group}/scorecards/{scorecard}` API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_identifier": schema.StringAttribute{
				MarkdownDescription: "The identifier of the scorecard group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scorecard_identifier": schema.StringAttribute{
				MarkdownDescription: "The identifier of the scorecard to add to the group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"override_group": schema.BoolAttribute{
				MarkdownDescription: "When true, moves the scorecard to this group even if it already belongs to another group.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func membershipID(groupIdentifier, scorecardIdentifier string) string {
	return fmt.Sprintf("%s/%s", groupIdentifier, scorecardIdentifier)
}

func (r *ScorecardGroupScorecardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ScorecardGroupScorecardModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.portClient.AddScorecardToGroup(
		ctx,
		state.GroupIdentifier.ValueString(),
		state.ScorecardIdentifier.ValueString(),
		state.OverrideGroup.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to add scorecard to group", err.Error())
		return
	}

	state.ID = types.StringValue(membershipID(state.GroupIdentifier.ValueString(), state.ScorecardIdentifier.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ScorecardGroupScorecardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScorecardGroupScorecardModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, statusCode, err := r.portClient.ReadScorecardGroup(ctx, state.GroupIdentifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read scorecard group", err.Error())
		return
	}

	state.ID = types.StringValue(membershipID(state.GroupIdentifier.ValueString(), state.ScorecardIdentifier.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ScorecardGroupScorecardResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"update not supported",
		"Changing group membership requires replacing the resource (destroy and recreate).",
	)
}

func (r *ScorecardGroupScorecardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScorecardGroupScorecardModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.portClient.RemoveScorecardFromGroup(
		ctx,
		state.GroupIdentifier.ValueString(),
		state.ScorecardIdentifier.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("failed to remove scorecard from group", err.Error())
		return
	}
}

func (r *ScorecardGroupScorecardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"invalid import ID",
			"Expected import ID in the form `<group_identifier>/<scorecard_identifier>`",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_identifier"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scorecard_identifier"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("override_group"), false)...)
}
