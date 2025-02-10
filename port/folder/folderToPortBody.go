package folder

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func FolderToPortRequest(fm *FolderModel) (*cli.Folder, error) {
	fb := &cli.Folder{
		Identifier: fm.Identifier.ValueString(),
		// Sidebar:    fm.Sidebar.ValueString(),

		Title:  fm.Title.ValueString(),
		After:  fm.After.ValueString(),
		Parent: fm.Parent.ValueString(),
		//Matan 10-02

		// Title:  fm.Title.ValueStringPointer(),
		// After:  fm.After.ValueStringPointer(),
		// Parent: fm.Parent.ValueStringPointer(),
	}
	return fb, nil
}
