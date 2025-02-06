package folder

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func FolderToPortBody(fm *FolderModel) (*cli.Folder, error) {
	fb := &cli.Folder{
		Identifier: fm.FolderIdentifier.ValueString(),
		// Sidebar:    fm.SidebarIdentifier.ValueString(),
		Title:  fm.Title.ValueStringPointer(),
		After:  fm.After.ValueStringPointer(),
		Parent: fm.Parent.ValueStringPointer(),
	}
	return fb, nil
}
