package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestMergeSpecSectionExplicitAppSpec(t *testing.T) {
	remote := map[string]any{
		"liveEventsEnabled":        true,
		"incrementalSyncEnabled":   false,
		"scheduledResyncInterval":  "12h",
		"liveEventsUuid":           "abc",
		"liveEventsIngestHostname": "ingest.example.com",
	}
	prior := map[string]any{
		"liveEventsEnabled":       false,
		"scheduledResyncInterval": "6h",
	}

	got := mergeSpecSectionExplicit(remote, prior, consts.IsServerManagedAppSpecKey)

	assert.Equal(t, map[string]any{
		"liveEventsEnabled":       false,
		"scheduledResyncInterval": "6h",
	}, got)
	assert.NotContains(t, got, "incrementalSyncEnabled")
	assert.NotContains(t, got, "liveEventsUuid")
	assert.NotContains(t, got, "liveEventsIngestHostname")
}

func TestMergeSpecSectionImplicitAppSpec(t *testing.T) {
	remote := map[string]any{
		"liveEventsEnabled":        true,
		"incrementalSyncEnabled":   false,
		"liveEventsUuid":           "abc",
		"liveEventsIngestHostname": "ingest.example.com",
	}

	got := mergeSpecSectionImplicit(remote, consts.IsServerManagedAppSpecKey)

	assert.Equal(t, map[string]any{
		"liveEventsEnabled":      true,
		"incrementalSyncEnabled": false,
	}, got)
}

func TestMergeSpecSectionExplicitPreservesPriorFalse(t *testing.T) {
	remote := map[string]any{"liveEventsEnabled": true}
	prior := map[string]any{"liveEventsEnabled": false}

	assert.Equal(t, map[string]any{"liveEventsEnabled": false},
		mergeSpecSectionExplicit(remote, prior, consts.IsServerManagedAppSpecKey))
}

func TestMergeSpecSectionExplicitIntegrationSpec(t *testing.T) {
	remote := map[string]any{
		"githubToken":        "",
		"authenticationMode": "GitHub App",
	}
	prior := map[string]any{"githubToken": "my-secret-ref"}

	got := mergeSpecSectionExplicit(remote, prior, nil)

	assert.Equal(t, map[string]any{"githubToken": "my-secret-ref"}, got)
	assert.NotContains(t, got, "authenticationMode")
}

func TestMergeSpecSectionImplicitIntegrationSpec(t *testing.T) {
	remote := map[string]any{
		"githubToken":        "resolved-token",
		"authenticationMode": "Personal Access Token",
	}

	assert.Equal(t, remote, mergeSpecSectionImplicit(remote, nil))
}

func TestMergeSpecPreservesAppSpecToggles(t *testing.T) {
	state := types.StringValue(`{"integrationSpec":{"token":"my-secret"},"appSpec":{"liveEventsEnabled":false}}`)
	remote := &cli.IntegrationClientSpec{
		IntegrationSpec: map[string]any{"token": ""},
		AppSpec: map[string]any{
			"liveEventsEnabled":       true,
			"scheduledResyncInterval": "12h",
			"liveEventsUuid":          "abc",
		},
	}

	got := mergeSpec(state, remote, false)
	assert.JSONEq(t,
		`{"integrationSpec":{"token":"my-secret"},"appSpec":{"liveEventsEnabled":false}}`,
		got.ValueString(),
	)
}

func TestMergeSpecImplicitAppSpec(t *testing.T) {
	state := types.StringValue(`{"integrationSpec":{"token":"my-secret"}}`)
	remote := &cli.IntegrationClientSpec{
		IntegrationSpec: map[string]any{"token": ""},
		AppSpec: map[string]any{
			"liveEventsEnabled":        true,
			"liveEventsUuid":           "abc",
			"liveEventsIngestHostname": "ingest.example.com",
		},
	}

	got := mergeSpec(state, remote, false)
	assert.JSONEq(t,
		`{"integrationSpec":{"token":"my-secret"},"appSpec":{"liveEventsEnabled":true}}`,
		got.ValueString(),
	)
}

func TestRefreshIntegrationState_SaasOAuth2StripsIntegrationSpec(t *testing.T) {
	r := &IntegrationResource{portClient: &cli.PortClient{}}
	state := &IntegrationModel{
		InstallationId: types.StringValue("github-oauth"),
	}
	installationType := consts.InstallationTypeSaasOAuth2
	remote := &cli.Integration{
		InstallationId:   "github-oauth",
		InstallationType: &installationType,
		Spec: &cli.IntegrationClientSpec{
			IntegrationSpec: map[string]any{
				"githubAppPrivateKey": "oauth-secret-name",
				"githubAppId":         "123",
			},
			AppSpec: map[string]any{
				"scheduledResyncInterval":  "12h",
				"liveEventsUuid":           "abc",
				"liveEventsIngestHostname": "ingest.example.com",
			},
		},
	}

	r.refreshIntegrationState(state, remote, "github-oauth")

	assert.Equal(t, consts.InstallationTypeSaasOAuth2, state.InstallationType.ValueString())
	assert.JSONEq(t, `{"appSpec":{"scheduledResyncInterval":"12h"}}`, state.Spec.ValueString())
}

func TestRefreshIntegrationState_SaasOAuth2PreservesPriorAppSpec(t *testing.T) {
	r := &IntegrationResource{portClient: &cli.PortClient{}}
	state := &IntegrationModel{
		InstallationId: types.StringValue("github-oauth"),
		Spec:           types.StringValue(`{"appSpec":{"scheduledResyncInterval":"6h","liveEventsEnabled":false}}`),
	}
	installationType := consts.InstallationTypeSaasOAuth2
	remote := &cli.Integration{
		InstallationId:   "github-oauth",
		InstallationType: &installationType,
		Spec: &cli.IntegrationClientSpec{
			IntegrationSpec: map[string]any{"githubAppPrivateKey": "oauth-secret-name"},
			AppSpec: map[string]any{
				"scheduledResyncInterval": "12h",
				"liveEventsEnabled":       true,
			},
		},
	}

	r.refreshIntegrationState(state, remote, "github-oauth")

	assert.JSONEq(t,
		`{"appSpec":{"scheduledResyncInterval":"6h","liveEventsEnabled":false}}`,
		state.Spec.ValueString(),
	)
}
