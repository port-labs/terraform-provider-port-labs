package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestPlanSpec(t *testing.T) {
	state := types.StringValue(`{"integrationSpec":{"githubToken":"_OLD"},"appSpec":{"liveEventsEnabled":false}}`)

	tests := []struct {
		name   string
		config string
		want   string
	}{
		{
			name:   "inherits appSpec when only integrationSpec changes",
			config: `{"integrationSpec":{"githubToken":"_NEW"}}`,
			want:   `{"appSpec":{"liveEventsEnabled":false},"integrationSpec":{"githubToken":"_NEW"}}`,
		},
		{
			name:   "keeps state when the configuration is unchanged",
			config: `{"integrationSpec":{"githubToken":"_OLD"}}`,
			want:   state.ValueString(),
		},
		{
			name:   "configuration wins when it manages appSpec",
			config: `{"integrationSpec":{"githubToken":"_OLD"},"appSpec":{"liveEventsEnabled":true}}`,
			want:   `{"appSpec":{"liveEventsEnabled":true},"integrationSpec":{"githubToken":"_OLD"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := planSpec(types.StringValue(tt.config), state)
			assert.JSONEq(t, tt.want, got.ValueString())
		})
	}
}
