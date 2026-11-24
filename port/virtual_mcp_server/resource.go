package virtual_mcp_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &VirtualMcpServerResource{}
var _ resource.ResourceWithImportState = &VirtualMcpServerResource{}

func NewVirtualMcpServerResource() resource.Resource {
	return &VirtualMcpServerResource{}
}

type VirtualMcpServerResource struct {
	portClient *cli.PortClient
}

func (r *VirtualMcpServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_mcp_server"
}

func (r *VirtualMcpServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *VirtualMcpServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual MCP server under the `_virtual_mcp_server` system blueprint.\n\n" +
			"A virtual MCP server is a curated MCP endpoint that exposes selected Port and connector tools through " +
			"`virtual_mcp_tools` and `passthrough_tools`. The `url` attribute is the OAuth MCP endpoint agents should connect to.",
		Attributes: VirtualMcpServerSchema(),
	}
}

func (r *VirtualMcpServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state VirtualMcpServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entityBody, err := virtualMcpServerResourceToPortBody(&state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert virtual MCP server resource to body", err.Error())
		return
	}

	createdEntity, err := r.portClient.CreateEntity(ctx, entityBody, "", false)
	if err != nil {
		resp.Diagnostics.AddError("failed to create virtual MCP server", err.Error())
		return
	}

	identifier := createdEntity.Identifier
	if identifier == "" {
		identifier = state.Identifier.ValueString()
	}
	if err = r.refreshState(ctx, &state, identifier); err != nil {
		resp.Diagnostics.AddError("failed writing virtual MCP server fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMcpServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VirtualMcpServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, statusCode, err := r.portClient.ReadEntityIncludingCalculatedProperties(ctx, state.Identifier.ValueString(), VirtualMcpServerBlueprintIdentifier)
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read virtual MCP server", err.Error())
		return
	}

	if err = refreshVirtualMcpServerState(ctx, &state, entity); err != nil {
		resp.Diagnostics.AddError("failed writing virtual MCP server fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMcpServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state VirtualMcpServerModel
	var previousState VirtualMcpServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entityBody, err := virtualMcpServerResourceToPortBody(&state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert virtual MCP server resource to body", err.Error())
		return
	}

	_, err = r.portClient.UpdateEntity(
		ctx,
		previousState.Identifier.ValueString(),
		VirtualMcpServerBlueprintIdentifier,
		entityBody,
		"",
		false,
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to update virtual MCP server", err.Error())
		return
	}

	if err = r.refreshState(ctx, &state, previousState.Identifier.ValueString()); err != nil {
		resp.Diagnostics.AddError("failed writing virtual MCP server fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMcpServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VirtualMcpServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.portClient.DeleteEntity(ctx, state.Identifier.ValueString(), VirtualMcpServerBlueprintIdentifier); err != nil {
		resp.Diagnostics.AddError("failed to delete virtual MCP server", err.Error())
		return
	}
}

func (r *VirtualMcpServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), req.ID)...)
}

func (r *VirtualMcpServerResource) refreshState(ctx context.Context, state *VirtualMcpServerModel, identifier string) error {
	entity, _, err := r.portClient.ReadEntityIncludingCalculatedProperties(ctx, identifier, VirtualMcpServerBlueprintIdentifier)
	if err != nil {
		return err
	}

	return refreshVirtualMcpServerState(ctx, state, entity)
}