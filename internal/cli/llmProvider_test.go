package cli_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListLLMProviders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/llm-providers", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"ok": true,
			"providers": [
				{
					"provider": "anthropic-compatible",
					"enabled": true,
					"config": {
						"baseUrl": "https://api.example.com",
						"apiKeySecretName": "MY_KEY",
						"models": [{"name": "claude-sonnet"}]
					}
				}
			]
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := cli.New(server.URL)
	require.NoError(t, err)

	providers, statusCode, err := client.ListLLMProviders(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	require.Len(t, providers, 1)
	assert.Equal(t, "anthropic-compatible", providers[0].Provider)
	assert.True(t, providers[0].Enabled)
	assert.Equal(t, "https://api.example.com", providers[0].Config["baseUrl"])
}

func TestUpsertLLMProviderValidateConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/llm-providers", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("validate_connection"))
		w.Header().Set("Content-Type", "application/json")

		var body map[string]any
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "anthropic-compatible", body["provider"])

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(`{
			"ok": true,
			"provider": {
				"provider": "anthropic-compatible",
				"enabled": true,
				"config": {"baseUrl": "https://api.example.com"}
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := cli.New(server.URL)
	require.NoError(t, err)

	provider, err := client.UpsertLLMProvider(context.Background(), &cli.LLMProviderUpsert{
		Provider: "anthropic-compatible",
		Enabled:  true,
		Config: map[string]any{
			"baseUrl": "https://api.example.com",
		},
	}, true)
	require.NoError(t, err)
	assert.Equal(t, "anthropic-compatible", provider.Provider)
}
