package folder

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshFolderToState(fm *FolderModel, f *cli.Folder) error {
	fm.FolderIdentifier = types.StringValue(f.FolderIdentifier)
	fm.SidebarIdentifier = types.StringValue(f.SidebarIdentifier)
	fm.Title = ConvertStringPointerToTypesString(f.Title)
	fm.After = ConvertStringPointerToTypesString(f.After)
	fm.Parent = ConvertStringPointerToTypesString(f.Parent)
	fm.CreatedAt = types.StringValue(f.CreatedAt.String())
	fm.CreatedBy = types.StringValue(f.CreatedBy)
	fm.UpdatedAt = types.StringValue(f.UpdatedAt.String())
	fm.UpdatedBy = types.StringValue(f.UpdatedBy)
	return nil
}

func ConvertStringPointerToTypesString(s *string) *types.String {
	if s == nil {
		return nil
	}
	strValue := types.StringValue(*s)
	return &strValue
}
