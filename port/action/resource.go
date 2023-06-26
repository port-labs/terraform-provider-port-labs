package action

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/port/cli"
)

var _ resource.Resource = &ActionResource{}
var _ resource.ResourceWithImportState = &ActionResource{}

func NewEntityResource() resource.Resource {
	return &ActionResource{}
}

type ActionResource struct {
	portClient *cli.PortClient
}

func (r *ActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entity"
}

func (r *ActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *ActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <entity_id>:<blueprint_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), idParts[1])...)
}

func (r *ActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ActionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := data.Blueprint.ValueString()
	a, statusCode, err := r.portClient.ReadAction(ctx, data.Identifier.ValueString(), data.Blueprint.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading action", err.Error())
		return
	}

	writeActionFieldsToResource(ctx, data, a, blueprintIdentifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeActionFieldsToResource(ctx context.Context, data *ActionModel, a *cli.Action, blueprintIdentifier string) {
	data.Title = types.StringValue(a.Title)
	data.Icon = types.StringValue(a.Icon)
	data.Description = types.StringValue(a.Description)
}

func (r *ActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
func (r *ActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bp, _, err := r.portClient.ReadBlueprint(ctx, data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read blueprint", err.Error())
		return
	}

	action, err := actionResourceToBody(ctx, data, bp)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert entity resource to body", err.Error())
		return
	}

	a, err := r.portClient.CreateAction(ctx, bp.Identifier, action)
	if err != nil {
		resp.Diagnostics.AddError("failed to create action", err.Error())
		return
	}

	data.ID = types.StringValue(a.Identifier)
	data.Identifier = types.StringValue(a.Identifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func actionResourceToBody(ctx context.Context, data *ActionModel, bp *cli.Blueprint) (*cli.Action, error) {
	action := &cli.Action{
		Identifier:  data.Identifier.ValueString(),
		Title:       data.Title.ValueString(),
		Icon:        data.Icon.ValueString(),
		Description: data.Description.ValueString(),
	}

	return action, nil
}
