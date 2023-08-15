package webhook

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &WebhookResource{}
var _ resource.ResourceWithImportState = &WebhookResource{}

func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

type WebhookResource struct {
	portClient *cli.PortClient
}

func (r *WebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *WebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *WebhookModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	identifier := state.Identifier.ValueString()
	w, statusCode, err := r.portClient.ReadWebhook(ctx, identifier)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read webhook", err.Error())
		return
	}

	err = refreshWebhookState(ctx, state, w)
	if err != nil {
		resp.Diagnostics.AddError("failed writing webhook fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *WebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	w, err := webhookResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert webhook resource to body", err.Error())
		return
	}

	wp, err := r.portClient.CreateWebhook(ctx, w)
	if err != nil {
		resp.Diagnostics.AddError("failed to create webhook", err.Error())
		return
	}

	writeWebhookComputedFieldsToState(state, wp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func writeWebhookComputedFieldsToState(state *WebhookModel, wp *cli.Webhook) {
	state.ID = types.StringValue(wp.Identifier)
	state.Identifier = types.StringValue(wp.Identifier)
	state.CreatedAt = types.StringValue(wp.CreatedAt.String())
	state.CreatedBy = types.StringValue(wp.CreatedBy)
	state.UpdatedAt = types.StringValue(wp.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(wp.UpdatedBy)
	state.Url = types.StringValue(wp.Url)
	state.WebhookKey = types.StringValue(wp.WebhookKey)
}

func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *WebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	w, err := webhookResourceToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert webhook resource to body", err.Error())
		return
	}

	wp, err := r.portClient.UpdateWebhook(ctx, w.Identifier, w)
	if err != nil {
		resp.Diagnostics.AddError("failed to update the webhook", err.Error())
		return
	}

	writeWebhookComputedFieldsToState(state, wp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *WebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteWebhook(ctx, state.Identifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete webhook", err.Error())
		return
	}

}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
}
