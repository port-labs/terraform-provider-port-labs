package entity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDeleteTestClient(t *testing.T, status int, body string) *cli.PortClient {
	t.Helper()
	t.Setenv("PORT_RATE_LIMIT_DISABLED", "1")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/blueprints/service_v1/entities/my-service", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	c, err := cli.New(server.URL)
	require.NoError(t, err)
	c.Client.SetRetryCount(0)
	return c
}

// I8: after a blueprint rename the old blueprint is gone, so deleting from it returns 404; that must not fail the update.
func TestDeleteEntityFromPreviousBlueprint_NotFoundIsSuccess(t *testing.T) {
	c := newDeleteTestClient(t, http.StatusNotFound,
		`{"ok":false,"error":"not_found","message":"Blueprint with identifier \"service_v1\" was not found"}`)

	assert.NoError(t, deleteEntityFromPreviousBlueprint(context.Background(), c, "my-service", "service_v1"))
}

// Any other failure must still surface, even when the body happens to contain "not_found".
func TestDeleteEntityFromPreviousBlueprint_ServerErrorMentioningNotFoundFails(t *testing.T) {
	c := newDeleteTestClient(t, http.StatusInternalServerError,
		`{"ok":false,"error":"internal_error","message":"upstream returned not_found while deleting relations"}`)

	err := deleteEntityFromPreviousBlueprint(context.Background(), c, "my-service", "service_v1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete entity")
}

func TestDeleteEntityFromPreviousBlueprint_Success(t *testing.T) {
	c := newDeleteTestClient(t, http.StatusOK, `{"ok":true}`)

	assert.NoError(t, deleteEntityFromPreviousBlueprint(context.Background(), c, "my-service", "service_v1"))
}
func TestDeleteEntity_NotFoundErrorCarriesStatusCode(t *testing.T) {
	c := newDeleteTestClient(t, http.StatusNotFound, `{"ok":false,"error":"not_found"}`)

	err := c.DeleteEntity(context.Background(), "my-service", "service_v1")
	require.Error(t, err)
	assert.True(t, cli.IsNotFound(err))
}
