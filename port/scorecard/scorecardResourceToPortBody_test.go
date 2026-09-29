package scorecard

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestScorecardResourceToPortBodyPropertiesAndRelations(t *testing.T) {
	state := &ScorecardModel{
		Identifier: types.StringValue("readiness"),
		Title:      types.StringValue("Readiness"),
		Properties: types.StringValue(`{"owner":"platform-team","priority":1}`),
		Relations:  types.StringValue(`{"owner_team":"platform-team"}`),
		Rules: []Rule{
			{
				Identifier: types.StringValue("has-owner"),
				Title:      types.StringValue("Has Owner"),
				Level:      types.StringValue("Gold"),
				Query: &Query{
					Combinator: types.StringValue("and"),
					Conditions: []types.String{
						types.StringValue(`{"property":"$team","operator":"isNotEmpty"}`),
					},
				},
			},
		},
	}

	body, err := scorecardResourceToPortBody(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body.Properties["owner"] != "platform-team" {
		t.Fatalf("expected owner property, got %#v", body.Properties["owner"])
	}
	if body.Properties["priority"] != float64(1) {
		t.Fatalf("expected priority property, got %#v", body.Properties["priority"])
	}
	if body.Relations["owner_team"] != "platform-team" {
		t.Fatalf("expected owner_team relation, got %#v", body.Relations["owner_team"])
	}
}
