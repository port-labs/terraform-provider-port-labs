package search

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ datasource.DataSource = &SearchDataSource{}

func NewSearchDataSource() datasource.DataSource {
	return &SearchDataSource{}
}

type SearchDataSource struct {
	portClient *cli.PortClient
}

func (d *SearchDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.portClient = req.ProviderData.(*cli.PortClient)
}

func (d *SearchDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search"
}

func (d *SearchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SearchDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	searchRequest, err := searchResourceToPortBody(&data)
	if err != nil {
		resp.Diagnostics.AddError("failed to convert search data to port body", err.Error())
		return
	}

	searchResult, err := d.portClient.Search(ctx, searchRequest)
	if err != nil {
		resp.Diagnostics.AddError("failed to search", err.Error())
		return
	}

	data.ID = types.StringValue(data.GenerateID())
	data.MatchingBlueprints = goStringListToTFList(searchResult.MatchingBlueprints)

	blueprints := make(map[string]cli.Blueprint)
	for _, blueprint := range searchResult.MatchingBlueprints {
		b, _, err := d.portClient.ReadBlueprint(ctx, blueprint)
		if err != nil {
			resp.Diagnostics.AddError("failed to read blueprint", err.Error())
			return
		}
		blueprints[blueprint] = *b
	}

	for _, entity := range searchResult.Entities {
		matchingEntityBlueprint := blueprints[entity.Blueprint]
		e := refreshEntityState(ctx, &entity, &matchingEntityBlueprint)
		data.Entities = append(data.Entities, *e)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func goStringListToTFList(list []string) []types.String {
	var result = make([]types.String, len(list))
	for i, u := range list {
		result[i] = types.StringValue(u)
	}

	return result
}
