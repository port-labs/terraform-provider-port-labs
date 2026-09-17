package scorecard_group

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestScorecardGroupResourceToPortBodyExtendedFields(t *testing.T) {
	t.Parallel()

	state := &ScorecardGroupModel{
		Identifier:          types.StringValue("my-group"),
		Title:               types.StringValue("My Group"),
		GroupProperties:     types.StringValue(`{"category":"production"}`),
		ScorecardProperties: types.StringValue(`{"owner":"platform-team"}`),
		GroupRelations:      types.StringValue(`{"owner_team":"team-123"}`),
		ScorecardRelations:  types.StringValue(`{"stakeholder":["entity-a","entity-b"]}`),
		Blueprints:          []types.String{types.StringValue("svc")},
		Rules: []scorecard.Rule{
			{
				Identifier: types.StringValue("rule-1"),
				Title:      types.StringValue("Rule 1"),
				Level:      types.StringValue("Gold"),
				Query: &scorecard.Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{types.StringValue(`{"property":"x","operator":"isNotEmpty"}`)},
				},
			},
		},
	}

	group, err := scorecardGroupResourceToPortBody(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.GroupProperties["category"] != "production" {
		t.Fatalf("unexpected groupProperties: %v", group.GroupProperties)
	}
	if group.ScorecardProperties["owner"] != "platform-team" {
		t.Fatalf("unexpected scorecardProperties: %v", group.ScorecardProperties)
	}
	if group.GroupRelations["owner_team"] != "team-123" {
		t.Fatalf("unexpected groupRelations: %v", group.GroupRelations)
	}
	stakeholders, ok := group.ScorecardRelations["stakeholder"].([]interface{})
	if !ok || len(stakeholders) != 2 {
		t.Fatalf("unexpected scorecardRelations: %v", group.ScorecardRelations)
	}

	raw, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	payload := string(raw)
	for _, key := range []string{`"groupProperties"`, `"scorecardProperties"`, `"groupRelations"`, `"scorecardRelations"`} {
		if !strings.Contains(payload, key) {
			t.Fatalf("expected %s in JSON payload: %s", key, payload)
		}
	}
	if strings.Contains(payload, `"properties"`) {
		t.Fatalf("legacy properties field should not be present: %s", payload)
	}
}
