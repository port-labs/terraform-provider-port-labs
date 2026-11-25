package utils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertyToState(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		id, obj, err := PropertyToState(types.StringNull(), nil, true)
		require.NoError(t, err)
		assert.True(t, id.IsNull())
		assert.True(t, obj.IsNull())
	})

	t.Run("identifier", func(t *testing.T) {
		id, obj, err := PropertyToState(types.StringNull(), "$identifier", true)
		require.NoError(t, err)
		assert.Equal(t, types.StringValue("$identifier"), id)
		assert.True(t, obj.IsNull())
	})

	t.Run("object keeps prior formatting and HTML escaping", func(t *testing.T) {
		property := map[string]any{"context": "user", "property": "a&b<c>"}
		prior := types.StringValue(`{"context":"user","property":"a\u0026b\u003cc\u003e"}`)
		for _, escapeHTML := range []bool{true, false} {
			id, obj, err := PropertyToState(prior, property, escapeHTML)
			require.NoError(t, err)
			assert.True(t, id.IsNull())
			assert.Equal(t, prior, obj, "escapeHTML=%v", escapeHTML)
		}
	})

	t.Run("object without prior follows json_escape_html", func(t *testing.T) {
		property := map[string]any{"property": "a&b"}
		_, obj, err := PropertyToState(types.StringNull(), property, false)
		require.NoError(t, err)
		assert.Contains(t, obj.ValueString(), "a&b")
		_, obj, err = PropertyToState(types.StringNull(), property, true)
		require.NoError(t, err)
		assert.Contains(t, obj.ValueString(), `\u0026`)
	})

	t.Run("object that differs from prior is replaced", func(t *testing.T) {
		_, obj, err := PropertyToState(types.StringValue(`{"context":"other"}`), map[string]any{"context": "user"}, false)
		require.NoError(t, err)
		assert.JSONEq(t, `{"context":"user"}`, obj.ValueString())
	})
}

func TestPriorAt(t *testing.T) {
	prior := []string{"a", "b"}
	assert.Equal(t, "b", *PriorAt(prior, 1))
	assert.Nil(t, PriorAt(prior, 2))
	assert.Nil(t, PriorAt(prior, -1))
	assert.Nil(t, PriorAt([]string(nil), 0))
}
