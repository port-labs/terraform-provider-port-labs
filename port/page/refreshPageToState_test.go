package page

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveWidgetState(t *testing.T) {
	tests := []struct {
		name         string
		prevWidget   types.String
		remoteWidget string
		keepPrev     bool
	}{
		{
			name:         "API injects empty hiddenQuery when config omits it",
			prevWidget:   types.StringValue(`{"id":"details","type":"entity-info"}`),
			remoteWidget: `{"id":"details","type":"entity-info","hiddenQuery":[]}`,
			keepPrev:     true,
		},
		{
			name:         "config sets empty hiddenQuery and API returns it",
			prevWidget:   types.StringValue(`{"id":"details","hiddenQuery":[]}`),
			remoteWidget: `{"id":"details","hiddenQuery":[]}`,
			keepPrev:     true,
		},
		{
			name:         "non-empty hiddenQuery returned in different order",
			prevWidget:   types.StringValue(`{"id":"details","hiddenQuery":["url","status","author"]}`),
			remoteWidget: `{"id":"details","hiddenQuery":["author","status","url"]}`,
			keepPrev:     true,
		},
		{
			name:         "nested hiddenQuery inside grouper is order-insensitive",
			prevWidget:   types.StringValue(`{"type":"grouper","groups":[{"widgets":[{"hiddenQuery":["b","a"]}]}]}`),
			remoteWidget: `{"type":"grouper","groups":[{"widgets":[{"hiddenQuery":["a","b"]}]}]}`,
			keepPrev:     true,
		},
		{
			name:         "hiddenQuery string and number stay distinct",
			prevWidget:   types.StringValue(`{"id":"details","hiddenQuery":["1",1]}`),
			remoteWidget: `{"id":"details","hiddenQuery":[1,"1"]}`,
			keepPrev:     true,
		},
		{
			name:         "numbers, HTML characters and key order do not cause drift",
			prevWidget:   types.StringValue(`{"title":"a<b>&c","layout":[{"height":400}],"id":"details"}`),
			remoteWidget: `{"id":"details","layout":[{"height":400}],"title":"a<b>&c"}`,
			keepPrev:     true,
		},
		{
			name:         "hiddenQuery cleared remotely is drift",
			prevWidget:   types.StringValue(`{"id":"details","hiddenQuery":["url"]}`),
			remoteWidget: `{"id":"details","hiddenQuery":[]}`,
			keepPrev:     false,
		},
		{
			name:         "hiddenQuery with different elements is drift",
			prevWidget:   types.StringValue(`{"id":"details","hiddenQuery":["url","status"]}`),
			remoteWidget: `{"id":"details","hiddenQuery":["url","author"]}`,
			keepPrev:     false,
		},
		{
			name:         "reordering an ordered array is drift",
			prevWidget:   types.StringValue(`{"id":"board","propertiesSettings":{"order":["$title","status"]}}`),
			remoteWidget: `{"id":"board","propertiesSettings":{"order":["status","$title"]}}`,
			keepPrev:     false,
		},
		{
			name:         "other empty arrays are not ignored",
			prevWidget:   types.StringValue(`{"id":"graph"}`),
			remoteWidget: `{"id":"graph","hiddenBlueprints":[]}`,
			keepPrev:     false,
		},
		{
			name:         "null prev widget returns remote widget",
			prevWidget:   types.StringNull(),
			remoteWidget: `{"id":"details","hiddenQuery":[]}`,
			keepPrev:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var remote map[string]any
			require.NoError(t, json.Unmarshal([]byte(tt.remoteWidget), &remote))

			got, err := resolveWidgetState(tt.prevWidget, remote, false)
			require.NoError(t, err)

			if tt.keepPrev {
				assert.Equal(t, tt.prevWidget, got)
			} else {
				assert.JSONEq(t, tt.remoteWidget, got.ValueString())
			}
		})
	}
}
