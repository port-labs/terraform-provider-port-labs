package page

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func PageToPortBody(pm *PageModel) (*cli.Page, error) {
	pb := &cli.Page{
		Identifier:  pm.Identifier.ValueString(),
		Type:        pm.Type.ValueString(),
		Icon:        pm.Icon.ValueStringPointer(),
		Title:       pm.Title.ValueStringPointer(),
		Locked:      pm.Locked.ValueBoolPointer(),
		Blueprint:   pm.Blueprint.ValueStringPointer(),
		Parent:      pm.Parent.ValueStringPointer(),
		After:       pm.After.ValueStringPointer(),
		Description: pm.Description.ValueStringPointer(),
	}

	widgets, err := utils.TerraformJsonStringListToGoObjects(pm.Widgets)
	if err != nil {
		return nil, err
	}
	pb.Widgets = widgets

	pageFilters, err := utils.TerraformJsonStringListToGoObjects(pm.PageFilters)
	if err != nil {
		return nil, err
	}
	pb.PageFilters = pageFilters

	pageFilterPresets, err := utils.TerraformJsonStringListToGoObjects(pm.PageFilterPresets)
	if err != nil {
		return nil, err
	}
	pb.PageFilterPresets = pageFilterPresets

	return pb, nil
}
