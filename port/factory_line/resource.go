package factory_line

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var (
	_ resource.Resource                = &FactoryLineResource{}
	_ resource.ResourceWithImportState = &FactoryLineResource{}
	_ resource.ResourceWithConfigure   = &FactoryLineResource{}
)

func NewFactoryLineResource() resource.Resource {
	return &FactoryLineResource{}
}

type FactoryLineResource struct {
	portClient *cli.PortClient
}

func (r *FactoryLineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_factory_line"
}

func (r *FactoryLineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *FactoryLineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state FactoryLineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	line, err := modelToCreateLine(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError("failed to build factory line", err.Error())
		return
	}

	created, err := r.portClient.CreateFactoryLine(ctx, line)
	if err != nil {
		resp.Diagnostics.AddError("failed to create factory line", err.Error())
		return
	}

	resp.Diagnostics.Append(refreshLineToState(ctx, &state, created)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FactoryLineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FactoryLineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	line, statusCode, err := r.portClient.ReadFactoryLine(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read factory line", err.Error())
		return
	}

	resp.Diagnostics.Append(refreshLineToState(ctx, &state, line)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FactoryLineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state FactoryLineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	change, err := modelToLineChange(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError("failed to build factory line update", err.Error())
		return
	}

	updated, err := r.portClient.UpdateFactoryLine(ctx, state.Identifier.ValueString(), change)
	if err != nil {
		resp.Diagnostics.AddError("failed to update factory line", err.Error())
		return
	}

	resp.Diagnostics.Append(refreshLineToState(ctx, &state, updated)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FactoryLineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FactoryLineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.portClient.DeleteFactoryLine(ctx, state.Identifier.ValueString()); err != nil {
		resp.Diagnostics.AddError("failed to delete factory line", err.Error())
	}
}

func (r *FactoryLineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("identifier"), req, resp)
}
