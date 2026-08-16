package scorecard_group

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestValidateScorecardGroupConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		state   *ScorecardGroupModel
		wantErr bool
	}{
		{
			name: "shared rules mode",
			state: &ScorecardGroupModel{
				Blueprints: []types.String{types.StringValue("bp-1")},
				Rules: []scorecard.Rule{
					{Identifier: types.StringValue("rule-1")},
				},
			},
		},
		{
			name: "per blueprint mode",
			state: &ScorecardGroupModel{
				Scorecards: map[string]MemberSpecModel{
					"bp-1": {
						Rules: []scorecard.Rule{
							{Identifier: types.StringValue("rule-1")},
						},
					},
				},
			},
		},
		{
			name: "missing configuration",
			state: &ScorecardGroupModel{
				Identifier: types.StringValue("group-1"),
				Title:      types.StringValue("Group 1"),
			},
			wantErr: true,
		},
		{
			name: "shared rules missing rules",
			state: &ScorecardGroupModel{
				Blueprints: []types.String{types.StringValue("bp-1")},
			},
			wantErr: true,
		},
		{
			name: "mixed modes",
			state: &ScorecardGroupModel{
				Blueprints: []types.String{types.StringValue("bp-1")},
				Rules: []scorecard.Rule{
					{Identifier: types.StringValue("rule-1")},
				},
				Scorecards: map[string]MemberSpecModel{
					"bp-2": {},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateScorecardGroupConfiguration(tt.state)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateScorecardGroupConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
