package page

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestPageToPortBodyPageDateFilter(t *testing.T) {
	t.Run("sends the configured date filter", func(t *testing.T) {
		body, err := PageToPortBody(&PageModel{
			Identifier:     types.StringValue("dashboard"),
			Type:           types.StringValue("dashboard"),
			PageDateFilter: types.StringValue(`{"range":{"type":"quickRange","quickRange":"30d"}}`),
		})
		if err != nil {
			t.Fatal(err)
		}
		if body.PageDateFilter == nil {
			t.Fatal("expected pageDateFilter to be set")
		}
		rangeBody, ok := (*body.PageDateFilter)["range"].(map[string]any)
		if !ok || rangeBody["type"] != "quickRange" || rangeBody["quickRange"] != "30d" {
			t.Fatalf("unexpected pageDateFilter: %v", *body.PageDateFilter)
		}
	})

	t.Run("omits the date filter when it is not configured", func(t *testing.T) {
		body, err := PageToPortBody(&PageModel{
			Identifier:     types.StringValue("dashboard"),
			Type:           types.StringValue("dashboard"),
			PageDateFilter: types.StringNull(),
		})
		if err != nil {
			t.Fatal(err)
		}
		if body.PageDateFilter != nil {
			t.Fatalf("expected pageDateFilter to be omitted, got %v", *body.PageDateFilter)
		}
	})

	t.Run("fails on invalid JSON", func(t *testing.T) {
		_, err := PageToPortBody(&PageModel{
			Identifier:     types.StringValue("dashboard"),
			Type:           types.StringValue("dashboard"),
			PageDateFilter: types.StringValue(`{not json`),
		})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestRefreshPageToStatePageDateFilter(t *testing.T) {
	r := &PageResource{portClient: &cli.PortClient{}}
	apiFilter := map[string]any{"range": map[string]any{"type": "quickRange", "quickRange": "30d"}}

	t.Run("keeps the configured string when it is semantically equal", func(t *testing.T) {
		configured := `{ "range": { "quickRange": "30d", "type": "quickRange" } }`
		model := &PageModel{PageDateFilter: types.StringValue(configured)}
		if err := r.refreshPageToState(model, &cli.Page{Identifier: "dashboard", Type: "dashboard", PageDateFilter: &apiFilter}); err != nil {
			t.Fatal(err)
		}
		if model.PageDateFilter.ValueString() != configured {
			t.Fatalf("expected the configured string to be kept, got %s", model.PageDateFilter.ValueString())
		}
	})

	t.Run("uses the API value when it differs from the configured one", func(t *testing.T) {
		model := &PageModel{PageDateFilter: types.StringValue(`{"range":{"type":"quickRange","quickRange":"7d"}}`)}
		if err := r.refreshPageToState(model, &cli.Page{Identifier: "dashboard", Type: "dashboard", PageDateFilter: &apiFilter}); err != nil {
			t.Fatal(err)
		}
		if model.PageDateFilter.ValueString() != `{"range":{"quickRange":"30d","type":"quickRange"}}` {
			t.Fatalf("unexpected state: %s", model.PageDateFilter.ValueString())
		}
	})

	t.Run("is null when the page has no date filter", func(t *testing.T) {
		model := &PageModel{PageDateFilter: types.StringValue(`{"range":{"type":"quickRange","quickRange":"7d"}}`)}
		if err := r.refreshPageToState(model, &cli.Page{Identifier: "dashboard", Type: "dashboard"}); err != nil {
			t.Fatal(err)
		}
		if !model.PageDateFilter.IsNull() {
			t.Fatalf("expected null, got %s", model.PageDateFilter.ValueString())
		}
	})
}
