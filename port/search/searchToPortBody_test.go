package search

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSearchResourceToPortBodyWithFullTextSearch(t *testing.T) {
	state := &SearchDataModel{
		Query: types.StringValue(`{"combinator":"and","rules":[]}`),
		FullTextSearch: &FullTextSearchModel{
			Term:         types.StringValue("payments"),
			TargetFields: []types.String{types.StringValue("title"), types.StringValue("description")},
		},
	}

	body, err := searchResourceToPortBody(state)
	require.NoError(t, err)
	require.NotNil(t, body.FullTextSearch)
	require.Equal(t, "payments", body.FullTextSearch.Term)
	require.Equal(t, []string{"title", "description"}, body.FullTextSearch.TargetFields)
}
