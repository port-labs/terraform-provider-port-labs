package llm_provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &LLMProviderResource{}
var _ resource.ResourceWithImportState = &LLMProviderResource{}

type LLMProviderResource struct {
	portClient *cli.PortClient
}

func NewLLMProviderResource() resource.Resource {
	return &LLMProviderResource{}
}

func (r *LLMProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_llm_provider"
}

func (r *LLMProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *LLMProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *LLMProviderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body, err := llmProviderToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert llm provider resource to body", err.Error())
		return
	}

	validateConnection := !state.ValidateConnection.IsNull() && !state.ValidateConnection.IsUnknown() && state.ValidateConnection.ValueBool()
	provider, err := r.portClient.UpsertLLMProvider(ctx, body, validateConnection)
	if err != nil {
		resp.Diagnostics.AddError("failed to create llm provider", err.Error())
		return
	}

	if err := refreshLLMProviderState(ctx, state, provider); err != nil {
		resp.Diagnostics.AddError("failed writing llm provider fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LLMProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *LLMProviderModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	providerName := state.ProviderType.ValueString()
	provider, statusCode, err := r.portClient.ReadLLMProvider(ctx, providerName)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read llm provider", err.Error())
		return
	}

	if err := refreshLLMProviderState(ctx, state, provider); err != nil {
		resp.Diagnostics.AddError("failed writing llm provider fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LLMProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *LLMProviderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body, err := llmProviderToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert llm provider resource to body", err.Error())
		return
	}

	validateConnection := !state.ValidateConnection.IsNull() && !state.ValidateConnection.IsUnknown() && state.ValidateConnection.ValueBool()
	provider, err := r.portClient.UpsertLLMProvider(ctx, body, validateConnection)
	if err != nil {
		resp.Diagnostics.AddError("failed to update llm provider", err.Error())
		return
	}

	if err := refreshLLMProviderState(ctx, state, provider); err != nil {
		resp.Diagnostics.AddError("failed writing llm provider fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LLMProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *LLMProviderModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.portClient.DeleteLLMProvider(ctx, state.ProviderType.ValueString()); err != nil {
		resp.Diagnostics.AddError("failed to delete llm provider", err.Error())
		return
	}
}

func (r *LLMProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("provider_type"), req.ID,
	)...)
}
