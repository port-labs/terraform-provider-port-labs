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

func widgetsToPortBody(widgets []types.String) (*[]map[string]any, error) {
	if widgets == nil {
		return nil, nil
	}
	widgetsBody := make([]map[string]any, len(widgets))
	for i, w := range widgets {
		v, err := utils.TerraformJsonStringToGoObject(w.ValueStringPointer())

		if err != nil {
			return nil, err
		}

		widgetsBody[i] = *v
	}

	return &widgetsBody, nil
}

func pageFiltersToPortBody(pageFilters []types.String) (*[]map[string]any, error) {
	if pageFilters == nil {
		return nil, nil
	}
	pageFiltersBody := make([]map[string]any, len(pageFilters))
	for i, pf := range pageFilters {
		v, err := utils.TerraformJsonStringToGoObject(pf.ValueStringPointer())

		if err != nil {
			return nil, err
		}

		pageFiltersBody[i] = *v
	}

	return &pageFiltersBody, nil
}
