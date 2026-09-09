package folder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func FolderSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Folder state identifier",
			Computed:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the folder.",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The display title of the folder.",
			Optional:            true,
		},
		"sidebar": schema.StringAttribute{
			MarkdownDescription: "The identifier of the sidebar that contains the folder. Currently only `catalog` is supported.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"after": schema.StringAttribute{
			MarkdownDescription: "The identifier of the sibling item after which this folder appears. Omitted when the folder is the first item in its parent.",
			Optional:            true,
			Computed:            true,
		},
		"parent": schema.StringAttribute{
			MarkdownDescription: "The identifier of the parent folder. Omitted when the folder is at the root level of the sidebar.",
			Optional:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the folder",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the folder",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the folder",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the folder",
			Computed:            true,
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
