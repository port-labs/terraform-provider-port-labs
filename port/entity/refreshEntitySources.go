package entity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func encodeEntitySources(ctx context.Context, sources map[string]any) (types.Map, error) {
	if len(sources) == 0 {
		return types.MapNull(types.StringType), nil
	}

	encodedSources := make(map[string]string, len(sources))
	for identifier, sourceValues := range sources {
		encoded, err := json.Marshal(sourceValues)
		if err != nil {
			return types.MapNull(types.StringType), fmt.Errorf("failed to encode sources for %q: %w", identifier, err)
		}
		encodedSources[identifier] = string(encoded)
	}

	encodedState, diags := types.MapValueFrom(ctx, types.StringType, encodedSources)
	if diags.HasError() {
		return types.MapNull(types.StringType), fmt.Errorf("failed to convert sources to state: %s", diags.Errors()[0].Summary())
	}

	return encodedState, nil
}

func refreshRelationSourcesEntityState(ctx context.Context, state *EntityModel, e *cli.Entity) error {
	sources, err := encodeEntitySources(ctx, e.RelationSources)
	if err != nil {
		return err
	}
	state.RelationSources = sources
	return nil
}

func refreshTeamSourcesEntityState(ctx context.Context, state *EntityModel, e *cli.Entity) error {
	sources, err := encodeEntitySources(ctx, e.TeamSources)
	if err != nil {
		return err
	}
	state.TeamSources = sources
	return nil
}

func refreshTeamsEntityState(state *EntityModel, e *cli.Entity, blueprint *cli.Blueprint) {
	if blueprintHasUnionTeamOwnership(blueprint) || cli.EntityTeamIsUnionSlice(e.Team) || len(e.TeamSources) > 0 {
		state.Teams = nil
		return
	}

	teamIdentifiers := cli.EntityTeamIdentifiers(e.Team)
	if len(teamIdentifiers) == 0 {
		state.Teams = nil
		return
	}

	state.Teams = make([]types.String, len(teamIdentifiers))
	for i, teamID := range teamIdentifiers {
		state.Teams[i] = types.StringValue(teamID)
	}
}

func preserveUnionManyRelationSlices(state *EntityModel) map[string]UnionManyRelationSliceModel {
	if state.Relations == nil {
		return nil
	}

	return state.Relations.UnionManyRelations
}

func preserveUnionTeamSlice(state *EntityModel) *UnionTeamSliceModel {
	return state.UnionTeamSlice
}

func restoreUnionTeamSlice(state *EntityModel, unionTeamSlice *UnionTeamSliceModel) {
	state.UnionTeamSlice = unionTeamSlice
}
