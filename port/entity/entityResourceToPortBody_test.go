package entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestEntityResourceToBodyIncludesMappingWarnings(t *testing.T) {
	state := &EntityModel{
		Title:     types.StringValue("Service"),
		Blueprint: types.StringValue("service"),
		MappingWarnings: []types.String{
			types.StringValue("properties"),
			types.StringValue("team"),
		},
	}

	entity, err := entityResourceToBody(context.Background(), state, &cli.Blueprint{Identifier: "service"})
	if err != nil {
		t.Fatalf("entityResourceToBody returned error: %v", err)
	}

	want := []string{"properties", "team"}
	if len(entity.MappingWarnings) != len(want) {
		t.Fatalf("MappingWarnings = %#v, want %#v", entity.MappingWarnings, want)
	}
	for i, warning := range want {
		if entity.MappingWarnings[i] != warning {
			t.Fatalf("MappingWarnings[%d] = %q, want %q", i, entity.MappingWarnings[i], warning)
		}
	}
}
