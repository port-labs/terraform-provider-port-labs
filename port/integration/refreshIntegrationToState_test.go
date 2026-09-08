package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGithubExternalPropertiesNamespaceClaimsToState(t *testing.T) {
	t.Run("nil claims", func(t *testing.T) {
		result := githubExternalPropertiesNamespaceClaimsToState(nil)
		assert.True(t, result.IsNull())
	})

	t.Run("empty claims", func(t *testing.T) {
		result := githubExternalPropertiesNamespaceClaimsToState(map[string]bool{})
		assert.True(t, result.IsNull())
	})

	t.Run("populated claims", func(t *testing.T) {
		result := githubExternalPropertiesNamespaceClaimsToState(map[string]bool{
			"my-org":    true,
			"other-org": true,
		})

		assert.False(t, result.IsNull())
		elements := result.Elements()
		assert.Len(t, elements, 2)
		assert.Equal(t, types.BoolValue(true), elements["my-org"])
		assert.Equal(t, types.BoolValue(true), elements["other-org"])
	})
}
