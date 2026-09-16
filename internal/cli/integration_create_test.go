package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_CreatePortResourcesOriginJSON(t *testing.T) {
	origin := "Empty"
	body := Integration{
		InstallationId:            "github-prod",
		CreatePortResourcesOrigin: &origin,
	}

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"createPortResourcesOrigin":"Empty"`)
}

func TestIntegration_CreatePortResourcesOriginOmittedFromJSON(t *testing.T) {
	body := Integration{
		InstallationId: "github-prod",
	}

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "createPortResourcesOrigin")
}
