package factory_line

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithValidateConfig = &FactoryLineResource{}

func (r *FactoryLineResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var state FactoryLineModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	seen := map[string]int{}
	for i, stage := range state.Stages {
		stagePath := path.Root("stage").AtListIndex(i)
		if stage.Identifier.IsUnknown() {
			continue
		}
		id := stage.Identifier.ValueString()
		if prior, dup := seen[id]; dup {
			resp.Diagnostics.AddAttributeError(
				stagePath.AtName("identifier"),
				"Duplicate stage identifier",
				fmt.Sprintf("Stage identifier %q is also used at index %d.", id, prior),
			)
		}
		seen[id] = i

		if stage.Executor.IsUnknown() {
			continue
		}
		executor := stage.Executor.ValueString()
		validateExecutorSettings(resp, stagePath, executor, stage)
	}
}

func validateExecutorSettings(resp *resource.ValidateConfigResponse, stagePath path.Path, executor string, stage StageModel) {
	has := map[string]bool{
		"factory_ai_managed_agent": stage.FactoryAI != nil,
		"claude_managed_agent":     stage.Claude != nil,
		"devin_managed_agent":      stage.Devin != nil,
		"custom_workflow":          stage.CustomWorkflow != nil,
		"human":                    stage.Human != nil,
	}

	set := 0
	for _, present := range has {
		if present {
			set++
		}
	}
	if executor != "none" && set != 1 {
		resp.Diagnostics.Append(exactlyOneExecutorBlock(stagePath, executor, set)...)
		return
	}
	if executor == "none" && set > 0 {
		resp.Diagnostics.AddAttributeError(
			stagePath,
			"Invalid executor settings",
			"Executor `none` must not include an executor settings block.",
		)
	}
	if executor == "factory_ai_managed_agent" && stage.FactoryAI != nil {
		validateFactoryAI(resp, stagePath.AtName("factory_ai_managed_agent"), *stage.FactoryAI)
	}
}

func validateFactoryAI(resp *resource.ValidateConfigResponse, blockPath path.Path, settings FactoryAIExecutorSettingsModel) {
	envSet := !settings.EnvironmentID.IsNull() && !settings.EnvironmentID.IsUnknown() && settings.EnvironmentID.ValueString() != ""
	compSet := !settings.ComputerID.IsNull() && !settings.ComputerID.IsUnknown() && settings.ComputerID.ValueString() != ""
	if envSet == compSet {
		resp.Diagnostics.AddAttributeError(
			blockPath,
			"Invalid Factory.ai target",
			"Set exactly one of `environment_id` or `computer_id`.",
		)
	}
}

func exactlyOneExecutorBlock(stagePath path.Path, executor string, count int) diag.Diagnostics {
	var diags diag.Diagnostics
	if count == 0 && executor != "none" {
		diags.AddAttributeError(
			stagePath,
			"Missing executor settings",
			fmt.Sprintf("Executor %q requires exactly one matching executor settings block.", executor),
		)
		return diags
	}
	if count > 1 {
		diags.AddAttributeError(
			stagePath,
			"Invalid executor settings",
			"At most one executor settings block may be set per stage.",
		)
	}
	return diags
}
