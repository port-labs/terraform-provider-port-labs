package folder

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshFolderToState(fm *FolderModel, f *cli.Folder) error {
	fm.Identifier = types.StringValue(f.Identifier)
	fm.ID = types.StringValue(f.Identifier)

	if f.Title != "" {
		fm.Title = types.StringValue(f.Title)
	}

	if f.After != "" {
		fm.After = types.StringValue(f.After)
	}

	if f.Parent != "" {
		fm.Parent = types.StringValue(f.Parent)
	}

	return nil
}
