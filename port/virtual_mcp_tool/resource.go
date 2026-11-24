package virtual_mcp_tool

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &VirtualMcpToolResource{}
var _ resource.ResourceWithImportState = &VirtualMcpToolResource{}

func NewVirtualMcpToolResource() resource.Resource {
	return &VirtualMcpToolResource{}
}

type VirtualMcpToolResource struct {
	portClient *cli.PortClient
}

func (r *VirtualMcpToolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_mcp_tool"
}

func (r *VirtualMcpToolResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *VirtualMcpToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual MCP tool under the `_virtual_mcp_tool` system blueprint.\n\n" +
			"Virtual tools customize a single underlying Port or connector tool with a fixed name, description, and argument rules. " +
			"Attach them to a `port_virtual_mcp_server` through the server's `virtual_mcp_tools` relation.",
		Attributes: VirtualMcpToolSchema(),
	}
}

func (r *VirtualMcpToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state VirtualMcpToolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entityBody, err := virtualMcpToolResourceToPortBody(&state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert virtual MCP tool resource to body", err.Error())
		return
	}

	createdEntity, err := r.portClient.CreateEntity(ctx, entityBody, "", false)
	if err != nil {
		resp.Diagnostics.AddError("failed to create virtual MCP tool", err.Error())
		return
	}

	if err = refreshVirtualMcpToolState(ctx, &state, createdEntity); err != nil {
		resp.Diagnostics.AddError("failed writing virtual MCP tool fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMcpToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VirtualMcpToolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, statusCode, err := r.portClient.ReadEntity(ctx, state.Identifier.ValueString(), VirtualMcpToolBlueprintIdentifier)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read virtual MCP tool", err.Error())
		return
	}

	if err = refreshVirtualMcpToolState(ctx, &state, entity); err != nil {
		resp.Diagnostics.AddError("failed writing virtual MCP tool fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMcpToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state VirtualMcpToolModel
	var previousState VirtualMcpToolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entityBody, err := virtualMcpToolResourceToPortBody(&state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert virtual MCP tool resource to body", err.Error())
		return
	}

	updatedEntity, err := r.portClient.UpdateEntity(
		ctx,
		previousState.Identifier.ValueString(),
		VirtualMcpToolBlueprintIdentifier,
		entityBody,
		"",
		false,
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to update virtual MCP tool", err.Error())
		return
	}

	if err = refreshVirtualMcpToolState(ctx, &state, updatedEntity); err != nil {
		resp.Diagnostics.AddError("failed writing virtual MCP tool fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMcpToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VirtualMcpToolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.portClient.DeleteEntity(ctx, state.Identifier.ValueString(), VirtualMcpToolBlueprintIdentifier); err != nil {
		resp.Diagnostics.AddError("failed to delete virtual MCP tool", err.Error())
		return
	}
}

func (r *VirtualMcpToolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), req.ID)...)
}