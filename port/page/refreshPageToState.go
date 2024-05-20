package page

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshPageToState(pm *PageModel, b *cli.Page) error {
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

	pm.Widgets = make([]types.String, len(*b.Widgets))
	if b.Widgets != nil {
		// go over each widget and convert it to a string and store it in the widgets array
		for i, widget := range *b.Widgets {
			bWidget, err := json.Marshal(widget)
			if err != nil {
				return err
			}
			pm.Widgets[i] = types.StringValue(string(bWidget))
		}
	}
	return nil
}
