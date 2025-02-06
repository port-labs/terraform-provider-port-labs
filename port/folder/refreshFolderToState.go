package folder

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshFolderToState(fm *FolderModel, f *cli.Folder) error {
	fm.FolderIdentifier = types.StringValue(f.FolderIdentifier)
	fm.SidebarIdentifier = types.StringValue(f.SidebarIdentifier)
	fm.Title = types.StringValue(f.Title)
	fm.After = types.StringValue(f.After)
	fm.Parent = types.StringValue(f.Parent)
	return nil
}
