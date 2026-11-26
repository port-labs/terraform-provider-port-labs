package scorecard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOmitManagedRelationKeys(t *testing.T) {
	input := map[string]any{
		"owner_team": "team-a",
		"group":      "my-group",
	}
	got := omitManagedRelationKeys(input, scorecardManagedRelationKey)
	require.Equal(t, map[string]any{"owner_team": "team-a"}, got)
}