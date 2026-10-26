package page

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
)

func TestRefreshPageWidgets(t *testing.T) {
	widget := map[string]any{"id": "table", "type": "table-entities-explorer"}
	testCases := []struct {
		name        string
		stateNull   bool
		apiWidgets  *[]map[string]any
		wantNull    bool
		wantEntries int
	}{
		{name: "no widgets declared and none in Port stays undeclared", stateNull: true, apiWidgets: &[]map[string]any{}, wantNull: true},
		{name: "no widgets declared but some in Port shows them", stateNull: true, apiWidgets: &[]map[string]any{widget}, wantEntries: 1},
		{name: "declared widgets removed in Port shows the removal", stateNull: false, apiWidgets: &[]map[string]any{}, wantEntries: 0},
		{name: "missing widgets field is null", stateNull: false, apiWidgets: nil, wantNull: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &PageResource{portClient: &cli.PortClient{}}
			state := &PageModel{Widgets: types.ListNull(types.StringType)}
			if !tc.stateNull {
				state.Widgets = types.ListValueMust(types.StringType, []attr.Value{types.StringValue(`{"id":"table"}`)})
			}

			err := r.refreshPageToState(state, &cli.Page{Identifier: "services", Widgets: tc.apiWidgets})

			assert.NoError(t, err)
			assert.Equal(t, tc.wantNull, state.Widgets.IsNull())
			if !tc.wantNull {
				assert.Len(t, state.Widgets.Elements(), tc.wantEntries)
			}
		})
	}
}
