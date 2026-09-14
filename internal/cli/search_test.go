package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchRequestBodyUsesWrappedQueryWhenFullTextSearchSet(t *testing.T) {
	query := map[string]any{"combinator": "and", "rules": []any{}}
	body := &SearchRequestQuery{
		Query:          &query,
		FullTextSearch: &FullTextSearch{Term: "foo", TargetFields: []string{"title"}},
	}

	result := searchRequestBody(body)
	require.Equal(t, query, result["query"])
	fullTextSearch, ok := result["fullTextSearch"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "foo", fullTextSearch["term"])
	require.Equal(t, []string{"title"}, fullTextSearch["targetFields"])
}

func TestSearchRequestBodyUsesFlatQueryWithoutFullTextSearch(t *testing.T) {
	query := map[string]any{"combinator": "and", "rules": []any{}}
	body := &SearchRequestQuery{Query: &query}

	result := searchRequestBody(body)
	require.Equal(t, query, result)
}
