package folder

import (
	"context"
	"os"

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

func (r *FolderResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var state FolderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	betaFeaturesEnabledEnv := os.Getenv("PORT_BETA_FEATURES_ENABLED")
	if !(betaFeaturesEnabledEnv == "true") {
		resp.Diagnostics.AddError("Beta features are not enabled", "Folder resource is currently in beta and is subject to change in future versions. Use it by setting the Environment Variable PORT_BETA_FEATURES_ENABLED=true.")
		return
	}
}

var FolderResourceMarkdownDescription = `

# Folder resource

A full list of the available folder types and their identifiers can be found [here](https://docs.getport.io/customize-pages-dashboards-and-plugins/folder/catalog-folder).

~> **WARNING**
The folder resource is currently in beta and is subject to change in future versions.
Use it by setting the Environment Variable ` + "`PORT_BETA_FEATURES_ENABLED=true`" + `.
If this Environment Variable isn't specified, you won't be able to use the resource.

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
