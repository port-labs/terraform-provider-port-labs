package folder

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func FolderSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			Description: "The Identifier of the folder",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the folder",
			Required:    true,
		},
		"parent": schema.StringAttribute{
			Description: "The identifier of the parent folder",
			Optional:    true,
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
		"description": schema.StringAttribute{
			MarkdownDescription: "The folder description",
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
  identifier            = "example_folder"
  name                  = "Example Folder"
  description           = "This is an example folder"
}

` + "```" + `

### Folder with Parent

Create a folder inside another folder.

` + "```hcl" + `

resource "port_folder" "child_folder" {
  identifier            = "child_folder"
  name                  = "Child Folder"
  parent                = port_folder.example_folder.identifier
  description           = "This is a child folder"
}

` + "```" + `

`
