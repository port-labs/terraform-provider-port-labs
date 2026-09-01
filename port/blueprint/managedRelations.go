package blueprint

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

const managedRelationsKey = "managed_relations"

type privateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

func retainManagedRelations(ctx context.Context, req privateState, resp privateState, b *cli.Blueprint, fallbackManaged []string) diag.Diagnostics {
	var diags diag.Diagnostics

	managed := fallbackManaged

	if req != nil {
		raw, getDiags := req.GetKey(ctx, managedRelationsKey)
		diags.Append(getDiags...)
		if diags.HasError() {
			return diags
		}

		if raw != nil {
			var recorded []string
			if err := json.Unmarshal(raw, &recorded); err == nil {
				managed = recorded
			}
		}
	}

	diags.Append(setManagedRelations(ctx, resp, managed)...)
	if diags.HasError() {
		return diags
	}

	owned := make(map[string]struct{}, len(managed))
	for _, identifier := range managed {
		owned[identifier] = struct{}{}
	}
	for identifier := range b.Relations {
		if _, ok := owned[identifier]; !ok {
			delete(b.Relations, identifier)
		}
	}

	return diags
}

func setManagedRelations(ctx context.Context, private privateState, identifiers []string) diag.Diagnostics {
	var diags diag.Diagnostics

	if private == nil {
		return diags
	}

	sorted := append([]string(nil), identifiers...)
	sort.Strings(sorted)

	raw, err := json.Marshal(sorted)
	if err != nil {
		diags.AddError("failed to record managed relations", err.Error())
		return diags
	}

	diags.Append(private.SetKey(ctx, managedRelationsKey, raw)...)

	return diags
}

func setManagedRelationsFromState(ctx context.Context, private privateState, state *BlueprintModel) diag.Diagnostics {
	identifiers := make([]string, 0, len(state.Relations))
	for identifier := range state.Relations {
		identifiers = append(identifiers, identifier)
	}
	return setManagedRelations(ctx, private, identifiers)
}

func stateRelationIdentifiers(relations map[string]RelationModel) []string {
	identifiers := make([]string, 0, len(relations))
	for identifier := range relations {
		identifiers = append(identifiers, identifier)
	}
	sort.Strings(identifiers)
	return identifiers
}

func mergeUnmanagedRelations(desired map[string]cli.Relation, previouslyManaged map[string]cli.Relation, existing map[string]cli.Relation) map[string]cli.Relation {
	merged := make(map[string]cli.Relation, len(existing)+len(desired))

	for identifier, relation := range existing {
		if _, wasManaged := previouslyManaged[identifier]; wasManaged {
			continue
		}
		merged[identifier] = relation
	}

	for identifier, relation := range desired {
		merged[identifier] = relation
	}

	return merged
}
