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

func TestGetDatasourceEntitiesPaginatesWithoutExtraConditions(t *testing.T) {
	t.Parallel()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/blueprints/entities/datasource-entities", r.URL.Path)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "port-ocean/github/", body["datasource_prefix"])
		require.Equal(t, "/installation/resync", body["datasource_suffix"])
		require.NotContains(t, body, "extra_conditions")

		requestCount++
		if requestCount == 1 {
			next := "cursor-2"
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"entities": []map[string]string{
					{"identifier": "entity-1", "blueprint": "service"},
				},
				"next": next,
			}))
			return
		}

		require.Equal(t, "cursor-2", body["from"])
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"ok": true,
			"entities": []map[string]string{
				{"identifier": "entity-2", "blueprint": "service"},
			},
		}))
	}))
	t.Cleanup(server.Close)

	client, err := cli.New(server.URL)
	require.NoError(t, err)

	entities, err := client.GetDatasourceEntities(context.Background(), &cli.DatasourceEntitiesRequest{
		DatasourcePrefix: "port-ocean/github/",
		DatasourceSuffix: "/installation/resync",
	})
	require.NoError(t, err)
	require.Len(t, entities, 2)
	assert.Equal(t, "entity-1", entities[0].Identifier)
	assert.Equal(t, "entity-2", entities[1].Identifier)
	assert.Equal(t, 2, requestCount)
}
