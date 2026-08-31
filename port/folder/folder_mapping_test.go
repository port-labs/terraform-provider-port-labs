package folder

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestFolderModelToCLIIncludesIconWhenSet(t *testing.T) {
	state := &FolderModel{
		Identifier: types.StringValue("engineering"),
		Title:      types.StringValue("Engineering"),
		Icon:       types.StringValue("AWS"),
	}

	folder := FolderModelToCLI(state)

	require.Equal(t, "engineering", folder.Identifier)
	require.Equal(t, "Engineering", folder.Title)
	require.Equal(t, "AWS", folder.Icon)
}

func TestFolderModelToCLIOmitsIconWhenUnset(t *testing.T) {
	state := &FolderModel{
		Identifier: types.StringValue("engineering"),
		Title:      types.StringValue("Engineering"),
		Icon:       types.StringNull(),
	}

	folder := FolderModelToCLI(state)

	require.Empty(t, folder.Icon)
}

func TestRefreshFolderToStateMapsIcon(t *testing.T) {
	state := &FolderModel{
		Icon: types.StringNull(),
	}

	err := refreshFolderToState(state, &cli.Folder{
		Identifier: "engineering",
		Title:      "Engineering",
		Icon:       "AWS",
	})
	require.NoError(t, err)
	require.Equal(t, "AWS", state.Icon.ValueString())

	err = refreshFolderToState(state, &cli.Folder{
		Identifier: "engineering",
		Title:      "Engineering",
	})
	require.NoError(t, err)
	require.True(t, state.Icon.IsNull())
}
