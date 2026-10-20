package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

var testSchema = schema.Schema{Attributes: map[string]schema.Attribute{
	"id":         schema.StringAttribute{Computed: true},
	"identifier": schema.StringAttribute{Required: true},
	"blueprint":  schema.StringAttribute{Required: true},
}}

var testObjectType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"id": tftypes.String, "identifier": tftypes.String, "blueprint": tftypes.String,
}}

func strVal(v any) tftypes.Value { return tftypes.NewValue(tftypes.String, v) }

func object(id, identifier, blueprint any) tftypes.Value {
	return tftypes.NewValue(testObjectType, map[string]tftypes.Value{
		"id": strVal(id), "identifier": strVal(identifier), "blueprint": strVal(blueprint),
	})
}

// run plans the "id" attribute with prior state (nil for create) and a planned object whose id is unknown.
func run(t *testing.T, state *tftypes.Value, planIdentifier, planBlueprint any) types.String {
	t.Helper()
	ctx := context.Background()
	m := UseStateForUnknownUnlessChanged(path.Root("identifier"), path.Root("blueprint"))

	planRaw := object(tftypes.UnknownValue, planIdentifier, planBlueprint)
	req := planmodifier.StringRequest{
		Path:        path.Root("id"),
		ConfigValue: types.StringNull(),
		PlanValue:   types.StringUnknown(),
		Plan:        tfsdk.Plan{Schema: testSchema, Raw: planRaw},
		StateValue:  types.StringNull(),
		State:       tfsdk.State{Schema: testSchema, Raw: tftypes.NewValue(testObjectType, nil)},
	}
	if state != nil {
		req.State = tfsdk.State{Schema: testSchema, Raw: *state}
		var id types.String
		req.State.GetAttribute(ctx, path.Root("id"), &id)
		req.StateValue = id
	}
	resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
	m.PlanModifyString(ctx, req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	return resp.PlanValue
}

func TestUseStateForUnknownUnlessChanged_KeepsStateWhenTriggersUnchanged(t *testing.T) {
	state := object("svc:my-service", "my-service", "svc")
	assert.Equal(t, types.StringValue("svc:my-service"), run(t, &state, "my-service", "svc"))
}

func TestUseStateForUnknownUnlessChanged_UnknownWhenIdentifierChanges(t *testing.T) {
	state := object("svc:my-service", "my-service", "svc")
	assert.True(t, run(t, &state, "my-service-2", "svc").IsUnknown())
}

func TestUseStateForUnknownUnlessChanged_UnknownWhenBlueprintChanges(t *testing.T) {
	state := object("svc:my-service", "my-service", "svc")
	assert.True(t, run(t, &state, "my-service", "svc2").IsUnknown())
}

func TestUseStateForUnknownUnlessChanged_UnknownWhenTriggerUnknown(t *testing.T) {
	state := object("svc:my-service", "my-service", "svc")
	assert.True(t, run(t, &state, tftypes.UnknownValue, "svc").IsUnknown())
}

func TestUseStateForUnknownUnlessChanged_UnknownOnCreate(t *testing.T) {
	assert.True(t, run(t, nil, "my-service", "svc").IsUnknown())
}
