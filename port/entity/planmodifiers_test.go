package entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// planComputed runs the plan modifiers of the computed string attribute attrName against a prior state and a plan in
// which only the given string attributes (identifier, blueprint, ...) are set, and returns the planned value.
func planComputed(t *testing.T, attr schema.StringAttribute, attrName string, prior string, priorAttrs, plannedAttrs map[string]string) types.String {
	t.Helper()
	ctx := context.Background()
	attrs := map[string]schema.Attribute{attrName: schema.StringAttribute{Computed: true}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{attrName: tftypes.String}}
	for k := range priorAttrs {
		attrs[k] = schema.StringAttribute{Optional: true}
		objType.AttributeTypes[k] = tftypes.String
	}
	s := schema.Schema{Attributes: attrs}
	build := func(self tftypes.Value, vals map[string]string) tftypes.Value {
		m := map[string]tftypes.Value{attrName: self}
		for k, v := range vals {
			m[k] = tftypes.NewValue(tftypes.String, v)
		}
		return tftypes.NewValue(objType, m)
	}
	req := planmodifier.StringRequest{
		Path:        path.Root(attrName),
		ConfigValue: types.StringNull(),
		PlanValue:   types.StringUnknown(),
		Plan:        tfsdk.Plan{Schema: s, Raw: build(tftypes.NewValue(tftypes.String, tftypes.UnknownValue), plannedAttrs)},
		StateValue:  types.StringValue(prior),
		State:       tfsdk.State{Schema: s, Raw: build(tftypes.NewValue(tftypes.String, prior), priorAttrs)},
	}
	resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
	for _, m := range attr.PlanModifiers {
		m.PlanModifyString(ctx, req, resp)
		req.PlanValue = resp.PlanValue
	}
	if resp.Diagnostics.HasError() {
		t.Fatalf("plan modifier error: %v", resp.Diagnostics)
	}
	return resp.PlanValue
}

func TestEntityIdAndCreatedAtPlan(t *testing.T) {
	s := EntitySchema()
	id := s["id"].(schema.StringAttribute)
	createdAt := s["created_at"].(schema.StringAttribute)
	prior := map[string]string{"identifier": "my-service", "blueprint": "svc"}

	if got := planComputed(t, id, "id", "svc:my-service", prior, prior); got != types.StringValue("svc:my-service") {
		t.Errorf("no change: want prior id kept, got %v", got)
	}
	if got := planComputed(t, createdAt, "created_at", "t0", prior, prior); got != types.StringValue("t0") {
		t.Errorf("no change: want prior created_at kept, got %v", got)
	}
	// TestAccPortEntityUpdateIdentifier: identifier changes, the entity is renamed in place.
	renamed := map[string]string{"identifier": "my-service-2", "blueprint": "svc"}
	if got := planComputed(t, id, "id", "svc:my-service", prior, renamed); !got.IsUnknown() {
		t.Errorf("identifier change: want unknown id, got %v", got)
	}
	// TestAccPortEntityUpdateBlueprintIdentifier: the blueprint changes, the entity is recreated under it.
	moved := map[string]string{"identifier": "my-service", "blueprint": "svc2"}
	if got := planComputed(t, id, "id", "svc:my-service", prior, moved); !got.IsUnknown() {
		t.Errorf("blueprint change: want unknown id, got %v", got)
	}
	if got := planComputed(t, createdAt, "created_at", "t0", prior, moved); !got.IsUnknown() {
		t.Errorf("blueprint change: want unknown created_at, got %v", got)
	}
}
