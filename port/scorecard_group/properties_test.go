package scorecard_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestPropertiesToCLI(t *testing.T) {
	properties := types.DynamicValue(types.MapValueMust(types.DynamicType, map[string]attr.Value{
		"owner": types.DynamicValue(types.StringValue("platform-team")),
		"score": types.DynamicValue(types.Float64Value(42)),
	}))

	result, err := propertiesToCLI(context.Background(), properties)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["owner"] != "platform-team" {
		t.Fatalf("expected owner platform-team, got %v", result["owner"])
	}
	if result["score"] != float64(42) {
		t.Fatalf("expected score 42, got %v", result["score"])
	}
}

func TestPropertiesToCLIRejectsReservedKeys(t *testing.T) {
	properties := types.DynamicValue(types.MapValueMust(types.DynamicType, map[string]attr.Value{
		"blueprint": types.DynamicValue(types.StringValue("svc")),
	}))

	_, err := propertiesToCLI(context.Background(), properties)
	if err == nil {
		t.Fatalf("expected reserved key validation error")
	}
}

func TestPropertiesFromCLI(t *testing.T) {
	properties := map[string]any{
		"owner": "platform-team",
		"score": float64(42),
	}

	result, diags := propertiesFromCLI(context.Background(), properties)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	underlying, ok := result.UnderlyingValue().(types.Map)
	if !ok {
		t.Fatalf("expected properties map, got %T", result.UnderlyingValue())
	}

	owner, ok := underlying.Elements()["owner"].(types.Dynamic)
	if !ok {
		t.Fatalf("unexpected owner value type: %T", underlying.Elements()["owner"])
	}
	ownerValue, ok := owner.UnderlyingValue().(types.String)
	if !ok || ownerValue.ValueString() != "platform-team" {
		t.Fatalf("unexpected owner value: %v", owner.UnderlyingValue())
	}
}

func TestScorecardGroupResourceToPortBodyWithProperties(t *testing.T) {
	properties := types.DynamicValue(types.MapValueMust(types.DynamicType, map[string]attr.Value{
		"owner": types.DynamicValue(types.StringValue("platform-team")),
	}))

	state := &ScorecardGroupModel{
		Identifier: types.StringValue("group-1"),
		Title:      types.StringValue("Group 1"),
		Properties: properties,
		Blueprints: []types.String{types.StringValue("bp-1")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"$team","operator":"isNotEmpty"}`),
					},
				},
			},
		},
	}

	group, err := scorecardGroupResourceToPortBody(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.Properties["owner"] != "platform-team" {
		t.Fatalf("expected properties owner=platform-team, got %v", group.Properties)
	}
}
