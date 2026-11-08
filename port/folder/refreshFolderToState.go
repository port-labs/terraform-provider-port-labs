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

	if f.Sidebar != "" {
		fm.Sidebar = types.StringValue(f.Sidebar)
	}

	if f.After != "" {
		fm.After = types.StringValue(f.After)
	} else if !fm.After.IsNull() {
		fm.After = types.StringNull()
	}

	if f.Parent != "" {
		fm.Parent = types.StringValue(f.Parent)
	}

	return nil
}

func refreshFolderMetaToState(fm *FolderModel, f *cli.Folder) {
	if f.CreatedAt != nil {
		fm.CreatedAt = types.StringValue(f.CreatedAt.String())
	}

	if f.CreatedBy != "" {
		fm.CreatedBy = types.StringValue(f.CreatedBy)
	}

	if f.UpdatedAt != nil {
		fm.UpdatedAt = types.StringValue(f.UpdatedAt.String())
	}

	if f.UpdatedBy != "" {
		fm.UpdatedBy = types.StringValue(f.UpdatedBy)
	}
}
