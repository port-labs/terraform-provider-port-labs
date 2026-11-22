package organization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshOrganizationState(ctx context.Context, state *OrganizationModel, org *cli.Organization) error {
	state.ID = types.StringValue(org.Name)
	state.Name = types.StringValue(org.Name)

	tabs := []string{}
	if org.Settings != nil && org.Settings.HiddenGlobalSearchTabs != nil {
		tabs = org.Settings.HiddenGlobalSearchTabs
	}

	tabSet, diags := types.SetValueFrom(ctx, types.StringType, tabs)
	if diags.HasError() {
		return fmt.Errorf("failed to set hidden_global_search_tabs: %s", diags.Errors()[0].Detail())
	}
	state.HiddenGlobalSearchTabs = tabSet

	return nil
}
