package factory_line

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func modelToCreateLine(ctx context.Context, state *FactoryLineModel) (*cli.FactoryLine, error) {
	stages, err := stagesFromModel(state.Stages)
	if err != nil {
		return nil, err
	}

	line := &cli.FactoryLine{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
		Stages:     stages,
	}

	if !state.Description.IsNull() && !state.Description.IsUnknown() {
		desc := state.Description.ValueString()
		line.Description = &desc
	}
	if !state.Team.IsNull() && !state.Team.IsUnknown() {
		team, err := stringListFromTerraform(ctx, state.Team)
		if err != nil {
			return nil, err
		}
		line.Team = team
	}
	if !state.MaxRunning.IsNull() && !state.MaxRunning.IsUnknown() {
		max := int(state.MaxRunning.ValueInt64())
		line.MaxRunning = &max
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		status := state.Status.ValueString()
		line.Status = &status
	}
	if !state.StatusReason.IsNull() && !state.StatusReason.IsUnknown() {
		reason := state.StatusReason.ValueString()
		line.StatusReason = &reason
	}
	return line, nil
}

func modelToLineChange(ctx context.Context, state *FactoryLineModel) (*cli.FactoryLineChange, error) {
	stages, err := stagesFromModel(state.Stages)
	if err != nil {
		return nil, err
	}

	change := &cli.FactoryLineChange{
		Title:  ptrString(state.Title),
		Stages: stages,
	}

	if !state.Description.IsNull() && !state.Description.IsUnknown() {
		change.Description = ptrString(state.Description)
	}
	if !state.Team.IsNull() && !state.Team.IsUnknown() {
		team, err := stringListFromTerraform(ctx, state.Team)
		if err != nil {
			return nil, err
		}
		change.Team = team
	}
	if !state.MaxRunning.IsNull() && !state.MaxRunning.IsUnknown() {
		max := int(state.MaxRunning.ValueInt64())
		change.MaxRunning = &max
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		change.Status = ptrString(state.Status)
	}
	if !state.StatusReason.IsNull() && !state.StatusReason.IsUnknown() {
		change.StatusReason = ptrString(state.StatusReason)
	}
	return change, nil
}

func stagesFromModel(stages []StageModel) ([]cli.FactoryLineStage, error) {
	out := make([]cli.FactoryLineStage, 0, len(stages))
	for _, stage := range stages {
		settings, err := executorSettingsFromStage(stage)
		if err != nil {
			return nil, fmt.Errorf("stage %q: %w", stage.Identifier.ValueString(), err)
		}

		checkpoints, err := checkpointsJSON(stage.Checkpoints)
		if err != nil {
			return nil, err
		}

		s := cli.FactoryLineStage{
			Identifier:       stage.Identifier.ValueString(),
			Executor:         stage.Executor.ValueString(),
			ExecutorSettings: settings,
			Checkpoints:      checkpoints,
			MaxTries:         int(stage.MaxTries.ValueInt64()),
			AllowTakeOver:    stage.AllowTakeOver.ValueBool(),
		}

		if !stage.Title.IsNull() && !stage.Title.IsUnknown() {
			title := stage.Title.ValueString()
			s.Title = &title
		}
		if !stage.Description.IsNull() && !stage.Description.IsUnknown() {
			desc := stage.Description.ValueString()
			s.Description = &desc
		}
		if !stage.Instructions.IsNull() && !stage.Instructions.IsUnknown() {
			instructions := stage.Instructions.ValueString()
			s.Instructions = &instructions
		}
		if !stage.Team.IsNull() && !stage.Team.IsUnknown() {
			team, err := stringListFromTerraform(context.Background(), stage.Team)
			if err != nil {
				return nil, err
			}
			s.Team = team
		}
		if !stage.Timeout.IsNull() && !stage.Timeout.IsUnknown() {
			timeout := stage.Timeout.ValueString()
			s.Timeout = &timeout
		}
		out = append(out, s)
	}
	return out, nil
}

