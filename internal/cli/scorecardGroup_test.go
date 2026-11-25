package cli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateScorecardGroupNullResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok": true, "scorecardGroup": null}`))
	}))
	t.Cleanup(server.Close)

	client := &PortClient{Client: resty.New().SetBaseURL(server.URL)}
	group, err := client.UpdateScorecardGroup(context.Background(), "my-group", &ScorecardGroup{Identifier: "my-group"})
	require.NoError(t, err)
	assert.Nil(t, group)
}

func TestUpdateScorecardGroupReturnsGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok": true, "scorecardGroup": {"identifier": "my-group", "title": "My Group"}}`))
	}))
	t.Cleanup(server.Close)

	client := &PortClient{Client: resty.New().SetBaseURL(server.URL)}
	group, err := client.UpdateScorecardGroup(context.Background(), "my-group", &ScorecardGroup{Identifier: "my-group"})
	require.NoError(t, err)
	require.NotNil(t, group)
	assert.Equal(t, "my-group", group.Identifier)
}