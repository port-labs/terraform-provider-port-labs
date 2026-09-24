package mcp_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &McpServerResource{}
var _ resource.ResourceWithImportState = &McpServerResource{}

func NewMcpServerResource() resource.Resource {
	return &McpServerResource{}
}

type McpServerResource struct {
	portClient *cli.PortClient
}

func (r *McpServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mcp_server"
}

func (r *McpServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *McpServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *McpServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entityBody, err := mcpServerResourceToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert MCP server resource to body", err.Error())
		return
	}

	createdEntity, err := r.portClient.CreateEntity(ctx, entityBody, "", false)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MCP server", err.Error())
		return
	}

	err = refreshMcpServerState(ctx, state, createdEntity)
	if err != nil {
		resp.Diagnostics.AddError("failed writing MCP server fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *McpServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *McpServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entity, statusCode, err := r.portClient.ReadEntity(ctx, state.Identifier.ValueString(), McpServerBlueprintIdentifier)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read MCP server", err.Error())
		return
	}

	err = refreshMcpServerState(ctx, state, entity)
	if err != nil {
		resp.Diagnostics.AddError("failed writing MCP server fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *McpServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *McpServerModel
	var previousState *McpServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entityBody, err := mcpServerResourceToPortBody(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert MCP server resource to body", err.Error())
		return
	}

	updatedEntity, err := r.portClient.UpdateEntity(
		ctx,
		previousState.Identifier.ValueString(),
		McpServerBlueprintIdentifier,
		entityBody,
		"",
		false,
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to update MCP server", err.Error())
		return
	}

	err = refreshMcpServerState(ctx, state, updatedEntity)
	if err != nil {
		resp.Diagnostics.AddError("failed writing MCP server fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *McpServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *McpServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteEntity(ctx, state.Identifier.ValueString(), McpServerBlueprintIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete MCP server", err.Error())
		return
	}
}

func (r *McpServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
}
