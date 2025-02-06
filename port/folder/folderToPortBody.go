package folder

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func FolderToPortBody(fm *FolderModel) (*cli.Folder, error) {
	fb := &cli.Folder{
		FolderIdentifier:  fm.FolderIdentifier.ValueString(),
		SidebarIdentifier: fm.SidebarIdentifier.ValueString(),
		Title:             fm.Title.ValueString(),
		After:             fm.After.ValueString(),
		Parent:            fm.Parent.ValueString(),
	}
	return fb, nil
}
