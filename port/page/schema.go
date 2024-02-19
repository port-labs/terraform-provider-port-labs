package page

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
)

func PageSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			Description: "The Identifier of the page",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the page, can be one of \"blueprint-entities\", \"dashboard\" or \"home\"",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"blueprint-entities",
					"dashboard",
					"home",
				),
			},
		},
		"parent": schema.StringAttribute{
			Description: "The identifier of the folder in which the page is in, default is the root of the sidebar",
			Optional:    true,
		},
		"after": schema.StringAttribute{
			Description: "The identifier of the page/folder after which the page should be placed",
			Optional:    true,
		},
		"icon": schema.StringAttribute{
			Description: "The icon of the page",
			Optional:    true,
		},
		"title": schema.StringAttribute{
			Description: "The title of the page",
			Optional:    true,
		},
		"locked": schema.BoolAttribute{
			Description: "Whether the page is locked, if true, viewers will not be able to edit the page widgets and filters",
			Optional:    true,
		},
		"blueprint": schema.StringAttribute{
			Description: "The blueprint for which the page is created, relevant only for pages of type \"blueprint-entities\"",
			Optional:    true,
		},
		"widgets": schema.ListAttribute{
			Description: "The widgets of the page",
			Optional:    true,
			ElementType: types.StringType,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the page",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the page",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the page",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the page",
			Computed:            true,
		},
	}
}

func (r *PageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: PageResourceMarkdownDescription,
		Attributes:          PageSchema(),
	}
}

func (r *PageResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var state PageModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	betaFeaturesEnabledEnv := os.Getenv("PORT_BETA_FEATURES_ENABLED")
	if !(betaFeaturesEnabledEnv == "true") {
		resp.Diagnostics.AddError("Beta features are not enabled", "Page resource is currently in beta and is subject to change in future versions. Use it by setting the Environment Variable PORT_BETA_FEATURES_ENABLED=true.")
		return
	}
}

var PageResourceMarkdownDescription = `

# Page resource

Docs about the different page types can be found [here](https://docs.getport.io/customize-pages-dashboards-and-plugins/page/catalog-page).

~> **WARNING**  
The page resource is currently in beta and is subject to change in future versions.  
Use it by setting the Environment Variable ` + "`PORT_BETA_FEATURES_ENABLED=true`" + `.  
If this Environment Variable isn't specified, you won't be able to use the resource. 

## Example Usage

### Blueprint Entities Page

` + "```hcl" + `

resource "port_page" "microservice_blueprint_page" {
  identifier            = "microservice_blueprint_page"
  title                 = "Microservices"
  type                  = "blueprint-entities"
  icon                  = "Microservice"
  blueprint             = port_blueprint.base_blueprint.identifier
  widgets               = [
    jsonencode(
      {
        "id" : "microservice-table-entities",
        "type" : "table-entities-explorer",
        "dataset" : {
          "combinator" : "and",
          "rules" : [
            {
              "operator" : "=",
              "property" : "$blueprint",
              "value" : ` + "{{`\"{{blueprint}}\"`}}" + `
            }
          ]
        }
      }
    )
  ]
}

` + "```" + `

### Dashboard Page

` + "```hcl" + `

resource "port_page" "microservice_dashboard_page" {
  identifier            = "microservice_dashboard_page"
  title                 = "Microservices"
  icon                  = "GitHub"
  type                  = "dashboard"
  widgets               = [
    jsonencode(
      {
        "id" : "dashboardWidget",
        "layout" : [
          {
            "height" : 400,
            "columns" : [
              {
                "id" : "microserviceGuide",
                "size" : 12
              }
            ]
          }
        ],
        "type" : "dashboard-widget",
        "widgets" : [
          {
            "title" : "Microservices Guide",
            "icon" : "BlankPage",
            "markdown" : "# This is the new Microservice Dashboard",
            "type" : "markdown",
            "description" : "",
            "id" : "microserviceGuide"
          }
        ],
      }
    )
  ]
}

` + "```" + `


### Page with parent

Create a page inside a folder.

` + "```hcl" + `

resource "port_page" "microservice_dashboard_page" {
  identifier            = "microservice_dashboard_page"
  title                 = "Microservices"
  icon                  = "GitHub"
  type                  = "dashboard"
  parent                = "microservices-folder"
  widgets               = [
    jsonencode(
      {
        "id" : "dashboardWidget",
        "layout" : [
          {
            "height" : 400,
            "columns" : [
              {
                "id" : "microserviceGuide",
                "size" : 12
              }
            ]
          }
        ],
        "type" : "dashboard-widget",
        "widgets" : [
          {
            "title" : "Microservices Guide",
            "icon" : "BlankPage",
            "markdown" : "# This is the new Microservice Dashboard",
            "type" : "markdown",
            "description" : "",
            "id" : "microserviceGuide"
          }
        ],
      }
    )
  ]
}

` + "```" + `


### Page with after

Create a page after another page.

` + "```hcl" + `

resource "port_page" "microservice_dashboard_page" {
  identifier            = "microservice_dashboard_page"
  title                 = "Microservices"
  icon                  = "GitHub"
  type                  = "dashboard"
  after                 = "microservices_entities_page"
  widgets               = [
    jsonencode(
      {
        "id" : "dashboardWidget",
        "layout" : [
          {
            "height" : 400,
            "columns" : [
              {
                "id" : "microserviceGuide",
                "size" : 12
              }
            ]
          }
        ],
        "type" : "dashboard-widget",
        "widgets" : [
          {
            "title" : "Microservices Guide",
            "icon" : "BlankPage",
            "markdown" : "# This is the new Microservice Dashboard",
            "type" : "markdown",
            "description" : "",
            "id" : "microserviceGuide"
          }
        ],
      }
    )
  ]
}

` + "```" + `

### Home Page

` + "```hcl" + `

resource "port_page" "home_page" {
  identifier            = "$home"
  title                 = "Home"
  type                  = "home"
  widgets               = [
    jsonencode(
      {
        "type" : "dashboard-widget",
        "id" : "azkLJD6wLk6nJSvA",
        "layout" : [
          {
            "columns" : [
              {
                "id" : "markdown",
                "size" : 6
              },
              {
                "id" : "overview",
                "size" : 6
              },
            ],
            "height" : 648
          }
        ],
        "widgets" : [
          {
            "type" : "markdown",
            "markdown" : "## Welcome to your internal developer portal",
            "icon" : "port",
            "title" : "About developer Portal",
            "id" : "markdown"
          },
          {
            "type" : "iframe-widget",
            "id" : "overview",
            "title" : "Overview",
            "icon" : "Docs",
            "url" : "https://www.youtube.com/embed/ggXL2ZsPVQM?si=xj6XtV0faatoOhss",
            "urlType" : "public"
          }
        ]
      }
    )
  ]
}

` + "```" + `

The home page is a special page, which is created by default when you create a new organization.

- When deleting the home page resource using terraform, the home page will not be deleted from Port as it isn't deletable page, instead, the home page will be removed from the terraform state. 
- Due to only having one home page you'll have to import the state of the home page manually.

` + "```" + `
terraform import port_page.home_page "\$home"
` + "```" + `
 
`
