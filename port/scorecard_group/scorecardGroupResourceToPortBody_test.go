package scorecard_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScorecardGroupResourceToPortBodyExtendedFields(t *testing.T) {
	ctx := context.Background()
	state := &ScorecardGroupModel{
		Identifier:          types.StringValue("my-group"),
		Title:               types.StringValue("My Group"),
		GroupProperties:     types.StringValue(`{"category":"production"}`),
		ScorecardProperties: types.StringValue(`{"owner":"platform-team"}`),
		GroupRelations:      types.StringValue(`{"owning_team":"team-1"}`),
		ScorecardRelations:  types.StringValue(`{"owner_team":["team-1","team-2"]}`),
		Blueprints:          []types.String{types.StringValue("service")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"author","operator":"isNotEmpty"}`),
					},
				},
			},
		},
	}

	body, err := scorecardGroupResourceToPortBody(ctx, state)
	require.NoError(t, err)
	assert.Equal(t, "my-group", body.Identifier)
	assert.Equal(t, "production", body.GroupProperties["category"])
	assert.Equal(t, "platform-team", body.ScorecardProperties["owner"])
	assert.Equal(t, "team-1", body.GroupRelations["owning_team"])
	assert.Equal(t, []interface{}{"team-1", "team-2"}, body.ScorecardRelations["owner_team"])
}

func TestScorecardGroupResourceToPortBodyLegacyPropertiesAlias(t *testing.T) {
	ctx := context.Background()
	state := &ScorecardGroupModel{
		Identifier: types.StringValue("my-group"),
		Title:      types.StringValue("My Group"),
		Properties: types.StringValue(`{"owner":"legacy"}`),
		Blueprints: []types.String{types.StringValue("service")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"author","operator":"isNotEmpty"}`),
					},
				},
			},
		},
	}

	body, err := scorecardGroupResourceToPortBody(ctx, state)
	require.NoError(t, err)
	assert.Equal(t, "legacy", body.ScorecardProperties["owner"])
	assert.Nil(t, body.GroupProperties)
}

func TestScorecardGroupResourceToPortBodyPrefersScorecardPropertiesOverLegacy(t *testing.T) {
	ctx := context.Background()
	state := &ScorecardGroupModel{
		Identifier:          types.StringValue("my-group"),
		Title:               types.StringValue("My Group"),
		Properties:          types.StringValue(`{"owner":"legacy"}`),
		ScorecardProperties: types.StringValue(`{"owner":"explicit"}`),
		Blueprints:          []types.String{types.StringValue("service")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"author","operator":"isNotEmpty"}`),
					},
				},
			},
		},
	}

	body, err := scorecardGroupResourceToPortBody(ctx, state)
	require.NoError(t, err)
	assert.Equal(t, "explicit", body.ScorecardProperties["owner"])
}
