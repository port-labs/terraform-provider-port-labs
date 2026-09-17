package scorecard_group

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ datasource.DataSource = &ScorecardGroupsDataSource{}

func NewScorecardGroupsDataSource() datasource.DataSource {
	return &ScorecardGroupsDataSource{}
}

type ScorecardGroupsDataSource struct {
	portClient *cli.PortClient
}

func (d *ScorecardGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.portClient = req.ProviderData.(*cli.PortClient)
}

func (d *ScorecardGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard_groups"
}

func (d *ScorecardGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to list all scorecard groups in the organization. " +
			"This mirrors the Port API `GET /v1/scorecard-groups` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Stable identifier for this data source result.",
				Computed:            true,
			},
			"groups": schema.ListNestedAttribute{
				MarkdownDescription: "All scorecard groups returned by the Port API.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							MarkdownDescription: "The scorecard group identifier.",
							Computed:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "The scorecard group title.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

type scorecardGroupsDataSourceModel struct {
	ID     types.String                 `tfsdk:"id"`
	Groups []scorecardGroupSummaryModel `tfsdk:"groups"`
}

type scorecardGroupSummaryModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
}

func (d *ScorecardGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scorecardGroupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := d.portClient.ListScorecardGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list scorecard groups", err.Error())
		return
	}

	state.Groups = make([]scorecardGroupSummaryModel, 0, len(groups))
	identifiers := make([]string, 0, len(groups))
	for _, group := range groups {
		state.Groups = append(state.Groups, scorecardGroupSummaryModel{
			Identifier: types.StringValue(group.Identifier),
			Title:      types.StringValue(group.Title),
		})
		identifiers = append(identifiers, group.Identifier)
	}

	sort.Strings(identifiers)
	hash := sha256.Sum256([]byte(strings.Join(identifiers, "\n")))
	state.ID = types.StringValue(hex.EncodeToString(hash[:]))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
