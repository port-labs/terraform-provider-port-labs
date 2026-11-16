package page

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// parseJSONState decodes a JSON string from state, returning nil when it is null, unknown or not valid JSON.
func parseJSONState(s types.String) interface{} {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s.ValueString()), &v); err != nil {
		return nil
	}
	return v
}

// goObjectToTerraformStringPreferExistingWithHiddenQueryNorm compares a widget from state (preferred)
// with an API response widget (v), normalizing hiddenQuery arrays in both before comparison.
// If the normalized values are semantically equal, the original preferred string is returned to
// prevent diffs caused only by hiddenQuery array reordering.
func goObjectToTerraformStringPreferExistingWithHiddenQueryNorm(preferred types.String, v map[string]any, jsonEscapeHTML bool) (types.String, error) {
	encoded, err := utils.GoObjectToTerraformString(v, jsonEscapeHTML)
	if err != nil {
		return types.StringNull(), err
	}

	if preferred.IsNull() || preferred.IsUnknown() {
		return encoded, nil
	}

	// Parse and normalize the preferred (state) value
	oldParsed := parseJSONState(preferred)
	if oldParsed == nil {
		return encoded, nil
	}
	oldMap, ok := oldParsed.(map[string]interface{})
	if !ok {
		return encoded, nil
	}

	// Normalize both for comparison
	normalizedOld := NormalizeWidgetForComparison(oldMap)
	normalizedNew := NormalizeWidgetForComparison(v)

	// Compare normalized versions
	normalizedOldJSON, err := utils.GoObjectToTerraformString(normalizedOld, jsonEscapeHTML)
	if err != nil {
		return encoded, nil
	}
	normalizedNewJSON, err := utils.GoObjectToTerraformString(normalizedNew, jsonEscapeHTML)
	if err != nil {
		return encoded, nil
	}

	equal, err := utils.JSONStringsSemanticallyEqual(normalizedOldJSON.ValueString(), normalizedNewJSON.ValueString(), jsonEscapeHTML)
	if err != nil || !equal {
		return encoded, nil
	}

	// Normalized values are equal, return original preferred to preserve user's formatting
	return preferred, nil
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
			oldWidget := utils.TerraformStringAtList(oldWidgets, i)
			// Normalize hiddenQuery arrays in both old state and API response to enable
			// order-insensitive comparison. The Port backend may reorder hiddenQuery arrays
			// asynchronously per organization, causing perpetual diffs without normalization.
			normalizedWidget := NormalizeWidgetForComparison(widget)
			bWidget, err := goObjectToTerraformStringPreferExistingWithHiddenQueryNorm(
				oldWidget, normalizedWidget, r.portClient.JSONEscapeHTML)
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
