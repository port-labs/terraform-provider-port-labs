package folder

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshFolderToState(fm *FolderModel, f *cli.Folder) error {
	fm.ID = types.StringValue(f.Identifier)
	fm.Sidebar = types.StringValue(f.Sidebar)
	fm.Parent = types.StringPointerValue(f.Parent)
	fm.After = types.StringPointerValue(f.After)
	fm.Title = types.StringPointerValue(f.Title)
	return nil
}
