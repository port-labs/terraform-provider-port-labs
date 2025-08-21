package cli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteEntity(t *testing.T) {
	tests := []struct {
		name              string
		entityID          string
		blueprint         string
		expectedOK        bool
		expectedError     bool
		responseBody      string
		responseStatus    int
		shouldHaveDelDeps bool
	}{
		{
			name:           "successful delete",
			entityID:       "test-entity",
			blueprint:      "test-blueprint",
			expectedOK:     true,
			expectedError:  false,
			responseBody:   `{"ok": true}`,
			responseStatus: http.StatusOK,
		},
		{
			name:           "failed delete",
			entityID:       "test-entity",
			blueprint:      "test-blueprint",
			expectedOK:     false,
			expectedError:  true,
			responseBody:   `{"ok": false, "error": "Entity not found"}`,
			responseStatus: http.StatusOK,
		},
		{
			name:           "server error",
			entityID:       "test-entity",
			blueprint:      "test-blueprint",
			expectedOK:     false,
			expectedError:  true,
			responseBody:   `{"error": "Internal server error"}`,
			responseStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedURL *url.URL
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedURL = r.URL

				expectedPath := "/v1/blueprints/" + tt.blueprint + "/entities/" + tt.entityID
				assert.Equal(t, expectedPath, r.URL.Path, "URL path should match expected pattern")

				assert.Equal(t, http.MethodDelete, r.Method)

				query := r.URL.Query()
				assert.False(t, query.Has("delete_dependents"), "delete_dependents query parameter should not be present")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := &PortClient{
				Client: resty.New().SetBaseURL(server.URL),
			}

			err := client.DeleteEntity(context.Background(), tt.entityID, tt.blueprint)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			require.NotNil(t, receivedURL)
		})
	}
}

func TestDeleteEntityWithDependents(t *testing.T) {
	tests := []struct {
		name              string
		entityID          string
		blueprint         string
		deleteDependents  bool
		expectedOK        bool
		expectedError     bool
		responseBody      string
		responseStatus    int
		shouldHaveDelDeps bool
	}{
		{
			name:              "successful delete without dependents",
			entityID:          "test-entity",
			blueprint:         "test-blueprint",
			deleteDependents:  false,
			expectedOK:        true,
			expectedError:     false,
			responseBody:      `{"ok": true}`,
			responseStatus:    http.StatusOK,
			shouldHaveDelDeps: false,
		},
		{
			name:              "successful delete with dependents",
			entityID:          "test-entity",
			blueprint:         "test-blueprint",
			deleteDependents:  true,
			expectedOK:        true,
			expectedError:     false,
			responseBody:      `{"ok": true}`,
			responseStatus:    http.StatusOK,
			shouldHaveDelDeps: true,
		},
		{
			name:              "failed delete with dependents",
			entityID:          "test-entity",
			blueprint:         "test-blueprint",
			deleteDependents:  true,
			expectedOK:        false,
			expectedError:     true,
			responseBody:      `{"ok": false, "error": "Cannot delete entity with dependents"}`,
			responseStatus:    http.StatusOK,
			shouldHaveDelDeps: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedURL *url.URL
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedURL = r.URL

				expectedPath := "/v1/blueprints/" + tt.blueprint + "/entities/" + tt.entityID
				assert.Equal(t, expectedPath, r.URL.Path, "URL path should match expected pattern")

				assert.Equal(t, http.MethodDelete, r.Method)

				query := r.URL.Query()
				if tt.shouldHaveDelDeps {
					assert.True(t, query.Has("delete_dependents"), "delete_dependents query parameter should be present")
					assert.Equal(t, "true", query.Get("delete_dependents"), "delete_dependents should be 'true'")
				} else {
					assert.False(t, query.Has("delete_dependents"), "delete_dependents query parameter should not be present")
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := &PortClient{
				Client: resty.New().SetBaseURL(server.URL),
			}

			err := client.DeleteEntityWithDependents(context.Background(), tt.entityID, tt.blueprint, tt.deleteDependents)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			require.NotNil(t, receivedURL)
		})
	}
}

func TestDeleteEntityBothFunctions(t *testing.T) {
	t.Run("wrapper vs explicit false should behave identically", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++

			query := r.URL.Query()
			assert.False(t, query.Has("delete_dependents"), "delete_dependents query parameter should not be present")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok": true}`))
		}))
		defer server.Close()

		client := &PortClient{
			Client: resty.New().SetBaseURL(server.URL),
		}

		err1 := client.DeleteEntity(context.Background(), "test-entity", "test-blueprint")
		assert.NoError(t, err1)

		err2 := client.DeleteEntityWithDependents(context.Background(), "test-entity", "test-blueprint", false)
		assert.NoError(t, err2)

		assert.Equal(t, 2, requestCount)
	})

	t.Run("explicit true should include delete_dependents parameter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			assert.True(t, query.Has("delete_dependents"), "delete_dependents query parameter should be present")
			assert.Equal(t, "true", query.Get("delete_dependents"), "delete_dependents should be 'true'")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok": true}`))
		}))
		defer server.Close()

		client := &PortClient{
			Client: resty.New().SetBaseURL(server.URL),
		}

		err := client.DeleteEntityWithDependents(context.Background(), "test-entity", "test-blueprint", true)
		assert.NoError(t, err)
	})
}
