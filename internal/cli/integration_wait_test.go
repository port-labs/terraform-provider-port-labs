package cli_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
)

// integrationAPI serves GET /v1/integration/my-gitlab with the given status, or 404 when status is empty
func integrationAPI(t *testing.T, status, message string) *cli.PortClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if status == "" {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": "not_found"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "integration": map[string]any{
			"installationId": "my-gitlab",
			"statusInfo":     map[string]any{"integrationStatus": map[string]any{"status": status, "message": message}},
		}})
	}))
	t.Cleanup(server.Close)
	return &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}
}

func TestWaitForIntegrationDeletedFailsFastWhenDeleteFailed(t *testing.T) {
	client := integrationAPI(t, "Error", "Fail to delete Integration")

	start := time.Now()
	err := client.WaitForIntegrationDeleted(context.Background(), "my-gitlab")

	var failedErr *cli.IntegrationOperationFailedError
	if assert.True(t, errors.As(err, &failedErr), "got %v", err) {
		assert.Equal(t, "delete", failedErr.Operation)
		assert.Equal(t, "Error", failedErr.Status)
		assert.Equal(t, "Fail to delete Integration", failedErr.Message)
	}
	assert.Less(t, time.Since(start), time.Second, "a failed delete must not be polled")
}

func TestWaitForIntegrationDeletedReturnsWhenGone(t *testing.T) {
	client := integrationAPI(t, "", "")

	start := time.Now()
	err := client.WaitForIntegrationDeleted(context.Background(), "my-gitlab")

	assert.NoError(t, err)
	assert.Less(t, time.Since(start), time.Second)
}

func TestWaitForIntegrationDeletedKeepsWaitingWhileDeleting(t *testing.T) {
	client := integrationAPI(t, "Deleting", "")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	err := client.WaitForIntegrationDeleted(ctx, "my-gitlab")

	var failedErr *cli.IntegrationOperationFailedError
	assert.False(t, errors.As(err, &failedErr), "a delete in progress is not a failure")
	assert.Error(t, err, "still deleting when the context ends")
}

func TestWaitForIntegrationReadyFailsFastOnError(t *testing.T) {
	client := integrationAPI(t, "Error", "Fail to create Integration")

	_, err := client.WaitForIntegrationReady(context.Background(), "my-gitlab")

	var failedErr *cli.IntegrationOperationFailedError
	if assert.True(t, errors.As(err, &failedErr), "got %v", err) {
		assert.Equal(t, "Fail to create Integration", failedErr.Message)
	}
}
