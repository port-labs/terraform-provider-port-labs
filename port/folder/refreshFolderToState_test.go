package folder

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestRefreshFolderToState(t *testing.T) {
	createdAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)

	state := &FolderModel{
		After:  types.StringValue("old-after"),
		Parent: types.StringValue("old-parent"),
	}

	folder := &cli.Folder{
		Meta: cli.Meta{
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
			CreatedBy: "creator@example.com",
			UpdatedBy: "updater@example.com",
		},
		Identifier: "my-folder",
		Sidebar:    "catalog",
		Title:      "My Folder",
	}

	err := refreshFolderToState(state, folder)
	require.NoError(t, err)

	require.Equal(t, "my-folder", state.ID.ValueString())
	require.Equal(t, "my-folder", state.Identifier.ValueString())
	require.Equal(t, "My Folder", state.Title.ValueString())
	require.Equal(t, "catalog", state.Sidebar.ValueString())
	require.True(t, state.After.IsNull())
	require.True(t, state.Parent.IsNull())
	require.Equal(t, createdAt.String(), state.CreatedAt.ValueString())
	require.Equal(t, "creator@example.com", state.CreatedBy.ValueString())
	require.Equal(t, updatedAt.String(), state.UpdatedAt.ValueString())
	require.Equal(t, "updater@example.com", state.UpdatedBy.ValueString())
}

func TestRefreshFolderToStateWithParentAndAfter(t *testing.T) {
	state := &FolderModel{}

	folder := &cli.Folder{
		Identifier: "child-folder",
		Sidebar:    "catalog",
		Title:      "Child Folder",
		Parent:     "parent-folder",
		After:      "sibling-folder",
	}

	err := refreshFolderToState(state, folder)
	require.NoError(t, err)

	require.Equal(t, "parent-folder", state.Parent.ValueString())
	require.Equal(t, "sibling-folder", state.After.ValueString())
}
