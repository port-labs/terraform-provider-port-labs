package folder

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func FolderToPortRequest(fm *FolderModel) (*cli.Folder, error) {
	fb := &cli.Folder{
		Identifier: fm.Identifier.ValueString(),
		Title:      fm.Title.ValueString(),
		After:      fm.After.ValueString(),
		Parent:     fm.Parent.ValueString(),
	}
	return fb, nil
}
