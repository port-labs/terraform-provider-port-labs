package entity

import (
	"context"
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestRefreshArrayEntityStateSkipsUnionProperties(t *testing.T) {
	union := true
	state := &EntityModel{
		Properties: &EntityPropertiesModel{},
	}
	blueprint := &cli.Blueprint{
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"vulnerabilities": {
					Type:  "array",
					Union: &union,
					Items: map[string]any{"type": "string"},
				},
				"tags": {
					Type:  "array",
					Items: map[string]any{"type": "string"},
				},
			},
		},
	}

	resource := &EntityResource{}
	resource.refreshArrayEntityState(context.Background(), state, map[string][]interface{}{
		"vulnerabilities": {"CVE-1", "CVE-2"},
		"tags":            {"alpha"},
	}, blueprint)

	if !state.Properties.ArrayProps.StringItems.IsNull() {
		elements := state.Properties.ArrayProps.StringItems.Elements()
		if _, ok := elements["vulnerabilities"]; ok {
			t.Fatalf("expected union property to be skipped on read, got %#v", elements)
		}
		if _, ok := elements["tags"]; !ok {
			t.Fatalf("expected non-union property to be populated, got %#v", elements)
		}
	}
}

func TestIsUnionArrayPropertyUsesUnknownFieldsFallback(t *testing.T) {
	prop := cli.BlueprintProperty{
		UnknownFields: map[string]any{
			"union": true,
		},
	}

	if !isUnionArrayProperty(prop) {
		t.Fatal("expected union property from unknown fields")
	}
}
