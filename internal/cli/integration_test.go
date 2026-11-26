package cli

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureIntegrationBody runs fn against a fake Port API and returns the JSON body it sent, decoded into a map.
func captureIntegrationBody(t *testing.T, fn func(c *PortClient) error) map[string]any {
	t.Helper()
	t.Setenv("PORT_RATE_LIMIT_DISABLED", "1")
	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"integration":{"installationId":"my-integration"}}`))
	}))
	t.Cleanup(server.Close)

	c, err := New(server.URL)
	require.NoError(t, err)
	require.NoError(t, fn(c))

	var got map[string]any
	require.NoError(t, json.Unmarshal(body, &got), "body: %s", body)
	return got
}

func strPtr(s string) *string { return &s }

// I7: a minimal custom integration must not send null for unset optional fields on create
// (the API rejects "installationAppType": null).
func TestCreateIntegration_OmitsUnsetOptionalFields(t *testing.T) {
	got := captureIntegrationBody(t, func(c *PortClient) error {
		_, err := c.CreateIntegration(context.Background(), &Integration{
			InstallationId:   "my-integration",
			InstallationType: strPtr("Custom"),
		})
		return err
	})

	assert.Equal(t, "my-integration", got["installationId"])
	for _, key := range []string{"title", "version", "installationAppType"} {
		assert.NotContains(t, got, key)
	}
}

// The PATCH body carries unset title and version as null, so removing title from config clears it. (version is
// Optional+Computed with UseStateForUnknown, so removing it from config keeps its prior value in the plan.)
func TestUpdateIntegration_SendsNullForUnsetTitleAndVersion(t *testing.T) {
	got := captureIntegrationBody(t, func(c *PortClient) error {
		_, err := c.UpdateIntegration(context.Background(), "my-integration", &Integration{
			InstallationId:   "my-integration",
			InstallationType: strPtr("Custom"),
		})
		return err
	})

	require.Contains(t, got, "title")
	assert.Nil(t, got["title"])
	require.Contains(t, got, "version")
	assert.Nil(t, got["version"])
	assert.NotContains(t, got, "installationAppType", "unset app type must still be omitted on update")
}

func TestUpdateIntegration_SendsSetTitleAndVersion(t *testing.T) {
	got := captureIntegrationBody(t, func(c *PortClient) error {
		_, err := c.UpdateIntegration(context.Background(), "my-integration", &Integration{
			InstallationId:      "my-integration",
			Title:               strPtr("My Integration"),
			Version:             strPtr("1.2.3"),
			InstallationAppType: strPtr("kafka"),
		})
		return err
	})

	assert.Equal(t, "My Integration", got["title"])
	assert.Equal(t, "1.2.3", got["version"])
	assert.Equal(t, "kafka", got["installationAppType"])
}
