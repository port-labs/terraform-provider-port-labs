package page

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func (r *PageResource) refreshPageToState(pm *PageModel, b *cli.Page) error {
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
		// go over each widget and convert it to a string and store it in the widgets array
		for i, widget := range *b.Widgets {
			bWidget, err := utils.GoObjectToTerraformString(widget, r.portClient.JSONEscapeHTML)
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
		// go over each page filter and convert it to a string and store it in the page filters array
		for i, pageFilter := range *b.PageFilters {
			bFilter, err := utils.GoObjectToTerraformString(pageFilter, r.portClient.JSONEscapeHTML)
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
