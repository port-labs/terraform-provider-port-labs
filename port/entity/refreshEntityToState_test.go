package entity

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestRefreshEntityStateClearsStaleCollectionsOnEmptyAPIResponse(t *testing.T) {
	staleRelation := "old-relation"
	state := &EntityModel{
		Identifier: types.StringValue("old"),
		Blueprint:  types.StringValue("old-blueprint"),
		Title:      types.StringValue("old-title"),
		Teams:      []types.String{types.StringValue("old-team")},
		Properties: &EntityPropertiesModel{
			StringProps: map[string]types.String{
				"name": types.StringValue("old-value"),
			},
		},
		Relations: &RelationModel{
			SingleRelation: map[string]*string{
				"owner": &staleRelation,
			},
		},
	}

	apiEntity := &cli.Entity{
		Meta: cli.Meta{
			CreatedAt: ptrTime(time.Now()),
			UpdatedAt: ptrTime(time.Now()),
		},
		Identifier: "new-identifier",
		Title:      "new-title",
		Blueprint:  "new-blueprint",
		Team:       nil,
		Properties: nil,
		Relations:  nil,
	}

	blueprint := &cli.Blueprint{
		Identifier: "new-blueprint",
	}

	resource := &EntityResource{}
	err := resource.refreshEntityState(context.Background(), state, apiEntity, blueprint)
	if err != nil {
		t.Fatalf("refreshEntityState returned error: %v", err)
	}

	if state.Teams != nil {
		t.Fatalf("expected teams to be nil when API returns empty, got: %#v", state.Teams)
	}

	if state.Properties != nil {
		t.Fatalf("expected properties to be nil when API returns empty, got: %#v", state.Properties)
	}

	if state.Relations != nil {
		t.Fatalf("expected relations to be nil when API returns empty, got: %#v", state.Relations)
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
