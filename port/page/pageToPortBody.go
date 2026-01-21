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

	widgets, err := widgetsToPortBody(pm.Widgets)
	if err != nil {
		return nil, err
	}
	pb.Widgets = widgets

	pageFilters, err := pageFiltersToPortBody(pm.PageFilters)
	if err != nil {
		return nil, err
	}
	pb.PageFilters = pageFilters

	return pb, nil
}

func widgetsToPortBody(widgets types.List) (*[]map[string]any, error) {
	if widgets.IsNull() || widgets.IsUnknown() {
		return nil, nil
	}
	widgetsBody := make([]map[string]any, len(widgets.Elements()))
	for i, w := range widgets.Elements() {
		strVal := w.(types.String)
		v, err := utils.TerraformJsonStringToGoObject(strVal.ValueStringPointer())

		if err != nil {
			return nil, err
		}

		widgetsBody[i] = *v
	}

	return &widgetsBody, nil
}

func pageFiltersToPortBody(pageFilters types.List) (*[]map[string]any, error) {
	if pageFilters.IsNull() || pageFilters.IsUnknown() {
		return nil, nil
	}
	pageFiltersBody := make([]map[string]any, len(pageFilters.Elements()))
	for i, pf := range pageFilters.Elements() {
		strVal := pf.(types.String)
		v, err := utils.TerraformJsonStringToGoObject(strVal.ValueStringPointer())

		if err != nil {
			return nil, err
		}

		pageFiltersBody[i] = *v
	}

	return &pageFiltersBody, nil
}
