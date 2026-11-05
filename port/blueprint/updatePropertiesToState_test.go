package blueprint

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
)

// strPtr and boolPtr are small helpers for constructing cli.Relation/Blueprint fixtures below.
func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

// TestAddRelationsToState_RemovesStaleKeys guards against a regression where addRelationsToState
// only ever ADDED keys found in the live Port blueprint into bm.Relations, and never cleared
// bm.Relations first. That meant a relation key already present in bm (e.g. because a caller
// reconstructs a Terraform "prior state" from the resource's desired configuration before the
// first real Read -- as crossplane-contrib/provider-port style Terraform Plugin Framework
// consumers do when GetObservation() returns an empty/missing status) would survive a refresh
// untouched even though the live Port blueprint (b) no longer -- or does not yet -- have it. That
// silently hides drift: a caller comparing the refreshed "observed" state against its desired
// config sees no difference and concludes an update is not needed, so a relation addition is
// dropped without any error. See the PR description for a concrete reproduction.
func TestAddRelationsToState_RemovesStaleKeys(t *testing.T) {
	bm := &BlueprintModel{
		Relations: map[string]RelationModel{
			"foo": {
				Target: types.StringValue("stale_target"),
			},
			// "bar" simulates a relation present in the caller's desired/prior state (e.g. just
			// added by the user) that Port has not actually persisted yet.
			"bar": {
				Target: types.StringValue("bar_blueprint"),
			},
		},
	}

	// Live Port blueprint only has "foo" -- "bar" was never actually applied.
	b := &cli.Blueprint{
		Relations: map[string]cli.Relation{
			"foo": {
				Target:   strPtr("foo_blueprint"),
				Required: boolPtr(true),
				Many:     boolPtr(false),
			},
		},
	}

	addRelationsToState(b, bm)

	assert.Len(t, bm.Relations, 1, "stale relation key absent from the live blueprint must be removed, not merely left untouched")
	_, hasBar := bm.Relations["bar"]
	assert.False(t, hasBar, "bar does not exist in Port and must not survive the refresh")
	assert.Equal(t, "foo_blueprint", bm.Relations["foo"].Target.ValueString(), "foo must be refreshed from the live blueprint, not left at its stale value")
}

// TestAddRelationsToState_AllRemoved guards the symmetric case: every relation was removed in
// Port. The old implementation was gated behind `if len(b.Relations) > 0` at the call site, so an
// empty b.Relations meant addRelationsToState was never even called and bm.Relations stayed
// stale.
func TestAddRelationsToState_AllRemoved(t *testing.T) {
	bm := &BlueprintModel{
		Relations: map[string]RelationModel{
			"foo": {Target: types.StringValue("foo_blueprint")},
		},
	}
	b := &cli.Blueprint{Relations: map[string]cli.Relation{}}

	addRelationsToState(b, bm)

	assert.Empty(t, bm.Relations, "bm.Relations must be cleared when the live blueprint has no relations left")
}

func TestAddMirrorPropertiesToState_RemovesStaleKeys(t *testing.T) {
	bm := &BlueprintModel{
		MirrorProperties: map[string]MirrorPropertyModel{
			"foo": {Path: types.StringValue("foo.$identifier")},
			"bar": {Path: types.StringValue("no_longer_real.$identifier")},
		},
	}
	b := &cli.Blueprint{
		MirrorProperties: map[string]cli.BlueprintMirrorProperty{
			"foo": {Path: "foo.$identifier"},
		},
	}

	addMirrorPropertiesToState(b, bm)

	assert.Len(t, bm.MirrorProperties, 1)
	_, hasBar := bm.MirrorProperties["bar"]
	assert.False(t, hasBar, "mirror property removed from Port must not survive the refresh")
}

func TestAddCalculationPropertiesToState_RemovesStaleKeys(t *testing.T) {
	bm := &BlueprintModel{
		CalculationProperties: map[string]CalculationPropertyModel{
			"bar": {Calculation: types.StringValue("1+1")},
		},
	}
	b := &cli.Blueprint{
		CalculationProperties: map[string]cli.BlueprintCalculationProperty{},
	}

	addCalculationPropertiesToState(context.Background(), b, bm)

	assert.Empty(t, bm.CalculationProperties, "calculation property removed from Port must not survive the refresh")
}
