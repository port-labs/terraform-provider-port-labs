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
	return newTeamTestClientWithSearch(t, featureFlags, func(int32) (int, string) { return searchStatus, searchBody })
}

// newTeamTestClientWithSearch is like newTeamTestClient, but lets the test pick the _team search response per call
// (calls are numbered from 1).
func newTeamTestClientWithSearch(t *testing.T, featureFlags []string, search func(call int32) (int, string)) (*PortClient, *atomic.Int32) {
	t.Helper()
	t.Setenv("PORT_RATE_LIMIT_DISABLED", "1")
	var searchCalls atomic.Int32

	prevDelay, prevJitter, prevMaxDelay := teamEntityRetryDelay, teamEntityRetryMaxJitter, teamEntityRetryMaxDelay
	teamEntityRetryDelay, teamEntityRetryMaxJitter, teamEntityRetryMaxDelay = time.Millisecond, time.Millisecond, time.Millisecond
	t.Cleanup(func() {
		teamEntityRetryDelay, teamEntityRetryMaxJitter, teamEntityRetryMaxDelay = prevDelay, prevJitter, prevMaxDelay
	})

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
		status, body := search(searchCalls.Add(1))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
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

func TestReadTeam_WithoutOwnershipV2_TeamBlueprintNotFound_FallsBackToNullIdentifier(t *testing.T) {
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

func TestReadTeam_WithoutOwnershipV2_SearchServerError_ReturnsError(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{}, http.StatusInternalServerError, `{"ok":false,"error":"internal_error"}`)

	team, _, err := c.ReadTeam(context.Background(), "Backend Team")
	require.Error(t, err, "a non-404 search failure must not silently become a null identifier")
	assert.Nil(t, team)
}

func TestReadTeam_WithoutOwnershipV2_SearchUnauthorized_ReturnsError(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{}, http.StatusForbidden, `{"ok":false,"error":"forbidden"}`)

	_, _, err := c.ReadTeam(context.Background(), "Backend Team")
	require.Error(t, err)
}

func TestReadTeam_WithoutOwnershipV2_MultipleTeamEntities_FallsBackToNullIdentifier(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{}, http.StatusOK,
		`{"ok":true,"entities":[{"identifier":"a","title":"Backend Team"},{"identifier":"b","title":"Backend Team"}]}`)

	team, _, err := c.ReadTeam(context.Background(), "Backend Team")
	require.NoError(t, err)
	assert.Nil(t, team.Identifier)
}

func TestEnrichTeamWithRetry_WithOwnershipV2_TeamEntityAppearsLate_RetriesUntilFound(t *testing.T) {
	c, searchCalls := newTeamTestClientWithSearch(t, []string{FeatureFlagUsersAndTeamsV2}, func(call int32) (int, string) {
		if call <= 2 {
			return http.StatusOK, `{"ok":true,"entities":[]}`
		}
		return http.StatusOK, `{"ok":true,"entities":[{"identifier":"backend_team","title":"Backend Team"}]}`
	})

	team, err := c.enrichTeamFromTeamEntityWithRetry(context.Background(), &PortTeam{Name: "Backend Team"})
	require.NoError(t, err)
	require.NotNil(t, team.Identifier)
	assert.Equal(t, "backend_team", *team.Identifier)
	assert.Equal(t, int32(3), searchCalls.Load())
}

func TestEnrichTeamWithRetry_WithOwnershipV2_TeamEntityNeverAppears_GivesUpWithMissingTeamEntityError(t *testing.T) {
	c, searchCalls := newTeamTestClient(t, []string{FeatureFlagUsersAndTeamsV2}, http.StatusOK, `{"ok":true,"entities":[]}`)

	_, err := c.enrichTeamFromTeamEntityWithRetry(context.Background(), &PortTeam{Name: "Backend Team"})
	assert.True(t, errors.Is(err, MissingTeamEntityError), "got %v", err)
	assert.Equal(t, int32(teamEntityRetryAttempts), searchCalls.Load())
}

func TestEnrichTeamWithRetry_WithOwnershipV2_SearchFails_DoesNotRetry(t *testing.T) {
	c, searchCalls := newTeamTestClient(t, []string{FeatureFlagUsersAndTeamsV2}, http.StatusInternalServerError, `{"ok":false,"error":"boom"}`)

	_, err := c.enrichTeamFromTeamEntityWithRetry(context.Background(), &PortTeam{Name: "Backend Team"})
	require.Error(t, err)
	assert.Equal(t, int32(1), searchCalls.Load(), "only a missing team entity is retried")
}

func TestEnrichTeamWithRetry_CancelledContext_ReturnsError(t *testing.T) {
	c, _ := newTeamTestClient(t, []string{}, http.StatusOK, `{"ok":true,"entities":[]}`)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	team, err := c.enrichTeamFromTeamEntityWithRetry(ctx, &PortTeam{Name: "Backend Team"})
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled), "got %v", err)
	assert.Nil(t, team)
}
