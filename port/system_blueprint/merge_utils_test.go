package system_blueprint

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/blueprint"
)

func ptr(s string) *string { return &s }

func TestRelationIsRuleResultTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    cli.Relation
		want bool
	}{
		{name: "nil type", r: cli.Relation{Type: nil}, want: false},
		{name: "empty type string", r: cli.Relation{Type: ptr("")}, want: false},
		{name: "other type", r: cli.Relation{Type: ptr("other")}, want: false},
		{name: "rule_result_target", r: cli.Relation{Type: ptr("rule_result_target")}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := RelationIsRuleResultTarget(tt.r); got != tt.want {
				t.Fatalf("RelationIsRuleResultTarget(%#v) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestMergeRelationsForSystemBlueprint_nilLive(t *testing.T) {
	structure := map[string]cli.Relation{
		"rule": {Title: ptr("Rule"), Target: ptr("_rule")},
	}
	state := map[string]blueprint.RelationModel{
		"vm": {Target: types.StringValue("microservice"), Title: types.StringValue("VM")},
	}
	got := MergeRelationsForSystemBlueprint("_team", nil, structure, state)
	if got["rule"].Target == nil || *got["rule"].Target != "_rule" {
		t.Fatalf("rule from structure: %#v", got["rule"])
	}
	if got["vm"].Target == nil || *got["vm"].Target != "microservice" {
		t.Fatalf("vm from state: %#v", got["vm"])
	}
}

func TestMergeRelationsForSystemBlueprint_stateAndLiveRuleResultTargets(t *testing.T) {
	live := map[string]cli.Relation{
		"_Env": {
			Title:  ptr("Env"),
			Target: ptr("Env"),
			Type:   ptr("rule_result_target"),
		},
	}
	structure := map[string]cli.Relation{
		"rule": {Title: ptr("Rule"), Target: ptr("_rule")},
	}
	state := map[string]blueprint.RelationModel{
		"vm": {Target: types.StringValue("microservice"), Title: types.StringValue("VM")},
	}
	got := MergeRelationsForSystemBlueprint("_rule_result", live, structure, state)
	if got["_Env"].Type == nil || *got["_Env"].Type != "rule_result_target" {
		t.Fatalf("expected _Env from live, got %#v", got["_Env"])
	}
	if got["vm"].Target == nil || *got["vm"].Target != "microservice" {
		t.Fatalf("expected vm from state, got %#v", got["vm"])
	}
	if got["rule"].Target == nil || *got["rule"].Target != "_rule" {
		t.Fatalf("expected rule from structure, got %#v", got["rule"])
	}
}

func TestMergeRelationsForSystemBlueprint_nonRuleResultBlueprintIgnoresLiveRuleResultTarget(t *testing.T) {
	live := map[string]cli.Relation{
		"_Env": {
			Title:  ptr("Env"),
			Target: ptr("Env"),
			Type:   ptr("rule_result_target"),
		},
	}
	structure := map[string]cli.Relation{
		"rule": {Title: ptr("Rule"), Target: ptr("_rule")},
	}

	got := MergeRelationsForSystemBlueprint("_team", live, structure, nil)
	if _, ok := got["_Env"]; ok {
		t.Fatalf("expected _Env to be ignored for non _rule_result blueprint, got %#v", got["_Env"])
	}
}

func TestMergeRelationsForSystemBlueprint_preservesRuleResultTargetRelations(t *testing.T) {
	live := map[string]cli.Relation{
		"rule": {
			Title:  ptr("Rule"),
			Target: ptr("_rule"),
		},
		"service_0b8f289d-2d70-4caa-856b-0af1d5e24f3a": {
			Title:  ptr("Service"),
			Target: ptr("service"),
			Type:   ptr("rule_result_target"),
		},
		"artifact": {
			Title:  ptr("Artifact"),
			Target: ptr("artifact"),
			Type:   ptr("rule_result_target"),
		},
	}
	structure := map[string]cli.Relation{
		"rule": {
			Title:  ptr("Rule"),
			Target: ptr("_rule"),
		},
	}
	got := MergeRelationsForSystemBlueprint("_rule_result", live, structure, nil)
	if got["rule"].Target == nil || *got["rule"].Target != "_rule" {
		t.Fatalf("expected rule from structure, got %#v", got["rule"])
	}
	for _, key := range []string{"service_0b8f289d-2d70-4caa-856b-0af1d5e24f3a", "artifact"} {
		if _, ok := got[key]; !ok {
			t.Fatalf("expected live relation %q to be preserved, got keys: %v", key, keysOf(got))
		}
		if got[key].Type == nil || *got[key].Type != "rule_result_target" {
			t.Fatalf("expected type rule_result_target on preserved relation %q", key)
		}
	}
}

func TestMergeRelationsForSystemBlueprint_liveWinsForRuleResultTargetType(t *testing.T) {
	live := map[string]cli.Relation{
		"svc_custom_key": {
			Target: ptr("from-live"),
			Type:   ptr("rule_result_target"),
		},
	}
	structure := map[string]cli.Relation{
		"svc_custom_key": {
			Target: ptr("from-structure"),
		},
	}
	got := MergeRelationsForSystemBlueprint("_rule_result", live, structure, nil)
	if got["svc_custom_key"].Target == nil || *got["svc_custom_key"].Target != "from-live" {
		t.Fatalf("expected live relation to overwrite merged for rule_result_target type, got %#v", got["svc_custom_key"])
	}
}

func TestMergeRelationsForSystemBlueprint_doesNotPreserveUntypedUUIDKey(t *testing.T) {
	live := map[string]cli.Relation{
		"svc_11111111-1111-4111-8111-111111111111": {
			Target: ptr("from-live"),
		},
	}
	structure := map[string]cli.Relation{
		"svc_11111111-1111-4111-8111-111111111111": {
			Target: ptr("from-structure"),
		},
	}
	got := MergeRelationsForSystemBlueprint("_rule_result", live, structure, nil)
	if got["svc_11111111-1111-4111-8111-111111111111"].Target == nil || *got["svc_11111111-1111-4111-8111-111111111111"].Target != "from-structure" {
		t.Fatalf("without type rule_result_target, structure merge should win, got %#v", got["svc_11111111-1111-4111-8111-111111111111"])
	}
}

func keysOf(m map[string]cli.Relation) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
