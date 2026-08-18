package organization

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestOrganizationResourceToPortBody_NameOnly(t *testing.T) {
	state := &OrganizationModel{
		Name: types.StringValue("my-org"),
	}

	update, err := organizationResourceToPortBody(context.Background(), state)
	require.NoError(t, err)
	require.NotNil(t, update.Name)
	require.Equal(t, "my-org", *update.Name)
	require.Nil(t, update.Settings)
}

func TestOrganizationResourceToPortBody_PortalIconURL(t *testing.T) {
	state := &OrganizationModel{
		PortalIcon: types.StringValue("https://example.com/icon.png"),
	}

	update, err := organizationResourceToPortBody(context.Background(), state)
	require.NoError(t, err)
	require.Nil(t, update.Name)
	require.NotNil(t, update.Settings)
	require.NotNil(t, update.Settings.PortalIcon)
	require.Equal(t, "https://example.com/icon.png", *update.Settings.PortalIcon)
}

func TestOrganizationResourceToPortBody_PortalIconAssetPath(t *testing.T) {
	state := &OrganizationModel{
		PortalIcon: types.StringValue("organizations/org-123/portal_icon/asset-456"),
	}

	update, err := organizationResourceToPortBody(context.Background(), state)
	require.NoError(t, err)
	require.NotNil(t, update.Settings)
	require.Equal(t, "organizations/org-123/portal_icon/asset-456", *update.Settings.PortalIcon)
}

func TestOrganizationResourceToPortBody_PortalIconEmptyString(t *testing.T) {
	state := &OrganizationModel{
		PortalIcon: types.StringValue(""),
	}

	update, err := organizationResourceToPortBody(context.Background(), state)
	require.NoError(t, err)
	require.NotNil(t, update.Settings)
	require.NotNil(t, update.Settings.PortalIcon)
	require.Equal(t, "", *update.Settings.PortalIcon)
}

func TestRefreshOrganizationState_PortalIcon(t *testing.T) {
	portalIcon := "organizations/org-123/portal_icon/asset-456"
	state := &OrganizationModel{}
	org := &cli.Organization{
		Name: "my-org",
		Settings: &cli.OrganizationSettings{
			PortalIcon: utils.PtrTo(portalIcon),
		},
	}

	err := refreshOrganizationState(context.Background(), state, org)
	require.NoError(t, err)
	require.Equal(t, "my-org", state.ID.ValueString())
	require.Equal(t, "my-org", state.Name.ValueString())
	require.Equal(t, portalIcon, state.PortalIcon.ValueString())
}

func TestRefreshOrganizationState_PortalIconNull(t *testing.T) {
	state := &OrganizationModel{
		PortalIcon: types.StringValue("stale-value"),
	}
	org := &cli.Organization{
		Name: "my-org",
	}

	err := refreshOrganizationState(context.Background(), state, org)
	require.NoError(t, err)
	require.True(t, state.PortalIcon.IsNull())
}
