package cli

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTeamTestClient starts a fake Port API that serves the organization feature flags, a team named "Backend Team",
// and a canned response for the _team entity search.
func newTeamTestClient(t *testing.T, featureFlags []string, searchStatus int, searchBody string) (*PortClient, *atomic.Int32) {
	t.Helper()
	t.Setenv("PORT_RATE_LIMIT_DISABLED", "1")
	var searchCalls atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/organization", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":           true,
			"organization": map[string]any{"name": "test-org", "featureFlags": featureFlags},
		})
	})
	mux.HandleFunc("GET /v1/teams/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":   true,
			"team": map[string]any{"name": r.PathValue("name"), "provider": "port"},
		})
	})
	mux.HandleFunc("POST /v1/entities/search", func(w http.ResponseWriter, _ *http.Request) {
		searchCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(searchStatus)
		_, _ = w.Write([]byte(searchBody))
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	c, err := New(server.URL)
	require.NoError(t, err)
	c.Client.SetRetryCount(0)
	return c, &searchCalls
}

func TestReadTeam_WithoutOwnershipV2_NoTeamEntity_FallsBackToNullIdentifier(t *testing.T) {
	c, searchCalls := newTeamTestClient(t, []string{}, http.StatusOK, `{"ok":true,"entities":[]}`)

	start := time.Now()
	team, _, err := c.ReadTeam(context.Background(), "Backend Team")
	require.NoError(t, err)
	assert.Nil(t, team.Identifier)
	assert.Equal(t, "Backend Team", team.Name)
	assert.LessOrEqual(t, searchCalls.Load(), int32(1), "must not retry the _team lookup for orgs without ownership v2")
	assert.Less(t, time.Since(start), 2*time.Second)
}

func TestReadTeam_WithoutOwnershipV2_SearchFails_FallsBackToNullIdentifier(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{}, http.StatusNotFound,
		`{"ok":false,"error":"not_found","message":"Blueprint with identifier \"_team\" was not found"}`)

	team, _, err := c.ReadTeam(context.Background(), "Backend Team")
	require.NoError(t, err)
	assert.Nil(t, team.Identifier)
}

func TestReadTeam_WithoutOwnershipV2_TeamEntityExists_SetsIdentifier(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{}, http.StatusOK,
		`{"ok":true,"entities":[{"identifier":"backend_team","title":"Backend Team"}]}`)

	team, _, err := c.ReadTeam(context.Background(), "Backend Team")
	require.NoError(t, err)
	require.NotNil(t, team.Identifier)
	assert.Equal(t, "backend_team", *team.Identifier)
}

func TestReadTeam_WithOwnershipV2_TeamEntityExists_SetsIdentifier(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{FeatureFlagUsersAndTeamsV2}, http.StatusOK,
		`{"ok":true,"entities":[{"identifier":"backend_team","title":"Backend Team"}]}`)

	team, _, err := c.ReadTeam(context.Background(), "Backend Team")
	require.NoError(t, err)
	require.NotNil(t, team.Identifier)
	assert.Equal(t, "backend_team", *team.Identifier)
}

func TestEnrichTeam_WithOwnershipV2_NoTeamEntity_ReturnsMissingTeamEntityError(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{FeatureFlagUsersAndTeamsV2}, http.StatusOK, `{"ok":true,"entities":[]}`)

	_, err := c.enrichTeamFromTeamEntity(context.Background(), &PortTeam{Name: "Backend Team"})
	assert.True(t, errors.Is(err, MissingTeamEntityError), "got %v", err)
}

func TestEnrichTeam_WithOwnershipV2_SearchFails_ReturnsError(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{FeatureFlagUsersAndTeamsV2}, http.StatusInternalServerError, `{"ok":false,"error":"boom"}`)

	_, err := c.enrichTeamFromTeamEntity(context.Background(), &PortTeam{Name: "Backend Team"})
	assert.Error(t, err)
}
