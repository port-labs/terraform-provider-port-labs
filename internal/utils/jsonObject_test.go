package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeJSONObject(t *testing.T) {
	object, err := DecodeJSONObject(`{"context":"user","property":"port_role"}`)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"context": "user", "property": "port_role"}, object)

	for _, invalid := range []string{`["a"]`, `"a"`, `1`, `null`, `{`, ``} {
		_, err := DecodeJSONObject(invalid)
		assert.Error(t, err, invalid)
	}
}

func TestJSONObjectValidator(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   types.String
		wantErr bool
	}{
		{"object", types.StringValue(`{"a":1}`), false},
		{"array", types.StringValue(`[1]`), true},
		{"invalid", types.StringValue(`{`), true},
		{"null is skipped", types.StringNull(), false},
		{"unknown is skipped", types.StringUnknown(), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &validator.StringResponse{}
			JSONObject().ValidateString(context.Background(), validator.StringRequest{Path: path.Root("p"), ConfigValue: tc.value}, resp)
			assert.Equal(t, tc.wantErr, resp.Diagnostics.HasError())
		})
	}
}

func TestPreferExistingKeepsPriorPropertyJson(t *testing.T) {
	property := map[string]any{"context": "user", "property": "a&b<c>"}

	// Escaped in the prior value, as `jsonencode()` writes it.
	prior := types.StringValue(`{"context":"user","property":"a&b<c>"}`)
	for _, escapeHTML := range []bool{true, false} {
		got, err := GoObjectToTerraformStringPreferExisting(prior, property, escapeHTML)
		require.NoError(t, err)
		assert.Equal(t, prior, got, "escapeHTML=%v", escapeHTML)
	}

	// Without a prior value the provider setting decides.
	got, err := GoObjectToTerraformStringPreferExisting(types.StringNull(), property, false)
	require.NoError(t, err)
	assert.JSONEq(t, `{"context":"user","property":"a&b<c>"}`, got.ValueString())
	assert.Contains(t, got.ValueString(), "a&b<c>")

	got, err = GoObjectToTerraformStringPreferExisting(types.StringNull(), property, true)
	require.NoError(t, err)
	assert.Contains(t, got.ValueString(), `\u0026`)

	// A prior value that differs semantically is replaced.
	got, err = GoObjectToTerraformStringPreferExisting(types.StringValue(`{"context":"other"}`), property, false)
	require.NoError(t, err)
	assert.Contains(t, got.ValueString(), "a&b<c>")
}

func TestNotJSONObjectValidator(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   types.String
		wantErr bool
	}{
		{"identifier", types.StringValue("$identifier"), false},
		{"jsonencoded object", types.StringValue(`{"context":"user","property":"port_role"}`), true},
		{"leading space", types.StringValue(` {"a":1}`), true},
		{"null is skipped", types.StringNull(), false},
		{"unknown is skipped", types.StringUnknown(), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &validator.StringResponse{}
			NotJSONObject("property_json").ValidateString(context.Background(), validator.StringRequest{Path: path.Root("p"), ConfigValue: tc.value}, resp)
			assert.Equal(t, tc.wantErr, resp.Diagnostics.HasError())
		})
	}
}
