package page

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

var _ resource.Resource = &PageResource{}
var _ resource.ResourceWithImportState = &PageResource{}

func NewPageResource() resource.Resource {
	return &PageResource{}
}

type PageResource struct {
	portClient *cli.PortClient
}

func (r *PageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_page"
}

func (r *PageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *PageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}

func (r *PageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *PageModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p, statusCode, err := r.portClient.GetPage(ctx, state.Identifier.ValueString())

	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to get page", err.Error())
		return
	}

	err = refreshPageToState(state, p)

	if err != nil {
		resp.Diagnostics.AddError("failed to write page fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *PageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *PageModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Identifier.ValueString() == "$home" {
		tflog.Debug(ctx, "$home page is not deletable, unregistering from state")
		resp.State.RemoveResource(ctx)
		return
	}
	statusCode, err := r.portClient.DeletePage(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to delete page", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *PageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *PageModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	page, err := PageToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert page resource to body", err.Error())
		return
	}

	p, err := r.portClient.CreatePage(ctx, page)
	if err != nil {
		resp.Diagnostics.AddError("failed to create page", err.Error())
		return
	}
	if p == nil {
		// if page is nil and err is nil this means that the page got created but the response body was empty
		// to be forward compatible we will query the page again
		p, _, err = r.portClient.GetPage(ctx, state.Identifier.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("failed to get page", err.Error())
			return
		}
	}

	state.ID = types.StringValue(p.Identifier)
	state.CreatedAt = types.StringValue(p.CreatedAt.String())
	state.CreatedBy = types.StringValue(p.CreatedBy)
	state.UpdatedAt = types.StringValue(p.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(p.UpdatedBy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *PageModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p, _, err := r.portClient.GetPage(ctx, state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get page", err.Error())
		return
	}

	page, err := PageToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert page resource to body", err.Error())
		return
	}

	_, err = r.portClient.UpdatePage(ctx, p.Identifier, page)

	if err != nil {
		resp.Diagnostics.AddError("failed to update page", err.Error())
		return
	}

	state.ID = types.StringValue(p.Identifier)
	state.CreatedAt = types.StringValue(p.CreatedAt.String())
	state.CreatedBy = types.StringValue(p.CreatedBy)
	state.UpdatedAt = types.StringValue(p.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(p.UpdatedBy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
