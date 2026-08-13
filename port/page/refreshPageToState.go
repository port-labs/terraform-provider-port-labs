package page

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *PageResource) refreshPageToState(pm *PageModel, b *cli.Page) error {
	oldWidgets := pm.Widgets
	oldPageFilters := pm.PageFilters

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
		var oldWidgetElements []attr.Value
		if !oldWidgets.IsNull() {
			oldWidgetElements = oldWidgets.Elements()
		}
		for i, widget := range *b.Widgets {
			var oldValue types.String
			if i < len(oldWidgetElements) {
				oldValue = oldWidgetElements[i].(types.String)
			}
			bWidget, err := utils.GoObjectToTerraformStringPreferExisting(oldValue, widget, r.portClient.JSONEscapeHTML)
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
		var oldFilterElements []attr.Value
		if !oldPageFilters.IsNull() {
			oldFilterElements = oldPageFilters.Elements()
		}
		for i, pageFilter := range *b.PageFilters {
			var oldValue types.String
			if i < len(oldFilterElements) {
				oldValue = oldFilterElements[i].(types.String)
			}
			bFilter, err := utils.GoObjectToTerraformStringPreferExisting(oldValue, pageFilter, r.portClient.JSONEscapeHTML)
			if err != nil {
				return err
			}
			filterAttrs[i] = bFilter
		}
		pm.PageFilters, _ = types.ListValue(types.StringType, filterAttrs)
	} else {
		pm.PageFilters = types.ListNull(types.StringType)
	}
	return nil
}
