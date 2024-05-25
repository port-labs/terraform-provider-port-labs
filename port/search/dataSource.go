package search

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &SearchDataSource{}

type SearchDataSource struct{}

type ThingDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
}

func (d *SearchDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search"
}

func (d *SearchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ThingDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Typically data sources will make external calls, however this example
	// hardcodes setting the id attribute to a specific value for brevity.
	data.ID = types.StringValue("example-id")
	data.Title = types.StringValue("blabla")
	data.Identifier = types.StringValue("blabla")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *SearchDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "###some schema",
		Attributes:          Schema(),
	}
}

func Schema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the blueprint",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The display name of the blueprint",
			Required:            true,
		},
	}

}
