package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringPropEncryptionToBody(t *testing.T) {
	t.Parallel()

	property, err := stringPropToBody(context.Background(), StringPropModel{
		Encryption: types.StringValue("aes256-gcm"),
	})
	require.NoError(t, err)
	assert.Equal(t, "aes256-gcm", property.Encryption)
}

func TestStringPropClientSideEncryptionToBody(t *testing.T) {
	t.Parallel()

	property, err := stringPropToBody(context.Background(), StringPropModel{
		ClientSideEncryption: &ClientSideEncryptionModel{
			Algorithm: types.StringValue("client-side"),
			Key:       types.StringValue("test-key"),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"algorithm": "client-side",
		"key":       "test-key",
	}, property.Encryption)
}

func TestObjectPropEncryptionToBody(t *testing.T) {
	t.Parallel()

	property, err := objectPropToBody(context.Background(), ObjectPropModel{
		Format:     types.StringValue("multi-line"),
		Encryption: types.StringValue("aes256-gcm"),
	})
	require.NoError(t, err)
	assert.Equal(t, "multi-line", *property.Format)
	assert.Equal(t, "aes256-gcm", property.Encryption)
}

func TestStringPropEncryptionToState(t *testing.T) {
	t.Parallel()

	prop, err := stringPropToState(context.Background(), cli.WorkflowInputProperty{
		Type:       "string",
		Encryption: "aes256-gcm",
	}, false)
	require.NoError(t, err)
	assert.Equal(t, "aes256-gcm", prop.Encryption.ValueString())
	assert.Nil(t, prop.ClientSideEncryption)
}

func TestStringPropClientSideEncryptionToState(t *testing.T) {
	t.Parallel()

	prop, err := stringPropToState(context.Background(), cli.WorkflowInputProperty{
		Type: "string",
		Encryption: map[string]string{
			"algorithm": "client-side",
			"key":       "test-key",
		},
	}, false)
	require.NoError(t, err)
	assert.True(t, prop.Encryption.IsNull())
	require.NotNil(t, prop.ClientSideEncryption)
	assert.Equal(t, "client-side", prop.ClientSideEncryption.Algorithm.ValueString())
	assert.Equal(t, "test-key", prop.ClientSideEncryption.Key.ValueString())
}

func TestArrayStringItemsConstraintsToBody(t *testing.T) {
	t.Parallel()

	property, err := arrayPropToBody(context.Background(), ArrayPropModel{
		StringItems: &StringItemsModel{
			Pattern:   types.StringValue("^[a-z]+$"),
			MinLength: types.Int64Value(2),
			MaxLength: types.Int64Value(10),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, property.Items)
	assert.Equal(t, "string", property.Items["type"])
	assert.Equal(t, "^[a-z]+$", property.Items["pattern"])
	assert.Equal(t, 2, property.Items["minLength"])
	assert.Equal(t, 10, property.Items["maxLength"])
}

func TestArrayStringItemsConstraintsToState(t *testing.T) {
	t.Parallel()

	prop, err := arrayPropToState(context.Background(), cli.WorkflowInputProperty{
		Type: "array",
		Items: map[string]any{
			"type":      "string",
			"pattern":   "^[a-z]+$",
			"minLength": float64(2),
			"maxLength": float64(10),
		},
	}, false)
	require.NoError(t, err)
	require.NotNil(t, prop.StringItems)
	assert.Equal(t, "^[a-z]+$", prop.StringItems.Pattern.ValueString())
	assert.Equal(t, int64(2), prop.StringItems.MinLength.ValueInt64())
	assert.Equal(t, int64(10), prop.StringItems.MaxLength.ValueInt64())
}

func TestEncryptionToStateRoundTrip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		value any
	}{
		{name: "server side", value: "aes256-gcm"},
		{name: "client side map", value: map[string]string{"algorithm": "client-side", "key": "pem-key"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			encryption, clientSide := encryptionToState(tc.value)
			assert.Equal(t, tc.value, encryptionToBody(encryption, clientSide))
		})
	}
}
