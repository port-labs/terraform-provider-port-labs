package search

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func searchResourceToPortBody(state *SearchDataModel) (*cli.SearchRequestQuery, error) {
	query, err := utils.TerraformJsonStringToGoObject(state.Query.ValueStringPointer())
	if err != nil {
		return nil, err
	}

	return &cli.SearchRequestQuery{
		Query:                       query,
		ExcludeCalculatedProperties: state.ExcludeCalculatedProperties.ValueBoolPointer(),
		Include:                     flex.TerraformStringListToGoArray(state.Include),
		Exclude:                     flex.TerraformStringListToGoArray(state.Exclude),
		AttachTitleToRelation:       state.AttachTitleToRelation.ValueBoolPointer(),
	}, nil
}
