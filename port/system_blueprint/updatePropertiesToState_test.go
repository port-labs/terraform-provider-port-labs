package system_blueprint

import (
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestAddRelationsToState_skipsRuleResultTargetRelations(t *testing.T) {
	envTarget := "Env"
	vmTarget := "microservice"
	b := &cli.Blueprint{
		Identifier: "_rule_result",
		Relations: map[string]cli.Relation{
			"_Env": {
				Title:  ptr("Env"),
				Target: &envTarget,
				Type:   ptr("rule_result_target"),
			},
			"vm": {
				Title:  ptr("VM"),
				Target: &vmTarget,
			},
		},
	}
	systemBp := &cli.Blueprint{Relations: map[string]cli.Relation{}}
	bm := &SystemBlueprintModel{}

	addRelationsToState(b, systemBp, bm)

	if _, ok := bm.Relations["_Env"]; ok {
		t.Fatal("rule_result_target relation _Env must not be written to Terraform state")
	}
	rm, ok := bm.Relations["vm"]
	if !ok {
		t.Fatal("expected vm relation in state")
	}
	if rm.Target.ValueString() != "microservice" {
		t.Fatalf("vm target: got %q", rm.Target.ValueString())
	}
}

func TestAddRelationsToState_nonRuleResultBlueprintKeepsRuleResultTargetRelations(t *testing.T) {
	envTarget := "Env"
	b := &cli.Blueprint{
		Identifier: "_team",
		Relations: map[string]cli.Relation{
			"_Env": {
				Title:  ptr("Env"),
				Target: &envTarget,
				Type:   ptr("rule_result_target"),
			},
		},
	}
	systemBp := &cli.Blueprint{Relations: map[string]cli.Relation{}}
	bm := &SystemBlueprintModel{}

	addRelationsToState(b, systemBp, bm)

	if _, ok := bm.Relations["_Env"]; !ok {
		t.Fatal("rule_result_target relation should be kept for non _rule_result blueprints")
	}
}
