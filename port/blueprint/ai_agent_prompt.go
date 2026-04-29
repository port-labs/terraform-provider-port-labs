package blueprint

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// NormalizeAIAgentPromptMaxLengths upgrades legacy maxLength on the _ai_agent
// blueprint from 2500 to 5000 when sending updates to Port. The system blueprint
// structure payload can still include the old cap while the API already accepts
// longer prompts; without this, applies that only touch other fields can keep
// resubmitting maxLength 2500 and block prompt values between 2501 and 5000 bytes.
func NormalizeAIAgentPromptMaxLengths(blueprintIdentifier string, props map[string]cli.BlueprintProperty) {
	if blueprintIdentifier != consts.SystemBlueprintAIAgentIdentifier || len(props) == 0 {
		return
	}
	for id := range props {
		prop := props[id]
		if prop.Type != "" && prop.Type != "string" {
			continue
		}
		if prop.MaxLength == nil || *prop.MaxLength != consts.LegacyAIAgentPromptMaxLength {
			continue
		}
		ml := consts.AIAgentPromptMaxLength
		prop.MaxLength = &ml
		props[id] = prop
	}
}
