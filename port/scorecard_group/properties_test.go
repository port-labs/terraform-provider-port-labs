package scorecard_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPropertiesToCLI(t *testing.T) {
	t.Parallel()

	properties, diags := types.MapValueFrom(context.Background(), types.StringType, map[string]string{
		"test_string": "shared-value",
		"enabled":     "true",
	})
	if diags.HasError() {
		t.Fatalf("failed to build properties map: %v", diags)
	}

	got, err := propertiesToCLI(context.Background(), properties)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["test_string"] != "shared-value" {
		t.Fatalf("unexpected test_string value: %v", got["test_string"])
	}
	if got["enabled"] != true {
		t.Fatalf("expected enabled to be true, got %v", got["enabled"])
	}
}

func TestPropertiesFromCLI(t *testing.T) {
	t.Parallel()

	got := propertiesFromCLI(map[string]any{
		"test_string": "shared-value",
		"test":        "https://example.com",
	}, false)

	elements := got.Elements()
	if len(elements) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(elements))
	}
	if elements["test_string"].Equal(types.StringValue("shared-value")) != true {
		t.Fatalf("unexpected test_string value: %v", elements["test_string"])
	}
	if elements["test"].Equal(types.StringValue("https://example.com")) != true {
		t.Fatalf("unexpected test value: %v", elements["test"])
	}
}