func executorSettingsFromStage(stage StageModel) (json.RawMessage, error) {
	executor := stage.Executor.ValueString()
	switch executor {
	case "factory_ai_managed_agent":
		if stage.FactoryAI == nil {
			return json.RawMessage("{}"), nil
		}
		m := map[string]any{}
		putString(m, "installation_id", stage.FactoryAI.InstallationID)
		putString(m, "environment_id", stage.FactoryAI.EnvironmentID)
		putString(m, "computer_id", stage.FactoryAI.ComputerID)
		putString(m, "model", stage.FactoryAI.Model)
		putString(m, "reasoning_effort", stage.FactoryAI.ReasoningEffort)
		return json.Marshal(m)
	case "claude_managed_agent":
		if stage.Claude == nil {
			return nil, fmt.Errorf("claude_managed_agent executor requires a claude_managed_agent settings block")
		}
		m := map[string]any{}
		putStringRequired(m, "agent_id", stage.Claude.AgentID)
		putStringRequired(m, "environment_id", stage.Claude.EnvironmentID)
		putString(m, "installation_id", stage.Claude.InstallationID)
		putString(m, "git_access_token", stage.Claude.GitAccessToken)
		if err := putJSONString(m, "session", stage.Claude.Session); err != nil {
			return nil, err
		}
		if err := putJSONString(m, "execution_identity", stage.Claude.ExecutionIdentity); err != nil {
			return nil, err
		}
		if _, ok := m["session"]; !ok {
			m["session"] = map[string]any{}
		}
		return json.Marshal(m)
	case "devin_managed_agent":
		if stage.Devin == nil {
			return json.RawMessage("{}"), nil
		}
		m := map[string]any{}
		putString(m, "installation_id", stage.Devin.InstallationID)
		putString(m, "playbook_id", stage.Devin.PlaybookID)
		if err := putJSONString(m, "execution_identity", stage.Devin.ExecutionIdentity); err != nil {
			return nil, err
		}
		if err := putJSONString(m, "session", stage.Devin.Session); err != nil {
			return nil, err
		}
		if !stage.Devin.Repos.IsNull() && !stage.Devin.Repos.IsUnknown() {
			repos := make([]string, 0, len(stage.Devin.Repos.Elements()))
			for _, elem := range stage.Devin.Repos.Elements() {
				if s, ok := elem.(types.String); ok {
					repos = append(repos, s.ValueString())
				}
			}
			if len(repos) > 0 {
				m["repos"] = repos
			}
		}
		if _, ok := m["session"]; !ok {
			m["session"] = map[string]any{}
		}
		return json.Marshal(m)
	case "custom_workflow":
		if stage.CustomWorkflow == nil {
			return nil, fmt.Errorf("custom_workflow executor requires a custom_workflow settings block")
		}
		m := map[string]any{}
		putStringRequired(m, "workflow", stage.CustomWorkflow.Workflow)
		return json.Marshal(m)
	case "human":
		if stage.Human == nil {
			return json.RawMessage("{}"), nil
		}
		m := map[string]any{}
		if err := putJSONString(m, "assignees", stage.Human.Assignees); err != nil {
			return nil, err
		}
		if err := putJSONString(m, "notify", stage.Human.Notify); err != nil {
			return nil, err
		}
		return json.Marshal(m)
	case "none":
		return json.RawMessage("{}"), nil
	default:
		return nil, fmt.Errorf("unsupported executor %q", executor)
	}
}

func checkpointsJSON(raw types.String) (json.RawMessage, error) {
	if raw.IsNull() || raw.IsUnknown() || stringsTrim(raw.ValueString()) == "" {
		return json.RawMessage("[]"), nil
	}
	if !json.Valid([]byte(raw.ValueString())) {
		return nil, fmt.Errorf("checkpoints must be valid JSON")
	}
	return json.RawMessage(raw.ValueString()), nil
}

func ptrString(s types.String) *string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	v := s.ValueString()
	return &v
}

func putString(m map[string]any, key string, v types.String) {
	if v.IsNull() || v.IsUnknown() || v.ValueString() == "" {
		return
	}
	m[key] = v.ValueString()
}

func putStringRequired(m map[string]any, key string, v types.String) {
	if v.IsNull() || v.IsUnknown() || v.ValueString() == "" {
		return
	}
	m[key] = v.ValueString()
}

func putJSONString(m map[string]any, key string, v types.String) error {
	if v.IsNull() || v.IsUnknown() || stringsTrim(v.ValueString()) == "" {
		return nil
	}
	if !json.Valid([]byte(v.ValueString())) {
		return fmt.Errorf("%s must be valid JSON", key)
	}
	var parsed any
	if err := json.Unmarshal([]byte(v.ValueString()), &parsed); err != nil {
		return fmt.Errorf("%s must be valid JSON: %w", key, err)
	}
	m[key] = parsed
	return nil
}

func stringsTrim(s string) string {
	return strings.TrimSpace(s)
}

func stringListFromTerraform(ctx context.Context, list types.List) ([]string, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	elems, err := utils.TerraformListToGoArray(ctx, list, "string")
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(elems))
	for _, elem := range elems {
		out = append(out, elem.(string))
	}
	return out, nil
}
