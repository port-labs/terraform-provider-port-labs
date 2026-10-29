package folder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func FolderSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Folder state identifier",
			Computed:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the folder",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the folder",
			Optional:            true,
		},
		"after": schema.StringAttribute{
			MarkdownDescription: "The identifier of the folder after which the folder should be placed",
			Optional:            true,
			Computed:            true,
		},
		"parent": schema.StringAttribute{
			MarkdownDescription: "The identifier of the parent folder",
			Optional:            true,
		},
	}
}

func (r *FolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: FolderResourceMarkdownDescription,
		Attributes:          FolderSchema(),
	}
}

var FolderResourceMarkdownDescription = `

# Folder resource

For more information about folders, see the [Port documentation](https://docs.port.io/customize-pages-dashboards-and-plugins/page/folders#folder-identifiers).

## Example Usage

### Basic Folder

` + "```hcl" + `

resource "port_folder" "example_folder" {
  identifier = "example_folder"
  title      = "Example Folder"
}

` + "```" + `

### Folder with Parent

Create a folder inside another folder.

` + "```hcl" + `

resource "port_folder" "child_folder" {
  identifier = "child_folder"
  parent     = port_folder.example_folder.identifier
  title      = "Child Folder"
}

` + "```" + `

### Folder with After

Create a folder after another folder.

` + "```hcl" + `

resource "port_folder" "another_folder" {
  identifier = "another_folder"
  after      = port_folder.example_folder.identifier
  title      = "Another Folder"
}

` + "```" + `

`
