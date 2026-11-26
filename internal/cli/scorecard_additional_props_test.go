package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScorecardRequestsAdditionalProperties(t *testing.T) {
	t.Run("nil scorecard", func(t *testing.T) {
		require.False(t, scorecardRequestsAdditionalProperties(nil))
	})

	t.Run("empty scorecard", func(t *testing.T) {
		require.False(t, scorecardRequestsAdditionalProperties(&Scorecard{}))
	})

	t.Run("scorecard properties", func(t *testing.T) {
		require.True(t, scorecardRequestsAdditionalProperties(&Scorecard{
			Properties: map[string]any{"owner": "team-a"},
		}))
	})

	t.Run("scorecard relations", func(t *testing.T) {
		require.True(t, scorecardRequestsAdditionalProperties(&Scorecard{
			Relations: map[string]any{"owner_team": "team-a"},
		}))
	})

	t.Run("rule properties", func(t *testing.T) {
		require.True(t, scorecardRequestsAdditionalProperties(&Scorecard{
			Rules: []Rule{{Properties: map[string]any{"sla": 1}}},
		}))
	})

	t.Run("rule relations", func(t *testing.T) {
		require.True(t, scorecardRequestsAdditionalProperties(&Scorecard{
			Rules: []Rule{{Relations: map[string]any{"owner": "team-a"}}},
		}))
	})
}