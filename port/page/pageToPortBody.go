package page

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	widgets, err := pageJSONListToPortBody(pm.Widgets)
	if err != nil {
		return nil, err
	}
	pb.Widgets = widgets

	pageFilters, err := pageJSONListToPortBody(pm.PageFilters)
	if err != nil {
		return nil, err
	}
	pb.PageFilters = pageFilters

	pageFilterPresets, err := pageJSONListToPortBody(pm.PageFilterPresets)
	if err != nil {
		return nil, err
	}
	pb.PageFilterPresets = pageFilterPresets

	return pb, nil
}

func pageJSONListToPortBody(items types.List) (*[]map[string]any, error) {
	if items.IsNull() || items.IsUnknown() {
		return nil, nil
	}
	body := make([]map[string]any, len(items.Elements()))
	for i, item := range items.Elements() {
		strVal := item.(types.String)
		v, err := utils.TerraformJsonStringToGoObject(strVal.ValueStringPointer())

		if err != nil {
			return nil, err
		}

		body[i] = *v
	}

	return &body, nil
}
