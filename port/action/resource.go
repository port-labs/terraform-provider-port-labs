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

func NewActionResource() resource.Resource {
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
	if a.Icon != nil {
		data.Icon = types.StringValue(*a.Icon)
	}
	if a.Description != nil {
		data.Description = types.StringValue(*a.Description)
	}

	if a.RequiredApproval != nil {
		data.RequiredApproval = types.BoolValue(*a.RequiredApproval)
	}

	if a.InvocationMethod.Type == "KAFKA" {
		data.KafkaMethod = types.MapNull(types.StringType)
	}

	if a.InvocationMethod.Type == "WEBHOOK" {
		data.WebhookMethod = &WebhookMethodModel{
			Url: types.StringValue(*a.InvocationMethod.Url),
		}
		if a.InvocationMethod.Agent != nil {
			data.WebhookMethod.Agent = types.BoolValue(*a.InvocationMethod.Agent)
		}
	}

	if a.InvocationMethod.Type == "GITHUB" {
		data.GithubMethod = &GithubMethodModel{
			Repo: types.StringValue(*a.InvocationMethod.Repo),
			Org:  types.StringValue(*a.InvocationMethod.Org),
		}

		if a.InvocationMethod.OmitPayload != nil {
			data.GithubMethod.OmitPayload = types.BoolValue(*a.InvocationMethod.OmitPayload)
		}

		if a.InvocationMethod.OmitUserInputs != nil {
			data.GithubMethod.OmitUserInputs = types.BoolValue(*a.InvocationMethod.OmitUserInputs)
		}

		if a.InvocationMethod.Workflow != nil {
			data.GithubMethod.Workflow = types.StringValue(*a.InvocationMethod.Workflow)
		}

		if a.InvocationMethod.Branch != nil {
			data.GithubMethod.Branch = types.StringValue(*a.InvocationMethod.Branch)
		}
	}

	if a.InvocationMethod.Type == "AZURE-DEVOPS" {
		data.AzureMethod = &AzureMethodModel{
			Org:     types.StringValue(*a.InvocationMethod.Org),
			Webhook: types.StringValue(*a.InvocationMethod.Webhook),
		}
	}

}

func (r *ActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.portClient.DeleteAction(ctx, data.Identifier.ValueString(), data.Blueprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete action", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
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

	a, err := r.portClient.UpdateAction(ctx, bp.Identifier, action.ID, action)
	if err != nil {
		resp.Diagnostics.AddError("failed to create action", err.Error())
		return
	}

	data.ID = types.StringValue(a.Identifier)
	data.Identifier = types.StringValue(a.Identifier)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func actionResourceToBody(ctx context.Context, data *ActionModel, bp *cli.Blueprint) (*cli.Action, error) {
	action := &cli.Action{
		Identifier: data.Identifier.ValueString(),
		Title:      data.Title.ValueString(),
	}

	if !data.Icon.IsNull() {
		icon := data.Icon.ValueString()
		action.Icon = &icon
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		action.Description = &description
	}

	action.InvocationMethod = invocationMethodToBody(data)

	return action, nil
}

func invocationMethodToBody(data *ActionModel) *cli.InvocationMethod {
	if data.AzureMethod != nil {
		org := data.AzureMethod.Org.ValueString()
		webhook := data.AzureMethod.Webhook.ValueString()
		return &cli.InvocationMethod{
			Type:    "AZURE-DEVOPS",
			Org:     &org,
			Webhook: &webhook,
		}
	}

	if data.GithubMethod != nil {
		org := data.GithubMethod.Org.ValueString()
		repo := data.GithubMethod.Repo.ValueString()
		githubInvocation := &cli.InvocationMethod{
			Type: "GITHUB",
			Org:  &org,
			Repo: &repo,
		}
		if !data.GithubMethod.Workflow.IsNull() {
			workflow := data.GithubMethod.Workflow.ValueString()
			githubInvocation.Workflow = &workflow
		}

		if !data.GithubMethod.OmitPayload.IsNull() {
			omitPayload := data.GithubMethod.OmitPayload.ValueBool()
			githubInvocation.OmitPayload = &omitPayload
		}

		if !data.GithubMethod.OmitUserInputs.IsNull() {
			omitUserInputs := data.GithubMethod.OmitUserInputs.ValueBool()
			githubInvocation.OmitUserInputs = &omitUserInputs
		}

		if !data.GithubMethod.ReportWorkflowStatus.IsNull() {
			reportWorkflowStatus := data.GithubMethod.ReportWorkflowStatus.ValueBool()
			githubInvocation.ReportWorkflowStatus = &reportWorkflowStatus
		}
		return githubInvocation
	}

	if !data.KafkaMethod.IsNull() {
		return &cli.InvocationMethod{
			Type: "KAFKA",
		}
	}

	if data.WebhookMethod != nil {
		url := data.WebhookMethod.Url.ValueString()
		webhookInvocation := &cli.InvocationMethod{
			Type: "WEBHOOK",
			Url:  &url,
		}
		if !data.WebhookMethod.Agent.IsNull() {
			agent := data.WebhookMethod.Agent.ValueBool()
			webhookInvocation.Agent = &agent
		}
		return webhookInvocation
	}
	return nil
}
