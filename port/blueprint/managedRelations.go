package blueprint

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// managedRelationsKey names the private state entry holding the relation identifiers this
// port_blueprint resource owns.
//
// Relations can also be managed one at a time by port_blueprint_relation. Without a record of
// which ones belong here, Read would pull those foreign relations into state and plan their
// removal on every run, and Update would delete them from the blueprint. Terraform state has no
// room for that distinction, so it is kept in private state.
const managedRelationsKey = "managed_relations"

type privateStateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateStateWriter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// resolveManagedRelations reports which relations this resource owns.
//
// When the private state entry is missing - state written by an older provider version, or a
// freshly imported resource - it falls back to every relation on the blueprint, which is exactly
// what the provider did before this entry existed. Read seeds the entry from that fallback, so
// the first refresh after upgrading behaves identically to previous versions and later refreshes
// have a baseline to filter against.
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
			// unreadable entry, fall through and rebuild it from the blueprint
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

// filterRelations drops relations this resource does not own, so they never reach Terraform state.
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

// mergeUnmanagedRelations builds the relations map to write to Port.
//
// Blueprint writes go through PUT, which replaces the blueprint wholesale and deletes any relation
// absent from the body. Relations owned by port_blueprint_relation must survive that, while
// relations dropped from this resource's configuration must still be deleted.
func mergeUnmanagedRelations(desired map[string]cli.Relation, previouslyManaged map[string]cli.Relation, existing map[string]cli.Relation) map[string]cli.Relation {
	merged := make(map[string]cli.Relation, len(existing)+len(desired))

	for identifier, relation := range existing {
		if _, wasManaged := previouslyManaged[identifier]; wasManaged {
			// owned here before; keep it only if it is still configured
			continue
		}
		merged[identifier] = relation
	}

	for identifier, relation := range desired {
		merged[identifier] = relation
	}

	return merged
}
