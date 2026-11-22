package workflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &WorkflowResource{}
var _ resource.ResourceWithImportState = &WorkflowResource{}

func NewWorkflowResource() resource.Resource {
	return &WorkflowResource{}
}

type WorkflowResource struct {
	portClient *cli.PortClient
}

func (r *WorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (r *WorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *WorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), req.ID)...)
}

func (r *WorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *WorkflowModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	w, statusCode, err := r.portClient.ReadWorkflow(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading workflow", err.Error())
		return
	}

	err = r.refreshWorkflowState(ctx, state, w)
	if err != nil {
		resp.Diagnostics.AddError("failed writing workflow fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *WorkflowModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	workflow, err := workflowStateToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert workflow resource to body", err.Error())
		return
	}

	w, err := r.portClient.CreateWorkflow(ctx, workflow)
	if err != nil {
		resp.Diagnostics.AddError("failed to create workflow", err.Error())
		return
	}

	if err := r.refreshWorkflowState(ctx, state, w); err != nil {
		resp.Diagnostics.AddError("failed writing workflow fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *WorkflowModel
	var previousState *WorkflowModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	workflow, err := workflowStateToPortBody(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert workflow resource to body", err.Error())
		return
	}

	w, err := r.portClient.UpdateWorkflow(ctx, previousState.Identifier.ValueString(), workflow)
	if err != nil {
		resp.Diagnostics.AddError("failed to update workflow", err.Error())
		return
	}

	if err := r.refreshWorkflowState(ctx, state, w); err != nil {
		resp.Diagnostics.AddError("failed writing workflow fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *WorkflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteWorkflow(ctx, state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete workflow", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
