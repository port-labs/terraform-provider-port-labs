package page

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// stripEmptyArrays recursively removes empty arrays from a map to prevent
// perpetual diffs when the API injects default empty arrays like "hiddenQuery: []"
// that weren't in the user's config.
func stripEmptyArrays(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			cleaned := stripEmptyArrays(value)
			// Skip empty arrays
			if arr, ok := cleaned.([]interface{}); ok && len(arr) == 0 {
				continue
			}
			result[key] = cleaned
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = stripEmptyArrays(item)
		}
		return result
	default:
		return data
	}
}

func (r *PageResource) refreshPageToState(pm *PageModel, b *cli.Page) error {
	oldWidgets := pm.Widgets
	oldPageFilters := pm.PageFilters
	oldPageFilterPresets := pm.PageFilterPresets

	pm.ID = types.StringValue(b.Identifier)
	pm.Identifier = types.StringValue(b.Identifier)
	pm.Type = types.StringValue(b.Type)
	pm.Icon = types.StringPointerValue(b.Icon)
	pm.Parent = types.StringPointerValue(b.Parent)
	pm.After = types.StringPointerValue(b.After)
	pm.Title = types.StringPointerValue(b.Title)
	pm.Locked = types.BoolPointerValue(b.Locked)
	pm.Blueprint = types.StringPointerValue(b.Blueprint)
	pm.Description = types.StringPointerValue(b.Description)

	if b.Widgets != nil {
		widgetAttrs := make([]attr.Value, len(*b.Widgets))
		for i, widget := range *b.Widgets {
			// Strip empty arrays like "hiddenQuery: []" that the API adds but aren't in user config
			cleanedWidget := stripEmptyArrays(widget)
			bWidget, err := utils.GoObjectToTerraformStringPreferExisting(
				utils.TerraformStringAtList(oldWidgets, i), cleanedWidget, r.portClient.JSONEscapeHTML)
			if err != nil {
				return err
			}
			widgetAttrs[i] = bWidget
		}
		pm.Widgets, _ = types.ListValue(types.StringType, widgetAttrs)
	} else {
		pm.Widgets = types.ListNull(types.StringType)
	}

	if b.PageFilters != nil {
		filterAttrs := make([]attr.Value, len(*b.PageFilters))
		for i, pageFilter := range *b.PageFilters {
			bFilter, err := utils.GoObjectToTerraformStringPreferExisting(
				utils.TerraformStringAtList(oldPageFilters, i), pageFilter, r.portClient.JSONEscapeHTML)
			if err != nil {
				return err
			}
			filterAttrs[i] = bFilter
		}
		pm.PageFilters, _ = types.ListValue(types.StringType, filterAttrs)
	} else {
		pm.PageFilters = types.ListNull(types.StringType)
	}

	if b.PageFilterPresets != nil {
		presetAttrs := make([]attr.Value, len(*b.PageFilterPresets))
		for i, preset := range *b.PageFilterPresets {
			bPreset, err := utils.GoObjectToTerraformStringPreferExisting(
				utils.TerraformStringAtList(oldPageFilterPresets, i), preset, r.portClient.JSONEscapeHTML)
			if err != nil {
				return err
			}
			presetAttrs[i] = bPreset
		}
		pm.PageFilterPresets, _ = types.ListValue(types.StringType, presetAttrs)
	} else {
		pm.PageFilterPresets = types.ListNull(types.StringType)
	}
	return nil
}
