package blueprint

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePrivateState map[string][]byte

func (f fakePrivateState) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f[key], nil
}

func (f fakePrivateState) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if len(value) == 0 {
		delete(f, key)
		return nil
	}
	f[key] = value
	return nil
}

func blueprintWithRelations(identifiers ...string) *cli.Blueprint {
	b := &cli.Blueprint{Identifier: "microservice", Relations: map[string]cli.Relation{}}
	for _, id := range identifiers {
		target := id
		b.Relations[id] = cli.Relation{Target: &target}
	}
	return b
}

func relationKeys(b *cli.Blueprint) []string {
	keys := make([]string, 0, len(b.Relations))
	for k := range b.Relations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func recordedManaged(t *testing.T, p fakePrivateState) []string {
	t.Helper()
	var got []string
	require.NoError(t, json.Unmarshal(p[managedRelationsKey], &got))
	return got
}

// A normal Read with a subset of relations in state must not adopt the others.
func TestRetainManagedRelations_ReadDoesNotAdoptUnmanaged(t *testing.T) {
	req := fakePrivateState{managedRelationsKey: []byte(`["environment"]`)}
	resp := fakePrivateState{}
	b := blueprintWithRelations("environment", "owner_team")

	diags := retainManagedRelations(context.Background(), req, resp, b, []string{"environment"})
	require.False(t, diags.HasError())

	assert.Equal(t, []string{"environment"}, relationKeys(b))
	assert.Equal(t, []string{"environment"}, recordedManaged(t, resp))
}

// I4 regression: a blueprint whose state has no relations and no recorded private state (e.g. state written by an older
// provider version, or a blueprint whose relations all live in port_blueprint_relation resources) must not adopt the
// relations managed elsewhere; adopting them would delete them on the next apply.
func TestRetainManagedRelations_ReadWithoutPrivateStateDoesNotAdopt(t *testing.T) {
	resp := fakePrivateState{}
	b := blueprintWithRelations("owner_team")

	diags := retainManagedRelations(context.Background(), fakePrivateState{}, resp, b, nil)
	require.False(t, diags.HasError())

	assert.Empty(t, b.Relations)
	assert.Empty(t, recordedManaged(t, resp))
}

// Import adopts every relation the blueprint has in Port, and the import marker is consumed so later Reads go back to
// only refreshing managed relations.
func TestRetainManagedRelations_ImportAdoptsAll(t *testing.T) {
	req := fakePrivateState{importAdoptRelationsKey: []byte(`true`)}
	resp := fakePrivateState{importAdoptRelationsKey: []byte(`true`)}
	b := blueprintWithRelations("environment", "owner_team")

	diags := retainManagedRelations(context.Background(), req, resp, b, nil)
	require.False(t, diags.HasError())

	assert.Equal(t, []string{"environment", "owner_team"}, relationKeys(b))
	assert.Equal(t, []string{"environment", "owner_team"}, recordedManaged(t, resp))
	assert.NotContains(t, resp, importAdoptRelationsKey)

	// The next Read (private state carried over from the import Read) keeps the adopted set and adopts nothing new.
	next := blueprintWithRelations("environment", "owner_team", "added_in_ui")
	nextResp := fakePrivateState{}
	diags = retainManagedRelations(context.Background(), resp, nextResp, next, []string{"environment", "owner_team"})
	require.False(t, diags.HasError())
	assert.Equal(t, []string{"environment", "owner_team"}, relationKeys(next))
}

func TestMarkImportForRelationAdoption(t *testing.T) {
	p := fakePrivateState{}
	diags := markImportForRelationAdoption(context.Background(), p)
	require.False(t, diags.HasError())
	assert.Contains(t, p, importAdoptRelationsKey)
}

// Update must keep relations that exist in Port but were never managed by this resource.
func TestMergeUnmanagedRelations_UpdateKeepsUnmanaged(t *testing.T) {
	desired := blueprintWithRelations("environment").Relations
	previouslyManaged := blueprintWithRelations("environment", "removed_from_config").Relations
	existing := blueprintWithRelations("environment", "removed_from_config", "owner_team").Relations

	merged := mergeUnmanagedRelations(desired, previouslyManaged, existing)

	assert.Contains(t, merged, "environment")
	assert.Contains(t, merged, "owner_team", "relation managed elsewhere must be kept")
	assert.NotContains(t, merged, "removed_from_config", "relation removed from this resource's config must be deleted")
}
