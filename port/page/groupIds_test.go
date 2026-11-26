package page

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

const (
	// Valid ids have the shape Port expects: "grp_" followed by 16 letters or numbers.
	overviewID = "grp_overview00000001"
	runsID     = "grp_runs000000000001"
)

func TestRefreshPageToStateAlignsEntityPageGroupIds(t *testing.T) {
	tests := []struct {
		name     string
		pageType string
		prior    []string // nil means there is no prior state (import)
		api      string
		want     string
	}{
		{
			name:     "removes ids that the configuration doesn't declare",
			pageType: "entity",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Overview","widgets":[]},{"title":"Runs","widgets":[]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]},{"id":"` + runsID + `","title":"Runs","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Overview","widgets":[]},{"title":"Runs","widgets":[]}]}]`,
		},
		{
			name:     "keeps ids that the configuration declares",
			pageType: "entity",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]},{"title":"Runs","widgets":[]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]},{"id":"` + runsID + `","title":"Runs","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]},{"title":"Runs","widgets":[]}]}]`,
		},
		{
			name:     "shows a changed id as drift",
			pageType: "entity",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + runsID + `","title":"Overview","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + runsID + `","title":"Overview","widgets":[]}]}]`,
		},
		{
			name:     "keeps declared ids while the API hides them",
			pageType: "entity",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Overview","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}]`,
		},
		{
			name:     "keeps the id of a tab that is not in the prior state",
			pageType: "entity",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Overview","widgets":[]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]},{"id":"` + runsID + `","title":"Added in the UI","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Overview","widgets":[]},{"id":"` + runsID + `","title":"Added in the UI","widgets":[]}]}]`,
		},
		{
			name:     "only changes the main grouper",
			pageType: "entity",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Related Entities","widgets":[{"id":"n","type":"grouper","displayMode":"switch","groups":[{"title":"Table","widgets":[]}]}]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Related Entities","widgets":[{"id":"n","type":"grouper","displayMode":"switch","groups":[{"id":"` + runsID + `","title":"Table","widgets":[]}]}]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Related Entities","widgets":[{"id":"n","type":"grouper","displayMode":"switch","groups":[{"id":"` + runsID + `","title":"Table","widgets":[]}]}]}]}]`,
		},
		{
			name:     "leaves hiddenGroups as returned",
			pageType: "entity",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Overview","widgets":[]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","hiddenGroups":["` + runsID + `"],"groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","hiddenGroups":["` + runsID + `"],"groups":[{"title":"Overview","widgets":[]}]}]`,
		},
		{
			name:     "leaves pages that are not entity pages as returned",
			pageType: "dashboard",
			prior:    []string{`{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"title":"Overview","widgets":[]}]}`},
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}]`,
		},
		{
			name:     "leaves the widgets as returned when there is no prior state",
			pageType: "entity",
			prior:    nil,
			api:      `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}]`,
			want:     `[{"id":"g","type":"grouper","displayMode":"tabs","groups":[{"id":"` + overviewID + `","title":"Overview","widgets":[]}]}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PageResource{portClient: &cli.PortClient{JSONEscapeHTML: true}}

			pm := &PageModel{Widgets: types.ListNull(types.StringType)}
			if tt.prior != nil {
				priorAttrs := make([]attr.Value, len(tt.prior))
				for i, widget := range tt.prior {
					priorAttrs[i] = canonicalWidget(t, decodeJSON[map[string]any](t, widget))
				}
				pm.Widgets = types.ListValueMust(types.StringType, priorAttrs)
			}

			apiWidgets := decodeJSON[[]map[string]any](t, tt.api)
			page := &cli.Page{Identifier: "microserviceEntity", Type: tt.pageType, Widgets: &apiWidgets}

			if err := r.refreshPageToState(pm, page); err != nil {
				t.Fatalf("refreshPageToState returned an error: %v", err)
			}

			wantWidgets := decodeJSON[[]map[string]any](t, tt.want)
			got := pm.Widgets.Elements()
			if len(got) != len(wantWidgets) {
				t.Fatalf("got %d widgets, want %d", len(got), len(wantWidgets))
			}
			for i, want := range wantWidgets {
				if got[i].(types.String).ValueString() != canonicalWidget(t, want).ValueString() {
					t.Errorf("widget %d:\n got: %s\nwant: %s", i, got[i].(types.String).ValueString(), canonicalWidget(t, want).ValueString())
				}
			}
		})
	}
}

func decodeJSON[T any](t *testing.T, s string) T {
	t.Helper()
	var v T
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Fatalf("invalid test JSON %q: %v", s, err)
	}
	return v
}

// canonicalWidget encodes a widget the same way the provider writes it to the state.
func canonicalWidget(t *testing.T, widget map[string]any) types.String {
	t.Helper()
	s, err := utils.GoObjectToTerraformString(widget, true)
	if err != nil {
		t.Fatalf("failed to encode widget: %v", err)
	}
	return s
}
