package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
)

func TestMergeAppSpec(t *testing.T) {
	tests := []struct {
		name   string
		remote map[string]any
		prior  map[string]any
		want   map[string]any
	}{
		{
			name:   "prior wins for overlapping toggles",
			remote: map[string]any{"liveEventsEnabled": true, "scheduledResyncInterval": "6h"},
			prior:  map[string]any{"liveEventsEnabled": false},
			want:   map[string]any{"liveEventsEnabled": false, "scheduledResyncInterval": "6h"},
		},
		{
			name:   "server defaults fill keys absent from state",
			remote: map[string]any{"liveEventsEnabled": true, "incrementalSyncEnabled": false},
			prior:  nil,
			want:   map[string]any{"liveEventsEnabled": true, "incrementalSyncEnabled": false},
		},
		{
			name:   "prior false is preserved",
			remote: map[string]any{"liveEventsEnabled": true},
			prior:  map[string]any{"liveEventsEnabled": false},
			want:   map[string]any{"liveEventsEnabled": false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mergeAppSpec(tt.remote, tt.prior))
		})
	}
}

func TestMergeSpecPreservesAppSpecToggles(t *testing.T) {
	state := types.StringValue(`{"integrationSpec":{"token":"my-secret"},"appSpec":{"liveEventsEnabled":false}}`)
	remote := &cli.IntegrationClientSpec{
		IntegrationSpec: map[string]any{"token": ""},
		AppSpec:         map[string]any{"liveEventsEnabled": true, "scheduledResyncInterval": "12h"},
	}

	got := mergeSpec(state, remote, false)
	assert.JSONEq(t,
		`{"integrationSpec":{"token":"my-secret"},"appSpec":{"liveEventsEnabled":false,"scheduledResyncInterval":"12h"}}`,
		got.ValueString(),
	)
}
