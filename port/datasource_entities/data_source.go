package datasource_entities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ datasource.DataSource = &DatasourceEntitiesDataSource{}

func NewDatasourceEntitiesDataSource() datasource.DataSource {
	return &DatasourceEntitiesDataSource{}
}

type DatasourceEntitiesDataSource struct {
	portClient *cli.PortClient
}

func (d *DatasourceEntitiesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.portClient = req.ProviderData.(*cli.PortClient)
}

func (d *DatasourceEntitiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datasource_entities"
}

func (d *DatasourceEntitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatasourceEntitiesDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	request := &cli.DatasourceEntitiesRequest{
		DatasourcePrefix: data.DatasourcePrefix.ValueString(),
		DatasourceSuffix: data.DatasourceSuffix.ValueString(),
	}

	if !data.Limit.IsNull() {
		limit := int(data.Limit.ValueInt64())
		request.Limit = &limit
	}

	if !data.Before.IsNull() {
		before := data.Before.ValueString()
		request.Before = &before
	}

	entities, err := d.portClient.GetDatasourceEntities(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("failed to get datasource entities", err.Error())
		return
	}

	data.ID = types.StringValue(data.GenerateID())
	data.Entities = make([]EntityModel, len(entities))
	for i, entity := range entities {
		data.Entities[i] = EntityModel{
			Identifier: types.StringValue(entity.Identifier),
			Blueprint:  types.StringValue(entity.Blueprint),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
