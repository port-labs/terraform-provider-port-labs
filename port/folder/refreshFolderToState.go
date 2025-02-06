package folder

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshFolderToState(fm *FolderModel, f *cli.Folder) error {
	fm.FolderIdentifier = types.StringValue(f.Identifier)
	fm.SidebarIdentifier = types.StringValue(f.Sidebar)
	fm.Title = types.StringPointerValue(f.Title)
	fm.After = types.StringPointerValue(f.After)
	fm.Parent = types.StringPointerValue(f.Parent)
	return nil
}
