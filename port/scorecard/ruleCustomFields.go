package scorecard

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func ApplyRuleCustomFieldsToCLI(stateRule Rule, rule *cli.Rule) error {
	if !stateRule.Properties.IsNull() && !stateRule.Properties.IsUnknown() {
		properties, err := utils.TerraformJsonStringToGoObject(stateRule.Properties.ValueStringPointer())
		if err != nil {
			return err
		}
		if properties != nil {
			rule.Properties = *properties
		}
	}

	if !stateRule.Relations.IsNull() && !stateRule.Relations.IsUnknown() {
		relations, err := utils.TerraformJsonStringToGoObject(stateRule.Relations.ValueStringPointer())
		if err != nil {
			return err
		}
		if relations != nil {
			rule.Relations = *relations
		}
	}

	return nil
}

func SyncRuleCustomFieldsFromAPI(stateRule *Rule, apiRule cli.Rule, jsonEscapeHTML bool, syncFromAPI bool) error {
	if err := utils.SyncJSONObjectState(&stateRule.Properties, apiRule.Properties, "rule properties", jsonEscapeHTML, syncFromAPI); err != nil {
		return err
	}
	if err := utils.SyncJSONObjectState(&stateRule.Relations, apiRule.Relations, "rule relations", jsonEscapeHTML, syncFromAPI); err != nil {
		return err
	}
	return nil
}
