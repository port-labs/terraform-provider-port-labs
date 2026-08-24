package blueprint

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestAddArrayPropertiesToState_unionFields(t *testing.T) {
	t.Parallel()

	union := true
	includeDuplicates := true
	prop := &cli.BlueprintProperty{
		Type: "array",
		Items: map[string]any{
			"type": "string",
		},
		Union:             &union,
		IncludeDuplicates: &includeDuplicates,
	}

	arrayProp := AddArrayPropertiesToState(context.Background(), prop, false, nil)

	if !arrayProp.Union.ValueBool() {
		t.Fatal("expected union to be true")
	}
	if !arrayProp.IncludeDuplicates.ValueBool() {
		t.Fatal("expected include_duplicates to be true")
	}
}

func TestBlueprintProperty_unionMarshalJSON(t *testing.T) {
	t.Parallel()

	union := true
	includeDuplicates := false
	prop := cli.BlueprintProperty{
		Type: "array",
		Items: map[string]any{
			"type": "number",
		},
		Union:             &union,
		IncludeDuplicates: &includeDuplicates,
	}

	data, err := json.Marshal(prop)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if result["union"] != true {
		t.Fatalf("expected union true, got %v", result["union"])
	}
	if result["includeDuplicates"] != false {
		t.Fatalf("expected includeDuplicates false, got %v", result["includeDuplicates"])
	}
}

func TestArrayPropResourceToBody_unionFields(t *testing.T) {
	t.Parallel()

	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"tags": {
				Union:             types.BoolValue(true),
				IncludeDuplicates: types.BoolValue(true),
				StringItems:       &StringItems{},
			},
		},
	}
	props := map[string]cli.BlueprintProperty{}

	if err := arrayPropResourceToBody(context.Background(), state, props, &[]string{}); err != nil {
		t.Fatalf("arrayPropResourceToBody: %v", err)
	}

	prop, ok := props["tags"]
	if !ok {
		t.Fatal("expected tags property")
	}
	if prop.Union == nil || !*prop.Union {
		t.Fatal("expected union to be true")
	}
	if prop.IncludeDuplicates == nil || !*prop.IncludeDuplicates {
		t.Fatal("expected include_duplicates to be true")
	}
}

func TestArrayPropResourceToBody_omitsUnionWhenFalse(t *testing.T) {
	t.Parallel()

	state := &PropertiesModel{
		ArrayProps: map[string]ArrayPropModel{
			"tags": {
				Union:       types.BoolValue(false),
				StringItems: &StringItems{},
			},
		},
	}
	props := map[string]cli.BlueprintProperty{}

	if err := arrayPropResourceToBody(context.Background(), state, props, &[]string{}); err != nil {
		t.Fatalf("arrayPropResourceToBody: %v", err)
	}

	prop := props["tags"]
	if prop.Union != nil {
		t.Fatal("expected union to be omitted when false")
	}
	if prop.IncludeDuplicates != nil {
		t.Fatal("expected include_duplicates to be omitted when union is false")
	}
}
