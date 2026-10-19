package llm_provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/llm_provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnthropicCompatibleToPortBody(t *testing.T) {
	state := &llm_provider.LLMProviderModel{
		ProviderType: types.StringValue("anthropic-compatible"),
		Enabled:  types.BoolValue(true),
		AnthropicCompatible: &llm_provider.AnthropicCompatibleModel{
			BaseUrl:          types.StringValue("https://api.example.com/anthropic"),
			ApiKeySecretName: types.StringValue("MY_KEY"),
			Models: []llm_provider.CustomModelModel{
				{
					Name:        types.StringValue("claude-sonnet"),
					DisplayName: types.StringValue("Claude Sonnet"),
					ContextWindow: types.Int64Value(200000),
					SupportedFeatures: &llm_provider.SupportedFeaturesModel{
						Caching:                types.BoolValue(true),
						NativeStructuredOutput: types.BoolValue(false),
					},
				},
			},
		},
	}

	body, err := llm_provider.LlmProviderToPortBodyForTest(state)
	require.NoError(t, err)
	assert.Equal(t, "anthropic-compatible", body.Provider)
	assert.True(t, body.Enabled)
	assert.Equal(t, "https://api.example.com/anthropic", body.Config["baseUrl"])
	assert.Equal(t, "MY_KEY", body.Config["apiKeySecretName"])

	models, ok := body.Config["models"].([]map[string]any)
	require.True(t, ok)
	require.Len(t, models, 1)
	assert.Equal(t, "claude-sonnet", models[0]["name"])
	assert.Equal(t, "Claude Sonnet", models[0]["displayName"])
	assert.Equal(t, int64(200000), models[0]["contextWindow"])

	features, ok := models[0]["supportedFeatures"].(map[string]bool)
	require.True(t, ok)
	assert.True(t, features["caching"])
	assert.False(t, features["nativeStructuredOutput"])
}

func TestAnthropicCompatibleRequiresExclusiveAuth(t *testing.T) {
	state := &llm_provider.LLMProviderModel{
		ProviderType: types.StringValue("anthropic-compatible"),
		AnthropicCompatible: &llm_provider.AnthropicCompatibleModel{
			BaseUrl:             types.StringValue("https://api.example.com/anthropic"),
			ApiKeySecretName:    types.StringValue("MY_KEY"),
			AuthTokenSecretName: types.StringValue("MY_TOKEN"),
			Models: []llm_provider.CustomModelModel{
				{Name: types.StringValue("claude-sonnet")},
			},
		},
	}

	_, err := llm_provider.LlmProviderToPortBodyForTest(state)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "either api_key_secret_name or auth_token_secret_name")
}

func TestAnthropicCompatibleAuthTokenSecretName(t *testing.T) {
	state := &llm_provider.LLMProviderModel{
		ProviderType: types.StringValue("anthropic-compatible"),
		AnthropicCompatible: &llm_provider.AnthropicCompatibleModel{
			BaseUrl:             types.StringValue("https://api.example.com/anthropic"),
			AuthTokenSecretName: types.StringValue("MY_TOKEN"),
			Models: []llm_provider.CustomModelModel{
				{Name: types.StringValue("claude-sonnet")},
			},
		},
	}

	body, err := llm_provider.LlmProviderToPortBodyForTest(state)
	require.NoError(t, err)
	assert.Equal(t, "MY_TOKEN", body.Config["authTokenSecretName"])
	assert.NotContains(t, body.Config, "apiKeySecretName")
}
