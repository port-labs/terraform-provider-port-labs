package page

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

const entityPageType = "entity"

// alignEntityPageGroupIds changes the widgets returned by the API, in place, so the tab (group) ids of an
// entity page match how the prior state tracks them.
//
// Port gives an id to every tab of an entity page's main grouper, but a configuration doesn't have to declare
// them. Tabs are matched to the prior state by title:
//   - The prior tab has no id and the API returns one: the id is removed, since the configuration doesn't
//     manage it.
//   - The prior tab has an id and the API returns none (Port hides ids while tab management is off for the
//     organization): the prior id is kept.
//   - Both have an id: the API's id is kept, so a changed id shows as drift.
//
// Tabs that aren't in the prior state are left as they are, so a tab added outside Terraform shows as drift.
// Nothing changes when there is no prior state (for example, on import).
func alignEntityPageGroupIds(pageType string, apiWidgets []map[string]any, priorWidgets types.List) {
	if pageType != entityPageType || priorWidgets.IsNull() || priorWidgets.IsUnknown() {
		return
	}

	apiGroups := mainGrouperGroups(apiWidgets)
	if len(apiGroups) == 0 {
		return
	}

	priorGroupsByTitle := make(map[string]map[string]any)
	for _, group := range mainGrouperGroups(decodeWidgets(priorWidgets)) {
		if title, ok := group["title"].(string); ok {
			priorGroupsByTitle[title] = group
		}
	}

	for _, group := range apiGroups {
		title, _ := group["title"].(string)
		priorGroup, ok := priorGroupsByTitle[title]
		if !ok {
			continue
		}

		priorID, priorHasID := priorGroup["id"]
		_, apiHasID := group["id"]

		switch {
		case !priorHasID && apiHasID:
			delete(group, "id")
		case priorHasID && !apiHasID:
			group["id"] = priorID
		}
	}
}

// mainGrouperGroups returns the groups of the first top-level grouper, the same grouper Port gives ids to.
// Nested groupers (such as the one in the Related Entities tab) are not included.
func mainGrouperGroups(widgets []map[string]any) []map[string]any {
	for _, widget := range widgets {
		if widget["type"] != "grouper" {
			continue
		}

		rawGroups, _ := widget["groups"].([]any)
		groups := make([]map[string]any, 0, len(rawGroups))
		for _, rawGroup := range rawGroups {
			if group, ok := rawGroup.(map[string]any); ok {
				groups = append(groups, group)
			}
		}
		return groups
	}

	return nil
}

// decodeWidgets decodes the widget JSON strings stored in the state. Widgets that can't be decoded are skipped.
func decodeWidgets(widgets types.List) []map[string]any {
	decoded := make([]map[string]any, 0, len(widgets.Elements()))
	for _, element := range widgets.Elements() {
		str, ok := element.(types.String)
		if !ok || str.IsNull() || str.IsUnknown() {
			continue
		}

		widget, err := utils.TerraformJsonStringToGoObject(str.ValueStringPointer())
		if err != nil || widget == nil {
			continue
		}
		decoded = append(decoded, *widget)
	}

	return decoded
}
