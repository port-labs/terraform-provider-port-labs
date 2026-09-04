package entity

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestWriteEntityComputedFieldsToStateIncludesMappingWarnings(t *testing.T) {
	now := time.Now()
	state := &EntityModel{}
	entity := &cli.Entity{
		Blueprint:  "service",
		Identifier: "svc-1",
		Meta: cli.Meta{
			CreatedAt: &now,
			UpdatedAt: &now,
			CreatedBy: "user",
			UpdatedBy: "user",
		},
		MappingWarnings: []string{"properties", "team"},
	}

	writeEntityComputedFieldsToState(state, entity)

	want := []string{"properties", "team"}
	if len(state.MappingWarnings) != len(want) {
		t.Fatalf("MappingWarnings = %#v, want %#v", state.MappingWarnings, want)
	}
	for i, warning := range want {
		if state.MappingWarnings[i].ValueString() != warning {
			t.Fatalf("MappingWarnings[%d] = %q, want %q", i, state.MappingWarnings[i].ValueString(), warning)
		}
	}
}

func TestWriteEntityComputedFieldsToStateClearsMappingWarningsWhenEmpty(t *testing.T) {
	now := time.Now()
	state := &EntityModel{
		MappingWarnings: []types.String{types.StringValue("properties")},
	}
	entity := &cli.Entity{
		Blueprint:  "service",
		Identifier: "svc-1",
		Meta: cli.Meta{
			CreatedAt: &now,
			UpdatedAt: &now,
			CreatedBy: "user",
			UpdatedBy: "user",
		},
	}

	writeEntityComputedFieldsToState(state, entity)

	if state.MappingWarnings != nil {
		t.Fatalf("expected mapping_warnings to be nil when API returns empty, got: %#v", state.MappingWarnings)
	}
}
