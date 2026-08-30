package blueprint

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

const managedRelationsKey = "managed_relations"

type privateStateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateStateWriter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

func resolveManagedRelations(ctx context.Context, private privateStateReader, b *cli.Blueprint) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if private != nil {
		raw, getDiags := private.GetKey(ctx, managedRelationsKey)
		diags.Append(getDiags...)
		if diags.HasError() {
			return nil, diags
		}

		if raw != nil {
			var managed []string
			if err := json.Unmarshal(raw, &managed); err == nil {
				return managed, diags
			}
		}
	}

	managed := make([]string, 0, len(b.Relations))
	for identifier := range b.Relations {
		managed = append(managed, identifier)
	}
	sort.Strings(managed)

	return managed, diags
}

func setManagedRelations(ctx context.Context, private privateStateWriter, identifiers []string) diag.Diagnostics {
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

func setManagedRelationsFromState(ctx context.Context, private privateStateWriter, state *BlueprintModel) diag.Diagnostics {
	identifiers := make([]string, 0, len(state.Relations))
	for identifier := range state.Relations {
		identifiers = append(identifiers, identifier)
	}
	return setManagedRelations(ctx, private, identifiers)
}

func filterRelations(relations map[string]cli.Relation, managed []string) map[string]cli.Relation {
	owned := make(map[string]struct{}, len(managed))
	for _, identifier := range managed {
		owned[identifier] = struct{}{}
	}

	filtered := make(map[string]cli.Relation, len(relations))
	for identifier, relation := range relations {
		if _, ok := owned[identifier]; ok {
			filtered[identifier] = relation
		}
	}

	return filtered
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
