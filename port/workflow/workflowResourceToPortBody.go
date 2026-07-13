package workflow

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func workflowResourceToPortBody(ctx context.Context, state *WorkflowModel) (*cli.Workflow, error) {
	w := &cli.Workflow{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
	}

	if !state.Icon.IsNull() {
		w.Icon = state.Icon.ValueStringPointer()
	}
	if !state.Description.IsNull() {
		w.Description = state.Description.ValueStringPointer()
	}
	if !state.AllowAnyoneToViewRuns.IsNull() {
		w.AllowAnyoneToViewRuns = state.AllowAnyoneToViewRuns.ValueBoolPointer()
	}
	if !state.Tags.IsNull() && !state.Tags.IsUnknown() {
		tagValues, err := utils.TerraformListToGoArray(ctx, state.Tags, "string")
		if err != nil {
			return nil, err
		}
		tagArray := utils.InterfaceToStringArray(tagValues)
		w.Tags = &tagArray
	}

	nodes, err := utils.TerraformStringToGoType[[]map[string]any](state.Nodes)
	if err != nil {
		return nil, err
	}
	w.Nodes = nodes

	connections, err := utils.TerraformStringToGoType[[]map[string]any](state.Connections)
	if err != nil {
		return nil, err
	}
	w.Connections = connections

	return w, nil
}

func refreshWorkflowState(ctx context.Context, state *WorkflowModel, w *cli.Workflow) error {
	state.ID = types.StringValue(w.Identifier)
	state.Identifier = types.StringValue(w.Identifier)
	state.Title = types.StringValue(w.Title)
	state.Icon = flex.GoStringToFramework(w.Icon)
	state.Description = flex.GoStringToFramework(w.Description)
	state.AllowAnyoneToViewRuns = flex.GoBoolToFramework(w.AllowAnyoneToViewRuns)

	if w.CreatedAt != nil {
		state.CreatedAt = types.StringValue(w.CreatedAt.String())
	}
	if w.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(w.UpdatedAt.String())
	}
	state.CreatedBy = types.StringValue(w.CreatedBy)
	state.UpdatedBy = types.StringValue(w.UpdatedBy)

	var tags []string
	if w.Tags != nil {
		tags = *w.Tags
	}
	state.Tags = flex.GoArrayStringToTerraformList(ctx, tags)

	nodesJSON, err := json.Marshal(w.Nodes)
	if err != nil {
		return err
	}
	state.Nodes = types.StringValue(string(nodesJSON))

	connectionsJSON, err := json.Marshal(w.Connections)
	if err != nil {
		return err
	}
	state.Connections = types.StringValue(string(connectionsJSON))

	return nil
}
