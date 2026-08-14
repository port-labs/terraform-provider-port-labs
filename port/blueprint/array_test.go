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
					Enum: types.ListNull(types.StringType),
				},
			},
		},
	}

	props := make(map[string]cli.BlueprintProperty)
	var required []string
	err := arrayPropResourceToBody(context.Background(), state, props, &required)
	if err != nil {
		t.Fatalf("arrayPropResourceToBody returned error: %v", err)
	}

	prop, ok := props["vulnerabilities"]
	if !ok {
		t.Fatal("expected vulnerabilities property to be set")
	}

	if prop.Union == nil || !*prop.Union {
		t.Fatalf("expected union to be true, got %#v", prop.Union)
	}

	if prop.IncludeDuplicates == nil || *prop.IncludeDuplicates {
		t.Fatalf("expected includeDuplicates to be false, got %#v", prop.IncludeDuplicates)
	}
}

func TestAddArrayPropertiesToStateUnionFields(t *testing.T) {
	union := true
	includeDuplicates := true
	prop := AddArrayPropertiesToState(context.Background(), &cli.BlueprintProperty{
		Type:              "array",
		Union:             &union,
		IncludeDuplicates: &includeDuplicates,
		Items: map[string]any{
			"type": "string",
		},
	}, false)

	if prop.Union.IsNull() || !prop.Union.ValueBool() {
		t.Fatalf("expected union to be true, got %#v", prop.Union)
	}

	if prop.IncludeDuplicates.IsNull() || !prop.IncludeDuplicates.ValueBool() {
		t.Fatalf("expected include_duplicates to be true, got %#v", prop.IncludeDuplicates)
	}
}
