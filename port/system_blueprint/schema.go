package system_blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/blueprint"
)

type Resource struct {
	client *cli.PortClient
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*cli.PortClient)
}

func (r *Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	blueprintSchemas := blueprint.BlueprintSchema()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Port System Blueprint Resource. This resource is used to extend system blueprints with additional properties and relations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the system blueprint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Identifier of the system blueprint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"properties":             blueprintSchemas["properties"],
			"relations":              blueprintSchemas["relations"],
			"mirror_properties":      blueprintSchemas["mirror_properties"],
			"calculation_properties": blueprintSchemas["calculation_properties"],
		},
	}
}
