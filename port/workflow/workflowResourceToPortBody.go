package workflow

import (
	"context"
	"encoding/json"
	"fmt"

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
	if !state.Tags.IsNull() {
		tagValues, err := utils.TerraformListToGoArray(ctx, state.Tags, "string")
		if err != nil {
			return nil, err
		}
		w.Tags = make([]string, 0, len(tagValues))
		for _, tag := range tagValues {
			if tagString, ok := tag.(string); ok {
				w.Tags = append(w.Tags, tagString)
			}
		}
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

	if len(w.Tags) > 0 {
		tags, diags := types.ListValueFrom(ctx, types.StringType, w.Tags)
		if diags.HasError() {
			return fmt.Errorf("failed to convert workflow tags: %s", diags)
		}
		state.Tags = tags
	} else {
		state.Tags = types.ListNull(types.StringType)
	}

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
