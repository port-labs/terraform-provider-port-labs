package blueprint

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestArrayPropResourceToBodyUnionFields(t *testing.T) {
	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"vulnerabilities": {
				Title:             types.StringValue("Vulnerabilities"),
				Union:             types.BoolValue(true),
				IncludeDuplicates: types.BoolValue(false),
				StringItems: &StringItems{
					Default: types.ListNull(types.StringType),
					Enum:    types.ListNull(types.StringType),
				},
			},
		},
	}

	props := map[string]cli.BlueprintProperty{}
	required := []string{}

	err := arrayPropResourceToBody(context.Background(), state, props, &required)
	if err != nil {
		t.Fatalf("arrayPropResourceToBody returned error: %v", err)
	}

	prop, ok := props["vulnerabilities"]
	if !ok {
		t.Fatal("expected vulnerabilities property in body")
	}

	if prop.Union == nil || !*prop.Union {
		t.Fatalf("expected union=true, got %#v", prop.Union)
	}

	if prop.IncludeDuplicates == nil || *prop.IncludeDuplicates {
		t.Fatalf("expected includeDuplicates=false, got %#v", prop.IncludeDuplicates)
	}

	if prop.Items == nil || prop.Items["type"] != "string" {
		t.Fatalf("expected string items, got %#v", prop.Items)
	}
}

func TestAddArrayPropertiesToStateUnionFields(t *testing.T) {
	union := true
	includeDuplicates := true

	prop := AddArrayPropertiesToState(context.Background(), &cli.BlueprintProperty{
		Type: "array",
		Items: map[string]any{
			"type": "string",
		},
		Union:             &union,
		IncludeDuplicates: &includeDuplicates,
	}, false)

	if prop.Union.IsNull() || !prop.Union.ValueBool() {
		t.Fatalf("expected union=true in state, got %#v", prop.Union)
	}

	if prop.IncludeDuplicates.IsNull() || !prop.IncludeDuplicates.ValueBool() {
		t.Fatalf("expected include_duplicates=true in state, got %#v", prop.IncludeDuplicates)
	}
}
