package page

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func normalizeHiddenQueries(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for key, value := range t {
			if arr, ok := value.([]any); ok && key == "hiddenQuery" {
				if len(arr) == 0 {
					continue
				}
				keys := make([]string, len(arr))
				for i, item := range arr {
					b, _ := json.Marshal(item)
					keys[i] = string(b)
				}
				idx := make([]int, len(arr))
				for i := range idx {
					idx[i] = i
				}
				slices.SortFunc(idx, func(a, b int) int { return strings.Compare(keys[a], keys[b]) })
				sorted := make([]any, 0, len(arr))
				var prevKey string
				for i, j := range idx {
					if i > 0 && keys[j] == prevKey {
						continue
					}
					sorted = append(sorted, arr[j])
					prevKey = keys[j]
				}
				out[key] = sorted
				continue
			}
			out[key] = normalizeHiddenQueries(value)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = normalizeHiddenQueries(item)
		}
		return out
	default:
		return v
	}
}

func resolveWidgetState(prevWidget types.String, remoteWidget map[string]any, jsonEscapeHTML bool) (types.String, error) {
	if !prevWidget.IsNull() && !prevWidget.IsUnknown() {
		var prev map[string]any
		if json.Unmarshal([]byte(prevWidget.ValueString()), &prev) == nil &&
			reflect.DeepEqual(normalizeHiddenQueries(prev), normalizeHiddenQueries(remoteWidget)) {
			return prevWidget, nil
		}
	}
	return utils.GoObjectToTerraformString(remoteWidget, jsonEscapeHTML)
}

func (r *PageResource) refreshPageToState(pm *PageModel, b *cli.Page) error {
	oldWidgets := pm.Widgets
	oldPageFilters := pm.PageFilters
	oldPageFilterPresets := pm.PageFilterPresets
	oldPageDateFilter := pm.PageDateFilter

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
		// pm.Widgets still holds the prior state here, so it can be compared with the API's widgets.
		alignEntityPageGroupIds(b.Type, *b.Widgets, pm.Widgets)

		widgetAttrs := make([]attr.Value, len(*b.Widgets))
		for i, widget := range *b.Widgets {
			bWidget, err := resolveWidgetState(
				utils.TerraformStringAtList(oldWidgets, i), widget, r.portClient.JSONEscapeHTML)
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

	pageDateFilter, err := utils.GoObjectToTerraformStringPreferExisting(oldPageDateFilter, b.PageDateFilter, r.portClient.JSONEscapeHTML)
	if err != nil {
		return err
	}
	pm.PageDateFilter = pageDateFilter
	return nil
}
