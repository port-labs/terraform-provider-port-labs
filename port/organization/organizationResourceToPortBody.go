package organization

import (
	"context"
	"fmt"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func organizationResourceToPortBody(ctx context.Context, state *OrganizationModel) (*cli.OrganizationUpdate, error) {
	update := &cli.OrganizationUpdate{}

	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		name := state.Name.ValueString()
		update.Name = &name
	}

	if !state.HiddenGlobalSearchTabs.IsNull() && !state.HiddenGlobalSearchTabs.IsUnknown() {
		var tabs []string
		diags := state.HiddenGlobalSearchTabs.ElementsAs(ctx, &tabs, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to read hidden_global_search_tabs: %s", diags.Errors()[0].Detail())
		}
		if tabs == nil {
			tabs = []string{}
		}
		update.Settings = &cli.OrganizationSettings{
			HiddenGlobalSearchTabs: tabs,
		}
	}

	return update, nil
}
