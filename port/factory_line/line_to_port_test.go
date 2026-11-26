package factory_line

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestExecutorSettingsFactoryAI(t *testing.T) {
	stage := StageModel{
		Executor: types.StringValue("factory_ai_managed_agent"),
		FactoryAI: &FactoryAIExecutorSettingsModel{
			EnvironmentID:   types.StringValue("env-123"),
			Model:           types.StringValue("droid-v2"),
			ReasoningEffort: types.StringValue("medium"),
		},
	}
	raw, err := executorSettingsFromStage(stage)
	require.NoError(t, err)

	var body map[string]any
	require.NoError(t, json.Unmarshal(raw, &body))
	require.Equal(t, "env-123", body["environment_id"])
	require.Equal(t, "droid-v2", body["model"])
	require.Equal(t, "medium", body["reasoning_effort"])
}
