package factory_line

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func refreshLineToState(ctx context.Context, state *FactoryLineModel, line *cli.FactoryLineResponse) diag.Diagnostics {
	state.ID = types.StringValue(line.Identifier)
	state.Identifier = types.StringValue(line.Identifier)
	state.Title = types.StringValue(line.Title)
	state.Status = types.StringValue(line.Status)
	state.MaxRunning = types.Int64Value(int64(line.MaxRunning))

	if line.Description != nil {
		state.Description = types.StringValue(*line.Description)
	} else {
		state.Description = types.StringNull()
	}
	if line.StatusReason != nil {
		state.StatusReason = types.StringValue(*line.StatusReason)
	} else {
		state.StatusReason = types.StringNull()
	}
	if line.CreatedAt != nil {
		state.CreatedAt = types.StringValue(*line.CreatedAt)
	}
	if line.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(*line.UpdatedAt)
	}

	team, diags := types.ListValueFrom(ctx, types.StringType, line.Team)
	if diags.HasError() {
		return diags
	}
	state.Team = team

	order, diags := types.ListValueFrom(ctx, types.StringType, line.StageOrder)
	if diags.HasError() {
		return diags
	}
	state.StageOrder = order

	stages := make([]StageModel, 0, len(line.Stages))
	for _, stage := range line.Stages {
		m, diags := stageToModel(ctx, stage)
		if diags.HasError() {
			return diags
		}
		stages = append(stages, m)
	}
	state.Stages = stages
	return nil
}

func stageToModel(ctx context.Context, stage cli.FactoryLineStage) (StageModel, diag.Diagnostics) {
	m := StageModel{
		Identifier:    types.StringValue(stage.Identifier),
		Executor:      types.StringValue(stage.Executor),
		MaxTries:      types.Int64Value(int64(stage.MaxTries)),
		AllowTakeOver: types.BoolValue(stage.AllowTakeOver),
	}
	if stage.Title != nil {
		m.Title = types.StringValue(*stage.Title)
	}
	if stage.Description != nil {
		m.Description = types.StringValue(*stage.Description)
	}
	if stage.Instructions != nil {
		m.Instructions = types.StringValue(*stage.Instructions)
	}
	if stage.Timeout != nil {
		m.Timeout = types.StringValue(*stage.Timeout)
	}
	if len(stage.Team) > 0 {
		team, diags := types.ListValueFrom(ctx, types.StringType, stage.Team)
		if diags.HasError() {
			return m, diags
		}
		m.Team = team
	}
	if len(stage.Checkpoints) > 0 {
		cp, err := utils.GoObjectToTerraformString(json.RawMessage(stage.Checkpoints), true)
		if err != nil {
			var diags diag.Diagnostics
			diags.AddError("failed to map checkpoints", err.Error())
			return m, diags
		}
		m.Checkpoints = cp
	} else {
		m.Checkpoints = types.StringValue("[]")
	}
	if err := applyExecutorSettingsToModel(&m, stage.Executor, stage.ExecutorSettings); err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to map executor settings", err.Error())
		return m, diags
	}
	return m, nil
}

func applyExecutorSettingsToModel(stage *StageModel, executor string, raw json.RawMessage) error {
	if len(raw) == 0 {
		raw = json.RawMessage("{}")
	}
	switch executor {
	case "factory_ai_managed_agent":
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			return err
		}
		settings := &FactoryAIExecutorSettingsModel{
			InstallationID:  stringField(body, "installation_id"),
			EnvironmentID:   stringField(body, "environment_id"),
			ComputerID:      stringField(body, "computer_id"),
			Model:           stringField(body, "model"),
			ReasoningEffort: stringField(body, "reasoning_effort"),
		}
		stage.FactoryAI = settings
	case "claude_managed_agent":
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			return err
		}
		settings := &ClaudeExecutorSettingsModel{
			AgentID:        stringField(body, "agent_id"),
			EnvironmentID:  stringField(body, "environment_id"),
			InstallationID: stringField(body, "installation_id"),
			GitAccessToken: stringField(body, "git_access_token"),
		}
		if session, ok := body["session"]; ok {
			s, err := utils.GoObjectToTerraformString(session, true)
			if err != nil {
				return err
			}
			settings.Session = s
		}
		if identity, ok := body["execution_identity"]; ok {
			s, err := utils.GoObjectToTerraformString(identity, true)
			if err != nil {
				return err
			}
			settings.ExecutionIdentity = s
		}
		stage.Claude = settings
	case "devin_managed_agent":
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			return err
		}
		settings := &DevinExecutorSettingsModel{
			InstallationID: stringField(body, "installation_id"),
			PlaybookID:     stringField(body, "playbook_id"),
		}
		if repos, ok := body["repos"].([]any); ok {
			repoStrings := make([]string, 0, len(repos))
			for _, r := range repos {
				if s, ok := r.(string); ok {
					repoStrings = append(repoStrings, s)
				}
			}
			list, diags := types.ListValueFrom(context.Background(), types.StringType, repoStrings)
			if diags.HasError() {
				return fmt.Errorf("repos: %s", diags.Errors()[0].Summary())
			}
			settings.Repos = list
		}
		if session, ok := body["session"]; ok {
			s, err := utils.GoObjectToTerraformString(session, true)
			if err != nil {
				return err
			}
			settings.Session = s
		}
		if identity, ok := body["execution_identity"]; ok {
			s, err := utils.GoObjectToTerraformString(identity, true)
			if err != nil {
				return err
			}
			settings.ExecutionIdentity = s
		}
		stage.Devin = settings
	case "custom_workflow":
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			return err
		}
		stage.CustomWorkflow = &CustomWorkflowExecutorSettingsModel{
			Workflow: stringField(body, "workflow"),
		}
	case "human":
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			return err
		}
		settings := &HumanExecutorSettingsModel{}
		if assignees, ok := body["assignees"]; ok {
			s, err := utils.GoObjectToTerraformString(assignees, true)
			if err != nil {
				return err
			}
			settings.Assignees = s
		}
		if notify, ok := body["notify"]; ok {
			s, err := utils.GoObjectToTerraformString(notify, true)
			if err != nil {
				return err
			}
			settings.Notify = s
		}
		stage.Human = settings
	}
	return nil
}

func stringField(body map[string]any, key string) types.String {
	v, ok := body[key]
	if !ok || v == nil {
		return types.StringNull()
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
