package scorecard_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ datasource.DataSource = &ScorecardGroupDataSource{}

func NewScorecardGroupDataSource() datasource.DataSource {
	return &ScorecardGroupDataSource{}
}

type ScorecardGroupDataSource struct {
	portClient *cli.PortClient
}

func (d *ScorecardGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard_group"
}

func (d *ScorecardGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.portClient = req.ProviderData.(*cli.PortClient)
}

func (d *ScorecardGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to read an existing scorecard group by identifier.",
		Attributes:          scorecardGroupDataSourceAttributes(),
	}
}

func (d *ScorecardGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ScorecardGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, statusCode, err := d.portClient.ReadScorecardGroup(ctx, config.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("scorecard group not found", config.Identifier.ValueString())
			return
		}
		resp.Diagnostics.AddError("failed to read scorecard group", err.Error())
		return
	}

	state := &ScorecardGroupModel{
		Identifier: config.Identifier,
	}
	jsonEscapeHTML := d.portClient != nil && d.portClient.JSONEscapeHTML
	if err := refreshScorecardGroupModel(state, group, jsonEscapeHTML, true); err != nil {
		resp.Diagnostics.AddError("failed to refresh scorecard group state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
