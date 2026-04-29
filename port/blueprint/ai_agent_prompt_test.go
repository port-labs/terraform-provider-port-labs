package blueprint

import (
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAIAgentPromptMaxLengths(t *testing.T) {
	t.Run("no op for other blueprints", func(t *testing.T) {
		ml := consts.LegacyAIAgentPromptMaxLength
		props := map[string]cli.BlueprintProperty{
			"prompt": {Type: "string", MaxLength: &ml},
		}
		NormalizeAIAgentPromptMaxLengths("service", props)
		require.Equal(t, consts.LegacyAIAgentPromptMaxLength, *props["prompt"].MaxLength)
	})

	t.Run("bumps legacy string maxLength on _ai_agent", func(t *testing.T) {
		ml := consts.LegacyAIAgentPromptMaxLength
		props := map[string]cli.BlueprintProperty{
			"prompt": {Type: "string", MaxLength: &ml},
		}
		NormalizeAIAgentPromptMaxLengths(consts.SystemBlueprintAIAgentIdentifier, props)
		require.Equal(t, consts.AIAgentPromptMaxLength, *props["prompt"].MaxLength)
	})

	t.Run("skips non string types", func(t *testing.T) {
		ml := consts.LegacyAIAgentPromptMaxLength
		props := map[string]cli.BlueprintProperty{
			"count": {Type: "number", MaxLength: &ml},
		}
		NormalizeAIAgentPromptMaxLengths(consts.SystemBlueprintAIAgentIdentifier, props)
		require.Equal(t, consts.LegacyAIAgentPromptMaxLength, *props["count"].MaxLength)
	})

	t.Run("does not change other maxLength values", func(t *testing.T) {
		ml := 4096
		props := map[string]cli.BlueprintProperty{
			"other": {Type: "string", MaxLength: &ml},
		}
		NormalizeAIAgentPromptMaxLengths(consts.SystemBlueprintAIAgentIdentifier, props)
		require.Equal(t, 4096, *props["other"].MaxLength)
	})
}
